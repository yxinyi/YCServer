package YUIDFactory

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func BuildUIDUint64()uint64{
	return rand.Uint64()
}