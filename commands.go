package main

import (
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"path/filepath"

	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/mattn/go-zglob"
	"github.com/urfave/cli"
	"go.uber.org/atomic"
	"golang.org/x/image/bmp"
)

var (
	// image quality
	quality = 0

	// ConvertCommands is image converter commands
	ConvertCommands = make([]cli.Command, 0)

	// source image file path list
	targetFilePaths = make([]string, 0)
)

// common convert encoder function interface
type convertEncodeFn = func(w io.Writer, m image.Image) error

// convert function
func convert(ctx *cli.Context, converter convertEncodeFn) error {
	partialFail := atomic.NewBool(false)

	Concurrent(jobs, targetFilePaths, func(srcPath string) {
		r, err := os.Open(srcPath)
		if err != nil {
			log.Println(fmt.Sprintf("error: %v [%s]", err, srcPath))
			partialFail.Store(true)
			return
		}
		defer r.Close()

		img, _, err := image.Decode(r)
		if err != nil {
			log.Println(fmt.Sprintf("error: %v [%s]", err, srcPath))
			partialFail.Store(true)
			return
		}

		destPath := generateDestinationPath(srcPath, ctx.Command.Name)
		w, err := os.OpenFile(destPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
		if err != nil {
			log.Println(fmt.Sprintf("error: %v [%s]", err, destPath))
			partialFail.Store(true)
			return
		}
		defer w.Close()

		err = converter(w, img)
		if err != nil {
			log.Println(fmt.Sprintf("error: %v [%s]", err, srcPath))
			partialFail.Store(true)
			return
		}
	})

	if partialFail.Load() {
		return cli.NewExitError("warn: some images failed to convert", 1)
	}
	return nil
}

// generate destination image file path
func generateDestinationPath(srcPath string, destExt string) string {
	ext := filepath.Ext(srcPath)
	if len(ext) > 0 {
		srcPath = srcPath[0 : len(srcPath)-len(ext)]
	}
	return srcPath + "." + destExt
}

// fetch source file path list from command-line arguments
func fetchSourceFilePaths(ctx *cli.Context) error {
	// check input files [need at least one]
	if len(ctx.Args()) == 0 {
		return cli.NewExitError("need at least one source file", 1)
	}

	// expand wildcard path
	targetFilePaths = make([]string, 0)
	for _, path := range ctx.Args() {
		paths := expandFilePath(path)
		targetFilePaths = append(targetFilePaths, paths...)
	}
	if len(targetFilePaths) == 0 {
		return cli.NewExitError("need at least one source file", 1)
	}

	return nil
}

// preprocess command flags & args
func preprocessCommandArgs(ctx *cli.Context) error {
	if err := validateFlags(ctx); err != nil {
		return err
	}

	if err := fetchSourceFilePaths(ctx); err != nil {
		return err
	}

	return nil
}

// validate jpeg flags
func validateJpegFlags(ctx *cli.Context) error {
	// check quality range
	if quality < 1 || quality > 100 {
		return cli.NewExitError("invalid jpeg quality [%d]", quality)
	}

	return nil
}

// expand file path (use wildcard)
func expandFilePath(path string) []string {
	mathces, err := zglob.Glob(path)
	if err != nil {
		log.Printf("warn: %v\n", err)
	}

	return mathces
}

func init() {
	ConvertCommands = []cli.Command{
		{
			// bitmap encoder command
			Name:  "bmp",
			Usage: "convert to BMP format",
			Action: func(ctx *cli.Context) error {
				if err := preprocessCommandArgs(ctx); err != nil {
					return err
				}

				encoder := func(w io.Writer, m image.Image) error {
					return bmp.Encode(w, m)
				}
				return convert(ctx, encoder)
			},
		},
		{
			// PNG encoder command
			Name:  "png",
			Usage: "convert to PNG format",
			Action: func(ctx *cli.Context) error {
				if err := preprocessCommandArgs(ctx); err != nil {
					return err
				}

				encoder := func(w io.Writer, m image.Image) error {
					return png.Encode(w, m)
				}
				return convert(ctx, encoder)
			},
		},
		{
			// GIF encoder command
			Name:  "gif",
			Usage: "convert to gif format",
			Action: func(ctx *cli.Context) error {
				if err := preprocessCommandArgs(ctx); err != nil {
					return err
				}

				opts := &gif.Options{}
				encoder := func(w io.Writer, m image.Image) error {
					return gif.Encode(w, m, opts)
				}
				return convert(ctx, encoder)
			},
		},
		{
			// JPEG encoder command
			Name:  "jpg",
			Usage: "convert to JPEG format",
			Action: func(ctx *cli.Context) error {
				if err := preprocessCommandArgs(ctx); err != nil {
					return err
				}
				if err := validateJpegFlags(ctx); err != nil {
					return err
				}

				opts := &jpeg.Options{
					Quality: quality,
				}
				encoder := func(w io.Writer, m image.Image) error {
					return jpeg.Encode(w, m, opts)
				}
				return convert(ctx, encoder)
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "q",
					Usage:       "jpeg quality",
					Value:       90,
					Destination: &quality,
				},
			},
		},
	}
}
