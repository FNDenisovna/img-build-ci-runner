package model

import (
	"time"
)

type ImgPkg struct {
	// key Image string `json:"image"`
	// value Packages list string `json:"package"`
	MapImgPkg map[string][]string
}

// Model of info package from local service db
type SqlPack struct {
	Changed time.Time
	Name    string
	Version string
	Release string
	Branch  string
	Id      int
	Epoch   int
}

// Model of info package from basealt site api
type SiteVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Release string `json:"release"`
	Changed string `json:"changed"`
	Branch  string `json:"branch"`
}

type PackInfoByName struct {
	Versions []PackVersionByName `json:"versions"`
	Name     string              `json:"name"`
	ByBinary bool                `json:"by_binary"`
}

type PackVersionByName struct {
	Version string `json:"version"`
	Release string `json:"release"`
	Deleted bool   `json:"deleted"`
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
	Images    []WfInputsImagesInfo `json:"-"`
	Branch    string               `json:"branch"`
	ImagesStr string               `json:"images"`
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
