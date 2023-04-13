package search

// type SubIndex struct {
// 	// strictly speaking we do not all of this data
// 	// TODO: replace cluster with strictly needed information
// 	Cluster []sub.SubWithJoins
// 	IDhash  map[uuid.UUID]int
// 	SmHash  sub.GroupedSm
// }

// func (ctx *Ctx) NewSubIndexCluster(dr helpers.DateRange) ([]sub.SubWithJoins, error) {
// 	return ctx.Rsv.Sub.DalReadSubWithJoins(sub.ReadSubRequest{
// 		DateRange: dr,
// 	})
// }

// func NewSubIndex(cluster []sub.SubWithJoins) SubIndex {
// 	var subIndex SubIndex

// 	subIndex.IDhash = make(map[uuid.UUID]int, len(cluster))
// 	for i := range cluster {
// 		subIndex.IDhash[cluster[i].ID] = i
// 	}

// 	subIndex.Cluster = cluster
// 	subIndex.SmHash = sub.NewGroupedSm(cluster)
// 	return subIndex
// }

// func (ctx *Ctx) DiffUpdateSubIndex(notSubIDs SingleTableIDmap) error {

// 	if len(notSubIDs) == 0 {
// 		return nil
// 	}

// 	ids := convertIDmapToStrArr(notSubIDs)
// 	modifiedSubs, err := ctx.Rsv.Sub.DalReadSubWithJoins(sub.ReadSubRequest{
// 		DateRange: NewCacheDateRange(),
// 		IDs:       ids,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	newCluster := ctx.Cache.SubIndex.Cluster

// 	modifiedIDs := make(map[uuid.UUID]struct{})
// 	for i := range modifiedSubs {
// 		id := modifiedSubs[i].ID
// 		if ix, f := ctx.Cache.SubIndex.IDhash[id]; f {
// 			newCluster[ix] = modifiedSubs[i]
// 		} else {
// 			newCluster = append(newCluster, modifiedSubs[i])
// 		}
// 		modifiedIDs[id] = struct{}{}
// 	}

// 	subIxsToBeRemoved := make([]int, 0, 2)

// 	for notID := range notSubIDs {
// 		if _, f := modifiedIDs[notID]; f {
// 			continue
// 		}
// 		ix := ctx.Cache.SubIndex.IDhash[notID]
// 		subIxsToBeRemoved = append(subIxsToBeRemoved, ix)
// 	}

// 	sort.Ints(subIxsToBeRemoved)

// 	for i := len(subIxsToBeRemoved) - 1; i >= 0; i-- {
// 		el := subIxsToBeRemoved[i]
// 		newCluster[el] = newCluster[len(newCluster)-1]
// 		newCluster = newCluster[:len(newCluster)-1]
// 	}

// 	if len(modifiedSubs) > 0 || len(subIxsToBeRemoved) > 0 {
// 		ctx.Cache.SubIndex = NewSubIndex(newCluster)
// 	}

// 	return nil
// }
