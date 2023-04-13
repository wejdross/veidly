package search

import (
	"fmt"
	"os"
	"sport/helpers"
	"sport/lang"
	"sport/train"
	"time"
)

type CacheBuildTime struct {
	ClusterBuildTime PrettyJsonDuration
	IndexBuildTime   PrettyJsonDuration
}

type CacheGenerationPerf struct {
	TotalBuildTime    PrettyJsonDuration
	TrainingBuildTime CacheBuildTime
	RsvBuildTime      CacheBuildTime
	InstrBuildTime    CacheBuildTime
}

type TrainingUpdateInternals struct {
	ReadTrainingsTime PrettyJsonDuration
	ReadGroupTime     PrettyJsonDuration
	ReadSmTime        PrettyJsonDuration
	TotalReadTime     PrettyJsonDuration
	MergeTrainingTime PrettyJsonDuration
	RmTrainingTime    PrettyJsonDuration
}

type CacheUpdatePerf struct {
	TrainUpdateTime         PrettyJsonDuration
	TrainingUpdateInternals TrainingUpdateInternals
	InstrUpdateTime         PrettyJsonDuration
	RsvUpdateTime           PrettyJsonDuration
	NotificationData        NotificationData
}

type SearchCache struct {
	Tag *lang.TagCtx

	TrainingIndex TrainingIndex
	InstrIndex    InstrIndex
	//SubIndex      SubIndex
	RsvIndex RsvIndex
	// for subs and rsvs
	DateRange helpers.DateRange
	//
	BuildPerf          CacheGenerationPerf
	UnlockedUpdatePerf CacheUpdatePerf
	LockedUpdatePerf   CacheUpdatePerf

	ignoreInvalidInstructors bool
}

func (ctx *Ctx) RegenerateCache() error {
	sc, err := ctx.NewSearchCache()
	if err != nil {
		return err
	}
	ctx.Cache = sc
	return nil
}

func (ctx *Ctx) RunCacheGenerationAgent(comm chan struct{}) {

	if err := ctx.RegenerateCache(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		// TODO: notify admin
		return
	}

	if comm != nil {
		comm <- struct{}{}
	}

	for {
		/* done updating cache */
		ctx.RWLock.Unlock()

	CONTINUE_NO_UNLOCK:

		if comm != nil {
			select {
			case <-comm:
				return
			case <-time.After(ctx.CacheRefreshInterval):
				break
			}
		} else {
			time.Sleep(ctx.CacheRefreshInterval)
		}

		// no need for locking for this
		d, err := ctx.GetPgChanges()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			goto CONTINUE_NO_UNLOCK
		}

		if len(d) == 0 {
			goto CONTINUE_NO_UNLOCK
		}

		req, err := ctx.NewUpdateSearchCacheRequest(d)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			goto CONTINUE_NO_UNLOCK
		}

		/* locking search - now writing */
		ctx.RWLock.Lock()

		// update cache
		if err := req.UpdateSearchCache(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	// for {
	// 	go ctx.RunTokenCollectionLoop(quitVal, quitValConfirm)
	// 	go ctx.RunTokenDistrubutionLoop(quit, quitVal)

	// 	time.Sleep(ctx.CacheRefreshInterval)

	// 	// no need for locking for this
	// 	d, err := ctx.GetPgChanges()
	// 	if err != nil {
	// 		fmt.Fprintln(os.Stderr, err)
	// 		continue
	// 	}

	// 	quit <- struct{}{}
	// 	c := <-quitValConfirm
	// 	if c != 0 {
	// 		fmt.Fprintln(os.Stderr, "redeemed too many tokens")
	// 	}

	// 	// this is thread safe area
	// 	if err := ctx.UpdateSearchCache(d); err != nil {
	// 		fmt.Fprintln(os.Stderr, err)
	// 		continue
	// 	}
	// }

	//ctx.PgListenForChangesCb(ctx.UpdateSearchCache)
}

var baseTrainReq = train.DalReadTrainingsRequest{
	WithOccs:      true,
	WithSubModels: true,
	WithGroups:    true,
}

func NewCacheDateRange() helpers.DateRange {
	s := time.Now().In(time.UTC)

	dr := helpers.DateRange{
		Start: s,
		// only lookup up to 160 days
		End: s.Add(time.Hour * 24 * 160),
	}

	return dr
}

func (ctx *Ctx) NewSearchCache() (*SearchCache, error) {

	s := time.Now()

	dr := NewCacheDateRange()

	sc := &SearchCache{
		DateRange:                dr,
		Tag:                      ctx.User.LangCtx.Tag,
		ignoreInvalidInstructors: ctx.IgnoreInvalidInstructors,
	}

	err := helpers.RunFunctionsInParallel([]func() error{
		func() error {
			sb := time.Now()
			var err error
			sc.TrainingIndex, err = ctx.NewTrainingIndex()
			if err != nil {
				return err
			}
			sc.BuildPerf.TrainingBuildTime.IndexBuildTime = PrettyJsonDuration(time.Since(sb))
			return nil
		},
		func() error {
			sb := time.Now()
			var err error
			sc.RsvIndex, err = ctx.NewRsvIndex(dr)
			sc.BuildPerf.RsvBuildTime.IndexBuildTime = PrettyJsonDuration(time.Since(sb))
			return err
		},
		// func() error {
		// 	cluster, err := ctx.NewSubIndexCluster(dr)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	sc.SubIndex = NewSubIndex(cluster)
		// 	return nil
		// },
		func() error {
			sb := time.Now()
			instrCluster, err := ctx.NewInstrIndexCluster()
			if err != nil {
				return err
			}
			sc.BuildPerf.InstrBuildTime.ClusterBuildTime = PrettyJsonDuration(time.Since(sb))
			sb = time.Now()
			sc.InstrIndex, err = ctx.NewInstrIndex(instrCluster)
			sc.BuildPerf.InstrBuildTime.IndexBuildTime = PrettyJsonDuration(time.Since(sb))
			return err
		},
	})

	sc.BuildPerf.TotalBuildTime = PrettyJsonDuration(time.Since(s))

	if err != nil {
		return sc, err
	}

	return sc, nil
}

type UpdateSearchCacheRequest struct {
	UpdateTrainingIxRequest UpdateTrainingIxRequest
	UpdateInstrIxRequest    UpdateInstrIxRequest
	UpdateRsvIxRequest      UpdateRsvIxRequest
	//
	Nd  NotificationData
	ctx *Ctx
}

func (ctx *Ctx) NewUpdateSearchCacheRequest(nd NotificationData) (UpdateSearchCacheRequest, error) {
	req := UpdateSearchCacheRequest{
		Nd:  nd,
		ctx: ctx,
	}
	up := &req.ctx.Cache.UnlockedUpdatePerf

	err := helpers.RunFunctionsInParallel([]func() error{
		func() error {
			s := time.Now()
			r, err := ctx.NewUpdateInstrIxRequest(nd)
			req.UpdateInstrIxRequest = r
			up.TrainUpdateTime = PrettyJsonDuration(time.Since(s))
			return err
		}, func() error {
			s := time.Now()
			r, err := ctx.NewUpdateRsvIxRequest(nd)
			req.UpdateRsvIxRequest = r
			up.TrainUpdateTime = PrettyJsonDuration(time.Since(s))
			return err
		}, func() error {
			s := time.Now()
			r, err := ctx.NewUpdateTrainingIxRequest(nd, &up.TrainingUpdateInternals)
			req.UpdateTrainingIxRequest = r
			up.TrainUpdateTime = PrettyJsonDuration(time.Since(s))
			return err
		},
	})

	return req, err
}

func (req *UpdateSearchCacheRequest) UpdateSearchCache() error {

	if req.ctx.DbgSlowdownCacheRefresh > 0 {
		time.Sleep(req.ctx.DbgSlowdownCacheRefresh)
	}

	up := &req.ctx.Cache.LockedUpdatePerf

	up.NotificationData = req.Nd

	if err := helpers.RunFunctionsInParallel([]func() error{
		// trainings
		func() error {
			s := time.Now()
			err := req.ctx.DiffUpdateTrainingIndex(req.UpdateTrainingIxRequest, &up.TrainingUpdateInternals)
			up.TrainUpdateTime = PrettyJsonDuration(time.Since(s))
			return err
		},
		// instructors
		func() error {
			s := time.Now()
			if err := req.ctx.DiffUpdateInstrIndex(req.UpdateInstrIxRequest); err != nil {
				return err
			}
			up.InstrUpdateTime = PrettyJsonDuration(time.Since(s))
			return nil
		},
		// rsvs
		func() error {
			s := time.Now()
			if err := req.ctx.DiffUpdateRsvIndex(req.UpdateRsvIxRequest); err != nil {
				return err
			}
			up.RsvUpdateTime = PrettyJsonDuration(time.Since(s))
			return nil
		},
	}); err != nil {
		return err
	}

	// subs
	// if err := ctx.DiffUpdateSubIndex(elems[Subs]); err != nil {
	// 	return err
	// }
	return nil
}

func (ctx *Ctx) UpdateSearchCache(nd NotificationData) error {
	req, err := ctx.NewUpdateSearchCacheRequest(nd)
	if err != nil {
		return err
	}
	return req.UpdateSearchCache()
}
