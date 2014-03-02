package porcelain

import (
	"io"
	"os"
	"path"

	"github.com/kourge/ggit/config"
)

const (
	defaultPerm os.FileMode = 0755
	defaultDesc string      = `Unnamed repository; edit this file 'description' to name the repository.`
)

var (
	repoDirs []string = []string{
		"branches",
		"hooks",
		"info",
		"objects", "objects/info", "objects/pack",
		"refs", "refs/heads", "refs/tags",
	}

	defaultConfig config.Config = config.Config{
		"core": {"core", config.Dict{
			"repositoryformatversion": 0,
			"filemode":                true,
			"ignorecase":              true,
			"precomposeunicode":       false,
		}},
	}
)

// InitOptions contains all the possible options for InitRepo.
//
// Dir is a string that is a path to the directory under which a repository
// should be initialized. If left blank, it defaults to the current working
// directory.
//
// Bare is a bool that indicates if this repository should be a bare repository,
// which is a repo that has no working tree but acts as an object storage.
type InitOptions struct {
	Dir  string
	Bare bool
}

// InitRepo initializes a repo, given o as its options. Equivalent to
// `git init`. See the documentation on InitOptions for more details.
func InitRepo(o InitOptions) error {
	if o.Dir == "" {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		o.Dir = dir
	}
	dir := o.Dir
	if o.Bare {
		defaultConfig["core"].Dict["bare"] = true
	} else {
		dir = path.Join(o.Dir, ".git")
		defaultConfig["core"].Dict["logallrefupdates"] = true
		defaultConfig["core"].Dict["bare"] = false
	}

	if err := os.MkdirAll(dir, defaultPerm); err != nil {
		return err
	} else if err := os.Chdir(dir); err != nil {
		return err
	}

	for _, dir := range repoDirs {
		if err := os.Mkdir(dir, defaultPerm); err != nil {
			return err
		}
	}

	for _, item := range []fileCreation{
		{"HEAD", func(f *os.File) error {
			_, err := f.Write([]byte("ref: refs/heads/master\n"))
			return err
		}},
		{"config", func(f *os.File) error {
			_, err := io.Copy(f, defaultConfig.Reader())
			return err
		}},
		{"description", func(f *os.File) error {
			_, err := f.Write([]byte(defaultDesc))
			return err
		}},
		{"info/exclude", fileNop},
	} {
		if file, err := os.Create(item.Name); err != nil {
			return err
		} else {
			defer file.Close()
			if err := item.Do(file); err != nil {
				return err
			}
		}
	}

	return nil
}
