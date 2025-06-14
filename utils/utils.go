package utils

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"
)

const maxRetry = 5

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

var (
	ErrCommandNotFound = errors.New("file not found")
)

func MakesureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
func DownloadFile(url string, path string, cmd string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		return ErrCommandNotFound
	}

	defer resp.Body.Close()

	content := make([]byte, 0)
	reader := bufio.NewReader(resp.Body)
	for {
		line, _, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		content = append(content, line...)
		content = append(content, []byte("\n")...)
	}
	err = os.WriteFile(path, content, 0666)
	if err != nil {
		return err
	}
	return nil
}

func RetryDownloadFile(url string, path string, cmd string) error {
	for j := 0; j < maxRetry; j++ {
		if err := DownloadFile(url, path, cmd); err != nil {
			continue
		}
		break
	}
	return nil
}
