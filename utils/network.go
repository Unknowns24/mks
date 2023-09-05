package utils

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/unknowns24/mks/config"
)

// chech if is an url or not (accepts http and https)
func IsUrl(url string) bool {
	return strings.HasPrefix(url, config.NETWORK_HTTP_PREFIX) || strings.HasPrefix(url, config.NETWORK_HTTPS_PREFIX)
}

// chech if is an github url  (package style like github.com/unknowns24/mks)
func IsGithubUrl(url string) bool {
	return strings.HasPrefix(url, config.NETWORK_GITHUB_DOMAIN)
}

// Download file from url to destination
func DownloadFile(url string, destPath string) error {
	// Create file
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Download file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
