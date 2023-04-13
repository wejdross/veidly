package search

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PgNotificationChannel string

const (
	Trainings        PgNotificationChannel = "trainings"
	Instructors                            = "instructors"
	Occurrences                            = "occurrences"
	InstrVacations                         = "instr_vacations"
	Reservations                           = "reservations"
	SecondaryOccs                          = "secondary_occs"
	SubModels                              = "sub_models"
	SubModelBindings                       = "sub_model_bindings"
	Subs                                   = "subs"
	TrainingGroups                         = "training_groups"
	TrainingsVgroups                       = "trainings_v_groups"
	Users                                  = "users"
	Reviews                                = "reviews"
)

// im using map to avoid duplicates
type SingleTableIDmap map[uuid.UUID]struct{}
type NotificationData map[PgNotificationChannel]SingleTableIDmap

type NotificationCallback func(elems NotificationData) error

// func (ctx *Ctx) NewListener() *pq.Listener {
// 	l := pq.NewListener(ctx.Dal.Cs, time.Second, time.Second, func(event pq.ListenerEventType, err error) {
// 		if err != nil {
// 			fmt.Fprintln(os.Stderr, err)
// 			// TODO: notify admin
// 		}
// 	})
// 	chans := []PgNotificationChannel{
// 		Trainings,
// 		Instructors,
// 		Occurrences,
// 		InstrVacations,
// 		Reservations,
// 		SecondaryOccs,
// 		SubModels,
// 		SubModelBindings,
// 		//Subs,
// 		TagsVtrainings,
// 		TrainingGroups,
// 		TrainingsVgroups,
// 		Users,
// 	}

// 	for i := range chans {
// 		if err := l.Listen(string(chans[i])); err != nil {
// 			fmt.Fprintln(os.Stderr, err)
// 			// TODO: notify admin
// 		}
// 	}

// 	return l
// }

func (d *NotificationData) AddChange(t string, id uuid.UUID) {
	channel := PgNotificationChannel(t)
	cd := (*d)[channel]
	if cd == nil {
		cd = make(SingleTableIDmap)
		cd[id] = struct{}{}
		(*d)[channel] = cd
	} else {
		(*d)[channel][id] = struct{}{}
	}
}

// func ProcessPgNotification(n *pq.Notification, d *NotificationData) error {
// 	id, err := uuid.Parse(n.Extra)
// 	if err != nil {
// 		return err
// 	}
// 	d.AddChange(n.Extra, id)
// 	return nil
// 	// channel := PgNotificationChannel(n.Channel)
// 	// cd := (*d)[channel]
// 	// if cd == nil {
// 	// 	cd = make(SingleTableIDmap)
// 	// 	cd[id] = struct{}{}
// 	// 	(*d)[channel] = cd
// 	// } else {
// 	// 	(*d)[channel][id] = struct{}{}
// 	// }
// 	// return nil
// }

// func (ctx *Ctx) PgListenForChangesCb(nc NotificationCallback) {

// 	l := ctx.NewListener()

// 	q := make(chan bool)

// 	go func() {
// 		for {
// 			time.Sleep(time.Second)
// 			q <- true
// 		}
// 	}()

// 	for {

// 		bf := make(NotificationData)
// 	L:
// 		for {
// 			select {
// 			case n := <-l.Notify:
// 				if err := ProcessPgNotification(n, &bf); err != nil {
// 					continue
// 				}
// 			case <-q:
// 				break L
// 			}
// 		}

// 		if len(bf) == 0 {
// 			continue
// 		}

// 		if err := nc(bf); err != nil {
// 			fmt.Fprintln(os.Stderr, err)
// 			// TODO: notify admin
// 		}
// 	}
// }

// type TestListenerMapKey struct {
// 	Channel PgNotificationChannel
// 	ID      uuid.UUID
// }

// type TestListenerMap map[TestListenerMapKey]int

// func DrainListenerQueue(l *pq.Listener, expectedData TestListenerMap) error {
// 	foundData := make(TestListenerMap)
// F:
// 	for {
// 		select {
// 		case e := <-l.Notify:
// 			id, err := uuid.Parse(e.Extra)
// 			if err != nil {
// 				return err
// 			}
// 			key := TestListenerMapKey{
// 				Channel: PgNotificationChannel(e.Channel),
// 				ID:      id,
// 			}
// 			foundData[key]++
// 			break
// 		case <-time.After(time.Millisecond * 50):
// 			break F
// 		}
// 	}

// 	for x := range expectedData {
// 		if e, found := foundData[x]; found {
// 			if e != expectedData[x] {
// 				return fmt.Errorf("Unexpected repeats of elem (%v). Expected %d, got %d", x, expectedData[x], e)
// 			}
// 		} else {
// 			return fmt.Errorf("Didnt find expected item in the queue: %v.", x)
// 		}
// 	}
// 	for x := range foundData {
// 		if _, found := expectedData[x]; !found {
// 			return fmt.Errorf("Received unexpected notification: %v (%d times)", x, foundData[x])
// 		}
// 	}

// 	return nil
// }

func AssertNotificationData(got, expected NotificationData) error {
	for x := range got {
		if e := expected[x]; e != nil {
			for gotid := range got[x] {
				if _, f := e[gotid]; !f {
					return fmt.Errorf("got unexpected id: %v", gotid)
				}
			}
		} else {
			return fmt.Errorf("got unexpected channel: %v", x)
		}
	}
	for x := range expected {
		if e := got[x]; e != nil {
			for expid := range expected[x] {
				if _, f := e[expid]; !f {
					return fmt.Errorf("expected id: %v not found", expid)
				}
			}
		} else {
			return fmt.Errorf("expected channel: %v not found", x)
		}
	}
	return nil
}

func (ctx *Ctx) GetPgChanges() (NotificationData, error) {
	const q = "select channel, id from table_changes"
	r, err := ctx.Dal.Db.Query(q)
	if err != nil {
		return nil, err
	}
	var t string
	var id uuid.UUID
	nd := make(NotificationData)
	ids := make([]string, 0, 4)
	for r.Next() {
		if err := r.Scan(&t, &id); err != nil {
			return nil, err
		}
		channel := PgNotificationChannel(t)
		channelData := nd[channel]
		added := false
		if channelData == nil {
			channelData = make(SingleTableIDmap)
			channelData[id] = struct{}{}
			nd[channel] = channelData
			added = true
		} else {
			if _, f := nd[channel][id]; !f {
				nd[channel][id] = struct{}{}
				added = true
			}
		}
		if added {
			ids = append(ids, id.String())
		}
		//nd.AddChange(t, id)
	}

	if len(nd) > 0 {
		_, err := ctx.Dal.Db.Query("delete from table_changes where id = any($1) ",
			pq.StringArray(ids))
		if err != nil {
			return nil, err
		}
	}

	return nd, nil
}
