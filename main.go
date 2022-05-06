package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	cp "github.com/otiai10/copy"
)

func LoadManifest(folder string) Manifest {
	manifestFile := path.Join(folder, "manifest.apkg.json")
	_, err := os.Stat(manifestFile)
	if err != nil {
		log.Fatalln("manifest file not found")
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
	return manifest
}

func Compile(folder string, manifest Manifest) {
	log.Printf(
		"Building %s from %s\n",
		manifest.PackageName,
		manifest.Project.Maintainer)
	log.Println("1. Running build script...")
	cmd := exec.Command(manifest.Build.Command, manifest.Build.Args...)
	cmd.Dir = folder

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	err := cmd.Run()
	if err != nil {
		log.Fatalln("build script failed")
	}
	log.Println(stdBuffer.String())
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

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("no arguments passed")
	}

	folder := os.Args[1]
	if folder[0] != '/' && folder[1] != ':' {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatalln("cannot get working directory")
		}
		folder = path.Join(pwd, folder)
	}
	buildFolder := path.Join(os.TempDir(), "apkg-build")

	manifest := LoadManifest(folder)
	Compile(folder, manifest)
	PreparePackage(folder, buildFolder, manifest)
	GenPkgInfo(buildFolder, manifest)
	outputFile := path.Join(folder, manifest.PackageName+".lcpkg")
	Compress(buildFolder, outputFile)

	log.Printf("Done. Package archive saved to %s.\n", outputFile)
}
