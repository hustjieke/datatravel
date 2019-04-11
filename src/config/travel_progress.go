package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"xbase"

	"github.com/pkg/errors"
)

const (
	// travelProgressFile file path that stores travel progress
	travelProgressFile = "travelProgress.json"
)

// TravelProgress used to store progress rate msg during travel data
type TravelProgress struct {
	DumpProgressRate string `json:"dump-progress"`
	DumpRemainTime   string `json:"remain-time"`
	PositionBehinds  string `json:"position-behinds"`
	SynGTID          string `json:target syncer gtid`
	MasterGTID       string `json:src master gtid`
}

// UpdateTravelProgress used to update progress rate during travel data
func UpdateTravelProgress(travelProgress *TravelProgress, metaDir string) error {
	b, err := json.Marshal(travelProgress)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := os.MkdirAll(metaDir, os.ModePerm); err != nil {
		return err
	}
	file := filepath.Join(metaDir, travelProgressFile)
	return xbase.WriteFile(file, b)
}

// ReadTravelProgress used to read the config version from the file.
func ReadTravelProgress(metaDir string) (string, error) {
	if err := os.MkdirAll(metaDir, os.ModePerm); err != nil {
		return "", err
	}
	file := filepath.Join(metaDir, travelProgressFile)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Sprintf("%v", err), errors.WithStack(err)
	}

	return string(data), nil
}
