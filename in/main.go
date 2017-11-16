package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/regevbr/keyval-resource/models"
	"fmt"
	"bufio"
	"sort"
)

func main() {
	if len(os.Args) < 2 {
		fatalNoErr("usage: " + os.Args[0] + " <destination>")
	}

	destination := os.Args[1]

	log("creating destination dir " + destination)
	err := os.MkdirAll(destination, 0755)
	if err != nil {
		fatal("creating destination", err)
	}

	output := filepath.Join(destination, "keyval.properties")
	log("creating output file " + destination)
	file, err := os.Create(output)
	if err != nil {
		fatal("creating output file", err)
	}

	defer file.Close()

	var request models.InRequest

	err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	var inVersion = request.Version

	w := bufio.NewWriter(file)

	var keys []string
	for k := range inVersion {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	log("writing " + strconv.Itoa(len(keys)) + " keys to output file")
	for _, k := range keys {
		fmt.Fprintf(w, "%s=%s\n", k, inVersion[k])
	}

	err = w.Flush()

	if err != nil {
		fatal("writing output file", err)
	}

	json.NewEncoder(os.Stdout).Encode(models.InResponse{
		Version:  inVersion,
	})

	log("Done")
}

func fatal(doing string, err error) {
	fmt.Fprintln(os.Stderr, "error " + doing + ": " + err.Error())
	os.Exit(1)
}

func log(doing string) {
	fmt.Fprintln(os.Stderr, doing)
}

func fatalNoErr(doing string) {
	log(doing)
	os.Exit(1)
}


