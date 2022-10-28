package version_control

import (
	"github.com/google/uuid"
	"hash/fnv"
	"io/fs"
	"log"
	"os"
	"strconv"
)

type Version struct {
	UUID            string
	PreviousVersion *Version
	source          string
	target          string
}

func NewVersion(basePath string, source string, versionTag string) *Version {
	extendedUuid := versionTag + "-" + uuid.New().String()

	return &Version{
		UUID:   extendedUuid,
		source: source,
		//target: target,
	}
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
