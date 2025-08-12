package model

import "time"

type ImgPkg struct {
	// key Image string `json:"image"`
	// value Packages list string `json:"package"`
	MapImgPkg map[string][]string
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
	Changed string `json:"changed"`
	Branch  string `json:"branch"`
}

type WfTag struct {
	Message string `json:"message"`
	TagName string `json:"tag_name"`
	Target  string `json:"target"`
	Image   string
	Version string
	Branch  string
	Org     string
}

type WfInputDataImages struct {
	Inputs WfInputsImages `json:"inputs"`
	Ref    string         `json:"ref"`
}

type WfInputsImages struct {
	Images []WfInputsImagesInfo `json:"images"`
	Branch string               `json:"branch"`
}

type WfInputsImagesInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type WfInputDataGroup struct {
	Inputs WfInputsGroup `json:"inputs"`
	Ref    string        `json:"ref"`
}

type WfInputsGroup struct {
	Group  string `json:"group"`
	Branch string `json:"branch"`
}
