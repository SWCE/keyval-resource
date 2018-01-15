package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SWCE/keyval-resource/models"
)

func main() {
	var request models.CheckRequest
	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err.Error())
		os.Exit(1)
	}

	versions := []models.Version{}

	versions = append(versions, request.Version)

	json.NewEncoder(os.Stdout).Encode(versions)

	//json.NewEncoder(os.Stdout).Encode(models.CheckResponse{
	//	Version:  request.Version["updated"],
	//})
	//versions := []models.EmptyVersion{}
	//json.NewEncoder(os.Stdout).Encode(versions)
}
