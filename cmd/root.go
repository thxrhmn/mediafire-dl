package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cheggaaa/pb/v3"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mediafire-dl [flags] [URLs]",
	Short: "A mediafire downloader",
	Long:  `A simple CLI tool to download files from mediafire URLs with a progress bar.`,
	Run:   run,
}

var file string

func init() {
	rootCmd.Flags().StringVarP(&file, "file", "f", "", "Path to a file containing URLs (one per line)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetDownloadableURL(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	// Find the downloadable URL
	downloadableURL, isExist := doc.Find(".input.popsok").Attr("href")
	if !isExist {
		return "", errors.New("downloadable URL not found")
	}

	return downloadableURL, nil
}

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the size of the file
	contentLength := resp.Header.Get("Content-Length")
	size, err := strconv.Atoi(contentLength)
	if err != nil {
		return err
	}

	// Create a progress bar and set it up
	bar := pb.Full.Start64(int64(size))
	defer bar.Finish()

	// Create a proxy reader
	reader := bar.NewProxyReader(resp.Body)

	// Copy the data from the response to the file
	buf := make([]byte, 512*1024) // 512 KB chunks
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			if _, err := out.Write(buf[:n]); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func run(cmd *cobra.Command, args []string) {
	figure := figure.NewFigure("Mediafire-dl", "smslant", false)
	figure.Print()
	fmt.Print("\nCreated by: @thxrhmn")
	fmt.Print("\nCreated at: 2024/06/10\n\n")

	var urls []string

	if file != "" {
		file, err := os.Open(file)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			urls = append(urls, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading file: %v", err)
		}
	} else {
		urls = args
	}

	if len(urls) == 0 {
		log.Fatalf("No URLs provided")
	}

	for _, downloadPageURL := range urls {
		downloadURL, err := GetDownloadableURL(downloadPageURL)
		if err != nil {
			log.Printf("Error fetching downloadable URL for %s: %v", downloadPageURL, err)
			continue
		}

		// Extract the filename from the download URL
		filename := path.Base(strings.Split(downloadURL, "/")[len(strings.Split(downloadURL, "/"))-1])

		fmt.Println("Downloading File:", filename)

		err = DownloadFile(filename, downloadURL)
		if err != nil {
			log.Printf("Error downloading file %s: %v", filename, err)
			continue
		}

		fmt.Printf("File downloaded successfully: %s\n\n", filename)
	}
}
