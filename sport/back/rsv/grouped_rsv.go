package rsv

import (
	"sort"
	"sport/train"
	"time"

	"github.com/google/uuid"
)

type GroupHashRecord struct {
	DateStart, DateEnd time.Time
	People             int
	Trainings          map[uuid.UUID]struct{}
}

type GroupHash map[uuid.UUID][]GroupHashRecord

func AddToGroupHash(
	h GroupHash,
	groups train.GroupArr,
	trainingID *uuid.UUID,
	dateStart, dateEnd time.Time,
) {
	for i := range groups {
		g := &groups[i]
		ixEntry := GroupHashRecord{
			People:    1,
			DateStart: dateStart,
			DateEnd:   dateEnd,
		}
		if trainingID != nil {
			ixEntry.Trainings = map[uuid.UUID]struct{}{
				*trainingID: {},
			}
		}
		if record := h[g.ID]; record == nil {
			h[g.ID] = []GroupHashRecord{ixEntry}
		} else {
			for z := range record {
				if !overlaps(record[z].DateStart, record[z].DateEnd, dateStart, dateEnd) {
					continue
				}
				record[z].People++
				ixEntry.People++
				if trainingID != nil {
					if _, ok := record[z].Trainings[*trainingID]; !ok {
						record[z].Trainings[*trainingID] = struct{}{}
					}
					if _, ok := ixEntry.Trainings[*trainingID]; !ok {
						ixEntry.Trainings[*trainingID] = struct{}{}
					}
				}
			}
			h[g.ID] = append(h[g.ID], ixEntry)
		}
	}
}

func AddShortRsvToGroupHash(h GroupHash, sr *ShortenedRsv) {
	AddToGroupHash(h, sr.Groups, sr.TrainingID, sr.DateStart, sr.DateEnd)
}

func RmFromGroupHash(
	h GroupHash,
	groups train.GroupArr,
	trainingID *uuid.UUID,
	dateStart, dateEnd time.Time,
) {
	for i := range groups {
		g := &groups[i]
		if records := h[g.ID]; records == nil {
			continue
		} else {
			ixsToBeRemoved := make([]int, 0, 2)
			for z := range records {
				if !overlaps(records[z].DateStart, records[z].DateEnd, dateStart, dateEnd) {
					continue
				}
				if records[z].People > 1 {
					records[z].People--
				} else {
					// delete record[z]
					ixsToBeRemoved = append(ixsToBeRemoved, z)
				}
				if trainingID != nil {
					delete(records[z].Trainings, *trainingID)
				}
			}
			if len(ixsToBeRemoved) > 0 {
				sort.Ints(ixsToBeRemoved)
				for z := len(ixsToBeRemoved) - 1; z >= 0; z-- {
					ix := ixsToBeRemoved[z]
					records[ix] = records[len(records)-1]
					records = records[:len(records)-1]
				}
				if len(records) == 0 {
					// delete the key
					delete(h, g.ID)
				} else {
					h[g.ID] = records
				}
			}
		}
	}
}

func RmShortRsvFromGroupHash(h GroupHash, sr *ShortenedRsv) {
	RmFromGroupHash(h, sr.Groups, sr.TrainingID, sr.DateStart, sr.DateEnd)
}

func NewGroupHashFromShortRsvs(rr []ShortenedRsv) GroupHash {
	ix := make(GroupHash)
	for i := range rr {
		r := &rr[i]
		AddToGroupHash(ix, r.Groups, r.TrainingID, r.DateStart, r.DateEnd)
	}
	return ix
}

func NewGroupHashFromRsvs(rr []RsvWithInstr) GroupHash {
	ix := make(GroupHash)
	for i := range rr {
		r := &rr[i]
		if !r.IsConfirmed {
			continue
		}
		AddToGroupHash(ix, r.Groups, r.TrainingID, r.DateStart, r.DateEnd)
	}
	return ix
}
