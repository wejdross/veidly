package search

import (
	"fmt"
	"sport/api"
	"sport/dal"
	"sport/instr"
	"sport/rsv"
	"sport/schedule"
	"sport/train"
	"sport/user"
	"sync"
	"time"
)

// config
type Config struct {
	MaxQueryLength  int `yaml:"max_query_length"`
	GoogleTranslate struct {
		Token   string
		Enabled bool
	} `yaml:"google_translate"`
	LinguaApiUrl string `yaml:"lingua_api_url"`
}

func (r *Config) ValidationErr(m string) error {
	return fmt.Errorf("Validate Config: " + m)
}

func (r *Config) Validate() error {
	if r.MaxQueryLength <= 0 {
		return r.ValidationErr("invalid max_query_length")
	}
	if r.GoogleTranslate.Enabled {
		if r.GoogleTranslate.Token == "" {
			return r.ValidationErr("invalid google_translate.token")
		}
	}
	if r.LinguaApiUrl == "" {
		return r.ValidationErr("invalid lingua_api_url")
	}
	return nil
}

type Token struct{}

type Ctx struct {
	Api   *api.Ctx
	Dal   *dal.Ctx
	User  *user.Ctx
	Train *train.Ctx
	Rsv   *rsv.Ctx
	Instr *instr.Ctx
	Sched *schedule.Ctx
	//
	Conf *Config
	//
	//cacheLock               sync.Mutex
	Cache *SearchCache

	RWLock sync.RWMutex

	CacheRefreshInterval    time.Duration
	DbgSlowdownCacheRefresh time.Duration

	IgnoreInvalidInstructors bool
}

func RegisterAllHandlers(ctx *Ctx) {
	ctx.Api.AnonGroup.POST("/search", ctx.SearchHandler())
}

func NewCtx(
	apiCtx *api.Ctx,
	dalCtx *dal.Ctx,
	userCtx *user.Ctx,
	instrCtx *instr.Ctx,
	trainCtx *train.Ctx,
	rsvCtx *rsv.Ctx,
	schedCtx *schedule.Ctx,
) *Ctx {

	_ = *apiCtx
	_ = *dalCtx
	_ = *userCtx
	_ = *instrCtx
	_ = *trainCtx
	_ = *rsvCtx

	ctx := new(Ctx)
	ctx.Conf = new(Config)
	apiCtx.Config.UnmarshalKeyPanic("search", ctx.Conf, ctx.Conf.Validate)
	ctx.Api = apiCtx
	ctx.Dal = dalCtx
	ctx.User = userCtx
	ctx.Train = trainCtx
	ctx.Rsv = rsvCtx
	ctx.Instr = instrCtx
	ctx.Sched = schedCtx

	ctx.CacheRefreshInterval = time.Second

	ctx.RWLock = sync.RWMutex{}

	/* start in locked state until cache gets updated */
	ctx.RWLock.Lock()

	RegisterAllHandlers(ctx)

	return ctx
}
