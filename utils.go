package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// copyFile copies the src file to dest.
func copyFile(src string, dest string) error {
	var srcFile, err = ioutil.ReadFile(src)

	if err != nil {
		return err
	}
	return writeFile(dest, &srcFile)
}

// deleteFile removes the given filename.
func deleteFile(filename string) error {
	if err := os.Remove(filename); err != nil {
		return err
	}
	return nil
}

// writeFile writes the given file name using buf content.
func writeFile(filename string, buf *[]byte) error {
	if err := ioutil.WriteFile(filename, *buf, 0644); err != nil {
		return err
	}
	return nil
}

// zipFile compress the given src file and write it to dest.
func zipFile(src string, dest string) error {
	zFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("archive file creation error: %v", err)
	}
	defer zFile.Close()

	zWriter := zip.NewWriter(zFile)
	defer zWriter.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("file open error: %v", err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("file stat error: %v", err)
	}

	srcHeader, err := zip.FileInfoHeader(srcInfo)
	if err != nil {
		return fmt.Errorf("file header info error: %v", err)
	}
	srcHeader.Method = zip.Deflate

	zWriter2, err := zWriter.CreateHeader(srcHeader)
	if err != nil {
		return fmt.Errorf("archive file append error: %v", err)
	}
	_, err = io.Copy(zWriter2, srcFile)
	return err
}
