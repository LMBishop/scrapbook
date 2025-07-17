package site

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
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
	site.Handler = NewSiteFileServer(http.Dir(path.Join(dir, "default")), config)
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

func (s *Site) DeleteDataOnDisk() error {
	return os.RemoveAll(s.Path)
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

	slices.Sort(versions)
	slices.Reverse(versions)

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
	if s.SiteConfig.Host == "" {
		return "inactive"
	}
	stat, err := os.Stat(s.GetCurrentPath())
	if err != nil || !stat.IsDir() {
		return "inactive"
	}
	if s.SiteConfig.Flags&config.FlagDisable != 0 {
		return "inactive"
	}

	return "live"
}

func (s *Site) EvaluateSiteStatusReason() string {
	if s.SiteConfig.Host == "" {
		return "This site is not served by scrapbook"
	}
	stat, err := os.Stat(s.GetCurrentPath())
	if err != nil || !stat.IsDir() {
		return "This site is inacessible because no version is active"
	}
	if s.SiteConfig.Flags&config.FlagDisable != 0 {
		return "This site is inacessible because it is disabled"
	}

	return "This site is live"
}

func (s *Site) ConvertFlagsToString() string {
	var bits []string
	bitNames := []string{"D", "T", "I", "P", "R"}

	for i := 0; i < len(bitNames); i++ {
		if s.SiteConfig.Flags&(1<<i) != 0 {
			bits = append(bits, bitNames[i])
		} else {
			bits = append(bits, "-")
		}
	}

	return strings.Join(bits, "")
}
