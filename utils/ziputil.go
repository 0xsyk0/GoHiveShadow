package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
)

type InternalFile struct {
	Name string
	Path string
}

func ZipFiles(filename string, files []InternalFile) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, f InternalFile) error {

	fileToZip, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = f.Name
	header.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func CreateZipArchive(outputPath string) bool {
	var samPath = path.Join(outputPath, "SAM")
	var symPath = path.Join(outputPath, "SYSTEM")
	var secPath = path.Join(outputPath, "SECURITY")
	var archPath = path.Join(outputPath, "samael.zip")
	fmt.Println("[+] Creating archive")
	var files = []InternalFile{
		InternalFile{Name: "SAM", Path: samPath},
		InternalFile{Name: "SYSTEM", Path: symPath},
		InternalFile{Name: "SECURITY", Path: secPath},
	}

	if err := ZipFiles(archPath, files); err != nil {
		panic(err)
	}
	fmt.Println("[+] Cleanup")
	os.Remove(samPath)
	os.Remove(symPath)
	os.Remove(secPath)
	fmt.Println(fmt.Sprintf("[+] Archive available at path %s", archPath))
	return true
}
