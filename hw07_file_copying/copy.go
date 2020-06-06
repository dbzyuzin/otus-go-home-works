package main

import (
	"errors"
	"io"
	"os"

	"github.com/mitchellh/ioprogress"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) (err error) {
	err = checkPaths(fromPath, toPath)
	if err != nil {
		return err
	}

	in, out, closeF, err := openFiles(fromPath, toPath)
	if err != nil {
		return err
	}
	defer func() {
		cerr := closeF()
		if cerr != nil {
			err = cerr
		}
	}()

	inr, err := getReaderWithOffset(in, offset)
	if err != nil {
		return err
	}

	if limit > 0 {
		_, err = io.CopyN(out, inr, limit)
	} else {
		_, err = io.Copy(out, inr)
	}
	if err != nil && err != io.EOF {
		return err
	}

	err = out.Sync()
	if err != nil {
		return err
	}
	return nil
}

func getReaderWithOffset(file *os.File, offset int64) (io.Reader, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if offset > stat.Size() {
		return nil, ErrOffsetExceedsFileSize
	}
	if offset > 0 {
		_, err = file.Seek(offset, io.SeekStart)
		if err != nil {
			return nil, err
		}
	}

	inr := &ioprogress.Reader{
		Reader: file,
		Size:   stat.Size() - offset,
	}

	return inr, nil
}

func openFiles(src, dst string) (in, out *os.File, closer func() error, err error) {
	in, err = os.Open(src)
	if err != nil {
		return
	}
	out, err = os.Create(dst)
	if err != nil {
		return
	}

	return in, out, func() error {
		cerr := in.Close()
		if cerr != nil {
			return cerr
		}
		cerr = out.Close()
		if cerr != nil {
			return cerr
		}
		return nil
	}, nil
}

func checkPaths(src, dst string) (err error) {
	sstat, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sstat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	dstat, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !dstat.Mode().IsRegular() {
			return ErrUnsupportedFile
		}
	}

	return nil
}
