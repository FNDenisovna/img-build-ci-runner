package git_getter

type InfoYaml struct {
	IsVersioned     bool     `yaml:"is_versioned"`
	VersionTemplate string   `yaml:"version_template"`
	SourcePackages  []string `yaml:"source_packages"`
}
