package site

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/LMBishop/scrapbook/pkg/config"
)

const versionRegex = "[0-9a-f]{64}"

type Site struct {
	Name   string
	Path   string
	Flags  config.SiteFlag
	Config *config.SiteConfig

	handler http.Handler
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

	cfg, strCfg := config.NewSiteConfig(host)
	site := &Site{
		Name:   name,
		Path:   dir,
		Flags:  0,
		Config: &cfg,
	}

	err = os.Mkdir(dir, 0o755)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	err = os.WriteFile(path.Join(dir, "config"), []byte(strCfg), 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to write site config: %w", err)
	}

	return site, nil
}

func (s *Site) Initialise() {
	s.handler = nil
	if s.Config != nil {
		s.handler = NewSiteFileServer(http.Dir(path.Join(s.GetCurrentPath(), "webroot")), slog.With("site", s.Name), s.Flags, *s.Config)
	}
}

func (s *Site) Handler() http.Handler {
	return s.handler
}

func (s *Site) GetCurrentPath() string {
	return path.Join(s.Path, "current")
}

func (s *Site) GetCurrentVersion() (string, error) {
	dir, err := filepath.EvalSymlinks(path.Join(s.Path, "current"))
	if err != nil {
		return "", err
	}

	return filepath.Base(dir), nil
}

func (s *Site) ReadVersionMetadata(version string) (string, error) {
	config, err := os.ReadFile(path.Join(s.Path, version, "metadata"))
	return string(config), err
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

func (s *Site) DeleteVersion(version string) error {
	if version == "" {
		return fmt.Errorf("not removing empty version")
	}

	return os.RemoveAll(path.Join(s.Path, version))
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

func (s *Site) CreateNewVersion(version string, size int64, filecount uint, via string) error {
	t := time.Now()
	newVersionDir := path.Join(s.Path, version)

	_, err := os.Stat(newVersionDir)
	if err == nil {
		return errors.New("version already exists")
	}

	err = os.MkdirAll(newVersionDir, os.FileMode(0o755))
	if err != nil {
		return err
	}

	_, versionMeta := config.NewVersionMeta(t, version, size, filecount, via)
	err = os.WriteFile(path.Join(newVersionDir, "metadata"), []byte(versionMeta), 0o644)
	if err != nil {
		return err
	}

	return nil
}

func (s *Site) EvaluateSiteStatus() (string, string) {
	if s.Config == nil {
		return "error", "This site does not have a valid configuration"
	}
	if s.Config.Host == "" {
		return "inactive", "This site is not served by scrapbook"
	}
	stat, err := os.Stat(s.GetCurrentPath())
	if err != nil || !stat.IsDir() {
		return "disabled", "This site is disabled because no version is active"
	}
	if s.Flags&config.FlagDisable != 0 {
		return "disabled", "This site is disabled because the disable flag is set"
	}
	if s.Flags&config.FlagPassword != 0 && s.Config.Password == "" {
		return "error", "This site is live, but inacessible because the password is empty and the password protect flag is set"
	}

	return "live", "This site is live"
}

func (s *Site) ConvertFlagsToString() string {
	var bits []string
	bitNames := []string{"D", "T", "I", "P", "R"}

	for i := 0; i < len(bitNames); i++ {
		if s.Flags&(1<<i) != 0 {
			bits = append(bits, bitNames[i])
		} else {
			bits = append(bits, "-")
		}
	}

	return strings.Join(bits, "")
}

func (s *Site) ReadSiteConfigFile() (string, error) {
	config, err := os.ReadFile(path.Join(s.Path, "config"))
	return string(config), err
}

func (s *Site) WriteSiteConfigFile(config string) error {
	return os.WriteFile(path.Join(s.Path, "config"), []byte(config), 0o644)
}

func (s *Site) ReadSiteFlags() config.SiteFlag {
	flags, _ := os.ReadFile(path.Join(s.Path, "flags"))
	val, _ := strconv.ParseUint(string(flags), 10, 64)
	return config.SiteFlag(val)
}

func (s *Site) WriteSiteFlags(flags config.SiteFlag) error {
	return os.WriteFile(path.Join(s.Path, "flags"), []byte(strconv.Itoa(int(flags))), 0o644)
}
