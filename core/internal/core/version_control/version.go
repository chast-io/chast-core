package version_control

import (
	"hash/fnv"
	"io/fs"
	"log"
	"os"
	"strconv"
	"time"
)

type Version struct {
	source string
	target string
}

func NewVersion(basePath string, source string, versionTag string) *Version {
	now := time.Now()
	timeStamp := strconv.FormatInt(now.UnixMilli(), 10)
	versionTagHash := hash(versionTag)
	target := basePath + timeStamp + versionTagHash
	return &Version{source: source, target: target}
}

func (v *Version) CreateTargetFolder() {
	if err := os.MkdirAll(v.target, fs.ModeDir); err != nil {
		log.Fatal(err)
	}
}

func hash(s string) string {
	h := fnv.New32()
	_, err := h.Write([]byte(s))
	if err != nil {
		panic(err)
	}
	return strconv.FormatUint(uint64(h.Sum32()), 10)
}
