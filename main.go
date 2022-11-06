package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
)

const (
	baseDownloadURL = "https://kpolyakov.spb.ru/cms/files/ege-proc/"
	baseTaskURL     = "https://kpolyakov.spb.ru/school/ege/gen.php?action=viewTopic&topicId"
)

var client = http.DefaultClient

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		topic := scanner.Text()
		if topic == "" {
			panic("invalid topic")
		}
		var cmd *exec.Cmd
		url := fmt.Sprintf("%s=%s", baseTaskURL, topic)

		switch runtime.GOOS {
		case "windows":
			winargs := []string{"cmd", "/c", "start"}
			cmd = exec.Command(winargs[0], append(winargs[1:], url)...)
		case "linux":
			cmd = exec.Command("xdg-open", url)
		}

		err := cmd.Run()
		if err != nil {
			log.Fatalf("could not start browser: %v", err)
		}

		content, err := fetchPage(url)
		if err != nil {
			log.Fatal(err.Error())
		}

		rxp, err := regexp.Compile(`22-\d+\.xls`)
		if err != nil {
			log.Fatalf("could not compile regexp: %v", err)
		}

		inSiteFileName := rxp.FindString(content)

		downloadURL := fmt.Sprintf("%s/%s", baseDownloadURL, inSiteFileName)

		fileData, err := fetchFile(downloadURL)
		if err != nil {
			log.Fatalf("could not fetch a file: %v", err)
		}

		fname := fmt.Sprintf("%s_%s", topic, inSiteFileName)

		saveFile(fname, fileData)

		fmt.Printf("file has saved under: %s", fname)
	}

}

func fetchFile(url string) ([]byte, error) {

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not fetch a file: %w", err)
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read file's body: %w", err)
	}

	return raw, nil
}

func saveFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0666)
}

func fetchPage(url string) (string, error) {
	res, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("could not fetch a page: %w", err)
	}

	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("could not read body: %w", err)
	}

	return string(raw), err
}
