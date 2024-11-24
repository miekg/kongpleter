package main

import (
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/miekg/kongpleter"
)

var cli struct {
	Debug bool `help:"Debug mode."`

	Complete struct {
		Something string `completion:"echo 1 2 3 4" help:"more"`

		Bla  string `arg:"" completion:"bloap"`
		Name string `arg:"" completion:"Name"`
	} `cmd:"" help:"Complete." cmdaliases:"ListComplete"`
}

func main() {
	parser, err := kong.New(&cli,
		kong.Name("shell"),
		kong.Description("A shell-like example app."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))
	kongpleter.Walk("shell", parser)

	ctx, err := parser.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	switch ctx.Command() {
	case "rm <path>":

	case "ls":
	}
}
