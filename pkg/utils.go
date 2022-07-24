package pkg

import (
	"os"
	"path/filepath"

	"github.com/alexcoder04/arrowprint"
)

func GetFilesList(folder string) []string {
	filesList := []string{}
	err := filepath.Walk(folder, func(p string, f os.FileInfo, err error) error {
		if p == folder {
			return nil
		}
		stat, err := os.Stat(p)
		if err != nil {
			arrowprint.Err0("cannot stat file")
			os.Exit(1)
		}
		if !stat.IsDir() {
			filesList = append(filesList, p[len(folder)+1:])
		}
		return nil
	})
	if err != nil {
		arrowprint.Err0("cannot list files")
		os.Exit(1)
	}
	return filesList
}

func GetPackageOS() string {
	if os.Getenv("APKG_BUILDER_OS") == "" {
		return "lnx64"
	}
	return os.Getenv("APKG_BUILDER_OS")
}
