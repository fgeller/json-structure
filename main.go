package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	flatten := flag.Bool("flatten", false, "flatten schemata in arrays and combine objects")
	js := flag.Bool("schema", false, "output a json schema")
	flag.Parse()

	var d any
	err := json.NewDecoder(os.Stdin).Decode(&d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode json err=%v\n", err)
		os.Exit(1)
	}

	s := schema(d, *flatten)
	s.Schema = "https://json-schema.org/draft/2020-12/schema"
	res := any(s)
	if !*js {
		res = s.toSimple()
	}

	err = json.NewEncoder(os.Stdout).Encode(&res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode to json err=%v\n", err)
		os.Exit(1)
	}
}
