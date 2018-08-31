package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func isExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func TestConvertCommands(t *testing.T) {
	tests := []struct {
		args       []string
		outputPath string
	}{
		{[]string{"kelp", "bmp", "testdata/dummy_png"}, "testdata/dummy_png.bmp"},
		{[]string{"kelp", "png", "testdata/dummy_png"}, "testdata/dummy_png.png"},
		{[]string{"kelp", "gif", "testdata/dummy_png"}, "testdata/dummy_png.gif"},
		{[]string{"kelp", "jpg", "testdata/dummy_png"}, "testdata/dummy_png.jpg"},
	}

	defer func() {
		for _, tt := range tests {
			os.Remove(tt.outputPath)
		}
	}()

	app := newApp()
	app.Writer = ioutil.Discard
	progressWriter = ioutil.Discard // supress progress bar
	for _, tt := range tests {
		err := app.Run(tt.args)
		assert.Nil(t, err)
		assert.Equal(t, isExists(tt.outputPath), true)
	}
}

func TestParallelConvert(t *testing.T) {
	tests := []struct {
		args        []string
		outputPaths []string
	}{
		{[]string{"kelp", "-j", "2", "png", "testdata/*.tmp"}, []string{"testdata/dummy_png.png"}},
		{[]string{"kelp", "-j", "2", "png", "testdata/**/*.tmp"}, []string{"testdata/dummy_png.png", "testdata/dir1/dummy_png.png", "testdata/dir1/dir2/dummy_png.png"}},
	}

	defer func() {
		for _, tt := range tests {
			for _, opath := range tt.outputPaths {
				os.Remove(opath)
			}
		}
	}()

	app := newApp()
	app.Writer = ioutil.Discard
	progressWriter = ioutil.Discard // supress progress bar
	for _, tt := range tests {
		err := app.Run(tt.args)
		assert.Nil(t, err)

		for _, opath := range tt.outputPaths {
			assert.Equal(t, isExists(opath), true)
			os.Remove(opath)
		}
	}
}

func TestExpandWildcardPath(t *testing.T) {
	tests := []struct {
		args        []string
		outputPaths []string
	}{
		{[]string{"kelp", "png", "testdata/*.tmp"}, []string{"testdata/dummy_png.png"}},
		{[]string{"kelp", "png", "testdata/**/*.tmp"}, []string{"testdata/dummy_png.png", "testdata/dir1/dummy_png.png", "testdata/dir1/dir2/dummy_png.png"}},
	}

	defer func() {
		for _, tt := range tests {
			for _, opath := range tt.outputPaths {
				os.Remove(opath)
			}
		}
	}()

	app := newApp()
	app.Writer = ioutil.Discard
	progressWriter = ioutil.Discard // supress progress bar
	for _, tt := range tests {
		err := app.Run(tt.args)
		assert.Nil(t, err)

		for _, opath := range tt.outputPaths {
			assert.Equal(t, isExists(opath), true)
			os.Remove(opath)
		}
	}
}

func TestForbiddenOverwrite(t *testing.T) {
	args := []string{"kelp", "png", "testdata/dummy_png"}
	outputPath := "testdata/dummy_png.png"

	defer func() {
		os.Remove(outputPath)
		log.SetOutput(os.Stdout)
		cli.OsExiter = os.Exit
		cli.ErrWriter = os.Stderr
	}()

	exitcode := 0
	cli.OsExiter = func(code int) {
		exitcode = code
	}
	cli.ErrWriter = ioutil.Discard

	app := newApp()
	app.Writer = ioutil.Discard
	log.SetOutput(ioutil.Discard)   // supress log output
	progressWriter = ioutil.Discard // supress progress bar

	err := app.Run(args)
	assert.Nil(t, err)

	err = app.Run(args)
	assert.NotNil(t, err)
	assert.NotEqual(t, exitcode, 0)
}
