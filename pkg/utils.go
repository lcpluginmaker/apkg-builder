package pkg

import (
	"os"
)

func GetPackageOS() string {
	if os.Getenv("APKG_BUILDER_OS") == "" {
		return "lnx64"
	}
	return os.Getenv("APKG_BUILDER_OS")
}

func RecreateFolder(folder string) error {
	err := os.RemoveAll(folder)
	if err != nil {
		return err
	}
	err = os.MkdirAll(folder, 0700)
	if err != nil {
		return err
	}
	return nil
}
