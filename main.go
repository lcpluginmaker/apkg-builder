package main

import (
	"os"
	"path"

	"github.com/alexcoder04/LeoConsole-apkg-builder/pkg"
	"github.com/alexcoder04/arrowprint"
	"github.com/alexcoder04/friendly"
)

func main() {
	if len(os.Args) < 2 {
		arrowprint.Err0("no arguments passed")
		os.Exit(1)
	}

	folder := os.Args[1]

	if folder[0] != '/' && folder[1] != ':' {
		pwd, err := os.Getwd()
		if err != nil {
			arrowprint.Err0("cannot get working directory")
			os.Exit(1)
		}
		folder = path.Join(pwd, folder)
	}
	buildFolder := path.Join(os.TempDir(), "apkg-build")

	manifest := pkg.LoadManifest(folder)
	pkg.Compile(folder, manifest)
	pkg.PreparePackage(folder, buildFolder, manifest)
	pkg.GenPkgInfo(buildFolder, manifest)
	outputFile := path.Join(folder, manifest.PackageName+".lcp")
	friendly.CompressFolder(buildFolder, outputFile)

	arrowprint.Suc0("Done. Package archive saved to %s.", outputFile)
}
