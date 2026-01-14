package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const appFolderName = ".wft"
const configFileName = "config.json"
const stateFileName = "state.json"

type Config struct {
	CurrentUser string `json:"current_user"`
}

type State struct {
	PRs map[string]PRState `json:"prs"`
}

type PRState struct {
	Number             int       `json:"number"`
	Repo               string    `json:"repo"`
	HeadSHA            string    `json:"head_sha"`
	UpdatedAt          time.Time `json:"updated_at"`
	Etag               string    `json:"etag,omitempty"`
	Ignored            bool      `json:"ignored,omitempty"`
	WorkflowStatus     string    `json:"workflow_status"`
	WorkflowConclusion string    `json:"workflow_conclusion"`
	WorkflowRunID      int       `json:"workflow_run_id"`
	WorkflowURL        string    `json:"workflow_url"`
	LastCheckedSHA     string    `json:"last_checked_sha"`
}

func ReadConfig() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}

	defer func() {
		err := jsonFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	bytes, _ := io.ReadAll(jsonFile)

	var config Config

	err = json.Unmarshal([]byte(bytes), &config)
	if err != nil {
		return Config{}, err
	}

	log.Printf("loaded config")

	return config, nil
}

func (cfg *Config) WriteConfig() error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func SetUser(user string) error {
	return nil
}

func getConfigFilePath() (string, error) {
	path := filepath.Join(getHomeDir(), appFolderName, configFileName)
	return path, nil
}

func getStateFilePath() (string, error) {
	path := fmt.Sprintf("%s/%s/%s", getHomeDir(), appFolderName, stateFileName)
	return path, nil
}

func getHomeDir() string {
	var home string
	switch runtime.GOOS {
	case "windows":
		home = os.Getenv("USERPROFILE")
	default:
		home = os.Getenv("HOME")
	}
	return home
}

func ReadState() (State, error) {
	path, err := getStateFilePath()
	if err != nil {
		return State{}, err
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		return State{}, err
	}

	defer func() {
		err := jsonFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	b, _ := io.ReadAll(jsonFile)

	var state State
	err = json.Unmarshal([]byte(b), &state)
	if err != nil {
		return State{}, err
	}

	log.Printf("loaded state")

	return state, nil
}

func WriteState(s State) error {
	path, err := getStateFilePath()
	if err != nil {
		return err
	}

	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, b, 0600)
	if err != nil {
		return err
	}
	return nil
}
