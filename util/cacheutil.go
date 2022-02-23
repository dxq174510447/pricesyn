package util

import (
	"github.com/hashicorp/golang-lru"
)

type cacheUtil struct {
}

func (d *cacheUtil) FormatNow() string {
	c1, _ := lru.New(100)
	c1.Add("a", "b")

	lru.New2Q(100)

	return ""
}

var CacheUtil cacheUtil = cacheUtil{}
