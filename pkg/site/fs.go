package site

import (
	"crypto/subtle"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/LMBishop/scrapbook/pkg/auth"
	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/html"
)

type SiteFileServer struct {
	root          http.FileSystem
	logger        *slog.Logger
	flags         config.SiteFlag
	config        config.SiteConfig
	authenticator auth.Authenticator
}

func NewSiteFileServer(root http.FileSystem, logger *slog.Logger, flags config.SiteFlag, config config.SiteConfig) *SiteFileServer {
	return &SiteFileServer{root: root, logger: logger, flags: flags, config: config, authenticator: *auth.NewAuthenticator()}
}

func (fs *SiteFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if fs.flags&config.FlagDisable != 0 {
		w.WriteHeader(http.StatusForbidden)
		html.ForbiddenDisabledPage(fs.config.Host).Render(w)
		return
	}

	if fs.flags&config.FlagPassword != 0 {
		jwt, err := r.Cookie("session")
		if err != nil {
			goto deny
		}

		err = fs.authenticator.VerifyJwt(jwt.Value)
		if err != nil {
			goto deny
		}

		goto permit

	deny:
		if strings.HasPrefix(r.URL.Path, "/authenticate") {
			goto ask
		}
		http.Redirect(w, r, "/authenticate?redirect="+url.QueryEscape(r.URL.Path), 302)
		return

	ask:
		fs.handleAskPassword(w, r)
		return

	permit:
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
			if fs.flags&config.FlagIndex == 0 {
				goto notFound
			}
			files, err := fs.listFiles(path)
			if path != "/" {
				files = append([]html.File{{Name: "..", IsDir: true, Size: 0}}, files...)
			}
			if err != nil {
				html.IndexPage(path, true, files).Render(w)
				fs.logger.Error("could not list directory for index page generation", "path", path, "error", err)
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
	html.NotFoundUrlPage(path).Render(w)
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
			files = append(files, html.File{Name: entry.Name(), IsDir: false, Size: entry.Size(), Mtime: entry.ModTime()})
		} else {
			files = append(files, html.File{Name: entry.Name(), IsDir: true, Size: 0})
		}
	}

	return files, nil
}

func (fs *SiteFileServer) handleAskPassword(w http.ResponseWriter, r *http.Request) {
	redirect := r.URL.Query().Get("redirect")

	switch r.Method {
	case "GET":
		html.AuthenticateSitePage("", redirect).Render(w)
	case "POST":
		err := r.ParseForm()
		if err != nil {
			html.AuthenticateSitePage(err.Error(), redirect).Render(w)
			return
		}

		password := r.Form.Get("password")

		if len(fs.config.Password) == 0 || subtle.ConstantTimeCompare([]byte(password), []byte(fs.config.Password)) != 1 {
			html.AuthenticateSitePage("The password is incorrect", redirect).Render(w)
			return
		}

		jwt, err := fs.authenticator.NewJwt()
		if err != nil {
			html.AuthenticateSitePage(fmt.Errorf("Failed to create jwt: %w", err).Error(), redirect).Render(w)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: jwt,

			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
		})
		http.Redirect(w, r, redirect, 302)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
