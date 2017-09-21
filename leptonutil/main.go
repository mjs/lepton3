package main

import (
	"fmt"
	"log"

	arg "github.com/alexflint/go-arg"
	"periph.io/x/periph/host"

	"github.com/TheCacophonyProject/lepton3"
)

type Options struct {
	Frames int    `arg:"-f,help:number of frames to collect (default=all)"`
	Output string `arg:"positional,required,help:png or none"`
}

func procCommandLine() Options {
	opts := Options{}
	arg.MustParse(&opts)
	if opts.Output != "png" && opts.Output != "none" {
		log.Fatalf("invalid output type: %q", opts.Output)
	}
	return opts
}

func main() {
	err := runMain()
	if err != nil {
		log.Fatal(err)
	}
}

func runMain() error {
	opts := procCommandLine()

	_, err := host.Init()
	if err != nil {
		return err
	}

	camera := lepton3.New()
	err = camera.Open()
	if err != nil {
		return err
	}
	defer camera.Close()

	im := lepton3.NewFrameImage()
	i := 0
	for {
		err := camera.NextFrame(im)
		if err != nil {
			return err
		}
		fmt.Printf(".")

		if opts.Output == "png" {
			err := dumpToPNG(fmt.Sprintf("%05d.png", i), im)
			if err != nil {
				return nil
			}
		}

		i++
		if opts.Frames > 0 && i >= opts.Frames {
			break
		}
	}
	fmt.Println()

	return nil
}