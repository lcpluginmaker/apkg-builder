package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func CopyFile(source string, destin string) {
	bytesRead, err := ioutil.ReadFile(source)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(destin, bytesRead, 0600)
	if err != nil {
		log.Fatal(err)
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
			log.Fatalln("cannot stat file")
		}
		if !stat.IsDir() {
			filesList = append(filesList, p[len(folder)+1:])
		}
		return nil
	})
	if err != nil {
		log.Fatalln("cannot list files")
	}
	return filesList
}

func Compress(folder string, destin string) {
	log.Printf("6. Compressing %s...\n", folder)
	file, err := os.Create(destin)
	if err != nil {
		log.Fatalln("cannot create output file")
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		log.Printf("adding: %#v\n", path)
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
		log.Fatalln("cannot compress folder")
	}
}
