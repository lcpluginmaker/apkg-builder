package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
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
	err = copyLicense(folder, buildFolder, manifest.PackageName)
	if err != nil {
		arrowprint.Err0("error installing license")
		os.Exit(1)
	}
}

func copyLicense(folder string, buildFolder string, pkgname string) error {
	for _, f := range []string{"LICENSE", "LICENSE.txt", "LICENSE.md"} {
		stat, err := os.Stat(path.Join(folder, f))
		if err != nil || stat.IsDir() {
			continue
		}
		arrowprint.Suc0("5. Copying license to build folder...")
		docsFolder := path.Join(buildFolder, "share", "docs", pkgname)
		stat, err = os.Stat(docsFolder)
		if err != nil {
			if os.IsNotExist(err) {
				err := os.MkdirAll(docsFolder, 0700)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		if !stat.IsDir() {
			return errors.New("docs folder exists, but is not a folder")
		}
		CopyFile(path.Join(folder, f), path.Join(docsFolder, "LICENSE"))
		break
	}
	return nil
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
	pkginfo.CompatibleVersions = manifest.CompatibleVersions
	pkginfo.Depends = manifest.Depends
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
