package search

import (
	"encoding/binary"
	"math"
	"net"
	"sport/dal"
)

var emptyLocalization = IPData{}

type IPData struct {
	IP       string
	IsFilled bool
	Country  string
	Region   string
	City     string
	Language string
	Prob     float64
}

// return pointer to empty IPData struct if not found
// note that this will never return nil
func GetBestIP4Localization(probs []IPData) *IPData {
	var max *IPData = &emptyLocalization
	for i := 0; i < len(probs); i++ {
		if probs[i].Prob > max.Prob {
			max = &probs[i]
		}
	}
	return max
}

func GetPossibleIP4Localizations(dal *dal.Ctx, ip string) ([]IPData, error) {

	ipn := int64(binary.BigEndian.Uint32(net.ParseIP(ip).To4()))

	/*
		optionally you may add this stmt:
		(select ip_start, ip_end, country, region, city, 1 as type
			from ip
			where $1 >= ip_start and $1 <= ip_end
			limit 1)
		union all
	*/
	r, err := dal.Db.Query(`
		(select ip_start, ip_end, country, region, city, language, 2 as type
			from ip 
			where ip_start <= $1
			order by ip_start desc
			limit 1)
	`, ipn)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var tmp IPData
	ret := make([]IPData, 0, 3) // predicting up to 3 rows

	var ips, ipe int64
	var recordType int

	for r.Next() {
		if err := r.Scan(
			&ips,
			&ipe,
			&tmp.Country,
			&tmp.Region,
			&tmp.City,
			&tmp.Language,
			&recordType); err != nil {
			return nil, err
		}
		tmp.IP = ip
		tmp.IsFilled = true
		switch recordType {
		case 1:
			tmp.Prob = 1
			// if found ip in range then add it with prob = 1 (we cant get any better match than this)
			ret = append(ret, tmp)
			break
		case 2:
			dist := math.Abs(float64(ips - ipn))
			// if found in outside of range then calculate prob that this ip fits
			if dist <= 0xFFFFFF {
				tmp.Prob = 1 - (dist / 0xFFFFFF)
				ret = append(ret, tmp)
			}
			break
		}
	}
	return ret, nil
}
