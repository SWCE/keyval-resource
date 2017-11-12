package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/regevbr/keyval-resource/models"
)

func main() {
	var request models.CheckRequest
	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err.Error())
		os.Exit(1)
	}
	versions := []models.EmptyVersion{}
	versions = append(versions, models.EmptyVersion{})
	json.NewEncoder(os.Stdout).Encode(versions)
}
