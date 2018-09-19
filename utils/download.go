package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/gky360/atsrv/models"
)

func DownloadTestcases(contest *models.Contest, tasks []*models.Task) error {
	tempdir, err := ioutil.TempDir("", contest.ID+"_")
	if err != nil {
		return err
	}
	// defer os.RemoveAll(tempdir)

	filename := contest.ID + ".zip"
	zipPath := filepath.Join(tempdir, filename)
	if err := downloadFromUrl(contest.TestcasesURL, zipPath); err != nil {
		return err
	}

	if err := unzip(zipPath, tempdir); err != nil {
		return err
	}

	return nil
}

func downloadFromUrl(url string, fpath string) error {
	fmt.Println("Downloading from", url)
	fmt.Println("to", fpath, "...")

	out, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Server returned non-200 status: %v", resp.Status)
	}
	contentLength, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	bar := pb.New(contentLength).SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true
	bar.Start()
	reader := bar.NewProxyReader(resp.Body)

	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}
	bar.Finish()

	fmt.Println("Downloaded", fpath, ".")
	return nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fmt.Println(f.Name)
	}

	return nil
}
