package search

import (
	"fmt"
	"os"
	"sport/helpers"
	"sport/rsv"

	"github.com/google/uuid"
)

var isConfirmed = true
var isConfirmedPtr = &isConfirmed

type RsvCluster map[uuid.UUID]rsv.ShortenedRsv

type RsvIndex struct {
	Cluster      RsvCluster
	RsvCountHash rsv.RsvCount
	GroupHash    rsv.GroupHash
}

func (ctx *Ctx) NewRsvIndex(dr helpers.DateRange) (RsvIndex, error) {

	r, err := ctx.Rsv.ReadShortenedRsvs(rsv.ReadShortenedRsvsArgs{
		DateRange: dr,
	})
	if err != nil {
		return RsvIndex{}, err
	}
	cl := make(map[uuid.UUID]rsv.ShortenedRsv, len(r))
	for i := range r {
		cl[r[i].ID] = r[i]
	}

	var ix RsvIndex
	ix.Cluster = cl

	ix.RsvCountHash = ctx.Rsv.NewRsvCountFromShortRsvs(r)
	ix.GroupHash = rsv.NewGroupHashFromShortRsvs(r)

	return ix, nil
}

func AddElementToRsvCount(hash rsv.RsvCount, rsvToBeAdded *rsv.ShortenedRsv) {
	countKey := rsv.ShortRsvToCountKey(rsvToBeAdded)
	hash[countKey]++
}

func RemoveElementFromRsvCount(hash rsv.RsvCount, rsvToBeRemoved *rsv.ShortenedRsv) {
	countKey := rsv.ShortRsvToCountKey(rsvToBeRemoved)
	if c, e := hash[countKey]; !e {
		fmt.Fprintln(os.Stderr, "1: didnt find reservation in count - data may be out of sync")
	} else {
		if c > 1 {
			hash[countKey]--
		} else {
			delete(hash, countKey)
		}
	}
}

type UpdateRsvIxRequest struct {
	ModifiedRsvs []rsv.ShortenedRsv
	NotRsvIDs    SingleTableIDmap
}

func (ctx *Ctx) NewUpdateRsvIxRequest(elems NotificationData) (UpdateRsvIxRequest, error) {
	notRsvIDs := elems[Reservations]

	var ret = UpdateRsvIxRequest{
		NotRsvIDs: notRsvIDs,
	}

	if len(notRsvIDs) == 0 {
		return ret, nil
	}

	ids := helpers.IDMapToStringArr(notRsvIDs)
	modifiedRsvs, err := ctx.Rsv.ReadShortenedRsvs(rsv.ReadShortenedRsvsArgs{
		DateRange: NewCacheDateRange(),
		IDs:       ids,
	})
	ret.ModifiedRsvs = modifiedRsvs

	return ret, err
}

func (ctx *Ctx) DiffUpdateRsvIndex(req UpdateRsvIxRequest) error {

	rsvIndex := ctx.Cache.RsvIndex
	cluster := rsvIndex.Cluster

	modifiedIDs := make(map[uuid.UUID]struct{})

	// rsv was modified or added
	for i := range req.ModifiedRsvs {
		id := req.ModifiedRsvs[i].ID

		if oldRsv, f := cluster[id]; f {

			// replace records in rsv count hash
			RemoveElementFromRsvCount(rsvIndex.RsvCountHash, &oldRsv)
			AddElementToRsvCount(rsvIndex.RsvCountHash, &req.ModifiedRsvs[i])

			// replace records in grp hash
			rsv.RmShortRsvFromGroupHash(rsvIndex.GroupHash, &oldRsv)
			rsv.AddShortRsvToGroupHash(rsvIndex.GroupHash, &req.ModifiedRsvs[i])

			// replace modified rsv in hash
			cluster[id] = req.ModifiedRsvs[i]

		} else {
			// new rsv has been added
			cluster[id] = req.ModifiedRsvs[i]
			AddElementToRsvCount(rsvIndex.RsvCountHash, &req.ModifiedRsvs[i])
			rsv.AddShortRsvToGroupHash(rsvIndex.GroupHash, &req.ModifiedRsvs[i])
		}
		modifiedIDs[id] = struct{}{}
	}

	for notID := range req.NotRsvIDs {
		if _, f := modifiedIDs[notID]; f {
			continue
		}
		rsvToBeRemoved, f := cluster[notID]
		if !f {
			continue
		}
		delete(cluster, notID)
		RemoveElementFromRsvCount(rsvIndex.RsvCountHash, &rsvToBeRemoved)
		rsv.RmShortRsvFromGroupHash(rsvIndex.GroupHash, &rsvToBeRemoved)
	}

	return nil
}
