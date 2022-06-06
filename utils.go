package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/alexcoder04/arrowprint"
)

func CopyFile(source string, destin string) {
	bytesRead, err := ioutil.ReadFile(source)
	if err != nil {
		arrowprint.Err0(err.Error())
		os.Exit(1)
	}

	err = ioutil.WriteFile(destin, bytesRead, 0600)
	if err != nil {
		arrowprint.Err0(err.Error())
		os.Exit(1)
	}
}

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

func Compress(folder string, destin string) {
	arrowprint.Suc1("6. Compressing %s...", folder)
	file, err := os.Create(destin)
	if err != nil {
		arrowprint.Err0("cannot create output file")
		os.Exit(1)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		arrowprint.Info1("adding: %#v", path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		f, err := w.Create(path[len(folder)+1:])
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	err = filepath.Walk(folder, walker)
	if err != nil {
		arrowprint.Err0("cannot compress folder")
		os.Exit(1)
	}
}
