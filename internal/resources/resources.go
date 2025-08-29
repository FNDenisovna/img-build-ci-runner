package resources

import (
	"log"
	"os"
	"path/filepath"
)

var (
	appPath   = "img-build-ci-runner"
	localPath = "/.local/share"
)

// Create resources directories and files
func ManageResources(path, fileName string) string {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Can't get home directory of current host user to create resources. Error: %v\n", err)
		}

		path = filepath.Join(home, localPath, appPath)
	}

	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Can't parse resources directory. Path %s. Error: %s", path, err)
	}
	log.Printf("Parsed resources path: %s\n", path)

	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0764)
		if err != nil && !os.IsExist(err) {
			log.Fatalf("Can't create render-python script directory. Path %s. Error: %s", path, err)
		}
		log.Printf("Resources path: %s\n", path)
	}

	path = filepath.Join(path, fileName)

	return path
}
