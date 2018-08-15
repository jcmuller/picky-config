package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"

	"github.com/jcmuller/gozenity"
	"github.com/jcmuller/picky/rule"
	"gopkg.in/yaml.v2"
)

var (
	configFilePath = fmt.Sprintf("%s/.config/picky/config.yaml", os.Getenv("HOME"))
)

type defaultProfile struct {
	Base    string `yaml:"base"`
	Args    []string `yaml:"args"`
}

type config struct {
	Debug          bool            `yaml:"debug"`
	DefaultProfile *defaultProfile `yaml:"default"`
	Rules          []*rule.Rule    `yaml:"rules"`
}

type pickyConfig struct {
	config   *config
	inputUri string
	uri      string
	profile  string
	rule     *rule.Rule
}

func fileContents(path string) (configFile []byte, err error) {
	configFile, err = ioutil.ReadFile(path)

	if os.IsNotExist(err) {
		return nil, errors.New("config not found")
	}

	return
}

func (pc *pickyConfig) getProfile() {
	profileNames := []string{}
	for _, rule := range pc.config.Rules {
		profileNames = append(profileNames, rule.Label)
	}

	sort.Strings(profileNames)

	profileNames = append(profileNames, "Add new profile")

	profile, err := gozenity.List("Choose profile: ", profileNames...)

	if err != nil {
		if err, ok := err.(*gozenity.EmptySelectionError); !ok {
			panic(err)
		} else {
			fmt.Println("No profile selected")
			os.Exit(0)
		}
	}

	pc.profile = profile
}

func GetConfig() (c *config) {
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

func (pc *pickyConfig) confirm() bool {
	answer, err := gozenity.Question(fmt.Sprintf(`Add %s to config?`, pc.uri))

	if !answer {
		fmt.Println("Not adding url")
		return false
	}

	if err != nil {
		fmt.Println("Assuming no.")
		return false
	}

	return true
}

func (pc *pickyConfig) editURI() {
	returnURL, err := gozenity.Entry(fmt.Sprintf(`Edit %s?`, pc.uri), pc.uri)

	if returnURL == "" {
		fmt.Printf("Not editing %s\n", pc.uri)
		return
	}

	if err != nil {
		fmt.Println("Assuming no.")
		return
	}

	pc.uri = returnURL
}

func (pc *pickyConfig) getRule() {
	for _, r := range pc.config.Rules {
		if r.Label == pc.profile {
			pc.rule = r
			return
		}
	}

	pc.rule = &rule.Rule{
		Label:   "New profile",
		Command: pc.config.DefaultProfile.Base,
		Args:    []string{"CHANGE ME"},
	}
	pc.config.Rules = append(pc.config.Rules, pc.rule)
}

func (pc *pickyConfig) addNewUriToRule() {
	pc.rule.URIs = append(pc.rule.URIs, pc.uri)
}

func (pc *pickyConfig) saveFile() {
	newFile, err := yaml.Marshal(pc.config)

	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	_, err = f.Write(newFile)
	fmt.Printf("%s\n", newFile)

	if err != nil {
		panic(err)
	}
}

func (pc *pickyConfig) openURI() {
	command, args := pc.rule.GetCommand()
	args = append(args, pc.inputUri)
	err := exec.Command(command, args...).Run()

	if err != nil {
		panic(err)
	}
}

// New creates an instance of this thing
func New(config *config, url string) (pc *pickyConfig) {
	return &pickyConfig{
		config:   config,
		uri:      url,
		inputUri: url,
	}
}

func (pc *pickyConfig) savePath() {
	pc.editURI()
	pc.getProfile()
	pc.getRule()
	pc.addNewUriToRule()
	pc.saveFile()
}

func (pc *pickyConfig) openPath() {
	pc.getProfile()
	pc.getRule()
}

// Call runs this thing
func (pc *pickyConfig) Call() {
	if pc.confirm() {
		pc.savePath()
	} else {
		pc.openPath()
	}

	pc.openURI()
}
