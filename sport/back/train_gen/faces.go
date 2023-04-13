package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

const faceDir = "avatars"

func geturl(url, name string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	c, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(faceDir, name), c, 0600); err != nil {
		return err
	}
	return nil
}

func getFacePair(i int) error {
	womenUrl := fmt.Sprintf("https://randomuser.me/api/portraits/women/%d.jpg", i)
	womenFile := fmt.Sprintf("w%d.jpg", i)
	if err := geturl(womenUrl, womenFile); err != nil {
		return err
	}
	menUrl := fmt.Sprintf("https://randomuser.me/api/portraits/men/%d.jpg", i)
	menFile := fmt.Sprintf("m%d.jpg", i)
	if err := geturl(menUrl, menFile); err != nil {
		return err
	}
	return nil
}

func getFaces() error {

	_, err := os.Stat(faceDir)
	if err == nil {
		fmt.Printf("using faces from local cache (%s)\n", faceDir)
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	if err := os.Mkdir(faceDir, 0700); err != nil {
		return err
	}

	// get 50 women and 50 men
	max := 50

	fmt.Println("getting faces")

	maxThread := 16

	s := make(chan struct{}, maxThread)
	t := make(chan error, max)

	for i := 1; i <= max; i++ {
		go func(i int) {
			s <- struct{}{}
			t <- getFacePair(i)
			<-s
		}(i)
	}

	var e error

	for i := 0; i < max; i++ {
		err := <-t
		if err != nil && e == nil {
			e = err
		}
	}

	if e == nil {
		fmt.Printf("got %d faces", max*2)
	}

	return e
}
