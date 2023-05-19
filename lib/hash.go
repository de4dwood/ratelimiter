package lib

import (
	"hash/fnv"
	"strconv"
)

func Hash(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return strconv.FormatUint(h.Sum64(), 36)
}
