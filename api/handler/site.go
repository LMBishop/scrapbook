package handler

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/pkg/upload"
)

func UploadSiteVersion(mainConfig *config.MainConfig, index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix("Bearer ", r.Header.Get("Authorization"))

		if len(mainConfig.Command.Secret) == 0 || subtle.ConstantTimeCompare([]byte(token), []byte(mainConfig.Command.Secret)) != 1 {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		site := r.PathValue("site")
		reader, err := r.MultipartReader()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "could not read stream: %s", err.Error())
			return
		}

		version, err := upload.HandleUpload(site, reader, index)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}

		fmt.Fprintf(w, "version created: %s", version)
	}
}
