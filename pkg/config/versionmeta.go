package config

import (
	"fmt"
	"strings"
	"time"

	"codeberg.org/emersion/go-scfg"
)

type VersionMeta struct {
	Created uint64 `scfg:"created"`
	Hash    string `scfg:"hash"`
	Size    int64  `scfg:"size"`
	Files   uint   `scfg:"filecount"`
	Via     string `scfg:"via"`
}

const DefaultVersionMeta = `created %d
hash "%s"
size %d
filecount %d
via "%s"
`

func NewVersionMeta(created time.Time, hash string, size int64, filecount uint, via string) (VersionMeta, string) {
	var config VersionMeta
	strConfig := fmt.Sprintf(DefaultVersionMeta, created.Unix(), hash, size, filecount, via)
	scfg.NewDecoder(strings.NewReader(strConfig)).Decode(&config)
	return config, strConfig
}

func VersionMetaFromString(str string) (VersionMeta, error) {
	var config VersionMeta
	err := scfg.NewDecoder(strings.NewReader(str)).Decode(&config)
	return config, err
}
