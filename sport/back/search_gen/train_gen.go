package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"sport/helpers"
	"sport/train"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/html"
)

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func hread(n *html.Node, tags *[]string) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "a" {
			for _, x := range c.Attr {
				if x.Key == "href" && strings.HasPrefix(x.Val, "https://www.britannica.com/sports/") {
					if c.FirstChild.Data != "" {
						*tags = append(*tags, c.FirstChild.Data)
					}
				}
			}
		}
		hread(c, tags)
	}
}

func getdoc(u string) (*html.Node, error) {
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code: %d", res.StatusCode)
	}
	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// casually crawl some tags
func createTags() ([]string, error) {

	_, err := os.Stat("tags")
	if err == nil {
		fmt.Printf("using tags from local cache (./tags)\n")
		fc, err := ioutil.ReadFile("tags")
		if err != nil {
			return nil, err
		}
		return strings.Split(string(fc), "\n"), nil
	}

	fmt.Println("crawling tags...")

	url := "https://www.britannica.com/topic/list-of-sports-2038581"

	doc, err := getdoc(url)
	if err != nil {
		return nil, err
	}

	// if you want to test this function, then take downloaded doc save it to file and read file.
	// dont spam their server with your requests
	//x, err := os.OpenFile("1.html", os.O_RDONLY, 0600)
	// defer x.Close()
	// if err != nil {
	// 	return nil, err
	// }
	// doc, err := html.Parse(x)

	tags := make([]string, 0, 100)
	hread(doc, &tags)

	fmt.Println("saving tags to local cache")

	return tags, ioutil.WriteFile("tags", []byte(strings.Join(tags, "\n")), 0600)
}

type TrainingIDWithOccs struct {
	TID  uuid.UUID
	Occs []train.CreateOccRequest
}

func TrainGen(trainCtx *train.Ctx, maxTh, trainNo int, instrs []TokenWithInstrID) ([]TrainingIDWithOccs, error) {
	ts := make([]TrainingIDWithOccs, trainNo)

	possibleTags, err := createTags()
	if err != nil {
		return nil, err
	}

	if len(possibleTags) == 0 {
		return nil, err
	}

	err = helpers.Sem(maxTh, trainNo, []func(int) error{
		func(i int) error {
			r := (rand.Int() % (24 * 60)) + 1
			s := helpers.NowMin().Add(time.Duration(r) * time.Hour)
			d := (rand.Int() % (4)) + 1
			e := s.Add(time.Hour * time.Duration(d))
			tr := train.NewTestCreateTrainingRequest(s, e)

			// lat and lng are based on europe (north of stockholm and russia excluded)
			minLat := 35
			maxLat := 60
			minLng := -10
			maxLng := 30

			// randLat := rand.Intn(maxLat-minLat) + minLat
			// randLng := rand.Intn(maxLng-minLng) + minLng
			randLat := (rand.Float64() * float64(maxLat-minLat)) + float64(minLat)
			randLng := (rand.Float64() * float64(maxLng-minLng)) + float64(minLng)
			tr.Training.LocationLat = float64(randLat)
			tr.Training.LocationLng = float64(randLng)
			tr.Training.Capacity = 100

			{
				num := rand.Int() % 6
				tags := make([]string, num)

				tx := make(map[string]struct{})
				for j := 0; j < num; j++ {
					nt := possibleTags[rand.Int()%len(possibleTags)]
					if _, e := tx[nt]; e {
						j--
						continue
					}
					tx[nt] = struct{}{}
					tags[j] = nt
				}

				tr.Training.Tags = tags
			}

			tid, err := trainCtx.ApiCreateTraining(instrs[i%len(instrs)].Token, &tr)
			if err != nil {
				return err
			}
			ts[i] = TrainingIDWithOccs{
				TID:  tid,
				Occs: tr.Occurrences,
			}
			return nil
		},
	})

	return ts, err
}
