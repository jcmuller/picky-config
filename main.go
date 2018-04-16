package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/jcmuller/picky-config/config"
)

func getURL() (uri string) {
	uri, err := clipboard.ReadAll()

	if err != nil {
		panic(err)
	}

	_, err = url.ParseRequestURI(uri)

	if err != nil {
		fmt.Printf("Invalid url (%s)\n", uri)
		os.Exit(0)
	}

	return
}

func logURL(url string) {
	f, _ := os.OpenFile("/tmp/urls", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()

	f.WriteString(fmt.Sprintf("%#v\n", url))
}

func main() {
	var url string
	c := config.GetConfig()

	if len(os.Args) < 2 {
		url = getURL()
	} else {
		url = strings.Join(os.Args[1:], "")
	}

	logURL(url)

	config.New(c, url).Call()
}
