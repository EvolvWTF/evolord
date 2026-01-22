//go:build windows
// +build windows

package persistence

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const registryKey = `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
const registryValueName = "EvolordAgent"

func install(exePath string) error {

	appDataDir := os.Getenv("APPDATA")
	if appDataDir == "" {
		return fmt.Errorf("APPDATA environment variable not set")
	}

	evolordDir := filepath.Join(appDataDir, "Evolord")
	err := os.MkdirAll(evolordDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create Evolord directory: %w", err)
	}

	targetPath := filepath.Join(evolordDir, "agent.exe")

	if exePath != targetPath {

		srcFile, err := os.Open(exePath)
		if err != nil {
			return fmt.Errorf("failed to open source executable: %w", err)
		}
		defer srcFile.Close()

		dstFile, err := os.Create(targetPath)
		if err != nil {
			return fmt.Errorf("failed to create destination executable: %w", err)
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return fmt.Errorf("failed to copy executable: %w", err)
		}

		err = dstFile.Sync()
		if err != nil {
			return fmt.Errorf("failed to sync destination file: %w", err)
		}
	}

	k, err := registry.OpenKey(registry.CURRENT_USER, registryKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer k.Close()

	err = k.SetStringValue(registryValueName, targetPath)
	if err != nil {
		return fmt.Errorf("failed to set registry value: %w", err)
	}

	return nil
}

func uninstall() error {

	k, err := registry.OpenKey(registry.CURRENT_USER, registryKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer k.Close()

	err = k.DeleteValue(registryValueName)
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("failed to delete registry value: %w", err)
	}

	appDataDir := os.Getenv("APPDATA")
	if appDataDir != "" {
		targetPath := filepath.Join(appDataDir, "Evolord", "agent.exe")

		os.Remove(targetPath)
	}

	return nil
}
