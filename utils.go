package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func copyFile(src *string, dest *string) error {
	srcFile, err := ioutil.ReadFile(*src)
	if err != nil {
		log.Print(fmt.Sprintf("[ERROR]copyFile -> Unable to copy the file '%s' to '%s'.", *src, *dest), err)
		return err
	}
	return writeFile(dest, &srcFile)
}

func delFile(filename *string) error {
	err := os.Remove(*filename)
	if err != nil {
		log.Print(fmt.Sprintf("[ERROR]delFile -> Unable to delete the file '%s'.", *filename), err)
	}
	return err
}

func writeFile(filename *string, buf *[]byte) error {
	err := ioutil.WriteFile(*filename, *buf, 0644)
	if err != nil {
		log.Print(fmt.Sprintf("[ERROR]writeFile -> Unable to write the file: '%s'.", *filename), err)
		return err
	}
	return nil
}

func zipFile(src *string, dest *string) error {
	zFile, err := os.Create(*dest)
	if err != nil {
		log.Print(fmt.Sprintf("[ERROR]zipFile -> Unable to create the zip archive: '%s'.", *dest), err)
		return err
	}
	defer zFile.Close()

	zWriter := zip.NewWriter(zFile)
	defer zWriter.Close()

	srcFile, err := os.Open(*src)
	if err != nil {
		log.Print(fmt.Sprintf("[ERROR]zipFile -> Unable to open the source file: '%s'.", *src), err)
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		log.Print(fmt.Sprintf("[ERROR]zipFile -> Unable to retrieve information from the source file: '%s'.", *src), err)
		return err
	}

	srcHeader, err := zip.FileInfoHeader(srcInfo)
	if err != nil {
		log.Print("[ERROR]zipFile -> Unable to create the file header info from the source file info.", err)
		return err
	}
	srcHeader.Method = zip.Deflate

	zWriter2, err := zWriter.CreateHeader(srcHeader)
	if err != nil {
		log.Print(fmt.Sprintf("[ERROR]zipFile -> Unable to add the file '%s' into the archive '%s'.", *src, *dest), err)
		return err
	}
	_, err = io.Copy(zWriter2, srcFile)
	return err
}
