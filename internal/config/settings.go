package config

// AltSiteUrl - url for checking packages version
// GiteaUrl - url images configs repo
type AppSettings struct {
	Branches          []string `json:"branches"`
	VersCronGroupExp  string   `json:"vers_cronexp"`
	DepsCronGroupExp  string   `json:"deps_cronexp"`
	GiteaToken        string   `json:"gitea_token"`
	AltSiteUrl        string   `json:"alt_site_url"`
	GiteaWfUrl        string   `json:"gitea_wf_url"`
	GiteaRepoUrl      string   `json:"gitea_repo_url"`
	VersCheckImgGroup string   `json:"vers_check_img_group"`
	DepsCronImgGroup  string   `json:"deps_cron_omg_group"`
}
