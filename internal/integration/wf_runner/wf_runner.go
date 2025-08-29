package wf_runner

import (
	"encoding/json"
	"fmt"
	"img-build-ci-runner/internal/api"
	model "img-build-ci-runner/internal/model"
	"log"
)

const (
	urlCfgKey          = "wf_url"
	orgRepoCfgKey      = "wf_org_repo"
	wfImagesNameCfgKey = "wf_images_name"
	wfGroupNameCfgKey  = "wf_group_name"
	tokenCfgKey        = "wf_token"
)

type WfInputData interface {
	model.WfInputDataGroup | model.WfInputDataImages
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
/// V1 - nhit method create tag for running ci
func RunBuildImageV1(tag *model.WfTag, url, orgRepo, token string) error {
	endpoint := fmt.Sprintf("%sapi/v1/repos/%s/tags", url, orgRepo)
	headers := make(map[string]string, 3)
	headers["accept"] = "application/json"
	headers["authorization"] = fmt.Sprintf("Basic %s", token)
	headers["Content-Type"] = "application/json"

	tag.TagName = fmt.Sprintf("%s_%s_%s", tag.Branch, tag.Image, tag.Version)
	req := api.New(endpoint)
	req.Headers = headers

	body, err := json.Marshal(tag)
	if err != nil {
		return fmt.Errorf("Can't marsal struct %v. Error: %v\n", tag, err)
	}
	row, _, err := req.Post(body)
	if err != nil {
		return err
	}

	log.Printf("Tag to run workflow is pushed. Response: %x\n", row)
	return nil
}

/*
 * curl -X 'POST' \
    'https://gitea.basealt.ru/api/v1/repos/fedorovand/image-forge/actions/workflows/workflow_multiple.yaml/dispatches' \
    -H 'accept: application/json' \
    -H 'authorization: Basic <TOKEN>' \
    -H 'Content-Type: application/json' \
        -d '{
                "inputs": {
                    "images": "[
                    	{\"name\":\"alt/etcd\",\"version\":\"1.27.3\"},
                     	{\"name\":\"alt/golang\",\"version\":\"1.24.0\"},
                      	{\"name\":\"alt/python\",\"version\":\"3.14.7\"}
                    ]",
                    "branch": "p11"
                },
                "ref": "refs/heads/master"
            }'
*/
/// ref - git target branch which has running workflow
/// url - contains workflow ID (workflow_multiple.yaml), organization/repository (fedorovand/image-forge)
/// inputs - is getting from difinition of workflow
/// V1 - nhit method create tag for running ci
func RunBuildImage[T WfInputData](inputData *T, url, wfName, orgRepo, token string) error {
	endpoint := fmt.Sprintf("%sapi/v1/repos/%s/actions/workflows/%s/dispatches",
		url, orgRepo, wfName)
	log.Printf("Endpoint to run workflow: %v", endpoint)
	headers := make(map[string]string, 3)
	headers["accept"] = "application/json"
	headers["authorization"] = fmt.Sprintf("Basic %s", token)
	headers["Content-Type"] = "application/json"

	req := api.New(endpoint)
	req.Headers = headers

	body, err := json.Marshal(inputData)
	if err != nil {
		log.Printf("Can't marsal struct *model.WfInputData %v. Error: %w\n", inputData, err)
		return err
	}

	log.Printf("Request body: %v\n", string(body))

	row, statusCode, err := req.Post(body)
	if err != nil {
		log.Printf("Can't run workflow. Status code: %v, Error: %v, Response: %v\n", statusCode, err, string(row))
		return err
	}

	if statusCode >= 300 {
		log.Printf("Can't run workflow. Status code: %v, Error: %v, Response: %v\n", statusCode, err, string(row))
		return err
	}

	log.Printf("Building workflow is running. Status code: %v. Response: %v\n", statusCode, string(row))
	return nil
}
