package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		m, err := ReadDir("./testdata/env")
		fmt.Println(m)
		require.Nil(t, err)
	})
}
