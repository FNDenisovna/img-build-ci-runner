package altapi

import (
	"altpack-vers-checker/internal/api"
	"encoding/json"
	"fmt"

	model "altpack-vers-checker/internal/integration/model"
)

type AltApi struct {
	url string
}

func New(url string) *AltApi {
	a := &AltApi{}
	if url != "" {
		a.url = url
	} else {
		a.url = "https://rdb.altlinux.org/api/"
	}
	return a
}

func (a *AltApi) GetSitePackInfo(name, branch string) (packinfo model.SiteVersion, err error) {
	endpoint := fmt.Sprint(a.url, "site/package_versions_from_tasks")
	params := make(map[string]string, 2)
	params["name"] = name
	params["branch"] = branch

	req := api.New(endpoint)
	req.Params = params

	row, err := req.Get()
	if err != nil {
		return
	}

	var resp SitePackInfo
	err = json.Unmarshal(row, &resp)
	if err != nil {
		err = fmt.Errorf("Unmarshal response is failed. Error: %w", err)
		return
	}

	if resp.Message != "" {
		err = fmt.Errorf("No found package %v on branch %v. Error message: %w", name, branch, err)
		return
	}

	if resp.Length > 0 && resp.Message == "" {
		packinfo = resp.Versions[0]
		return
	}

	err = fmt.Errorf("Something wrong. No found package %v on branch %v", name, branch)
	return
}

type SitePackInfo struct {
	Length   int                 `json:"length"`
	Versions []model.SiteVersion `json:"versions"`
	Message  string              `json:"message"`
}

func (a *AltApi) GetPackInfo(name, branch string) (pack Pack, err error) {
	endpoint := fmt.Sprint(a.url, "package/package_info")
	params := make(map[string]string, 2)
	params["name"] = name
	params["branch"] = branch
	//params["arch"] = "x86_64"

	req := api.New(endpoint)
	req.Params = params

	row, err := req.Get()
	if err != nil {
		return
	}

	var resp PackInfo
	err = json.Unmarshal(row, &resp)
	if err != nil {
		err = fmt.Errorf("Unmarshal response is failed. Error: %w", err)
		return
	}

	if resp.Length > 0 {
		pack = resp.Packages[0]
		return
	}

	err = fmt.Errorf("No found package %v on branch %v", name, branch)
	return
}

type PackInfo struct {
	Length   int    `json:"length"`
	Packages []Pack `json:"packages"`
}

type Pack struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Release string `json:"release"`
	Epoch   string `json:"epoch"`
}
