package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Key    string `json:"key"`
	UserID string `json:"user_id"`
	AppID  string `json:"app_id"`
}

type VersionInfo struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

func createConfigTemplate(filePath string) error {
	configTemplate := Config{
		Key:    "your_key_here",
		UserID: "your_user_id_here",
		AppID:  "your_app_id_here",
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(configTemplate)
	if err != nil {
		return err
	}
	fmt.Printf("Template config file created at %s. Please update it with your values.\n", filePath)
	return nil
}

func loadConfig(filePath string) (Config, error) {
	var config Config
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := createConfigTemplate(filePath)
		if err != nil {
			return config, err
		}
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func getUploadURL(config Config, version string) (string, error) {
	url := "https://appho.st/api/get_upload_url"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Add("user_id", config.UserID)
	q.Add("app_id", config.AppID)
	q.Add("key", config.Key)
	q.Add("platform", "android")
	q.Add("version", version)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	return err == nil
}

func uploadFile(uploadURL string, filePath string) (int, error) {
	if !isValidURL(uploadURL) {
		return 0, fmt.Errorf("invalid upload URL: %s", uploadURL)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("PUT", uploadURL, bytes.NewBuffer(fileContent))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func getCurrentVersion(config Config) (VersionInfo, error) {
	var versionInfo VersionInfo
	url := "https://appho.st/api/get_current_version"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return versionInfo, err
	}
	q := req.URL.Query()
	q.Add("u", config.UserID)
	q.Add("a", config.AppID)
	q.Add("platform", "android")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return versionInfo, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return versionInfo, err
	}

	// Remove any extraneous characters like HTML tags
	bodyStr := strings.TrimSpace(string(body))
	if !isJSON(bodyStr) {
		return versionInfo, fmt.Errorf("response is not valid JSON: %s", bodyStr)
	}

	err = json.Unmarshal([]byte(bodyStr), &versionInfo)
	if err != nil {
		return versionInfo, err
	}

	return versionInfo, nil
}

func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: upload <file_path> <version>")
		os.Exit(1)
	}
	filePath := os.Args[1]
	version := os.Args[2]

	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	uploadURL, err := getUploadURL(config, version)
	if err != nil {
		fmt.Println("Error getting upload URL:", err)
		return
	}

	statusCode, err := uploadFile(uploadURL, filePath)
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return
	}

	if statusCode == http.StatusOK {
		fmt.Println("File uploaded successfully.")
		versionInfo, err := getCurrentVersion(config)
		if err != nil {
			fmt.Println("Error getting current version:", err)
			return
		}
		fmt.Printf("Current Version: %s\n", versionInfo.Version)
		fmt.Printf("Download URL: %s\n", versionInfo.URL)
	} else {
		fmt.Printf("Failed to upload file. Status code: %d\n", statusCode)
	}
}
