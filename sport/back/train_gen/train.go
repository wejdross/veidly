package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"runtime"
	"sport/api"
	"sport/config"
	"sport/dal"
	"sport/helpers"
	"sport/instr"
	"sport/static"
	"sport/train"
	"sport/user"
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

func createTrainingsAsync(ctx *train.Ctx, ic, tpc int) chan error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ret := make(chan error)

	go func() {

		faced, err := os.ReadDir(faceDir)
		if err != nil {
			ret <- err
			return
		}

		imgd, err := os.ReadDir(imgDir)
		if err != nil {
			ret <- err
			return
		}

		var instructorCount = ic
		var trainingsPerInstructor = tpc
		var userTokens = make([]string, instructorCount)
		var instructorIds = make([]uuid.UUID, instructorCount)

		n := []string{
			"Isaac",
			"Richard",
			"Andrzej",
			"Max",
			"Jeanette",
			"V",
			"Enrico",
			"Mike",
			"Roger",
			"Łukasz",
			"Kamil",
			"Katherine",
			"Sheldon",
			"James",
			"Clerk",
		}

		ln := []string{
			"Maxwell",
			"Dirac",
			"Fermi",
			"Faraday",
			"Planck",
			"Gauss",
			"Euler",
			"Schrodinger",
			"Bohr",
			"Rutheford",
			"Feynman",
			"Newton",
			"Bohr",
			"Rutheford",
			"Maxwell",
			"Dirac",
			"Fermi",
			"Faraday",
			"Planck",
			"Gauss",
			"Euler",
			"Schrodinger",
			"Gołota",
			"Tyson",
			"Federer",
			"Widera",
			"Zagórski",
			"Maulwurf",
			"Cooper",
		}

		domains := []string{
			"gmail.com",
			"protonmail.onion",
			"interia.pl",
			"wp.pl",
			"bazinga.com",
		}

		rand.Seed(time.Now().UTC().UnixNano())

		for i := 0; i < instructorCount; i++ {
			var err error
			name := n[rand.Int()%len(n)]
			namep := n[rand.Int()%len(n)]
			name2 := ln[rand.Int()%len(ln)]
			domain := domains[rand.Int()%len(domains)]
			password := "password"
			userTokens[i], err = ctx.User.ApiCreateAndLoginUser(&user.UserRequest{
				Email:    name + "_" + namep + "_" + name2 + "@" + domain,
				Password: password,
				UserData: user.UserData{
					Name:     name + " " + name2,
					Language: "pl",
					Country:  "PL",
					AboutMe:  "long text that describes user, just to catch more border cases, long text that describes user",
					Urls: []user.NamedUrl{
						{
							Name:   "instagram",
							Url:    "https://instagram.com",
							Avatar: "instagram",
						}, {
							Name:   "facebook",
							Url:    "https://facebook.com",
							Avatar: "facebook",
						}, {
							Name:   "twitter",
							Url:    "https://twitter.com",
							Avatar: "twitter",
						},
					},
				},
			})

			if err != nil {
				ret <- err
				return
			}

			if len(faced) != 0 {
				x := rand.Int() % len(faced)
				p := path.Join(faceDir, faced[x].Name())
				if err = ctx.User.ApiUploadAvatarFromPath(userTokens[i], p); err != nil {
					ret <- err
					return
				}
			}

			uid, err := ctx.Api.AuthorizeUserFromToken(userTokens[i])
			if err != nil {
				ret <- err
				return
			}

			if instructorIds[i], err = ctx.Instr.CreateTestInstructor(uid, nil); err != nil {
				ret <- err
				return
			}
		}

		possibleTags, err := createTags()
		if err != nil {
			ret <- err
			return
		}

		if len(possibleTags) == 0 {
			ret <- fmt.Errorf("no possible tags")
			return
		}

		// adding trainings
		//apiReq := make([]api.TestCase, len(instructorIds)*trainingsPerInstructor)

		for z := 0; z < len(instructorIds); z++ {
			for i := 0; i < trainingsPerInstructor; i++ {
				num := rand.Int() % 6
				tags := make([]string, num)
				tx := make(map[string]struct{})
				for j := 0; j < num; j++ {
					rand.Seed(time.Now().UTC().UnixNano())
					nt := possibleTags[rand.Int()%len(possibleTags)]
					if _, e := tx[nt]; e {
						j--
						continue
					}
					tx[nt] = struct{}{}
					tags[j] = nt
				}
				rand.Seed(time.Now().UTC().UnixNano())
				t := possibleTags[rand.Int()%len(possibleTags)]
				t += " with " + n[rand.Int()%len(n)]

				// this is completely random lat long
				// var latD, lngD int
				// latD = rand.Int() % 90
				// if rand.Int()%2 == 1 {
				// 	latD *= -1
				// }
				// lngD = rand.Int() % 180
				// if rand.Int()%2 == 1 {
				// 	lngD *= -1
				// }
				// lat := float64(latD) + rand.Float64()
				// lng := float64(lngD) + rand.Float64()

				// this is origin based lat lng
				lat := 52.2297
				lng := 21.017532
				latRad := float64((rand.Int()%300)-150) / 100
				lngRad := float64((rand.Int()%300)-150) / 100
				lat = lat + latRad
				lng = lng + lngRad

				var tid uuid.UUID

				occCount := rand.Intn(2) + 1
				ocr := make([]train.CreateOccRequest, occCount)
				cp := []string{
					"#0D4A6D",
					"#0884C1",
					"#24C084",
					"#5CD3AC",
					"#B00020",
					"#FAB5A4",
					"#BB8779",
				}
				for i := 0; i < occCount; i++ {
					d := rand.Intn(6) + 1
					h := rand.Intn(10) + 7
					dr := rand.Intn(3) + 1
					ocr[i] = train.CreateOccRequest{
						OccRequest: train.OccRequest{
							DateStart:  time.Date(2020, 02, d, h, 0, 0, 0, time.UTC),
							DateEnd:    time.Date(2020, 02, d, h+dr, 0, 0, 0, time.UTC),
							RepeatDays: 7,
							Remarks:    helpers.CRNG_stringPanic(128),
							Color:      cp[rand.Int()%len(cp)],
						},
					}
				}

				if err := ctx.Api.TestAssertCaseErr(&api.TestCase{
					RequestMethod: "POST",
					RequestUrl:    "/api/training",
					RequestReader: helpers.JsonMustSerializeReader(train.CreateTrainingRequest{
						Training: train.TrainingRequest{
							Title:           t,
							Capacity:        rand.Intn(100) + 1,
							Price:           rand.Intn(5000) + 100,
							Currency:        "PLN",
							LocationCountry: "PL",
							Tags:            tags,
							LocationLat:     lat,
							LocationLng:     lng,
							RequiredGear: []string{
								"rękwice",
								"kąpielówki",
								"zestaw do nurkowania",
								"ochraniacz na zęby",
								"mleko w proszku",
							},
							RecommendedGear: []string{
								"ładne rękawice",
							},
							InstructorGear: []string{
								"nie wim",
							},
							Description:   "Zajebisty trening, dla wszystkich, którzy mogą i którym się chce, będziesz Pan zdrowszy i szczęśliwszy",
							ManualConfirm: rand.Intn(10)+1 > 5 && true || false,
							Diff: []int32{
								int32(rand.Intn(5) + 1),
							},
							LocationText: "Warszawa gdzieśtam",
						},
						ReturnID:    true,
						Occurrences: ocr,
					}),
					ExpectedBodyVal: func(b []byte, i interface{}) error {
						var err error
						tid, err = uuid.Parse(string(b))
						return err
					},
					ExpectedStatusCode: 200,
					RequestHeaders: map[string]string{
						"Authorization": "Bearer " + userTokens[z],
					},
				}); err != nil {
					fmt.Println(err)
					ret <- err
					return
				}

				if len(imgd) != 0 {
					x := rand.Int() % len(imgd)
					p := path.Join(imgDir, imgd[x].Name())
					if err = ctx.ApiUploadImgFromPath(userTokens[z], p, tid); err != nil {
						ret <- err
						return
					}
				}

			}
		}
		//ret <- ctx.Api.TestAssertCasesSemaphoreErr(apiReq, runtime.NumCPU())
		ret <- nil
	}()

	return ret
}

func main() {

	if err := getFaces(); err != nil {
		log.Fatal(err)
	}

	if err := getImgs(); err != nil {
		log.Fatal(err)
	}

	dbname := "sportdb"
	apiCtx := api.NewApi(config.NewLocalCtx())
	dalCtx := dal.NewDal(apiCtx.Config, dbname)

	var sample bool
	flag.BoolVar(&sample, "sample", false, "")
	flag.Parse()

	staticCtx := static.NewCtx(apiCtx)
	userCtx := user.NewCtx(apiCtx, dalCtx, staticCtx, nil)
	instrCtx := instr.NewCtx(apiCtx, dalCtx, userCtx, nil)
	trainCtx := train.NewCtx(apiCtx, dalCtx, userCtx, staticCtx, instrCtx)
	ic := 20
	if sample {
		ic = 2
	}
	tpc := 10
	if sample {
		tpc = 4
	}
	errc := createTrainingsAsync(trainCtx, ic, tpc)
	for {
		select {
		case err := <-errc:
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
}
