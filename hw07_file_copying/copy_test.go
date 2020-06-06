package main

import "testing"

func TestCopy(t *testing.T) {
	err := Copy("/Volumes/Данные (hdd)/dima/Downloads/Охота.2012.BDRip.(1080p).mkv", "/Volumes/Данные (hdd)/dima/Downloads/test.mkv", 0, 1000<<20)
	if err != nil {
		t.Error(err)
	}
}
