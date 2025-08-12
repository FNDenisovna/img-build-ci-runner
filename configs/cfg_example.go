package configs

const Cfg_example = `{
	"branches": "p10 p11 sisyphus",
	"vers_cronexp": "@daily",
	"period_cronexp": "@monthly",
	"alt_site_url": "https://rdb.altlinux.org/api/",
	"wf_url": "https://gitea.basealt.ru/",
	"wf_org_repo": "fedorovand/image-forge",
	"wf_ref_repo": "master",
	"wf_token": "",
	"wf_images_name": "workflow_multiple.yaml",
	"wf_group_name": "wf_full.yaml",
	"images_templ_repo_url": "https://gitea.basealt.ru/alt/image-forge",
	"storage_path": "~/.local/share/img-build-ci-runner",
	"vers_check_img_group": "alt k8s",
	"period_cron_img_group": "base kubevirt"
}`
