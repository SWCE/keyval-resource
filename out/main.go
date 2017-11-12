package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/regevbr/keyval-resource/models"
	"sort"
	"github.com/magiconair/properties"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
	}

	destination := os.Args[1]

	var request models.OutRequest

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	var data = properties.MustLoadFile(filepath.Join(destination, "keyval.properties"), properties.UTF8).Map();
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var metadata = models.Metadata{}

	for _, k := range keys {
		metadata = append(metadata, models.MetadataField{
			Name:  k,
			Value: data[k],
		})
	}

	json.NewEncoder(os.Stdout).Encode(models.OutResponse{
		Version:  data,
		Metadata: metadata,
	})
}

func fatal(doing string, err error) {
	println("error " + doing + ": " + err.Error())
	os.Exit(1)
}
