package chat

import (
	"fmt"
	"sport/helpers"
	"sport/lang"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

func NewClusterConfig(config *CassandraConfig) *gocql.ClusterConfig {
	cluster := gocql.NewCluster(config.Hosts...)
	cluster.Consistency = gocql.Consistency(config.Consistency)
	cluster.ProtoVersion = config.ProtoVersion
	cluster.ConnectTimeout = time.Second * 5
	cluster.Timeout = time.Second * 2
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.Username,
		Password: config.Password,
	}
	if config.CaPath != "" {
		cluster.SslOpts = &gocql.SslOptions{
			EnableHostVerification: true,
			CaPath:                 config.CaPath,
		}
	}
	// cluster.Keyspace = config.Keyspace
	return cluster
}

func (ctx *Ctx) DeployDdl(
	config *gocql.ClusterConfig, keyspace string, forceRecreateDdl bool) error {

	var conn *gocql.Session
	var err error

	retries := 15
	for i := 0; i < retries; i++ {

		conn, err = gocql.NewSession(*config)
		if err == nil {
			break
		}

		time.Sleep(time.Second * 6)
	}

	if conn == nil {
		return fmt.Errorf("failed to establish cass connection")
	}

	defer conn.Close()

	if forceRecreateDdl {
		if err := conn.Query(fmt.Sprintf("drop keyspace if exists %s", keyspace)).Exec(); err != nil {
			return err
		}
	}

	for i := range ctx.Config.Cassandra.Ddl {
		q := ctx.Config.Cassandra.Ddl[i]
		err := conn.Query(fmt.Sprintf(q, keyspace)).Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

type ChatMemberCol int

const (
	Email        ChatMemberCol = 0x1
	DisplayName  ChatMemberCol = 0x2
	LastNotified ChatMemberCol = 0x4
	ServerID     ChatMemberCol = 0x8
	ChatRoomName ChatMemberCol = 0x10
	IconRelpath  ChatMemberCol = 0x20
	LastReadMsg  ChatMemberCol = 0x40
	Language     ChatMemberCol = 0x80

	All ChatMemberCol = (1 << 8) - 1
)

func (ctx *Ctx) DalCreateChatRoomMember(cr *ChatRoomMember, mask ChatMemberCol) error {

	var qb strings.Builder

	if mask == 0 {
		return fmt.Errorf("invalid mask")
	}

	var args = make([]interface{}, 0, 7)

	qb.WriteString(`insert into chat_room_members (
		chat_room_id, user_id `)

	args = append(args, cr.ChatRoomID[:], cr.UserID[:])

	if mask&Email != 0 {
		qb.WriteString(`,email`)
		args = append(args, cr.Data.Email)
	}

	if mask&DisplayName != 0 {
		qb.WriteString(`,display_name`)
		args = append(args, cr.Data.DisplayName)
	}

	if mask&LastNotified != 0 {
		qb.WriteString(`,last_notified`)
		args = append(args, cr.LastNotified)
	}

	if mask&ServerID != 0 {
		qb.WriteString(`,server_id`)
		args = append(args, cr.ServerID)
	}

	if mask&ChatRoomName != 0 {
		qb.WriteString(`,chat_room_name`)
		args = append(args, cr.Data.ChatRoomName)
	}

	if mask&IconRelpath != 0 {
		qb.WriteString(`,icon_relpath`)
		args = append(args, cr.Data.IconRelpath)
	}

	if mask&LastReadMsg != 0 {
		qb.WriteString(`,last_read_msg`)
		args = append(args, cr.LastReadMsg)
	}

	if mask&Language != 0 {
		qb.WriteString(`,language`)
		args = append(args, cr.Data.Language)
	}

	qb.WriteString(`) values (?,?`)

	for mask > 0 {
		if mask%2 != 0 {
			qb.WriteString(",?")
		}
		mask /= 2
	}

	qb.WriteString(`)`)

	return ctx.cass.Query(qb.String(), args...).Exec()
}

func (ctx *Ctx) DalReadChatRoomMembers(
	chatRoomID uuid.UUID, userID uuid.UUID) ([]ChatRoomMember, error) {

	if chatRoomID == uuid.Nil && userID == uuid.Nil {
		return nil, fmt.Errorf("scanning chat_room_members not supported")
	}

	var qb strings.Builder

	qb.WriteString(`select 
		chat_room_id,
		user_id,
		email,
		display_name,
		last_notified,
		last_read_msg,
		server_id,
		chat_room_name,
		icon_relpath,
		language
	from chat_room_members where `)

	args := make([]interface{}, 0, 3)

	if chatRoomID != uuid.Nil {
		qb.WriteString(` chat_room_id = ? `)
		args = append(args, chatRoomID[:])
	}

	if userID != uuid.Nil {
		if chatRoomID != uuid.Nil {
			qb.WriteString(` and `)
		}
		qb.WriteString(` user_id = ? `)
		args = append(args, userID[:])
	}

	if chatRoomID != uuid.Nil {
		qb.WriteString(` order by user_id`)
	}

	scanner := ctx.cass.Query(qb.String(), args...).Iter().Scanner()
	res := make([]ChatRoomMember, 0, 3)

	for scanner.Next() {
		res = append(res, ChatRoomMember{})
		i := len(res) - 1
		if err := scanner.Scan(
			(*[16]byte)(&res[i].ChatRoomID),
			(*[16]byte)(&res[i].UserID),
			&res[i].Data.Email,
			&res[i].Data.DisplayName,
			&res[i].LastNotified,
			&res[i].LastReadMsg,
			&res[i].ServerID,
			&res[i].Data.ChatRoomName,
			&res[i].Data.IconRelpath,
			&res[i].Data.Language); err != nil {

			_ = scanner.Err()
			return nil, err
		}
	}

	return res, scanner.Err()
}

/* updates last message timestamp for all chatroom members */
func (ctx *Ctx) dalUpdateChatRoomLMT(chatRoomID uuid.UUID, lmt int64) error {
	const q = `update chat_rooms set last_msg_timestamp = ? where chat_room_id = ?`
	return ctx.cass.Query(q, lmt, chatRoomID[:]).Exec()
}

func (ctx *Ctx) dalCreateChatMsg(m *ChatMsg) error {

	const q = `
		insert into messages (chat_room_id, timestamp, user_id, content) 
		values (?, ?, ?, ?)`

	return ctx.cass.Query(
		q,
		m.ChatRoomID[:], m.Timestamp, m.UserID[:], m.Content).Exec()
}

func (ctx *Ctx) EnqueueChanNotification(author *ChatRoomMember, m *ChatMsg) {

	if ctx.noReply == nil || ctx.langCtx == nil {
		return
	}

	crid := author.ChatRoomID

	cring := ctx.chanNotificationMap[crid]
	if cring == nil {
		ctx.chanNotificationMapLock.Lock()
		if cring == nil {
			cring = NewNotificationRing(ctx.Config.NotRingSize)
			ctx.chanNotificationMap[crid] = cring
		}
		ctx.chanNotificationMapLock.Unlock()
	}

	ts := MsgTimestampToTime(m.Timestamp)

	not := &NotificationItem{
		Author:       author.Data.DisplayName,
		MsgSummary:   helpers.CutString(m.Content, 64),
		Timestamp:    ts,
		TimestampStr: ts.Format(lang.MailDateFmt),
		AuthorID:     author.UserID,
	}

	cring.Mu.Lock()
	cring.Enqueue(not)
	cring.Mu.Unlock()
}

func (ctx *Ctx) CreateChatMsgNotify(author *ChatRoomMember, m *ChatMsg) error {

	var oerr = [2]error{nil, nil}
	var hasErr bool
	var nc = make(chan struct{})

	go func() {
		if err := ctx.dalCreateChatMsg(m); err != nil {
			oerr[0] = err
			hasErr = true
		}
		nc <- struct{}{}
	}()

	go func() {
		if err := ctx.dalUpdateChatRoomLMT(m.ChatRoomID, m.Timestamp); err != nil {
			oerr[1] = err
			hasErr = true
		}
		nc <- struct{}{}
	}()

	for i := 0; i < 2; i++ {
		<-nc
	}

	if hasErr {
		return fmt.Errorf("CreateChatMsg: %v, UpdateChatRoomLMT: %v", oerr[0], oerr[1])
	}

	ctx.EnqueueChanNotification(author, m)

	return nil
}

func (ctx *Ctx) DalReadChatMsgs(
	chatRoomID uuid.UUID, feedOpts *MsgFeedOptions, asc bool) ([]ChatMsgData, error) {

	const q = `
		select user_id, timestamp, content 
		from messages
		where chat_room_id = ? `

	var args = make([]interface{}, 0, 2)
	args = append(args, chatRoomID[:])

	var qb strings.Builder
	qb.WriteString(q)

	ord := ""
	if asc {
		ord = " order by timestamp asc "
	} else {
		ord = " order by timestamp desc "
	}

	if feedOpts != nil {

		start := feedOpts.Start
		end := feedOpts.End

		if start != 0 {
			qb.WriteString(" and timestamp >= ? ")
			args = append(args, start)
		}

		if end != 0 {
			qb.WriteString(" and timestamp <= ? ")
			args = append(args, end)
		}

		qb.WriteString(ord)

		qb.WriteString(" limit ? ")
		args = append(args, feedOpts.Limit)

	} else {
		qb.WriteString(ord)
	}

	scanner := ctx.cass.Query(qb.String(), args...).Iter().Scanner()

	var res = make([]ChatMsgData, 0, 4)
	for scanner.Next() {

		res = append(res, ChatMsgData{})
		i := len(res) - 1

		if err := scanner.Scan(
			(*[16]byte)(&res[i].UserID),
			&res[i].Timestamp,
			&res[i].Content); err != nil {

			_ = scanner.Err()
			return nil, err
		}
	}

	return res, scanner.Err()
}

func (ctx *Ctx) DalCreateChatRoom(cr *ChatRoom) error {
	const q = `
		insert into chat_rooms (chat_room_id, flags, last_msg_timestamp) 
		values (?,?,?)
	`
	return ctx.cass.Query(q, cr.ChatRoomID[:], cr.Flags, cr.LastMsgTimestmap).Exec()
}

func (ctx *Ctx) CreateChatRoomWithNotify(cr *ChatRoom, crm *ChatRoomMember) error {
	if err := ctx.DalCreateChatRoom(cr); err != nil {
		return err
	}
	ctx.SendEmailAboutNewChanAsync(crm)
	return nil
}

func (ctx *Ctx) DalDeleteChatRoom(chatRoomID uuid.UUID) error {
	return ctx.cass.Query(
		`delete from chat_rooms where id = ?`, chatRoomID[:]).Exec()
}

func (ctx *Ctx) DalReadChatRooms(chatRoomIDs []uuid.UUID) ([]ChatRoom, error) {
	if len(chatRoomIDs) == 0 {
		return nil, fmt.Errorf("invalid chatRoomIDs")
	}
	qb := strings.Builder{}
	qb.WriteString(`select 
		chat_room_id, flags, last_msg_timestamp 
		from chat_rooms where chat_room_id in (`)

	args := make([]interface{}, 0, 2)
	for i := range chatRoomIDs {
		if i == 0 {
			qb.WriteString("?")
		} else {
			qb.WriteString(",?")
		}
		args = append(args, chatRoomIDs[i][:])
	}

	qb.WriteString(")")
	scanner := ctx.cass.Query(qb.String(), args...).Iter().Scanner()

	var res = make([]ChatRoom, 0, 4)

	for scanner.Next() {
		res = append(res, ChatRoom{})
		i := len(res) - 1
		if err := scanner.Scan(
			(*[16]byte)(&res[i].ChatRoomID),
			&res[i].Flags,
			&res[i].LastMsgTimestmap,
		); err != nil {
			_ = scanner.Err()
			return nil, err
		}
	}

	return res, scanner.Err()
}

func (ctx *Ctx) DalReadChatRoom(chatRoomID uuid.UUID) (*ChatRoom, error) {
	res, err := ctx.DalReadChatRooms([]uuid.UUID{chatRoomID})
	if err != nil {
		return nil, err
	}
	if len(res) != 1 {
		return nil, fmt.Errorf("invalid number of rooms read")
	}
	return &res[0], nil
}

func (ctx *Ctx) DalCreateAccessToken(it *AccessToken) error {
	const q = `
		insert into access_tokens (
			chat_room_id,
			creator_id,
			token_value,
			user_id,
			expires_on
		) values (
			?, ?, ?, ?, ?
		)`
	return ctx.cass.Query(q,
		it.ChatRoomID[:],
		it.CreatorID[:],
		it.TokenValue[:],
		it.UserID[:],
		it.ExpiresOn).Exec()
}

func (ctx *Ctx) DalDeleteAccessToken(
	chatRoomID, token uuid.UUID) error {

	const q = `
		delete from access_tokens
		where chat_room_id = ? and token_value = ?`

	return ctx.cass.Query(q, chatRoomID[:], token[:]).Exec()
}

func (ctx *Ctx) DalUpdateAccessTokenUser(chatRoomID, token, userID uuid.UUID) error {
	const q = `
		update access_tokens set
			user_id = ?
		where chat_room_id = ? and token_value = ?
	`
	return ctx.cass.Query(q, userID[:], chatRoomID[:], token[:]).Exec()
}

func (ctx *Ctx) DalReadAccessTokens(
	chatRoomID, token uuid.UUID,
) ([]AccessToken, error) {

	if chatRoomID == uuid.Nil {
		return nil, fmt.Errorf("invalid chatRoomID")
	}

	qb := strings.Builder{}

	qb.WriteString(`
		select 
			chat_room_id,
			creator_id,
			token_value,
			user_id,
			expires_on
		from access_tokens 
		where chat_room_id = ? `)

	args := make([]interface{}, 0, 3)
	args = append(args, chatRoomID[:])

	if token != uuid.Nil {
		qb.WriteString("and token_value = ? ")
		args = append(args, token[:])
	}

	var res = make([]AccessToken, 0, 2)

	var scanner = ctx.cass.Query(qb.String(), args...).Iter().Scanner()

	for scanner.Next() {
		res = append(res, AccessToken{})
		i := len(res) - 1
		if err := scanner.Scan(
			(*[16]byte)(&res[i].ChatRoomID),
			(*[16]byte)(&res[i].CreatorID),
			(*[16]byte)(&res[i].TokenValue),
			(*[16]byte)(&res[i].UserID),
			&res[i].ExpiresOn,
		); err != nil {
			_ = scanner.Err()
			return nil, err
		}
	}

	return res, scanner.Err()
}
