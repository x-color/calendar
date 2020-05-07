package testutils

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func IgnoreKey(key string) cmp.Option {
	return cmpopts.IgnoreMapEntries(func(k string, t interface{}) bool {
		return k == key
	})
}
