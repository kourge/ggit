package plumbing

import (
	"path/filepath"

	"github.com/kourge/ggit/format"
)

func (repo *Repository) Ignores() (*format.GlobTable, error) {
	path := filepath.Join(repo.path, "..", ".gitignore")
	return format.GlobTableAtPath(path)
}

func (repo *Repository) Excludes() (*format.GlobTable, error) {
	path := filepath.Join(repo.path, "info", "exclude")
	return format.GlobTableAtPath(path)
}
