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

func GetAPODImageHref() (string, error) {
	baseURL := "https://apod.nasa.gov/apod/"

	// get the page
	resp, err := http.Get(baseURL + "astropix.html")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// read html
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// original <a href>
	aHrefSel := doc.Find("html > body > center").First().
		Find("p").Eq(1).
		Find("a").First()

	href, exists := aHrefSel.Attr("href")
	if !exists {
		return "", fmt.Errorf("href not found")
	}

	// corresponding <img src>
	imgSel := aHrefSel.Find("img").First()
	src, srcExists := imgSel.Attr("src")
	if !srcExists {
		return baseURL + href, nil // fallback
	}

	// normalize URLs
	hrefURL := ensureAbsoluteURL(baseURL, href)
	srcURL := ensureAbsoluteURL(baseURL, src)

	// size check
	hrefSize, err := getContentLength(hrefURL)
	if err != nil {
		fmt.Println("warning: could not get size of href:", err)
	}

	srcSize, err := getContentLength(srcURL)
	if err != nil {
		fmt.Println("warning: could not get size of src:", err)
	}

	// choose the beefier image
	if srcSize > hrefSize {
		fmt.Printf("%s > %s\n", srcURL, hrefURL)
		return srcURL, nil
	}
	fmt.Printf("%s < %s\n", srcURL, hrefURL)
	return hrefURL, nil
}

// downloads a file over http, extracts it to userCacheDir().
func DownloadFile(href, path string) error {
	// check if file already exists
	if _, err := os.Stat(path); err == nil {
		fmt.Println("already downloaded", path)
		return nil
	} else if !os.IsNotExist(err) {
		return err // ????
	}

	// get the page
	resp, err := http.Get(href)
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

	fmt.Printf("downloaded %s to %s\n", href, path)
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
