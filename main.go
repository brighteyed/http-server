package main

import (
	"flag"
)

func main() {
	root := flag.String("d", ".", "the directory of files to host")
	port := flag.String("p", "8100", "port to serve on")
	idle := flag.Uint("t", 0, "duration before shutdown while inactive (0 – disable)")
	flag.Parse()

	run(root, idle, port)
}
