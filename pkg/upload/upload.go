package upload

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/index"
	"golang.org/x/mod/sumdb/dirhash"
)

func HandleUpload(siteName, source, via string, reader *multipart.Reader, index *index.SiteIndex) (string, error) {
	s := index.GetSite(siteName)
	if s == nil {
		return "", fmt.Errorf("no such site: %s", siteName)
	}

	if s.Flags&config.FlagReadOnly != 0 {
		return "", fmt.Errorf("site is read only: %s", siteName)
	}

	uploadedZip, err := os.CreateTemp(os.TempDir(), "scrapbookupload")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer func() {
		uploadedZip.Close()
		os.Remove(uploadedZip.Name())
	}()

	var zipName string

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", fmt.Errorf("failed to read multipart stream: %w", err)
		}
		if part.FormName() == "upload" {
			zipName = part.FileName()
			io.Copy(uploadedZip, part)
		}
	}

	tmpDest, err := os.MkdirTemp(s.Path, "tmp")
	defer func() {
		os.Remove(tmpDest)
	}()
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}

	zipReader, err := zip.OpenReader(uploadedZip.Name())
	if err != nil {
		return "", fmt.Errorf("failed to open zip reader: %w", err)
	}
	defer zipReader.Close()

	count, size, err := unzipSource(uploadedZip.Name(), tmpDest)
	if err != nil {
		return "", fmt.Errorf("failed to unzip archive: %w", err)
	}

	versionHash, err := dirhash.HashDir(tmpDest, "", notHash1)
	if err != nil {
		return "", fmt.Errorf("failed to hash directory: %w", err)
	}

	err = s.CreateNewVersion(versionHash, "ZIPArchive", zipName, size, count, source, via)
	if err != nil {
		return "", fmt.Errorf("failed to create new version: %w", err)
	}
	webroot := path.Join(s.Path, versionHash, "webroot")

	err = os.Rename(tmpDest, webroot)
	if err != nil {
		return "", fmt.Errorf("failed to move tmp directory to webroot: %w", err)
	}

	err = s.UpdateVersion(versionHash)
	if err != nil {
		return "", fmt.Errorf("failed to update version: %w", err)
	}

	s.ProcessRetention()

	return versionHash, nil
}

// https://gosamples.dev/unzip-file/
// (adapted to add mtimes)

func unzipSource(source, destination string) (uint, int64, error) {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return 0, 0, err
	}
	defer reader.Close()

	destination, err = filepath.Abs(destination)
	if err != nil {
		return 0, 0, err
	}

	var count uint
	var totalSize int64

	for _, f := range reader.File {
		n, err := unzipFile(f, destination)
		if err != nil {
			return 0, 0, err
		}
		if n > 0 {
			count++
			totalSize = totalSize + n
		}
	}

	return count, totalSize, nil
}

func unzipFile(f *zip.File, destination string) (int64, error) {
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return 0, fmt.Errorf("invalid file path: %s", filePath)
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return 0, err
		}
		if modTime := f.FileHeader.Modified; !modTime.IsZero() {
			os.Chtimes(filePath, modTime, modTime)
		}
		return 0, nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return 0, err
	}

	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return 0, err
	}
	defer destinationFile.Close()

	zippedFile, err := f.Open()
	if err != nil {
		return 0, err
	}
	defer zippedFile.Close()

	n, err := io.Copy(destinationFile, zippedFile)
	if err != nil {
		return 0, err
	}

	if modTime := f.FileHeader.Modified; !modTime.IsZero() {
		os.Chtimes(filePath, modTime, modTime)
	}

	return n, nil
}

// This is copied almost verbaitm from dirhash.Hash1 except modified to return as
// hex and remove h1: prefix
func notHash1(files []string, open func(string) (io.ReadCloser, error)) (string, error) {
	h := sha256.New()
	files = append([]string(nil), files...)
	slices.Sort(files)
	for _, file := range files {
		if strings.Contains(file, "\n") {
			return "", errors.New("dirhash: filenames with newlines are not supported")
		}
		r, err := open(file)
		if err != nil {
			return "", err
		}
		hf := sha256.New()
		_, err = io.Copy(hf, r)
		r.Close()
		if err != nil {
			return "", err
		}
		fmt.Fprintf(h, "%x  %s\n", hf.Sum(nil), file)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
