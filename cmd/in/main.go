package main

import (
	"encoding/json"
	"os"

	resource "github.com/dcsg/bugsnag-build-resource"
)

func main() {
	var (
		inRequest  resource.InRequest
		inResponse resource.InResponse
	)

	err := json.NewDecoder(os.Stdin).Decode(&inRequest)
	if err != nil {
		panic(err)
	}

	inResponse = resource.InResponse{
		Version: inRequest.Version,
	}

	err = json.NewEncoder(os.Stdout).Encode(inResponse)
	if err != nil {
		panic(err)
	}
}
