package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

/*
Support passing a flag to move files instead of copying, useful if storage space is limited.
Passing -x flag will extract files from their original location to the FOUND directory.
*/
var (
	xFlag bool
)

func init() {
	flag.BoolVar(&xFlag, "x", false, "Extract files to output folder (original file will be deleted) Use this option if you do not have room to copy all files found.")
	flag.Parse()
}

func moveFile(source, dest string) error {
	err := os.Rename(source, dest)
	if err != nil {
		return err
	}
	return nil
}

func copyFile(source, dest string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func seekFiles(sourceDir, destDir string, filetypes []string) error {
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil //continue walking
		}

		for _, ext := range filetypes {
			if strings.EqualFold(filepath.Ext(path), ext) {
				destFile := filepath.Join(destDir, filepath.Base(path))
				if xFlag {
					err = moveFile(path, destFile)
				} else {
					err = copyFile(path, destFile)
				}
				if err != nil {
					return err
				}
				if xFlag {
					fmt.Printf("Extracted: %s\n", path)
				} else {
					fmt.Printf("Copied: %s\n", path)
				}
			}
		}

		return nil
	})

	return err
}

func main() {
	var fileTypes []string
	fileTypes = append(fileTypes, os.Args[1:]...)

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error detecting working directory: %s\n", err)
		return
	}

	foundDir := filepath.Join(currentDir, "FOUND")
	err = os.MkdirAll(foundDir, 0755)
	if err != nil {
		fmt.Printf("Error creating FOUND directory: %s\n", err)
		return
	}

	err = seekFiles(currentDir, foundDir, fileTypes)
	if err != nil {
		fmt.Printf("Error copying or extracting file: %s\n", err)
	}

	fmt.Println("All files have been copied or moved to the FOUND directory.")

	fmt.Println("Press enter key to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
