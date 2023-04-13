package chat

import (
	"fmt"
	"sport/api"
	"sport/helpers"
	"sport/lang"
	"sport/notify"
	"strings"
	"sync"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type CassandraConfig struct {
	Hosts        []string
	Username     string
	Password     string
	Consistency  int
	ProtoVersion int      `yaml:"proto_version"`
	Ddl          []string `yaml:"ddl"`
	Keyspace     string
	CaPath       string `yaml:"ca_path"`
}

func (c *CassandraConfig) Validate() error {
	const _fmt = "validate CassandraConfig: %s"

	if len(c.Hosts) == 0 {
		return fmt.Errorf(_fmt, "invalid hosts")
	}

	if c.Keyspace == "" {
		return fmt.Errorf(_fmt, "invalid keyspace")
	}

	return nil
}

type Config struct {
	Cassandra        CassandraConfig
	UiJoinChatUrlFmt string           `yaml:"ui_join_chat_url_fmt"`
	ServerID         string           `yaml:"server_id"`
	OpenUrlFmt       string           `yaml:"open_url_fmt"`
	OpenPathFmt      string           `yaml:"open_path_fmt"`
	NotRingSize      int              `yaml:"not_ring_size"`
	EmailNotAfter    helpers.Duration `yaml:"email_not_after"`
	Jwt              api.JwtRequest
}

func (c *Config) Validate() error {
	const _fmt = "validate Config: %s"

	if c.UiJoinChatUrlFmt == "" {
		return fmt.Errorf(_fmt, "invalid ApiPubUrl")
	}

	if c.ServerID == "" {
		return fmt.Errorf(_fmt, "invalid ServerID")
	}

	if c.OpenUrlFmt == "" {
		return fmt.Errorf(_fmt, "invalid ServerUrlFmt")
	}

	if c.OpenPathFmt == "" {
		return fmt.Errorf(_fmt, "invalid OpenPathFmt")
	}

	if err := c.Jwt.Validate(); err != nil {
		return fmt.Errorf(_fmt, "invalid chat jwt")
	}

	if c.NotRingSize <= 0 {
		return fmt.Errorf(_fmt, "invalid not_ring_size")
	}

	if c.EmailNotAfter <= 0 {
		return fmt.Errorf(_fmt, "invalid email_not_after")
	}

	return c.Cassandra.Validate()
}

type Ctx struct {
	Config *Config
	cass   *gocql.Session
	jwt    *api.Jwt
	//
	wsConnMapLock sync.RWMutex
	wsConnMap     WssConnMap

	chanNotificationMap     map[uuid.UUID]*NotificationRing
	chanNotificationMapLock sync.Mutex

	// map[chatRoomID]map[UserID]
	// allMembersMap     map[uuid.UUID]map[uuid.UUID]*ChatRoomMember
	// allMembersMapLock sync.RWMutex

	langCtx *lang.Ctx
	noReply notify.EmailSender

	wsTokenCache *WsTokenCache
}

func (ctx *Ctx) GetOpenPath() string {
	openPath := ""
	if strings.Contains(ctx.Config.OpenPathFmt, "%s") {
		openPath = fmt.Sprintf(ctx.Config.OpenPathFmt, ctx.Config.ServerID)
	} else {
		openPath = ctx.Config.OpenPathFmt
	}
	return openPath
}

func NewCtx(
	apiCtx *api.Ctx,
	langCtx *lang.Ctx,
	noReply notify.EmailSender,
	overrideKeyspace string,
	forceRecreateDdl bool) *Ctx {

	var err error

	_ = *apiCtx
	_ = *langCtx

	ctx := new(Ctx)

	ctx.langCtx = langCtx
	ctx.noReply = noReply
	ctx.Config = new(Config)
	apiCtx.Config.UnmarshalKeyPanic("chat", ctx.Config, ctx.Config.Validate)
	ctx.chanNotificationMap = make(map[uuid.UUID]*NotificationRing)
	ctx.wsTokenCache = NewWsTokenCache()
	// ctx.allMembersMap = make(map[uuid.UUID]map[uuid.UUID]*ChatRoomMember)

	if overrideKeyspace != "" {
		ctx.Config.Cassandra.Keyspace = overrideKeyspace
	}

	clusterConfig := NewClusterConfig(&ctx.Config.Cassandra)

	if err := ctx.DeployDdl(clusterConfig, ctx.Config.Cassandra.Keyspace, forceRecreateDdl); err != nil {
		panic(err)
	}

	clusterConfig.Keyspace = ctx.Config.Cassandra.Keyspace
	ctx.cass, err = gocql.NewSession(*clusterConfig)
	if err != nil {
		panic(err)
	}

	ctx.wsConnMap = make(WssConnMap)

	ctx.jwt, err = ctx.Config.Jwt.NewJwt()
	if err != nil {
		panic(err)
	}

	apiCtx.AnonGroup.GET("/chat/token/validate", ctx.ValidateTokenHandler())
	apiCtx.AnonGroup.POST("/chat/room", ctx.CreateChatRoomHandler())
	apiCtx.AnonGroup.POST("/chat/token/ws", ctx.ObtainWsTokenHandler())
	apiCtx.AnonGroup.GET(ctx.GetOpenPath(), ctx.OpenChatRoomHandler())
	apiCtx.AnonGroup.GET("/chat/room", ctx.GetChatRoomsHandler())
	apiCtx.AnonGroup.POST("/chat/room/join", ctx.JoinChatRoomHandler())
	apiCtx.AnonGroup.GET("/chat/room/access_token", ctx.GetAccessTokensHandler())
	apiCtx.AnonGroup.POST("/chat/room/access_token", ctx.CreateAccessTokenHandler())
	apiCtx.AnonGroup.GET("/chat/notify/open", ctx.OpenNotificationChan())

	go ctx.StartChatroomRunner()

	return ctx
}
