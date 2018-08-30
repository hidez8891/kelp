package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var (
	// job thread number
	jobs = 1
)

// newApp returns *cli.App [use testing]
func newApp() *cli.App {
	app := cli.NewApp()

	// application infomation
	app.Name = _Name
	app.Usage = _Description
	app.Version = _Version

	// convert commands
	app.Commands = ConvertCommands
	app.HideHelp = true

	// global options
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "jobs, j",
			Usage:       "convert job thread number",
			Value:       1,
			Destination: &jobs,
		},
	}

	return app
}

// validate flags
func validateFlags(ctx *cli.Context) error {
	// check jobs range
	if jobs < 1 {
		msg := fmt.Sprintf("invalid job thread number [%d]", jobs)
		return cli.NewExitError(msg, 1)
	}

	return nil
}

func main() {
	app := newApp()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	cli.AppHelpTemplate = `
{{- "NAME:"}}
  {{.Name}} - {{.Usage}}
  {{- if .Version}}
    {{- " [version "}}
    {{- .Version}}
    {{- "]"}}
  {{- end}}

USAGE:
  {{.HelpName}}
    {{- if .Commands}}
      {{- " format [format options]"}}
    {{- end}}
    {{- if .ArgsUsage}}
      {{- .ArgsUsage}}
    {{- else}}
      {{- " [file...]"}}
    {{- end}}
  {{- "\n"}}

{{- if .Commands}}
FORMAT:
  {{- "\n"}}
  {{- range .Commands}}
    {{- if not .HideHelp}}
      {{- "  "}}
      {{- join .Names ", "}}
      {{- "\t"}}
      {{- .Usage}}
      {{- "\n"}}
    {{- end}}
  {{- end}}
{{- end}}

{{- if .Commands}}
FORMAT OPTIONS:
  {{- "\n"}}
  {{- range .Commands}}
    {{- if .VisibleFlags}}
      {{- "  "}}
      {{- join .Names ", "}}
      {{- "\n"}}
      {{- range .VisibleFlags}}
        {{- "    "}}
        {{- .}}
        {{- "\n"}}
      {{- end}}
    {{- end}}
  {{- end}}
{{- end}}
`
}
