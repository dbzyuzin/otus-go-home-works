package main

import (
	"errors"
	"github.com/mitchellh/ioprogress"
	"io"
	"os"
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

	in, out, err, closeF := openFiles(fromPath, toPath)
	if err != nil {
		return err
	}
	defer func() {
		cerr := closeF()
		if cerr != nil {
			err = cerr
		}
	}()

	stat, err := in.Stat()
	if err != nil {
		return err
	}
	if offset > stat.Size() {
		return ErrOffsetExceedsFileSize
	}
	if offset > 0 {
		offset, err = in.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	inr := &ioprogress.Reader{
		Reader: in,
		Size:   stat.Size() - offset,
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

func openFiles(src, dst string) (in, out *os.File, err error, closer func() error) {
	in, err = os.Open(src)
	if err != nil {
		return
	}
	out, err = os.Create(dst)
	if err != nil {
		return
	}

	return in, out, nil, func() error {
		cerr := in.Close()
		if cerr != nil {
			return cerr
		}
		cerr = out.Close()
		if cerr != nil {
			return cerr
		}
		return nil
	}
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
