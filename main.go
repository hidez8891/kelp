package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

var ()

func main() {
	app := cli.NewApp()

	// application infomation
	app.Name = "kelp"
	app.Usage = "simple image converter"
	app.Version = "0.1.0"

	// convert commands
	app.Commands = ConvertCommands
	app.HideHelp = true

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
