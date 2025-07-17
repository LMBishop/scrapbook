package handler

import (
	"crypto/subtle"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/pkg/upload"
)

func UploadSiteVersion(mainConfig *config.MainConfig, index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		if len(mainConfig.Command.Secret) == 0 || subtle.ConstantTimeCompare([]byte(token), []byte(mainConfig.Command.Secret)) != 1 {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "forbidden")
			return
		}

		site := r.PathValue("site")
		reader, err := r.MultipartReader()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "could not read stream: %s", err.Error())
			slog.Error("could not read stream", "remoteAddr", r.RemoteAddr, "error", err)
			return
		}

		version, err := upload.HandleUpload(site, reader, index)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			slog.Error("could not handle upload", "remoteAddr", r.RemoteAddr, "error", err)
			return
		}

		slog.Info("new version created", "site", site, "version", version, "remoteAddr", r.RemoteAddr)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "version created: %s", version)
	}
}
