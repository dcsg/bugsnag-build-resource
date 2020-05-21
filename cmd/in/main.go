package main

import (
	"encoding/json"
	"os"

	resource "github.com/dcsg/bugsnag-build-resource"
)

func main() {
	response := resource.InResponse{
		Version: resource.Version{
			Id: "unsupported",
		},
	}

	err := json.NewEncoder(os.Stdout).Encode(response)
	if err != nil {
		panic(err)
	}
}
