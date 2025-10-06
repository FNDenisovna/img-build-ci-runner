package alt_api

import (
	"encoding/json"
	"fmt"
	"img-build-ci-runner/internal/api"
	model "img-build-ci-runner/internal/model"
	"log"
	"regexp"
	"sort"
	"time"
)

type SitePackInfo struct {
	Versions []model.SiteVersion `json:"versions"`
	Message  string              `json:"message"`
	Length   int                 `json:"length"`
}

type PackInfo struct {
	Versions []model.SiteVersion `json:"versions"`
	Message  string              `json:"message"`
}

type PackListByName struct {
	Packages []model.PackInfoByName `json:"packages"`
	Message  string                 `json:"message"`
	Length   int                    `json:"length"`
}

func GetTaskPackInfo(url, name, branch string) (packinfo model.SiteVersion, err error) {
	endpoint := fmt.Sprint(url, "site/package_versions_from_tasks")
	params := make(map[string]string, 2)
	params["name"] = name
	params["branch"] = branch

	req := api.New(endpoint)
	req.Params = params

	row, statusCode, err := req.Get()
	if statusCode == 429 {
		log.Printf("Can't get response: Too Many Requests. Sleep and try again.")
		time.Sleep(time.Second * 10)
		row, _, err = req.Get()
	}

	if err != nil {
		return
	}

	var resp SitePackInfo

	err = json.Unmarshal(row, &resp)
	if err != nil {
		log.Printf("Alt-api response: %v", string(row))
		err = fmt.Errorf("unmarshal response is failed. Error: %w", err)
		return
	}

	if resp.Message != "" {
		err = fmt.Errorf("no found package %v on branch %v. Error message: %s", name, branch, resp.Message)
		return
	}

	if resp.Length > 0 && resp.Message == "" {
		packinfo = resp.Versions[0]
		return
	}

	err = fmt.Errorf("something wrong. No found package %v on branch %v", name, branch)
	return
}

func GetPacksListByName(url, template, branch string) (packlist []model.PackInfoByName, err error) {
	endpoint := fmt.Sprint(url, "site/find_packages")
	params := make(map[string]string, 3)
	params["name"] = template
	params["arch"] = "x86_64"
	params["branch"] = branch

	req := api.New(endpoint)
	req.Params = params

	row, statusCode, err := req.Get()
	if statusCode == 429 {
		log.Printf("Can't get response: Too Many Requests. Sleep and try again.")
		time.Sleep(time.Second * 10)
		row, _, err = req.Get()
	}
	if err != nil {
		return
	}

	var resp PackListByName
	err = json.Unmarshal(row, &resp)
	if err != nil {
		err = fmt.Errorf("unmarshal response is failed. Error: %w", err)
		return
	}

	if resp.Message != "" {
		err = fmt.Errorf("no found packages by name-template %v on branch %v. Error message: %s", template, branch, resp.Message)
		return
	}

	if len(resp.Packages) > 0 {
		sort.Slice(resp.Packages, func(i, j int) bool { return resp.Packages[i].Name > resp.Packages[j].Name })
		log.Printf("Sort result packages by template: %+v", resp.Packages)

		regexpr := fmt.Sprintf("^%s.+", template)
		packlist = make([]model.PackInfoByName, 0, 3)

		for _, pack := range resp.Packages {
			matched, regerr := regexp.MatchString(regexpr, pack.Name)
			if regerr != nil {
				log.Printf("Can't regexp-parse package name %s by expression %s\n", pack.Name, regexpr)
			}
			if matched && !pack.ByBinary && !pack.Versions[0].Deleted {
				packlist = append(packlist, pack)
			}
			if len(packlist) >= 3 {
				log.Printf("Resulting packages list finding by name-template %s: %+v\n", template, packlist)
				return
			}
		}
	}

	err = fmt.Errorf("no found packages by name-template %v on branch %v", template, branch)
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

	row, statusCode, err := req.Get()
	if statusCode == 429 {
		log.Printf("Can't get response: Too Many Requests. Sleep and try again.")
		time.Sleep(time.Second * 10)
		row, _, err = req.Get()
	}
	if err != nil {
		return
	}

	var resp PackInfo
	err = json.Unmarshal(row, &resp)
	if err != nil {
		err = fmt.Errorf("unmarshal response is failed. Error: %w", err)
		return
	}

	if resp.Message != "" {
		err = fmt.Errorf("no found package %v on branch %v. Error message: %s", name, branch, resp.Message)
		return
	}

	if len(resp.Versions) > 0 {
		for _, ver := range resp.Versions {
			if ver.Branch == branch {
				packinfo = ver
				return
			}
		}
	}

	err = fmt.Errorf("no found package %v on branch %v", name, branch)
	return
}
