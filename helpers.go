package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// gets the cache directory
func UserCacheDir() (string, error) {
	// try env var first
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return xdg, nil
	}

	// if it fails then assemble manually
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache"), nil
}

// makes relative URLs absolute
func ensureAbsoluteURL(base, ref string) string {
	if strings.HasPrefix(ref, "http") {
		return ref
	}
	return base + ref
}

// does a HEAD request to get Content-Length
func getContentLength(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("bad status: %s", resp.Status)
	}

	cl := resp.ContentLength
	if cl <= 0 {
		// Fallback: attempt GET with no body read
		resp, err = http.Get(url)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()

		n, err := io.Copy(io.Discard, resp.Body)
		if err != nil {
			return 0, err
		}
		return n, nil
	}

	return cl, nil
}
