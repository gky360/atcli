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
	"strings"

	pb "gopkg.in/cheggaaa/pb.v2"

	"github.com/gky360/atsrv/models"
)

func DownloadTestcases(contest *models.Contest, tasks []*models.Task) error {
	tempdir, err := ioutil.TempDir("", "atcli_"+contest.ID+"_")
	if err != nil {
		return err
	}
	defer func() {
		os.RemoveAll(tempdir)
		fmt.Println("Cleaned", tempdir, ".")
	}()

	filename := contest.ID + ".zip"
	zipPath := filepath.Join(tempdir, filename)
	if err = downloadFromUrl(contest.TestcasesURL, zipPath); err != nil {
		return err
	}

	extractedDir := filepath.Join(tempdir, contest.ID)
	if err = unzip(zipPath, extractedDir); err != nil {
		return err
	}

	for _, task := range tasks {
		src := filepath.Join(extractedDir, strings.ToLower(task.Name))
		dest, err := TaskSampleDir(task.Name, true)
		if err != nil {
			return err
		}
		destRel, err := filepath.Rel(RootDir(), dest)
		if err != nil {
			return err
		}
		if _, err := os.Stat(dest); err == nil {
			// Already exists
			fmt.Println("Already exists:", destRel)
			continue
		}

		if err = os.Rename(src, dest); err != nil {
			return err
		}

		fmt.Println("Created folder:", destRel)
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

	bar := pb.Full.Start(contentLength)
	defer bar.Finish()

	reader := bar.NewProxyReader(resp.Body)

	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}

	fmt.Println("Downloaded", fpath, ".")
	return nil
}

func unzip(src, dest string) error {
	fmt.Println("Unziping", src)
	fmt.Println("to", dest, "...")

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	filesCnt := len(r.File)
	bar := pb.Full.Start(filesCnt)

	for _, f := range r.File {
		bar.Set("prefix", fnameToPrefix(f.Name)).
			Write().
			Increment()
		if err = extractAndWriteFile(f, dest); err != nil {
			return err
		}
	}

	bar.Finish()
	fmt.Println("Unzipped", src, ".")
	return nil
}

func fnameToPrefix(fname string) string {
	return fmt.Sprintf("%-20s", fname)[:20] + " "
}

func extractAndWriteFile(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	fpath := filepath.Join(dest, f.Name)

	if f.FileInfo().IsDir() {
		if err = os.MkdirAll(fpath, f.Mode()); err != nil {
			return err
		}
	} else {
		if err = os.MkdirAll(filepath.Dir(fpath), f.Mode()); err != nil {
			return err
		}

		out, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, rc)
		if err != nil {
			return err
		}
	}
	return nil
}
