package main

import (
	"encoding/json"
	"os"
	"os/exec"
)

func readTool(args string) (string, error) {
	var jsonArgs struct {
		FilePath string `json:"file_path"`
	}
	err := json.Unmarshal([]byte(args), &jsonArgs)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(jsonArgs.FilePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func writeTool(args string) error {
	var jsonArgs struct {
		FilePath string `json:"file_path"`
		Content  string `json:"content"`
	}

	err := json.Unmarshal([]byte(args), &jsonArgs)
	if err != nil {
		return err
	}

	err = os.WriteFile(jsonArgs.FilePath, []byte(jsonArgs.Content), 0o644)
	if err != nil {
		return err
	}

	return nil
}

func bashTool(args string) (string, error) {
	var jsonArgs struct {
		Command string `json:"command"`
	}

	err := json.Unmarshal([]byte(args), &jsonArgs)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("bash", "-c", jsonArgs.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
