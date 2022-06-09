package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/alexcoder04/arrowprint"
	cp "github.com/otiai10/copy"
)

func LoadManifest(folder string) Manifest {
	manifestFile := path.Join(folder, "manifest.json")
	_, err := os.Stat(manifestFile)
	if err != nil {
		arrowprint.Err0("manifest file not found")
		os.Exit(1)
	}
	content, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		arrowprint.Err0("cannot read manifest file")
		os.Exit(1)
	}
	manifest := Manifest{}
	err = json.Unmarshal(content, &manifest)
	if err != nil {
		arrowprint.Err0("cannot unmarshal manifest file")
		os.Exit(1)
	}
	return manifest
}

func Compile(folder string, manifest Manifest) {
	arrowprint.InfoC(
		"Building %s from %s",
		manifest.PackageName,
		manifest.Project.Maintainer)
	arrowprint.Suc0("1. Running build script...")
	cmd := exec.Command(manifest.Build.Command, manifest.Build.Args...)
	cmd.Dir = folder

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	err := cmd.Run()
	if err != nil {
		arrowprint.Err0("build script failed")
		os.Exit(1)
	}
	fmt.Println(stdBuffer.String())
}

func PreparePackage(folder string, buildFolder string, manifest Manifest) {
	arrowprint.Suc0("2. Creating build folder...")
	_, err := os.Stat(buildFolder)
	if err != nil {
		if !os.IsNotExist(err) {
			arrowprint.Err0("error while stat build folder")
			os.Exit(1)
		}
	} else {
		os.RemoveAll(buildFolder)
	}
	arrowprint.Suc0("3. Populating build folder with dlls...")
	err = os.MkdirAll(path.Join(buildFolder, "plugins"), 0700)
	if err != nil {
		arrowprint.Err0("error creating build folder")
		os.Exit(1)
	}
	for i, d := range manifest.Build.Dlls {
		arrowprint.Info1("3.%d. %s...", i+1, d)
		CopyFile(
			path.Join(folder, "bin", "Debug", "net6.0", d),
			path.Join(buildFolder, "plugins", d))
	}
	arrowprint.Suc0("4. Copying share files to build folder...")
	err = cp.Copy(
		path.Join(folder, manifest.Build.Share),
		path.Join(buildFolder, "share"))
	if err != nil {
		arrowprint.Err0("cannot copy shared files")
		os.Exit(1)
	}
}

func GetPackageOS() string {
	if os.Getenv("APKG_BUILDER_OS") == "" {
		return "lnx64"
	}
	return os.Getenv("APKG_BUILDER_OS")
}

func GenPkgInfo(buildFolder string, manifest Manifest) {
	arrowprint.Suc0("5. Generating PKGINFO...")
	pkginfo := PKGINFO{}
	pkginfo.ManifestVersion = 1.1
	pkginfo.PackageName = manifest.PackageName
	pkginfo.PackageVersion = manifest.PackageVersion
	pkginfo.PackageOS = GetPackageOS()
	pkginfo.Files = GetFilesList(buildFolder)
	pkginfo.Project = manifest.Project

	res, err := json.Marshal(pkginfo)
	if err != nil {
		arrowprint.Err0("cannot marschal manifest")
		os.Exit(1)
	}
	err = ioutil.WriteFile(path.Join(buildFolder, "PKGINFO.json"), res, 0600)
	if err != nil {
		arrowprint.Err0("cannot write pkginfo")
		os.Exit(1)
	}
}

func main() {
	var folder string

	switch len(os.Args) {
	case 0:
		arrowprint.Err0("no arguments passed")
		os.Exit(1)
		break
	case 1:
		folder = os.Args[1]
		break
	case 2:
		// Args[1] is IData
		folder = os.Args[2]
	}

	if folder[0] != '/' && folder[1] != ':' {
		pwd, err := os.Getwd()
		if err != nil {
			arrowprint.Err0("cannot get working directory")
			os.Exit(1)
		}
		folder = path.Join(pwd, folder)
	}
	buildFolder := path.Join(os.TempDir(), "apkg-build")

	manifest := LoadManifest(folder)
	Compile(folder, manifest)
	PreparePackage(folder, buildFolder, manifest)
	GenPkgInfo(buildFolder, manifest)
	outputFile := path.Join(folder, manifest.PackageName+".lcp")
	Compress(buildFolder, outputFile)

	arrowprint.Suc0("Done. Package archive saved to %s.", outputFile)
}
