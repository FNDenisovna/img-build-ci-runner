package git_getter

import (
	"fmt"
	"io"
	"io/fs"
	_ "os"
	"regexp"
	"slices"
	"strings"

	"log"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"gopkg.in/yaml.v2"
)

type GitGetter struct {
	virtfs   billy.Filesystem
	GitUrl   string
	ImgGroup string
}

func New(giturl string, imgGroup string) *GitGetter {
	virtfs := memfs.New()
	return &GitGetter{
		GitUrl:   giturl,
		ImgGroup: imgGroup,
		virtfs:   virtfs,
	}
}

// Return map with image key and value as base image of this image
func (g *GitGetter) GetImgWithDeps() map[string]string {
	orgdirs, err := g.getImgGroups()
	if err != nil {
		return nil
	}
	log.Printf("Got groups from git repo: %v\n", orgdirs)

	imgDeps := make(map[string]string, 0)
	imgGroupList := strings.Split(g.ImgGroup, " ")
	for _, orgd := range *orgdirs {
		if !slices.Contains(imgGroupList, orgd.Name()) {
			continue
		}

		//TODO
		// Read every image Dockerfile.template
		// Pasre base image from <FROM> construction
	}
	return imgDeps
}

// return map with image key and value as package list in this image
func (g *GitGetter) GetImgPkgMap() map[string][]string {
	imgPkgMap := make(map[string][]string)

	orgdirs, err := g.getImgGroups()
	if err != nil {
		return nil
	}
	//log.Printf("Got groups from git repo: %v\n", orgdirs)

	log.Printf("Start reading mapping images-packes from repo %s...\n", g.GitUrl)
	imgGroupList := strings.Split(g.ImgGroup, " ")
	for _, orgd := range *orgdirs {
		if !slices.Contains(imgGroupList, orgd.Name()) {
			continue
		}

		imgs, err := g.virtfs.ReadDir("org/" + orgd.Name())
		if err != nil {
			//log.Fatalf("Can't read dir %s is not exist. Error: %v\n", "org/"+orgd.Name(), err)
			continue
		}

		for _, img := range imgs {
			filepath := fmt.Sprintf("org/%s/%s/info.yaml", orgd.Name(), img.Name())
			if _, err := g.virtfs.Lstat(filepath); err != nil {
				//log.Printf("File %s is not exist. Error: %v", filepath, err)
				continue
			}

			file, err := g.virtfs.Open(filepath)
			if err != nil {
				//log.Printf("Can't open file %s. Error: %v", filepath, err)
				continue
			}

			infodata, err := io.ReadAll(file)
			if err != nil {
				//log.Printf("Can't read info.yaml file. Error: %v\n", err)
				continue
			}

			var iy InfoYaml
			err = yaml.Unmarshal(infodata, &iy)
			if err != nil {
				//log.Printf("Can't read info.yaml file. Error: %v\n", err)
				continue
			}

			if iy.IsVersioned {
				packs := beautify(iy.SourcePackages[:])
				imgPkgMap[img.Name()] = packs
				log.Printf("image %s, packeges: %v\n", img.Name(), packs)
			}
		}
	}
	log.Printf("Finish reading mapping images-packes from repo %s\n", g.GitUrl)
	return imgPkgMap
}

func (g *GitGetter) getImgGroups() (*[]fs.FileInfo, error) {
	//"https://gitea.basealt.ru/alt/image-forge"
	_, err := git.Clone(memory.NewStorage(), g.virtfs, &git.CloneOptions{
		URL: g.GitUrl,
		//Progress:      os.Stdout,
		ReferenceName: plumbing.NewBranchReferenceName("master"),
	})

	if err != nil {
		log.Fatalf("Can't read git repo with images and inside it packages info. Giturl: %s. Error: %v\n", g.GitUrl, err)
		return nil, err
		//panic(err)
	}

	log.Printf("Repo %s is cloned", g.GitUrl)
	var orgdirs []fs.FileInfo
	orgdirs, err = g.virtfs.ReadDir("org")
	if err != nil {
		log.Fatalf("Can't read org dirs. Error: %v\n", err)
		return nil, err
		//panic(err)
	}

	return &orgdirs, nil
}

func beautify(asis []string) []string {
	tobe := make([]string, len(asis))
	templ := `\{%.*?%\}|\{{2}.*?\}{2}`
	for i, item := range asis {
		if ok, _ := regexp.MatchString(templ, item); ok {
			re, _ := regexp.Compile(templ)
			item = re.ReplaceAllString(item, "")
			item = strings.TrimSpace(item)
			item = strings.Split(item, " ")[0]
		}
		if item != "" || !slices.Contains(tobe, item) {
			tobe[i] = item
		}
	}

	return tobe
}
