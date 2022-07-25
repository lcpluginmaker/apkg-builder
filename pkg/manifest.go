package pkg

type ManifestBuild struct {
	Folders   []string            `json:"create"`
	Downloads []map[string]string `json:"downloads"`
	Command   string              `json:"command"`
	Args      []string            `json:"args"`
	Folder    string              `json:"folder"`
	Dlls      []string            `json:"dlls"`
	Share     string              `json:"share"`
}

type ManifestProject struct {
	Maintainer string `json:"maintainer"`
	Email      string `json:"email"`
	Homepage   string `json:"homepage"`
	BugTracker string `json:"bugTracker"`
}

type Manifest struct {
	ManifestVersion    float64         `json:"manifestVersion"`
	PackageName        string          `json:"packageName"`
	PackageVersion     string          `json:"packageVersion"`
	Depends            []string        `json:"depends"`
	CompatibleVersions []string        `json:"compatibleVersions"`
	Build              ManifestBuild   `json:"build"`
	Project            ManifestProject `json:"project"`
}

type PKGINFO struct {
	ManifestVersion    float64         `json:"manifestVersion"`
	PackageName        string          `json:"packageName"`
	PackageVersion     string          `json:"packageVersion"`
	PackageOS          string          `json:"packageOS"`
	CompatibleVersions []string        `json:"compatibleVersions"`
	Depends            []string        `json:"depends"`
	Files              []string        `json:"files"`
	Project            ManifestProject `json:"project"`
}
