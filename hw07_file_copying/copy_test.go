package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("not regular input", func(t *testing.T) {
		err := Copy("/dev/null", "", 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("offset largest then file size", func(t *testing.T) {
		tmpFile, closef, err := tempFile()
		if err != nil {
			t.Error(t)
			return
		}
		defer func() {
			err = closef()
			if err != nil {
				t.Error(err)
			}
		}()

		err = Copy(tmpFile.Name(), os.TempDir()+"filename", 1<<20, 0)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("simple copy", func(t *testing.T) {
		data := "data in temp file"
		tmpFile, closef, err := tempFile()
		if err != nil {
			t.Error(t)
			return
		}
		defer func() {
			err = closef()
			if err != nil {
				t.Error(err)
			}
		}()
		_, err = tmpFile.WriteString(data)
		if err != nil {
			t.Error(err)
			return
		}

		outName := os.TempDir() + "filename.data"
		err = Copy(tmpFile.Name(), outName, 0, 0)

		res, err := ioutil.ReadFile(outName)
		if err != nil {
			t.Error(err)
			return
		}
		require.Nil(t, err)
		require.Equal(t, data, string(res))
	})

	t.Run("sh test run", func(t *testing.T) {
		_, err := os.Stat("test.sh")
		if err != nil {
			return
		}
		cpath, err := filepath.Abs("./")
		if err != nil {
			t.Error(err)
		}
		cmd := exec.Command(cpath + "/test.sh")
		buf := new(bytes.Buffer)
		cmd.Stdout = buf
		if err = cmd.Run(); err != nil {
			t.Error(err)
		}
		fmt.Println(buf.String())
		ok := strings.HasSuffix(buf.String(), "PASS\n")
		require.True(t, ok, "sh tests not PASS")
	})
}

func tempFile() (*os.File, func() error, error) {
	tmpfile, err := ioutil.TempFile("", "temp.*.data")
	if err != nil {
		return nil, nil, err
	}
	closeF := func() error {
		err = tmpfile.Close()
		if err != nil {
			return err
		}
		err = os.Remove(tmpfile.Name())
		if err != nil {
			return err
		}
		return nil
	}

	return tmpfile, closeF, nil
}
