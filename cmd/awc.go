package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/bfontaine/lines/lines"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n\t$ %s [options...] <filename>\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// support '$ awc -l'
	flag.Bool("l", true, "Don't do anything; kept for compatibility with wc")
	flag.Parse()

	nargs := flag.NArg()
	if nargs < 1 {
		flag.Usage()
		os.Exit(1)
	}

	for _, filename := range flag.Args() {
		if err := processFilename(filename); err != nil {
			log.Fatal(err)
		}
	}
}

const (
	kb = 1024
	mb = 1024 * kb

	CHUNK_COUNT    = 100
	MIN_CHUNK_SIZE = 10 * mb
)

func processFilename(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	size := info.Size()
	chunkSize := size / CHUNK_COUNT

	readSize := size
	if chunkSize > MIN_CHUNK_SIZE {
		readSize = chunkSize
	}

	r := io.LimitReader(f, readSize)
	count, err := lines.CountFromReader(r)
	if err != nil {
		return err
	}

	count *= CHUNK_COUNT

	fmt.Printf("%d\n", count)
	return nil
}
