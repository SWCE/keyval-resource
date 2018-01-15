package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/SWCE/keyval-resource/models"
	"github.com/magiconair/properties"
	"github.com/google/uuid"
	"fmt"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fatalNoErr("usage: " + os.Args[0] + " <destination>")
	}

	destination := os.Args[1]

	var request models.OutRequest

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	if request.Params.File != "" {
		inputFile := filepath.Join(destination, request.Params.File)
		log("reading input file " + inputFile)
		var data = properties.MustLoadFile(inputFile, properties.UTF8).Map()
		data["UPDATED"] = time.Now().Format(time.RFC850)
		data["UUID"] = uuid.New().String()
		log("read " + strconv.Itoa(len(data)) + " keys from input file")

		json.NewEncoder(os.Stdout).Encode(models.OutResponse{
			Version:  data,
		})
	} else {
		fatalNoErr("no properties file specified")
	}

}

func fatal(doing string, err error) {
	println("error " + doing + ": " + err.Error())
	os.Exit(1)
}

func log(doing string) {
	fmt.Fprintln(os.Stderr, doing)
}

func fatalNoErr(doing string) {
	log(doing)
	os.Exit(1)
}
