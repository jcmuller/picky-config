package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"sort"

	yaml "gopkg.in/yaml.v2"

	"github.com/atotto/clipboard"
	"github.com/jcmuller/dmenu"
)

var (
	configFilePath = fmt.Sprintf("%s/.config/picky/config.yaml", os.Getenv("HOME"))
)

type defaultProfile struct {
	Base    string `yaml:"base"`
	Profile string `yaml:"profile"`
	Args    string `yaml:"args"`
}

type rule struct {
	Label   string   `yaml:"label"`
	Base    string   `yaml:"base"`
	Profile string   `yaml:"profile"`
	Args    string   `yaml:"args"`
	URIs    []string `yaml:"uris"`
}

type config struct {
	Debug          bool            `yaml:"debug"`
	DefaultProfile *defaultProfile `yaml:"default"`
	Rules          []*rule         `yaml:"rules"`
}

func fileContents(path string) (configFile []byte, err error) {
	configFile, err = ioutil.ReadFile(path)

	if os.IsNotExist(err) {
		return nil, errors.New("config not found")
	}

	return
}

func getProfile(c *config) (profile string, err error) {
	profileNames := []string{}
	for _, rule := range c.Rules {
		profileNames = append(profileNames, rule.Label)
	}

	sort.Strings(profileNames)

	profileNames = append(profileNames, "Add new profile")

	profile, err = dmenu.Popup("Choose profile: ", profileNames...)

	if err != nil {
		if err, ok := err.(*dmenu.EmptySelectionError); !ok {
			panic(err)
		} else {
			fmt.Println("No profile selected")
			os.Exit(0)
		}
	}

	return
}

func readConfig() (c *config) {
	c = &config{}
	configFile, err := fileContents(configFilePath)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configFile, c)

	if err != nil {
		panic(err)
	}

	return
}

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

func confirm(uri string) {
	yesNo := []string{"yes", "no"}
	answer, err := dmenu.Popup(fmt.Sprintf(`Add_%s_to_config?`, uri), yesNo...)

	if answer == "no" {
		fmt.Println("Not adding url")
		os.Exit(0)
	}

	if err != nil {
		if err, ok := err.(*dmenu.EmptySelectionError); !ok {
			panic(err)
		} else {
			fmt.Println("Assuming no.")
			os.Exit(0)
		}
	}
}

func getRule(c *config, profile string) (rr *rule) {
	for _, r := range c.Rules {
		if r.Label == profile {
			rr = r
			return
		}
	}

	rr = &rule{
		Label:   "New profile",
		Base:    c.DefaultProfile.Base,
		Profile: c.DefaultProfile.Profile,
		Args:    "CHANGE ME",
	}
	c.Rules = append(c.Rules, rr)

	return
}

func saveFile(c *config) {
	newFile, err := yaml.Marshal(c)

	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	_, err = f.Write(newFile)

	if err != nil {
		panic(err)
	}
}

func main() {
	config := readConfig()
	uri := getURL()
	confirm(uri)
	profile, _ := getProfile(config)
	rule := getRule(config, profile)

	rule.URIs = append(rule.URIs, uri)

	saveFile(config)
}
