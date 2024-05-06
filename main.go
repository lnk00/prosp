package main

import (
	"log"

	"github.com/lnk00/prosp/imap"
	"github.com/lnk00/prosp/parser"
	"github.com/lnk00/prosp/tui"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigFile("./conf.toml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error config file: %v", err)
	}

	imap := imap.New()
	imap.Login()
	messages := imap.FetchFrom("jobalerts-noreply@linkedin.com")
	jobs := parser.ParseAll(messages)
	tui.Render(jobs)
	imap.Logout()

}
