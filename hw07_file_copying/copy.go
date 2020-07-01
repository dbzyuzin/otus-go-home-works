package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) (err error) {
	err = Validate(fromPath, toPath, offset)
	if err != nil {
		return err
	}

	in, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	out, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		inerr := in.Close()
		outerr := out.Close()
		if inerr != nil {
			err = inerr
		}
		if outerr != nil {
			err = outerr
		}
	}()

	if offset > 0 {
		_, err = in.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}
	var inr io.Reader = in

	if limit > 0 {
		inr = io.LimitReader(in, limit)
	}

	bar := getPBar(in, offset, limit)
	defer bar.Finish()
	inr = bar.NewProxyReader(inr)

	_, err = io.Copy(out, inr)
	if err != nil {
		return err
	}

	return nil
}

func getPBar(in *os.File, offset, limit int64) *pb.ProgressBar {
	inStats, _ := in.Stat()
	inSize := inStats.Size()
	needCopy := inSize - offset
	if needCopy > limit {
		needCopy = limit
	}

	bar := pb.Simple.Start64(needCopy)

	return bar
}

func Validate(fromPath string, toPath string, offset int64) error {
	stats, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if stats.Size() < offset {
		return ErrOffsetExceedsFileSize
	}
	if !stats.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	stats, err = os.Stat(toPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if !stats.Mode().IsRegular() {
			return ErrUnsupportedFile
		}
	}
	return nil
}
