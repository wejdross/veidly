package adyen_sm

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type StateChangeEventSourceType string

const (
	Agent   StateChangeEventSourceType = "rsv_daemon"
	Webhook StateChangeEventSourceType = "rsv_webhook"
	Handler StateChangeEventSourceType = "rsv_handler"
	Manual  StateChangeEventSourceType = "rsv_manual"
)

type StateChangeEventSource struct {
	Type StateChangeEventSourceType
	Data string
}

var ManualSrc = StateChangeEventSource{Manual, ""}

type StateChangeEvent struct {
	ID                       uuid.UUID
	ObjKey                   uuid.UUID
	PreviousState, NextState string
	Timestamp                time.Time
	Success                  bool
	Error                    string
	Source                   StateChangeEventSource
}

func (ctx *Ctx) LogEventToFs(obj interface{}, id uuid.UUID, timestamp time.Time) {

	j, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		if _, err := fmt.Fprint(os.Stderr, err); err != nil {
			panic(err)
		}
		return
	}

	dir := path.Join(ctx.Config.SmLogDir, id.String())
	if err := os.MkdirAll(dir, 0700); err != nil {
		if _, err := fmt.Fprint(os.Stderr, err); err != nil {
			panic(err)
		}
		return
	}

	path := path.Join(dir, strconv.FormatInt(timestamp.Unix(), 10))
	if err := os.WriteFile(path, j, 0600); err != nil {
		if _, err := fmt.Fprint(os.Stderr, err); err != nil {
			panic(err)
		}
		return
	}
}

func (sce *StateChangeEvent) LogWithError(ctx *Ctx, err error, tx *sql.Tx) {
	sce.Success = false
	sce.Error = err.Error()
	ctx.LogStateChangeEvent(sce, nil)
}

func (sce *StateChangeEvent) Log(ctx *Ctx, tx *sql.Tx) {
	ctx.LogStateChangeEvent(sce, tx)
}

/*
	tries to log state into database,
		+ if dst state is error or dispute then emails message to support
	if failed logs into file
	if that failed then will print error to the stdout
	if that fails then panic
*/
func (ctx *Ctx) LogStateChangeEvent(event *StateChangeEvent, tx *sql.Tx) {

	if ctx.noReplyCtx != nil {
		if event.Success {
			if ctx.Config.NotifyOn[string(event.NextState)] {
				go ctx.MailSupportAboutRsvStateChanged(event)
			}
		} else {
			go ctx.MailSupportAboutRsvStateChanged(event)
		}
	}

	q := `insert into sm_log (
		id, 
		obj_key, 
		previous_state, 
		next_state, 
		timestamp,
		success,
		source,
		error,
		source_data,
		nstamp
	) values (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10
	)`
	params := []interface{}{
		event.ID,
		event.ObjKey,
		event.PreviousState,
		event.NextState,
		event.Timestamp,
		event.Success,
		event.Source.Type,
		event.Error,
		event.Source.Data,
		time.Now().UnixNano(),
	}
	var err error
	if tx == nil {
		_, err = ctx.dal.Db.Exec(q, params...)
	} else {
		_, err = tx.Exec(q, params...)
	}

	if err == nil {
		return
	}

	if _, err := fmt.Fprint(os.Stderr, err); err != nil {
		panic(err)
	}

	ctx.LogEventToFs(event, event.ID, event.Timestamp)
}
