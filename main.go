package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/PuerkitoBio/goquery"
)

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

	return href, nil
}

// DownloadImage downloads the image from APOD given a relative href.
func DownloadImage(href string) error {
	// check if file already exists
	if _, err := os.Stat(path.Base(href)); err == nil {
		fmt.Println("already downloaded", path.Base(href))
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
	out, err := os.Create(path.Base(href))
	if err != nil {
		return err
	}
	defer out.Close()

	// download
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("downloaded", path.Base(href))
	return nil
}

// SetWallpaper sets the desktop background using feh --bg-fill.
func SetWallpaper(filepath string) error {
	cmd := exec.Command("feh", "--bg-fill", filepath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to set wallpaper: %w", err)
	}
	return nil
}

func main() {
	href, err := GetAPODImageHref()
	if err != nil {
		panic(err)
	}

	if err := DownloadImage(href); err != nil {
		panic(err)
	}

	filename := path.Base(href)
	if err := SetWallpaper(filename); err != nil {
		panic(err)
	}
}
