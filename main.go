package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
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

	return app
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

{{- if .Version}}
VERSION:
   {{.Version}}
{{- end}}
`
}
