package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDynamicStyle(t *testing.T) {
	input := "abcdefghi"
	gota := GetDynamicStyle(input)
	gotb := GetDynamicStyle(input)
	// assert same style returned for the same string
	assert.Equal(t, gota.GetBackground(), gotb.GetBackground())
}
