package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
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

// gets the link for the latest APOD image
func GetAPODImageHref() (string, error) {
	// get the page
	resp, err := http.Get("https://apod.nasa.gov/apod/astropix.html")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// read html
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// find image link element
	selection := doc.Find("html > body > center").First().
		Find("p").Eq(1).
		Find("a").First()

	// get the link
	href, exists := selection.Attr("href")
	if !exists {
		return "", fmt.Errorf("href not found")
	}

	fmt.Println("found image link:", href)

	return href, nil
}

// downloads a file over http, extracts it to userCacheDir().
func DownloadFile(href string, path string) error {
	// check if file already exists
	if _, err := os.Stat(path); err == nil {
		fmt.Println("already downloaded", path)
		return nil
	} else if !os.IsNotExist(err) {
		return err // ????
	}

	// get the page
	resp, err := http.Get("https://apod.nasa.gov/apod/" + href)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// create file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// download
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("downloaded", path)
	return nil
}

// feh --bg-fill {filepath}
func SetWallpaper(path string) error {
	cmd := exec.Command("feh", "--bg-fill", path)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to set wallpaper: %w", err)
	}
	fmt.Println("wallpaper set to", path)
	return nil
}

func main() {
	// get latest apod image link
	href, err := GetAPODImageHref()
	if err != nil {
		panic(err)
	}

	// assembles filepath for the image
	cacheDir, err := UserCacheDir()
	if err != nil {
		panic(err)
	}
	path := filepath.Join(cacheDir, filepath.Base(href))

	// downloads the image
	if err := DownloadFile(href, path); err != nil {
		panic(err)
	}

	// sets the downloaded image as the background
	if err := SetWallpaper(path); err != nil {
		panic(err)
	}
}
