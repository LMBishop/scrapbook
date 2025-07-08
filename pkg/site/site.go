package site

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/LMBishop/scrapbook/pkg/config"
)

const versionRegex = "[0-9]{4}_[0-9]{2}_[0-9]{2}_[0-9]{2}_[0-9]{2}_[0-9]{2}"
const timeFormat = "2006_01_02_15_04_05"

type Site struct {
	Name       string
	Path       string
	Handler    http.Handler
	SiteConfig *config.SiteConfig
}

func NewSite(name string, dir string, config *config.SiteConfig) *Site {
	var site Site
	site.Name = name
	site.Path = dir
	site.SiteConfig = config
	site.Handler = http.FileServer(siteFS{http.Dir(path.Join(dir, "default"))})
	return &site
}

func CreateNewSite(name string, baseDir string, host string) (*Site, error) {
	dir := path.Join(baseDir, name)
	_, err := os.Stat(dir)
	if err == nil {
		return nil, fmt.Errorf("site with name already exists: %s", name)
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to check site uniqueness: %w", err)
	}

	cfg := &config.SiteConfig{
		Host: host,
	}
	site := NewSite(name, dir, cfg)

	err = os.Mkdir(dir, 0o755)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	err = config.WriteSiteConfig(path.Join(dir, "site.toml"), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to write site config: %w", err)
	}

	return site, nil
}

func (s *Site) GetCurrentPath() string {
	return path.Join(s.Path, "default")
}

func (s *Site) GetCurrentVersion() (string, error) {
	dir, err := filepath.EvalSymlinks(path.Join(s.Path, "default"))
	if err != nil {
		return "", err
	}

	return filepath.Base(dir), nil
}

func (s *Site) UpdateVersion(newVersion string) error {
	newVersionPath := path.Join(s.Path, newVersion)

	stat, err := os.Stat(newVersionPath)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("not a directory: %s", newVersionPath)
	}

	currentVersionPath := s.GetCurrentPath()

	os.Remove(currentVersionPath)
	return os.Symlink(newVersion, currentVersionPath)
}

func (s *Site) GetAllVersions() ([]string, error) {
	entries, err := os.ReadDir(s.Path)

	if err != nil {
		return nil, err
	}

	var versions []string

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		match, err := regexp.MatchString(versionRegex, entry.Name())
		if err != nil {
			return nil, err
		}

		if match {
			versions = append(versions, entry.Name())
		}
	}

	return versions, err
}

func (s *Site) CreateNewVersion() (string, error) {
	t := time.Now()
	dirName := t.Format("2006_01_02_15_04_05")
	newVersionDir := path.Join(s.Path, dirName)

	err := os.MkdirAll(newVersionDir, os.FileMode(0o755))
	if err != nil {
		return "", err
	}

	return dirName, nil
}

func (s *Site) EvaluateSiteStatus() string {
	stat, err := os.Stat(s.GetCurrentPath())
	if err != nil || !stat.IsDir() {
		return "inactive"
	}

	return "live"
}

func (s *Site) EvaluateSiteStatusReason() string {
	stat, err := os.Stat(s.GetCurrentPath())
	if err != nil || !stat.IsDir() {
		return "This site is inacessible because no version is active"
	}

	return "This site is live"
}
