package model

import "time"

type AppSettings struct {
	ImgPkgFileUrl string   `json:"img_pkg_url"`
	Branches      []string `json:"branches"`
	CronExp       string   `json:"cronexp"`
	AltSiteUrl    string   `json:"alt_site_url"`
	GiteaUrl      string   `json:"gitea_url"`
}

type ImgPkg struct {
	Package string `json:"package"`
	Image   string `json:"image"`
}

// Model of info package from local service db
type SqlPack struct {
	Id      int
	Name    string
	Version string
	Release string
	Epoch   int
	Changed time.Time
	Branch  string
}

// Model of info package from basealt site api
type SiteVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Release string `json:"release"`
	//"2024-09-16T16:22:17"
	Changed time.Time `json:"changed"`
}

type GiteaTag struct {
	Message string `json:"message"`
	TagName string `json:"tag_name"`
	Target  string `json:"target"`
	Image   string
	Version string
	Branch  string
}
