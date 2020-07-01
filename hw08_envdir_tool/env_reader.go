package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Environment map[string]string

var (
	ErrCantReadFile = errors.New("can't read file")
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	elems, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := make([]os.FileInfo, 0)
	for _, elem := range elems {
		if !elem.IsDir() {
			files = append(files, elem)
		}
	}

	env := make(Environment, len(files))
	for _, file := range files {
		line, err := readFirstLineFromFile(path.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		line = strings.TrimRight(line, "\t ")
		line = strings.Replace(line, string(rune(0x00)), "\n", -1)
		env[file.Name()] = line
	}

	return env, nil
}

func readFirstLineFromFile(path string) (line string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}
	if stat.Size() == 0 {
		return "", nil
	}
	snl := bufio.NewScanner(file)
	if ok := snl.Scan(); !ok || snl.Err() != nil {
		return "", ErrCantReadFile
	}

	return snl.Text(), nil
}
