package main

import (
	"flag"
	"log"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
}

func main() {
	flag.Parse()

	e := Copy(from, to, offset, limit)

	if e != nil {
		log.Fatal(e)
	}
}
