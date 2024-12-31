package maven

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Caedis/gtnh-updater/internal/utils"
)

type Downloader struct {
	CacheDir string
}

func getDefaultCacheDir() (string, error) {
	var cacheDir string
	switch runtime.GOOS {
	case "windows":
		// Windows: Use %LocalAppData%
		localAppData := os.Getenv("LocalAppData")
		if localAppData == "" {
			return "", fmt.Errorf("LocalAppData environment variable is not set")
		}
		cacheDir = filepath.Join(localAppData)
	case "darwin":
		// macOS: Use ~/Library/Caches
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		cacheDir = filepath.Join(homeDir, "Library", "Caches")
	default:
		// Linux: Use $XDG_CACHE_HOME or ~/.cache
		xdgCacheHome := os.Getenv("XDG_CACHE_HOME")
		if xdgCacheHome == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get user home directory: %w", err)
			}
			cacheDir = filepath.Join(homeDir, ".cache")
		} else {
			cacheDir = xdgCacheHome
		}
	}

	cacheDir = filepath.Join(cacheDir, "gtnh-updater", "mods")

	return cacheDir, nil
}

func (d *Downloader) Download(filename, url string) (string, error) {
	cachePath := filepath.Join(d.CacheDir, filename)

	// Download file to cache
	if !utils.FileExists(cachePath) {
		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("unexpected HTTP status: %s", resp.Status)
		}

		out, err := os.Create(cachePath)
		if err != nil {
			return "", fmt.Errorf("Failed creating cache file %s: %w", cachePath, err)
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return "", fmt.Errorf("Failed writing to cache file %s: %w", cachePath, err)
		}
	}

	return cachePath, nil
}

func NewDownloader() (*Downloader, error) {
	cacheDir, err := getDefaultCacheDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &Downloader{
		CacheDir: cacheDir,
	}, nil
}
