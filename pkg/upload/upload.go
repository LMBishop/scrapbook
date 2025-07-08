package upload

import (
	"archive/zip"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/LMBishop/scrapbook/pkg/index"
)

func HandleUpload(siteName string, reader *multipart.Reader, index *index.SiteIndex) (string, error) {
	s := index.GetSite(siteName)
	if s == nil {
		return "", fmt.Errorf("no such site: %s", siteName)
	}

	temp, err := os.CreateTemp(os.TempDir(), "scrapbook")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer func() {
		temp.Close()
		os.Remove(temp.Name())
	}()

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", fmt.Errorf("failed to read multipart stream: %w", err)
		}
		if part.FormName() == "upload" {
			io.Copy(temp, part)
		}
	}

	zipReader, err := zip.OpenReader(temp.Name())
	if err != nil {
		return "", fmt.Errorf("failed to open zip reader: %w", err)
	}
	defer zipReader.Close()

	version, err := s.CreateNewVersion()
	if err != nil {
		return "", fmt.Errorf("failed to create new version: %w", err)
	}
	versionDir := path.Join(s.Path, version)

	err = unzipSource(temp.Name(), versionDir)
	if err != nil {
		return "", fmt.Errorf("failed to unzip archive: %w", err)
	}

	err = s.UpdateVersion(version)
	if err != nil {
		return "", fmt.Errorf("failed to update version: %w", err)
	}

	return version, nil
}

// https://gosamples.dev/unzip-file/

func unzipSource(source, destination string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}

	for _, f := range reader.File {
		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(f *zip.File, destination string) error {
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}
