package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/logutils"
	"github.com/urfave/cli"
)

var (
	// allow overwrite destination file
	allowOverwrite = false

	// job thread number
	jobs = 1

	// is supress displayed progress bar
	hideProgress = false

	// output directory
	outDir = ""

	// use command pipe
	usePipe = false
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
		cli.BoolFlag{
			Name:        "force, f",
			Usage:       "allow file overwrite",
			Destination: &allowOverwrite,
		},
		cli.IntFlag{
			Name:        "jobs, j",
			Usage:       "convert job thread `n`umber",
			Value:       1,
			Destination: &jobs,
		},
		cli.BoolFlag{
			Name:        "no-progress",
			Usage:       "hide progress bar",
			Destination: &hideProgress,
		},
		cli.StringFlag{
			Name:        "outdir",
			Usage:       "output `dir`ectory",
			Value:       "",
			Destination: &outDir,
		},
		cli.BoolFlag{
			Name:        "pipe",
			Usage:       "input from stdin, output to stdout",
			Destination: &usePipe,
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

	// use pipe mode
	if usePipe {
		jobs = 1             // single job
		hideProgress = false // supress progress bar
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
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("WARN"),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)

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
    {{- if .VisibleFlags}}
      {{- " [options]"}}
    {{- end}}
    {{- if .Commands}}
      {{- " format [format options]"}}
    {{- end}}
    {{- if .ArgsUsage}}
      {{- .ArgsUsage}}
    {{- else}}
      {{- " [file...]"}}
    {{- end}}
  {{- "\n"}}

{{- if .VisibleFlags}}
OPTIONS:
  {{- "\n"}}
  {{- range .VisibleFlags}}
    {{- "  "}}
    {{- .}}
    {{- "\n"}}
  {{- end}}
{{- end}}

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
