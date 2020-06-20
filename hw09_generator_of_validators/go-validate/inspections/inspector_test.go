package inspections

import (
	"github.com/stretchr/testify/require"
	"go/parser"
	"go/token"
	"testing"
)

func TestPrivate(t *testing.T) {
	t.Run("decl struct without needed tag", func(t *testing.T) {
		fs := token.NewFileSet()
		file, err := parser.ParseFile(fs, "", "package main; type Name struct {I int64}", parser.AllErrors)
		if err != nil {
			t.Error(err)
			return
		}

		res := checkNode(file.Decls[0], nil)
		require.True(t, res)
	})

	t.Run("only one decl. Is not struct", func(t *testing.T) {
		fs := token.NewFileSet()
		file, err := parser.ParseFile(fs, "", "package main; type Name interface {f()}", parser.AllErrors)
		if err != nil {
			t.Error(err)
			return
		}

		res := checkNode(file.Decls[0], nil)
		require.True(t, res)
	})

	t.Run("correct decl of struct", func(t *testing.T) {
		fs := token.NewFileSet()
		file, err := parser.ParseFile(fs, "", "package main; type Name struct {I int64 `validate:\"min:18\"`}", parser.AllErrors)
		if err != nil {
			t.Error(err)
			return
		}

		res := checkNode(file.Decls[0], nil)
		require.False(t, res)
	})

	t.Run("", func(t *testing.T) {

	})
}
