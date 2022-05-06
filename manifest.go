package main

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
