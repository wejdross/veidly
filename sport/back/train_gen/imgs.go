package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"golang.org/x/net/html"
)

const imgDir = "imgs"

func iread(n *html.Node, us *[]string) {
	for x := n.FirstChild; x != nil; x = x.NextSibling {
		if x.Data == "img" {
			for _, a := range x.Attr {
				if a.Key == "src" {
					*us = append(*us, a.Val)
				}
			}
		}
		iread(x, us)
	}
}

func getimg(u string, outpath string) error {
	pr, err := pexelsRequest(u)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outpath, pr, 0600)
}

func pexelsRequest(u string) ([]byte, error) {

	var err error

	req, err := http.NewRequest("GET", u, nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	//req.Header.Set("Cookie", "__cf_bm=b5d9fbbcb17d1b0786cd79ddd02892c74f984ea9-1624289448-1800-AdglUEteXmrCQQcjsG2TECmL9N695H15TQGNONon7MDhcPZxUcjCVE+8d+XflJal1MGWmenmorlhotjB6GH+/K8=; ab.storage.sessionId.5791d6db-4410-4ace-8814-12c903a548ba=%7B%22g%22%3A%22ca64bc56-590e-de80-e376-bb66c96d66b5%22%2C%22e%22%3A1624291249214%2C%22c%22%3A1624289449145%2C%22l%22%3A1624289449214%7D; ab.storage.deviceId.5791d6db-4410-4ace-8814-12c903a548ba=%7B%22g%22%3A%22b9878311-4a8f-e327-17d8-75822727cd07%22%2C%22c%22%3A1624289449146%2C%22l%22%3A1624289449146%7D; locale=en-US; NEXT_LOCALE=en-US; _fbp=fb.1.1624289449801.2128540019; _ga=GA1.2.965276669.1624289450; _gid=GA1.2.1711658147.1624289450; _gaexp=GAX1.2.V6wdJBU3R5uIG-k6WbHtRg.18887.0!tR3-05irSjCuHGCWvr4mHw.18888.0; _gat=1; _hjTLDTest=1; _hjid=141a07dd-1cc9-48da-ab47-eb8f1d005d7c; _hjFirstSeen=1; _hjIncludedInSessionSample=1; _hjAbsoluteSessionInProgress=0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("TE", "Trailers")

	if err != nil {
		return nil, err
	}
	c := &http.Client{}
	c.Transport = cloudflarebp.AddCloudFlareByPass(c.Transport)
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		if c, err := ioutil.ReadAll(res.Body); err == nil {
			ioutil.WriteFile("errout.html", c, 0600)
		}
		return nil, fmt.Errorf("invalid status code: %d", res.StatusCode)
	}
	cc, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// doc, err := html.Parse(bytes.NewBuffer(cc))
	// if err != nil {
	// 	return nil, err
	// }

	return cc, nil
}

func getImgs() error {
	_, err := os.Stat(imgDir)
	if err == nil {
		fmt.Printf("using imgs from local cache (%s)\n", imgDir)
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	if err := os.Mkdir(imgDir, 0700); err != nil {
		return err
	}

	var c *html.Node

	var cc []byte

	if cc, err = ioutil.ReadFile("ic.html"); err != nil {
		fmt.Println("warn: gotta refresh cache - making sport img request")
		cc, err = pexelsRequest("https://www.pexels.com/search/sport/")
		if err != nil {
			return err
		}
		err = ioutil.WriteFile("ic.html", cc, 0600)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("using img html from local cache")
	}
	c, err = html.Parse(bytes.NewBuffer(cc))

	us := make([]string, 0, 100)

	iread(c, &us)

	if len(us) == 0 {
		return fmt.Errorf("no imgs found")
	}

	s := make(chan struct{}, 16)
	t := make(chan error, len(us))

	count := 0

	for i := 0; i < len(us); i++ {

		u, err := url.Parse(us[i])
		if err != nil {
			fmt.Printf("warning %d: not an url (%s)\n", i, us[i])
			continue
		}
		n := path.Base(u.Path)
		if n == "" {
			fmt.Printf("warning %d: empty name (%s)\n", i, us[i])
			continue
		}
		if !strings.HasSuffix(n, ".jpeg") {
			continue
		}
		count++
		go func(i int) {
			s <- struct{}{}
			t <- getimg(us[i], path.Join("imgs", fmt.Sprintf("%d_%s", i, n)))
			<-s
		}(i)
	}

	var e error

	for i := 0; i < count; i++ {
		err := <-t
		if err != nil && e == nil {
			e = err
		}
	}

	if e == nil {
		fmt.Printf("got %d imgs", len(us))
	}

	return e
}
