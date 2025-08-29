package config

// AltSiteUrl - url for checking packages version
// GiteaUrl - url images configs repo
type AppSettings struct {
	Branches           string `json:"branches"`
	VersCronGroupExp   string `json:"vers_cronexp"`
	PeriodCronGroupExp string `json:"period_cronexp"`
	WfToken            string `json:"wf_token"`
	AltSiteUrl         string `json:"alt_site_url"`
	WfUrl              string `json:"wf_url"`
	WfOrgRepo          string `json:"wf_org_repo"`
	WfRefRepo          string `json:"wf_ref_repo"`
	ImagesTemplRepoUrl string `json:"images_templ_repo_url"`
	VersCheckImgGroup  string `json:"vers_check_img_group"`
	PeriodCronImgGroup string `json:"period_cron_img_group"`
	StoragePath        string `json:"storage_path"`
}
