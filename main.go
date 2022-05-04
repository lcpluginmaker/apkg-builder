package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	cp "github.com/otiai10/copy"
)

type ManifestBuild struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Folder  string   `json:"folder"`
	Dlls    []string `json:"dlls"`
	Share   string   `json:"share"`
}

type ManifestProject struct {
	Maintainer string `json:"maintainer"`
	Email      string `json:"email"`
	Homepage   string `json:"homepage"`
	BugTracker string `json:"bugTracker"`
}

type Manifest struct {
	ManifestVersion float64         `json:"manifestVersion"`
	PackageName     string          `json:"packageName"`
	PackageVersion  string          `json:"packageVersion"`
	Build           ManifestBuild   `json:"build"`
	Project         ManifestProject `json:"project"`
}

type PKGINFO struct {
	PackageName    string          `json:"packageName"`
	PackageVersion string          `json:"packageVersion"`
	Files          []string        `json:"files"`
	Project        ManifestProject `json:"project"`
}

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

func Compile(folder string) Manifest {
	manifestFile := path.Join(folder, "manifest.apkg.json")
	_, err := os.Stat(manifestFile)
	if err != nil {
		log.Fatalln("manifest does not exist")
	}
	content, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		log.Fatalln("cannot read manifest file")
	}
	manifest := Manifest{}
	err = json.Unmarshal(content, &manifest)
	if err != nil {
		log.Fatalln("cannot unmarshal manifest file")
	}
	log.Printf("Building %s from %s\n", manifest.PackageName, manifest.Project.Maintainer)
	log.Println("1. Running build script...")
	cmd := exec.Command(manifest.Build.Command, manifest.Build.Args...)
	cmd.Dir = folder

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	err = cmd.Run()
	if err != nil {
		log.Fatalln("build script failed")
	}
	log.Println(stdBuffer.String())

	return manifest
}

func PreparePackage(folder string, buildFolder string, manifest Manifest) {
	log.Println("2. Creating build folder...")
	_, err := os.Stat(buildFolder)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalln("error while stat build folder")
		}
	} else {
		os.RemoveAll(buildFolder)
	}
	log.Println("3. Populating build folder with dlls...")
	err = os.MkdirAll(path.Join(buildFolder, "plugins"), 0700)
	if err != nil {
		log.Fatalln("error creating build folder")
	}
	for i, d := range manifest.Build.Dlls {
		log.Printf("3.%d. %s...\n", i+1, d)
		CopyFile(
			path.Join(folder, "bin", "Debug", "net6.0", d),
			path.Join(buildFolder, "plugins", d))
	}
	log.Println("4. Copying share files to build folder...")
	err = cp.Copy(
		path.Join(folder, manifest.Build.Share),
		path.Join(buildFolder, "share"))
	if err != nil {
		log.Fatalln("cannot copy shared files")
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

func GenPkgInfo(buildFolder string, manifest Manifest) {
	log.Println("5. Generating PKGINFO...")
	pkginfo := PKGINFO{}
	pkginfo.PackageName = manifest.PackageName
	pkginfo.PackageVersion = manifest.PackageVersion
	pkginfo.Files = GetFilesList(buildFolder)
	pkginfo.Project = manifest.Project

	res, err := json.Marshal(pkginfo)
	if err != nil {
		log.Fatalln("cannot marschal manifest")
	}
	err = ioutil.WriteFile(path.Join(buildFolder, "PKGINFO.json"), res, 0600)
	if err != nil {
		log.Fatalln("cannot write pkginfo")
	}
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

		// Ensure that `path` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute paths, but ensure your real-world code
		// transforms path into a zip-root relative path.
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

func main() {
	log.Println("starting")
	if len(os.Args) <= 1 {
		log.Fatalln("no arguments passed")
	}

	folder := os.Args[1]

	manifest := Compile(folder)
	buildFolder := path.Join(os.TempDir(), "apkg-build")
	PreparePackage(folder, buildFolder, manifest)
	GenPkgInfo(buildFolder, manifest)
	outputFile := path.Join(folder, manifest.PackageName+".lcpkg")
	Compress(buildFolder, outputFile)

	log.Printf("Done. Package archive saved to %s.\n", outputFile)
}
