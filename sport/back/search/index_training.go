package search

import (
	"sport/helpers"
	"sport/sub"
	"sport/train"
	"time"

	"github.com/google/btree"
	"github.com/google/uuid"
)

type TrainingCluster map[uuid.UUID]*train.TrainingWithJoins
type CountryHashValue struct {
	IDhash map[uuid.UUID]int
	Slice  []*train.TrainingWithJoins
}
type CountryHash map[string]*CountryHashValue
type TrainingMapForID map[uuid.UUID]map[uuid.UUID]int

type BTreeItem struct {
	Dg  float64
	Val []*train.TrainingWithJoins
}

func (bti *BTreeItem) Less(other btree.Item) bool {
	oi := other.(*BTreeItem)
	return bti.Dg < oi.Dg
}

type TrainingIndex struct {
	//Cluster []*train.TrainingWithJoins
	cluster TrainingCluster

	/*Lat,*/
	Lng *btree.BTree
	//
	CountryHash CountryHash

	// all training ids containing specified sub models
	SubModelHash TrainingMapForID
	// all training ids containing specified groups
	GroupIDHash TrainingMapForID
}

func (hash TrainingMapForID) AddElement(objID, trainingID uuid.UUID, elementIndex int) {
	if el, e := hash[objID]; e {
		el[trainingID] = elementIndex
	} else {
		el = make(map[uuid.UUID]int)
		el[trainingID] = elementIndex
		hash[objID] = el
	}
}

func AddAllTrainingSmsAndGrpsIntoMaps(t *train.TrainingWithJoins, smHash, grpHash TrainingMapForID) {
	trainingID := t.Training.ID
	for j := range t.Sms {
		smID := t.Sms[j].ID
		smHash.AddElement(smID, trainingID, j)
	}
	for j := range t.Groups {
		grpID := t.Groups[j].ID
		grpHash.AddElement(grpID, trainingID, j)
	}
}

func (hash TrainingMapForID) RmElement(objID, trainingID uuid.UUID) {
	if e, f := hash[objID]; f {
		delete(e, trainingID)
		if len(e) == 0 {
			delete(hash, objID)
		}
	}
}

func RmAllTrainingSmsAndGrpsFromMaps(t *train.TrainingWithJoins, smHash, grpHash TrainingMapForID) {
	trainingID := t.Training.ID
	for j := range t.Sms {
		smID := t.Sms[j].ID
		smHash.RmElement(smID, trainingID)
	}
	for j := range t.Groups {
		grpID := t.Groups[j].ID
		grpHash.RmElement(grpID, trainingID)
	}
}

func RmTrainingFromBtree(t *train.TrainingWithJoins, dg float64, btr *btree.BTree) {
	shouldDeleteElement := false
	shouldReturn := false
	req := &BTreeItem{Dg: dg}
	btr.AscendGreaterOrEqual(req, func(i btree.Item) bool {
		bi := i.(*BTreeItem)
		found := -1
		for i := range bi.Val {
			if bi.Val[i].Training.ID == t.Training.ID {
				found = i
				break
			}
		}
		if found == -1 {
			shouldReturn = true
			return false
		}
		if len(bi.Val) == 1 {
			shouldDeleteElement = true
			return false
		}
		lastIx := len(bi.Val) - 1
		bi.Val[found] = bi.Val[lastIx]
		bi.Val = bi.Val[:lastIx]
		return false
	})
	if shouldReturn {
		return
	}
	if shouldDeleteElement {
		btr.Delete(req)
	}
}

func AddTrainingToBtree(t *train.TrainingWithJoins, dg float64, btr *btree.BTree) {
	shouldAdd := true
	shouldReturn := false
	req := &BTreeItem{Dg: dg, Val: []*train.TrainingWithJoins{t}}
	btr.AscendGreaterOrEqual(req, func(i btree.Item) bool {
		bi := i.(*BTreeItem)
		if bi.Dg != dg {
			return false
		}
		found := -1
		for i := range bi.Val {
			if bi.Val[i].Training.ID == t.Training.ID {
				found = i
				break
			}
		}
		if found != -1 {
			shouldReturn = true
			return false
		}
		bi.Val = append(bi.Val, t)
		shouldAdd = false
		return false
	})
	if shouldReturn {
		return
	}
	if shouldAdd {
		btr.ReplaceOrInsert(req)
	}
}

// func RmTrainingFromIndexBtrees(t *train.TrainingWithJoins, lat, lng *btree.BTree) {
// 	RmTrainingFromBtree(t, t.Training.LocationLat, lat)
// 	RmTrainingFromBtree(t, t.Training.LocationLng, lng)
// }

// func AddTrainingToIndexBtrees(t *train.TrainingWithJoins, lat, lng *btree.BTree) {
// 	AddTrainingToBtree(t, t.Training.LocationLat, lat)
// 	AddTrainingToBtree(t, t.Training.LocationLng, lng)
// }

func (ctx *Ctx) NewTrainingCluster() ([]*train.TrainingWithJoins, error) {
	return ctx.Train.DalReadTrainings(baseTrainReq)
}

func (hash CountryHash) AddElement(t *train.TrainingWithJoins) {
	x, found := hash[t.Training.LocationCountry]
	if found {
		if _, e := x.IDhash[t.Training.ID]; e {
			return
		}
	} else {
		x = &CountryHashValue{
			IDhash: make(map[uuid.UUID]int),
			Slice:  make([]*train.TrainingWithJoins, 0, 4),
		}
	}
	x.IDhash[t.Training.ID] = len(x.Slice)
	x.Slice = append(x.Slice, t)
	hash[t.Training.LocationCountry] = x
}

func (hash CountryHash) RmElement(t *train.TrainingWithJoins) {
	v, found := hash[t.Training.LocationCountry]
	if !found {
		return
	}
	ix, found := v.IDhash[t.Training.ID]
	if !found {
		return
	}
	lastElemIx := len(v.Slice) - 1
	if ix == lastElemIx {
		delete(v.IDhash, t.Training.ID)
		v.Slice = v.Slice[:lastElemIx]
		return
	}
	v.Slice[ix] = v.Slice[lastElemIx]
	v.Slice = v.Slice[:lastElemIx]
	delete(v.IDhash, t.Training.ID)
	v.IDhash[v.Slice[ix].Training.ID] = ix
}

func (ctx *Ctx) NewTrainingIndex() (TrainingIndex, error) {

	t, err := ctx.Train.DalReadTrainings(baseTrainReq)
	if err != nil {
		return TrainingIndex{}, err
	}

	// sorted arr of trainings
	// one is lat asc,
	// other is lng asc
	degree := 128
	//latOrd := btree.New(degree)
	lngOrd := btree.New(degree)
	countryHash := make(CountryHash)

	smHash := make(TrainingMapForID)
	grpHash := make(TrainingMapForID)
	cluster := make(TrainingCluster)

	//latDupl := make(map[float64]*BTreeItem)
	lngDupl := make(map[float64]*BTreeItem)

	for i := 0; i < len(t); i++ {
		cluster[t[i].Training.ID] = t[i]

		AddAllTrainingSmsAndGrpsIntoMaps(t[i], smHash, grpHash)

		countryHash.AddElement(t[i])

		if t[i].Training.LocationLat == 0 && t[i].Training.LocationLng == 0 {
			continue
		}

		// lat := t[i].Training.LocationLat
		// if el, e := latDupl[lat]; e {
		// 	el.Val = append(el.Val, t[i])
		// } else {
		// 	latItem := &BTreeItem{
		// 		Dg:  lat,
		// 		Val: []*train.TrainingWithJoins{t[i]},
		// 	}
		// 	latDupl[lat] = latItem
		// 	latOrd.ReplaceOrInsert(latItem)
		// }

		lng := t[i].Training.LocationLng
		if el, e := lngDupl[lng]; e {
			el.Val = append(el.Val, t[i])
		} else {
			lngItem := &BTreeItem{
				Dg:  lng,
				Val: []*train.TrainingWithJoins{t[i]},
			}
			lngDupl[lng] = lngItem
			lngOrd.ReplaceOrInsert(lngItem)
		}
	}

	// _ = RunFunctionsInParallel([]func() error{
	// 	func() error {
	// 		sort.Slice(latOrd, func(i, j int) bool {
	// 			return latOrd[i].Training.LocationLat < latOrd[j].Training.LocationLat
	// 		})
	// 		return nil
	// 	}, func() error {
	// 		sort.Slice(lngOrd, func(i, j int) bool {
	// 			return lngOrd[i].Training.LocationLng < lngOrd[j].Training.LocationLng
	// 		})
	// 		return nil
	// 	},
	// })

	var sc TrainingIndex
	sc.cluster = cluster
	//sc.Lat = latOrd
	sc.Lng = lngOrd
	sc.CountryHash = countryHash
	sc.GroupIDHash = grpHash
	sc.SubModelHash = smHash

	return sc, nil
}

type UpdateTrainingIxRequest struct {
	ChangedTrainings []*train.TrainingWithJoins
	ChangedGroups    []train.Group
	ChangedSubModels []sub.SubModel
	//
	NotTrainingIDs SingleTableIDmap
}

func (ctx *Ctx) NewUpdateTrainingIxRequest(elems NotificationData, perf *TrainingUpdateInternals) (UpdateTrainingIxRequest, error) {

	notTrainingIDs := elems[Trainings]
	if notTrainingIDs == nil {
		notTrainingIDs = make(SingleTableIDmap)
	}
	addToMap(notTrainingIDs, elems[Occurrences])
	addToMap(notTrainingIDs, elems[SecondaryOccs])
	addToMap(notTrainingIDs, elems[SubModelBindings])
	addToMap(notTrainingIDs, elems[TrainingsVgroups])
	addToMap(notTrainingIDs, elems[Reviews])
	notSmIDs := elems[SubModels]
	notGroupIDs := elems[TrainingGroups]

	var (
		ret UpdateTrainingIxRequest = UpdateTrainingIxRequest{
			NotTrainingIDs: notTrainingIDs,
		}
		pfs []func() error = make([]func() error, 0, 3)
		//perf *TrainingUpdateInternals = &ctx.Cache.UnlockedUpdatePerf.TrainingUpdateInternals
	)

	// remove / insert / update trainings from notification
	if len(notTrainingIDs) > 0 {
		pfs = append(pfs, func() error {
			var err error
			s := time.Now()
			r := baseTrainReq
			r.TrainingIDs = make([]string, len(notTrainingIDs))
			i := 0
			for x := range notTrainingIDs {
				r.TrainingIDs[i] = x.String()
				i++
			}
			ret.ChangedTrainings, err = ctx.Train.DalReadTrainings(r)
			perf.ReadTrainingsTime = PrettyJsonDuration(time.Since(s))
			return err
		})
	}

	if len(notGroupIDs) > 0 {
		pfs = append(pfs, func() error {
			s := time.Now()
			ids := make([]string, len(notGroupIDs))
			var i int
			for x := range notGroupIDs {
				ids[i] = x.String()
				i++
			}
			_gs, err := ctx.Train.DalReadTrainingGroups(train.ReadTrainingGroupsRequest{
				IDs: ids,
			})
			ret.ChangedGroups = _gs
			perf.ReadGroupTime = PrettyJsonDuration(time.Since(s))
			return err
		})
	}

	if len(notSmIDs) > 0 {
		pfs = append(pfs, func() error {
			s := time.Now()
			ids := make([]string, len(notSmIDs))
			var i int
			for x := range notSmIDs {
				ids[i] = x.String()
				i++
			}
			_sm, err := ctx.Rsv.Sub.DalReadSubModels(sub.ReadSubModelRequest{
				IDs: ids,
			})
			ret.ChangedSubModels = _sm
			perf.ReadSmTime = PrettyJsonDuration(time.Since(s))
			return err
		})
	}

	s := time.Now()
	if err := helpers.RunFunctionsInParallel(pfs); err != nil {
		return ret, err
	}
	perf.TotalReadTime = PrettyJsonDuration(time.Since(s))

	return ret, nil
}

func (ctx *Ctx) DiffUpdateTrainingIndex(req UpdateTrainingIxRequest, perf *TrainingUpdateInternals) error {

	var (
		index                                       = ctx.Cache.TrainingIndex
		cluster          TrainingCluster            = index.cluster
		changedTrainings []*train.TrainingWithJoins = req.ChangedTrainings
		gs               []train.Group              = req.ChangedGroups
		sm               []sub.SubModel             = req.ChangedSubModels
		//perf             *TrainingUpdateInternals   = &ctx.Cache.LockedUpdatePerf.TrainingUpdateInternals
	)

	// note that following 2 changes (sub models and groups) only operate on training references
	// since every index contains not training copy but reference to cluster training
	// there is no need to update them

	// find trainings with specified sms and update them
	for i := range sm {
		trainings := index.SubModelHash[sm[i].ID]
		for tid := range trainings {
			if t, e := cluster[tid]; e {
				t.Sms[trainings[tid]] = sm[i]
			}
		}
	}

	// find trainings with specified groups and update them
	for i := range gs {
		trainings := index.GroupIDHash[gs[i].ID]
		for tid := range trainings {
			if t, e := cluster[tid]; e {
				t.Groups[trainings[tid]] = gs[i]
			}
		}
	}

	modifiedIDs := make(map[uuid.UUID]struct{})

	s := time.Now()

	// merge trainings
	for i := range changedTrainings {
		// add or replace training
		trainingID := changedTrainings[i].Training.ID
		oldTraining, exists := cluster[trainingID]

		newTraining := changedTrainings[i]
		cluster[trainingID] = newTraining

		if exists {
			index.CountryHash.RmElement(oldTraining)
			// one may validate if sms / grps really changed
			RmAllTrainingSmsAndGrpsFromMaps(
				oldTraining, index.SubModelHash, index.GroupIDHash)

			RmTrainingFromBtree(oldTraining, oldTraining.Training.LocationLng, index.Lng)
			//RmTrainingFromIndexBtrees(oldTraining, index.Lat, index.Lng)
		}

		index.CountryHash.AddElement(newTraining)

		AddAllTrainingSmsAndGrpsIntoMaps(newTraining, index.SubModelHash, index.GroupIDHash)
		//AddTrainingToIndexBtrees(newTraining, index.Lat, index.Lng)
		AddTrainingToBtree(newTraining, newTraining.Training.LocationLng, index.Lng)

		modifiedIDs[trainingID] = struct{}{}
	}

	perf.MergeTrainingTime = PrettyJsonDuration(time.Since(s))

	s = time.Now()

	// remove trainings
	for trainingID := range req.NotTrainingIDs {
		if _, found := modifiedIDs[trainingID]; found {
			continue
		}
		trainingToBeRemoved := cluster[trainingID]
		if trainingToBeRemoved == nil {
			continue
		}
		// Received notification with training ID but training doesnt exist any more => it was deleted
		delete(cluster, trainingID)
		index.CountryHash.RmElement(trainingToBeRemoved)
		RmAllTrainingSmsAndGrpsFromMaps(
			trainingToBeRemoved,
			index.SubModelHash, index.GroupIDHash)
		RmTrainingFromBtree(trainingToBeRemoved, trainingToBeRemoved.Training.LocationLng, index.Lng)
	}

	perf.RmTrainingTime = PrettyJsonDuration(time.Since(s))

	// rebuild indexes if training order / structure changed
	// if len(changedTrainings) > 0 || len(trainingIxsToBeRemoved) > 0 {
	// 	ctx.Cache.TrainingIndex = NewTrainingIndex(cluster)
	// }

	return nil
}
