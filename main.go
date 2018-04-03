package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/atotto/clipboard"
	"github.com/jcmuller/dmenu"
)

var (
	profilesMap = map[string]string{
		//"Personal": "Profile 3",
		"GH":      "Profile 4",
		"SHED":    "Profile 2",
		"Zendesk": "Profile 1",
		"Twitter": "Profile 5",
	}
	configFilePath = fmt.Sprintf("%s/.config/picky/config.yaml", os.Getenv("HOME"))
)

type rule struct {
	Base    string   `yaml:"base"`
	Profile string   `yaml:"profile"`
	Args    string   `yaml:"args"`
	URIs    []string `yaml:"uris"`
}

type config struct {
	Rules []*rule `yaml:"rules"`
}

func fileContents(path string) (configFile []byte, err error) {
	configFile, err = ioutil.ReadFile(path)

	if os.IsNotExist(err) {
		return nil, errors.New("config not found")
	}

	return
}

func getProfile() (profile string, err error) {
	profileNames := make([]string, 0, len(profilesMap))
	for name := range profilesMap {
		profileNames = append(profileNames, name)
	}

	profileName, err := dmenu.Popup("Choose profile: ", profileNames...)
	profile = profilesMap[profileName]

	if err != nil {
		panic(err)
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
		panic(err)
	}
}

func getRule(c *config, profile string) (rule *rule) {
	for _, r := range c.Rules {
		if r.Args == profile {
			rule = r
			return
		}
	}

	return nil
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
	profile, _ := getProfile()
	rule := getRule(config, profile)

	rule.URIs = append(rule.URIs, uri)

	saveFile(config)
}
