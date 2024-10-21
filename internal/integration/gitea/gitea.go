package gitea

import (
	"altpack-vers-checker/internal/api"
	model "altpack-vers-checker/internal/model"
	"encoding/json"
	"fmt"
	"log"
)

type GiteaApi struct {
	url string
}

func New(url string) *GiteaApi {
	a := &GiteaApi{}
	if url != "" {
		a.url = url
	} else {
		a.url = "https://gitea.basealt.ru/"
	}
	return a
}

/*
	curl -X 'POST' \
	  'https://gitea.basealt.ru/api/v1/repos/alt/image-forge/tags' \
	  -H 'accept: application/json' \
	  -H 'authorization: Basic ...' \
	  -H 'Content-Type: application/json' \
	  -d '{
	  "message": "building",
	  "tag_name": "p11_alt",
	  "target": "master"
	}'
*/
/// target - git branch for creating tag
func (g *GiteaApi) CreateTag(tag *model.GiteaTag, token string) error {
	endpoint := fmt.Sprint(g.url, "api/v1/repos/alt/image-forge/tags")
	headers := make(map[string]string, 3)
	headers["accept"] = "application/json"
	headers["authorization"] = fmt.Sprintf("Basic %s", token)
	headers["Content-Type"] = "application/json"

	tag.TagName = fmt.Sprintf("%s_%s_%s", tag.Branch, tag.Image, tag.Version)
	req := api.New(endpoint)
	req.Params = headers

	body, err := json.Marshal(tag)
	if err != nil {
		return fmt.Errorf("Can't marsal struct %v. Error: %v", tag, err)
	}
	row, err := req.Post(body)
	if err != nil {
		return err
	}

	log.Printf("Tag in gitea is created. Response: %x", row)
	return nil
}
