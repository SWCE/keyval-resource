package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/regevbr/keyval-resource/models"
	"fmt"
	"bufio"
	"sort"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
	}

	destination := os.Args[1]

	err := os.MkdirAll(destination, 0755)
	if err != nil {
		fatal("creating destination", err)
	}

	file, err := os.Create(filepath.Join(destination, "keyval.properties"))
	if err != nil {
		fatal("creating input file", err)
	}

	defer file.Close()

	var request models.InRequest

	err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	var inVersion = request.Version
	var metadata = models.Metadata{}

	w := bufio.NewWriter(file)

	delete(inVersion,"dummy")

	var keys []string
	for k := range inVersion {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "%s=%s\n", k, inVersion[k])
		metadata = append(metadata, models.MetadataField{
			Name:  k,
			Value: inVersion[k],
		})
	}

	err = w.Flush()

	if err != nil {
		fatal("writing input file", err)
	}

	json.NewEncoder(os.Stdout).Encode(models.InResponse{
		Version:  inVersion,
		Metadata: metadata,
	})
}

func fatal(doing string, err error) {
	println("error " + doing + ": " + err.Error())
	os.Exit(1)
}
