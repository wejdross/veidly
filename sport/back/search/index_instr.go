package search

import (
	"sport/helpers"
	"sport/instr"

	"github.com/google/uuid"
)

type InstrCluster map[uuid.UUID]*instr.InstructorWithUser

type InstrIndex struct {
	Cluster      InstrCluster
	VacationHash instr.VacationsGroupedByInstr
	//BuildTime    PrettyJsonDuration
}

func (ctx *Ctx) NewInstrIndexCluster() (InstrCluster, error) {
	_instr, err := ctx.Instr.DalReadInstructors(instr.DalReadInstructorsRequest{})
	if err != nil {
		return nil, err
	}

	iix := make(InstrCluster)
	for i := range _instr {
		iix[_instr[i].ID] = &_instr[i]
	}

	return iix, nil
}

func (ctx *Ctx) NewInstrIndex(
	cluster InstrCluster,
) (InstrIndex, error) {

	var iids = make([]string, len(cluster))
	var i int
	for instrID := range cluster {
		iids[i] = cluster[instrID].ID.String()
		i++
	}

	vh, err := ctx.Instr.DalCreateVacationIndex(iids)
	if err != nil {
		return InstrIndex{}, err
	}

	var ix InstrIndex
	ix.VacationHash = vh
	ix.Cluster = cluster

	return ix, nil
}

type UpdateInstrIxRequest struct {
	ModifiedInstrs     []instr.InstructorWithUser
	ModifiedUserInstrs []instr.InstructorWithUser
	ModifiedVac        instr.VacationsGroupedByInstr
	NotInstrIDs        SingleTableIDmap
	NotVacInstrIDs     SingleTableIDmap
}

func (ctx *Ctx) NewUpdateInstrIxRequest(elems NotificationData) (UpdateInstrIxRequest, error) {

	notInstrIDs := elems[Instructors]
	// if notInstrIDs == nil {
	// 	notInstrIDs = make(SingleTableIDmap)
	// }
	// addToMap(notInstrIDs, elems[InstrVacations])
	notInstrUserIDs := elems[Users]
	notVacInstrIDs := elems[InstrVacations]

	var ret = UpdateInstrIxRequest{
		NotInstrIDs:    notInstrIDs,
		NotVacInstrIDs: notVacInstrIDs,
	}

	if len(notInstrIDs) == 0 && len(notInstrUserIDs) == 0 && len(notVacInstrIDs) == 0 {
		return ret, nil
	}

	err := helpers.RunFunctionsInParallel(
		[]func() error{
			func() error {
				if len(notInstrIDs) == 0 {
					return nil
				}
				iids := helpers.IDMapToStringArr(notInstrIDs)
				x, err := ctx.Instr.DalReadInstructors(instr.DalReadInstructorsRequest{
					IDs: iids,
				})
				ret.ModifiedInstrs = x
				return err
			}, func() error {
				if len(notInstrUserIDs) == 0 {
					return nil
				}
				uids := helpers.IDMapToStringArr(notInstrUserIDs)
				x, err := ctx.Instr.DalReadInstructors(instr.DalReadInstructorsRequest{
					UserIDs: uids,
				})
				ret.ModifiedUserInstrs = x
				return err
			}, func() error {
				if len(notVacInstrIDs) == 0 {
					return nil
				}
				iids := helpers.IDMapToStringArr(notVacInstrIDs)
				v, err := ctx.Instr.DalCreateVacationIndex(iids)
				ret.ModifiedVac = v
				return err
			},
		},
	)

	return ret, err
}

func (ctx *Ctx) DiffUpdateInstrIndex(req UpdateInstrIxRequest) error {

	newcl := ctx.Cache.InstrIndex.Cluster

	// for user changes only modify existing instrs
	for i := range req.ModifiedUserInstrs {
		id := req.ModifiedUserInstrs[i].ID
		if _, f := newcl[id]; f {
			newcl[id] = &req.ModifiedUserInstrs[i]
		}
	}

	// merge

	modifiedIDs := make(map[uuid.UUID]struct{}, len(req.ModifiedInstrs))

	// for every modified instructor - replace it in cluster
	for i := range req.ModifiedInstrs {
		id := req.ModifiedInstrs[i].ID
		newcl[id] = &req.ModifiedInstrs[i]
		modifiedIDs[id] = struct{}{}
	}

	// for every instructor that we got in notification
	// but we didnt find in database - it doesnt exist any more - remove it
	for notID := range req.NotInstrIDs {
		if _, f := modifiedIDs[notID]; f {
			continue
		}
		//newcl[notID] = nil
		delete(newcl, notID)
		// also delete record from vacation index
		delete(ctx.Cache.InstrIndex.VacationHash, notID)
		//ctx.Cache.InstrIndex.VacationHash[notID] = nil
	}

	// alter vacation index
	for viid := range req.NotVacInstrIDs {
		ctx.Cache.InstrIndex.VacationHash[viid] = req.ModifiedVac[viid]
	}

	/*
		code below completely rebuilds index based on new cluster.
		This is very slow thats why im only doing diff
		its much, much faster, but
	*/

	// if structure of cluster changed - rebuild index
	// if len(modifiedIDs) != 0 || removedElem {
	// 	var err error
	// 	if ctx.Cache.InstrIndex, err = ctx.NewInstrIndex(newcl); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
