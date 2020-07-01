package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	_, err := os.Stat("test.sh")
	if err != nil {
		t.Error(err)
	}
	cpath, err := filepath.Abs("./")
	if err != nil {
		t.Error(err)
	}
	cmd := exec.Command(path.Join(cpath, "test.sh"))
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	if err = cmd.Run(); err != nil {
		t.Error(err)
	}
	ok := strings.HasSuffix(buf.String(), "PASS\n")
	require.True(t, ok, "sh tests not PASS")
}

func TestMergeEnvironments(t *testing.T) {
	t.Run("merge environments add new var", func(t *testing.T) {
		env := mergeEnvironments([]string{}, map[string]string{"NAME": "Value"})

		require.Len(t, env, 1)
		require.Equal(t, env[0], "NAME=Value")
	})
	t.Run("merge environments update var if exists", func(t *testing.T) {
		env := mergeEnvironments([]string{"NAME=Value"}, map[string]string{"NAME": "NEW_Value"})

		require.Len(t, env, 1)
		require.Equal(t, env[0], "NAME=NEW_Value")
	})
	t.Run("merge environments drop value", func(t *testing.T) {
		env := mergeEnvironments([]string{"NAME=Value"}, map[string]string{"NAME": ""})

		require.Len(t, env, 0)
	})
}
