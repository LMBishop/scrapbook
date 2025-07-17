package site

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/html"
)

type SiteFileServer struct {
	root       http.FileSystem
	siteConfig *config.SiteConfig
}

func NewSiteFileServer(root http.FileSystem, siteConfig *config.SiteConfig) *SiteFileServer {
	return &SiteFileServer{root: root, siteConfig: siteConfig}
}

func (fs *SiteFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if fs.siteConfig.Flags&config.FlagDisable != 0 {
		w.WriteHeader(http.StatusForbidden)
		html.ForbiddenDisabledPage(fs.siteConfig.Host).Render(w)
		return
	}

	path := filepath.Clean(r.URL.Path)

	var info os.FileInfo
	file, err := fs.root.Open(path)
	if err != nil {
		if strings.HasSuffix(path, ".html") {
			goto notFound
		}

		htmlPath := path + ".html"
		file, err = fs.root.Open(htmlPath)
		if err != nil {
			goto notFound
		}
	}
	defer file.Close()

	info, err = file.Stat()
	if err != nil {
		goto notFound
	}

	if info.IsDir() {
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, path+"/", http.StatusFound)
			return
		}
		indexPath := filepath.Join(path, "index.html")
		if file, err = fs.root.Open(indexPath); os.IsNotExist(err) {
			if fs.siteConfig.Flags&config.FlagIndex == 0 {
				goto notFound
			}
			files, err := fs.listFiles(path)
			if path != "/" {
				files = append([]html.File{{Name: "..", IsDir: true, Size: 0}}, files...)
			}
			if err != nil {
				html.IndexPage(path, true, files).Render(w)
				slog.Error("could not list directory for index page generation", "host", fs.siteConfig.Host, "path", path, "error", err)
			} else {
				html.IndexPage(path, false, files).Render(w)
			}
			return
		}

		info, err = file.Stat()
		if err != nil {
			goto notFound
		}
	}

	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
	return

notFound:
	w.WriteHeader(http.StatusNotFound)
	html.NotFoundUrlPage(path, fs.siteConfig.Host).Render(w)
}

func (fs *SiteFileServer) listFiles(dir string) ([]html.File, error) {
	file, err := fs.root.Open(dir)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries, err := file.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var files []html.File
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, html.File{Name: entry.Name(), IsDir: false, Size: entry.Size()})
		} else {
			files = append(files, html.File{Name: entry.Name(), IsDir: true, Size: 0})
		}
	}

	return files, nil
}
