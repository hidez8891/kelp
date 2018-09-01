package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

type testData struct {
	args        []string
	outputPaths []string
}

func isExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func setupApp() *cli.App {
	app := newApp()
	app.Writer = ioutil.Discard

	progressWriter = ioutil.Discard // supress progress bar
	log.SetOutput(ioutil.Discard)   // supress log output

	return app
}

func testHelperRunApp(t *testing.T, tests []testData) {
	t.Helper()

	app := setupApp()
	for _, tt := range tests {
		err := app.Run(tt.args)
		assert.Nil(t, err)

		for _, opath := range tt.outputPaths {
			assert.Equal(t, isExists(opath), true)
			os.Remove(opath)
		}
	}
}

func TestConvertCommands(t *testing.T) {
	tests := []testData{
		{
			[]string{"kelp", "bmp", "testdata/dummy_png"},
			[]string{"testdata/dummy_png.bmp"},
		},
		{
			[]string{"kelp", "png", "testdata/dummy_png"},
			[]string{"testdata/dummy_png.png"},
		},
		{
			[]string{"kelp", "gif", "testdata/dummy_png"},
			[]string{"testdata/dummy_png.gif"},
		},
		{
			[]string{"kelp", "jpg", "testdata/dummy_png"},
			[]string{"testdata/dummy_png.jpg"},
		},
	}

	testHelperRunApp(t, tests)
}

func TestParallelConvert(t *testing.T) {
	tests := []testData{
		{
			[]string{"kelp", "-j", "2", "png", "testdata/*.tmp"},
			[]string{"testdata/dummy_png.png"},
		},
		{
			[]string{"kelp", "-j", "2", "png", "testdata/**/*.tmp"},
			[]string{
				"testdata/dummy_png.png",
				"testdata/dir1/dummy_png.png",
				"testdata/dir1/dir2/dummy_png.png",
			},
		},
	}

	testHelperRunApp(t, tests)
}

func TestExpandWildcardPath(t *testing.T) {
	tests := []testData{
		{
			[]string{"kelp", "png", "testdata/*.tmp"},
			[]string{"testdata/dummy_png.png"},
		},
		{
			[]string{"kelp", "png", "testdata/**/*.tmp"},
			[]string{
				"testdata/dummy_png.png",
				"testdata/dir1/dummy_png.png",
				"testdata/dir1/dir2/dummy_png.png",
			},
		},
	}

	testHelperRunApp(t, tests)
}

func TestForbiddenOverwrite(t *testing.T) {
	args := []string{"kelp", "png", "testdata/dummy_png"}
	outputPath := "testdata/dummy_png.png"

	defer func() {
		os.Remove(outputPath)
		cli.OsExiter = os.Exit
		cli.ErrWriter = os.Stderr
	}()

	exitcode := 0
	cli.OsExiter = func(code int) {
		exitcode = code
	}
	cli.ErrWriter = ioutil.Discard

	app := setupApp()

	// create file
	err := app.Run(args)
	assert.Nil(t, err)

	// create duplicate file
	err = app.Run(args)
	assert.NotNil(t, err)
	assert.NotEqual(t, exitcode, 0)
}

func TestSetOutputDirectory(t *testing.T) {
	testdir := "./test_out"

	tests := []testData{
		{
			[]string{"kelp", "--outdir", testdir, "png", "testdata/*.tmp"},
			[]string{testdir + "/testdata/dummy_png.png"}},
		{
			[]string{"kelp", "--outdir", testdir, "png", "testdata/**/*.tmp"},
			[]string{
				testdir + "/testdata/dummy_png.png",
				testdir + "/testdata/dir1/dummy_png.png",
				testdir + "/testdata/dir1/dir2/dummy_png.png",
			},
		},
	}

	defer os.RemoveAll(testdir)
	testHelperRunApp(t, tests)
}
