package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	dedupe := flag.Bool("dedupe", false, "dedupe elements of in arrays")
	mergeObjects := flag.Bool("merge-objects", false, "merge objects in arrays")
	flag.Parse()

	var d any
	err := json.NewDecoder(os.Stdin).Decode(&d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode json err=%v\n", err)
		os.Exit(1)
	}
	pt := identify(d, *dedupe, *mergeObjects)

	err = json.NewEncoder(os.Stdout).Encode(&pt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode to json err=%v\n", err)
		os.Exit(1)
	}
}
