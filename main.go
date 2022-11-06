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
	"strings"
)

const (
	baseDownloadURL = "https://kpolyakov.spb.ru/cms/files/ege-proc/"
	baseTaskURL     = "https://kpolyakov.spb.ru/school/ege/gen.php?action=viewTopic&topicId"
)

var client = http.DefaultClient

func main() {
	run()
}

func run() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		topic := scanner.Text()
		if topic == "" {
			panic("invalid topic")
		}

		url := fmt.Sprintf("%s=%s", baseTaskURL, topic)

		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "windows":
			// cmd does not accept raw ampersants
			winURL := strings.ReplaceAll(url, "&", "^&")
			cmd = exec.Command("cmd", "/c", "start", winURL)
		case "linux":
			cmd = exec.Command("xdg-open", url)
		}

		err := cmd.Start()
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

		fmt.Printf("file has saved: %s\n", fname)
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
	return os.WriteFile(name, data, 0777)
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
