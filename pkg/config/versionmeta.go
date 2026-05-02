package config

import (
	"fmt"
	"strings"
	"time"

	"codeberg.org/emersion/go-scfg"
)

type VersionMeta struct {
	Created  uint64 `scfg:"created"`
	Hash     string `scfg:"hash"`
	Kind     string `scfg:"kind"`
	Original string `scfg:"original"`
	Size     int64  `scfg:"size"`
	Files    uint   `scfg:"filecount"`
	Source   string `scfg:"source"`
	Via      string `scfg:"via"`
}

const DefaultVersionMeta = `created %d
hash "%s"
kind "%s"
original "%s"
size %d
filecount %d
source "%s"
via "%s"
`

func NewVersionMeta(created time.Time, hash string, kind string, original string, size int64, filecount uint, source string, via string) (VersionMeta, string) {
	var config VersionMeta
	strConfig := fmt.Sprintf(DefaultVersionMeta, created.Unix(), hash, kind, original, size, filecount, source, via)
	scfg.NewDecoder(strings.NewReader(strConfig)).Decode(&config)
	return config, strConfig
}

func VersionMetaFromString(str string) (VersionMeta, error) {
	var config VersionMeta
	err := scfg.NewDecoder(strings.NewReader(str)).Decode(&config)
	return config, err
}
