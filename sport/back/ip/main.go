package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sport/config"
	"sport/dal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/lib/pq"
)

// lang index is map of ISO-3166 country : ISO 639-1 language
// example: PL : pl
func CreateLangIndex() (map[string]string, error) {

	var err error

	fp, err := os.Open("lang.tsv")
	if err != nil {
		log.Fatal(err)
	}

	cr := csv.NewReader(fp)
	cr.Comma = rune('\t')
	cr.Comment = rune('#')

	// csv common column indexes
	const LocaleIx = 15
	const CountryNameIx = 4
	const CountryCodeIx = 0

	ret := make(map[string]string)

	for {
		x, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		if len(x) != 19 {
			return nil, fmt.Errorf("unrecognized format of data file")
		}

		//fmt.Printf("%s - %s\n", x[15], x[4])

		// lc is "ar-AE,fa,en,hi,ur"
		lc := x[LocaleIx]

		// parts is ["ar-AE","fa","en","hi","ur"]
		parts := strings.Split(lc, ",")

		if len(parts) == 0 {
			continue
		}

		langCode := ""

		// check if first part containes dash
		parts = strings.Split(parts[0], "-")
		switch len(parts) {
		case 0:
			continue
		case 1:
			fallthrough
		case 2:
			langCode = parts[0]
			break
		default:
			continue
		}
		//fmt.Println(langCode)
		ret[x[CountryCodeIx]] = langCode

	}

	return ret, nil
}

func AppendInsertRecord(q *strings.Builder, x []string, i int64, langIx map[string]string) error {
	var err error
	if len(x) != 6 {
		return fmt.Errorf("unrecognized format of data file")
	}

	// find language code
	var lc string
	if lc = langIx[x[2]]; lc == "" {
		if x[2] != "-" {
			fmt.Printf("\033[33mWarning: Couldnt find language code for %s (%s)\033[0m\n", x[2], x[3])
		}
	}

	// parse IPs
	var ips, ipe int
	if ips, err = strconv.Atoi(x[0]); err != nil {
		return err
	}
	if ipe, err = strconv.Atoi(x[0]); err != nil {
		return err
	}

	if i == 0 {
		q.WriteString(fmt.Sprintf(`insert into ip values (
			%d, %d, %s, %s, %s, %s
			)`, ips, ipe, pq.QuoteLiteral(x[3]), pq.QuoteLiteral(x[4]), pq.QuoteLiteral(x[5]), pq.QuoteLiteral(lc)))
	} else {
		q.WriteString(fmt.Sprintf(`,(
			%d, %d, %s, %s, %s, %s
		)`, ips, ipe, pq.QuoteLiteral(x[3]), pq.QuoteLiteral(x[4]), pq.QuoteLiteral(x[5]), pq.QuoteLiteral(lc)))
	}

	return nil
}

// will run faster but will consume fuck load of memory
func uploadHighMemory(langIx map[string]string, fp *os.File, d *dal.Ctx) (int64, error) {

	fmt.Println("\033[33mWarning: running HIGH memory variant. This may consume lots of memory\033[0m")

	cr := csv.NewReader(fp)

	var i int64

	q := strings.Builder{}

	fmt.Println("forming query...")

	for {
		x, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return 0, err
			}
		}
		if err := AppendInsertRecord(&q, x, i, langIx); err != nil {
			return i, err
		}
		i++
	}

	fmt.Println("inserting...")

	r, err := d.Db.Exec(q.String())
	if err != nil {
		return -1, err
	}

	return r.RowsAffected()
}

// willl run slower but will consume a little memory
func uploadLowMemory(langIx map[string]string, fp *os.File, d *dal.Ctx) (int64, error) {

	fmt.Println("\033[33mWarning: running LOW memory variant. This may take a while\033[0m")

	cr := csv.NewReader(fp)

	var i, c int64

	q := strings.Builder{}

	for {
		x, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				if i != 0 {
					_, err = d.Db.Exec(q.String())
					if err != nil {
						return c, err
					}
				}
				break
			} else {
				return c, err
			}
		}

		if err := AppendInsertRecord(&q, x, i, langIx); err != nil {
			return i, err
		}

		c++
		i++

		if (i % 1000) == 0 {
			_, err = d.Db.Exec(q.String())
			if err != nil {
				return c, err
			}
			q.Reset()
			i = 0
		}

	}

	return c, nil
}

func getMem() uint64 {
	in := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(in)
	if err != nil {
		return 0
	}
	return uint64(in.Totalram) * uint64(in.Unit)
}

func main() {

	d := dal.NewDal(config.NewLocalCtx(), "sportdb")

	var err error

	fmt.Println("creating language index...")

	ix, err := CreateLangIndex()
	if err != nil {
		log.Fatal(err)
	}

	fp, err := os.Open("data.csv")
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}

	var start time.Time = time.Now()

	var x int64

	// megabytes
	var reqMem uint64 = 4000

	// if system memory is bigger than reqMem megs then run high memory version
	// otherwise fallback to slower low mem
	if getMem()/1024/1024 > reqMem {
		x, err = uploadHighMemory(ix, fp, d)
	} else {
		x, err = uploadLowMemory(ix, fp, d)
	}

	elapsed := time.Since(start).Seconds()

	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mTerminated with error/s: %v\033[0m", err)
	} else {
		fmt.Print("\033[32mSuccess\033[0m")
	}

	fmt.Printf("\n\tOperation took %f seconds", elapsed)
	fmt.Printf("\n\tInserted %d rows\n", x)

}
