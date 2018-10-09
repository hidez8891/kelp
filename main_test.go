package main

import (
	"io/ioutil"
	"io"
	"log"
	"os"
	"testing"
	"bytes"

	"image/jpeg"

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

func setupApp(writer io.Writer) *cli.App {
	app := newApp()
	app.Writer = writer

	progressWriter = writer
	log.SetOutput(writer)

	return app
}

func testHelperRunApp(t *testing.T, tests []testData) {
	t.Helper()

	app := setupApp(ioutil.Discard)
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
			[]string{"kelp", "bmp", "testdata/dummy_png.tmp"},
			[]string{"testdata/dummy_png.bmp"},
		},
		{
			[]string{"kelp", "png", "testdata/dummy_png.tmp"},
			[]string{"testdata/dummy_png.png"},
		},
		{
			[]string{"kelp", "gif", "testdata/dummy_png.tmp"},
			[]string{"testdata/dummy_png.gif"},
		},
		{
			[]string{"kelp", "jpg", "testdata/dummy_png.tmp"},
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
		{
			[]string{"kelp", "png", "testdata/**/*.{tmp,tmp2}"},
			[]string{
				"testdata/dummy_png.png",
				"testdata/dummy_jpg.png",
				"testdata/dir1/dummy_png.png",
				"testdata/dir1/dir2/dummy_png.png",
				"testdata/dir1/dir2/dummy_jpg.png",
			},
		},
	}

	testHelperRunApp(t, tests)
}

func TestForbiddenOverwrite(t *testing.T) {
	args := []string{"kelp", "png", "testdata/dummy_png.tmp"}
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

	app := setupApp(ioutil.Discard)

	// create file
	err := app.Run(args)
	assert.Nil(t, err)

	// forbid create duplicate file
	err = app.Run(args)
	assert.NotNil(t, err)
	assert.NotEqual(t, exitcode, 0)
}

func TestAllowOverwrite(t *testing.T) {
	args := []string{"kelp", "-f", "png", "testdata/dummy_png.tmp"}
	outputPath := "testdata/dummy_png.png"

	defer os.Remove(outputPath)
	app := setupApp(ioutil.Discard)

	// create file
	err := app.Run(args)
	assert.Nil(t, err)

	// allow create duplicate file
	err = app.Run(args)
	assert.Nil(t, err)
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

func TestPipeIO(t *testing.T) {
	args := []string{"kelp", "--pipe", "jpg"}
	src := "testdata/dummy_png.tmp"

	r, err := os.Open(src)
	assert.Nil(t, err)
	defer r.Close()

	w := new(bytes.Buffer)

	// watch progress/stderr output
	progress := new(bytes.Buffer)

	// pipe execute
	pipeStdin = r
	pipeStdout = w
	app := setupApp(progress)

	err = app.Run(args)
	assert.Nil(t, err)

	// check progress/stderr output
	assert.Equal(t, 0, progress.Len())

	// check
	_, err = jpeg.Decode(bytes.NewReader(w.Bytes()))
	assert.Nil(t, err)
}
