package keyring

import (
	"os"
	"path/filepath"
)

func Set(secret string) error {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "stack", "hosts.txt")
	err := ensureDirsExist(configPath)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, []byte(secret), 0644)
	if err != nil {
		return err
	}

	return nil
}

func Get() (string, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "stack", "hosts.txt")
	err := ensureDirsExist(configPath)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func Delete() error {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "stack", "hosts.txt")
	err := os.Remove(configPath)
	if err != nil {
		return err
	}

	return nil
}

func ensureDirsExist(filePath string) error {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
