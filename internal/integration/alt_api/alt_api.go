package alt_api

import (
	"encoding/json"
	"fmt"
	"img-build-ci-runner/internal/api"
	model "img-build-ci-runner/internal/model"
)

type SitePackInfo struct {
	Length   int                 `json:"length"`
	Versions []model.SiteVersion `json:"versions"`
	Message  string              `json:"message"`
}

type PackInfo struct {
	Versions []model.SiteVersion `json:"versions"`
	Message  string              `json:"message"`
}

func GetSitePackInfo(url, name, branch string) (packinfo model.SiteVersion, err error) {
	endpoint := fmt.Sprint(url, "site/package_versions_from_tasks")
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
		err = fmt.Errorf("Unmarshal response is failed. Error: %w\n", err)
		return
	}

	if resp.Message != "" {
		err = fmt.Errorf("No found package %v on branch %v. Error message: %s\n", name, branch, resp.Message)
		return
	}

	if resp.Length > 0 && resp.Message == "" {
		packinfo = resp.Versions[0]
		return
	}

	err = fmt.Errorf("Something wrong. No found package %v on branch %v\n", name, branch)
	return
}

func GetPackInfo(url, name, branch string) (packinfo model.SiteVersion, err error) {
	endpoint := fmt.Sprint(url, "site/package_versions")
	params := make(map[string]string, 3)
	params["name"] = name
	params["arch"] = "x86_64"
	params["package_type"] = "source"

	req := api.New(endpoint)
	req.Params = params

	row, err := req.Get()
	if err != nil {
		return
	}

	var resp PackInfo
	err = json.Unmarshal(row, &resp)
	if err != nil {
		err = fmt.Errorf("Unmarshal response is failed. Error: %w\n", err)
		return
	}

	if resp.Message != "" {
		err = fmt.Errorf("No found package %v on branch %v. Error message: %s\n", name, branch, resp.Message)
		return
	}

	if resp.Versions != nil && len(resp.Versions) > 0 {
		for _, ver := range resp.Versions {
			if ver.Branch == branch {
				packinfo = ver
				return
			}
		}
	}

	err = fmt.Errorf("No found package %v on branch %v\n", name, branch)
	return
}
