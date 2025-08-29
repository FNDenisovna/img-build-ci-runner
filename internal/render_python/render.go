package renderpython

import (
	"img-build-ci-runner/internal/resources"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	path       = ""
	scriptName = "render_python_template.py"
	script     = `#!/usr/bin/python3

import sys
from jinja2 import Environment, BaseLoader

#print("Command: ")
#print(sys.argv)
branch = sys.argv[2]
version = sys.argv[3]
template = sys.argv[1]
#Example template:
#'{%- if branch in [ "p10", "c10f1", "c10f2"] -%} php8.2 {%- else -%} php8.3 {%- endif -%}'
 
print(Environment(loader=BaseLoader()).from_string(template).render(branch=branch,version=version))`
)

func CreateScriptFile(_path string) (err error) {
	path = resources.ManageResources(_path, scriptName)
	log.Printf("Render-python script path: %s\n", path)

	scriptFile, err := os.Create(path)
	if err != nil {
		log.Fatalf("Can't create render-python script by path %s. Error: %s", path, err)
	}

	defer scriptFile.Close()
	_, err = scriptFile.WriteString(script)
	if err != nil {
		log.Fatalf("Can't write render-python script to file by path %s. Error: %s", path, err)
	}

	err = scriptFile.Chmod(0744)
	if err != nil {
		log.Fatalf("Can't set chmod 0744 to render-python script by path %s. Error: %s", path, err)
	}
	return
}

func CheckTemplate(packageName string) bool {
	return strings.ContainsAny(packageName, "{%}")
}

// Try render python template with package name
// If it fails, return template without changes
func RenderPackageName(packageTempl, branch, version string) string {
	cmd := exec.Command("/usr/bin/python3", path, packageTempl, branch, version)
	stdout, err := cmd.Output()

	if err != nil {
		log.Printf("Rendering package name by python script %s is failded. Branch: %s. Version: %s. \nTemplate: %s\n", path, branch, version, packageTempl)
		return packageTempl
	}

	res := string(stdout)
	res = strings.TrimSpace(res)
	log.Printf("Package name by python script is rendered: %s", res)
	return res
}
