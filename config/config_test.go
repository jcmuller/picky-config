package config_test

import (
	"github.com/jcmuller/picky-config/config"
)

func Example_1() {
	c := config.GetConfig()
	url := "https://github.com/mmatczuk/go-http-tunnel"

	config.New(c, url).Call()

	// Output: Foo
}
