package utils

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

func ReadJsonFromFile[T any](relativePath string) (T, error) {
	var jsonObj T

	filePath := GetFullPath(relativePath)

	// Open the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("error while opening file", "filename", filePath, "error", err)
		return jsonObj, err
	}
	defer file.Close()

	// Decode the JSON data
	err = json.NewDecoder(file).Decode(&jsonObj)
	if err != nil {
		slog.Error("error decoding json from file", "filename", relativePath, "error", err)
		return jsonObj, err
	}
	return jsonObj, nil
}

func WriteJsonToFile[T any](relativePath string, data T) error {
	// Marshal the data into a JSON byte slice
	
	filePath := GetFullPath(relativePath)
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		slog.Error("error while marshalling json", "filename", filePath, "error", err)
		return err
	}

	// Open the file for writing, overwriting any existing content
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if errors.Is(err, syscall.ENOTDIR) {
		file, err := os.Create(filePath)
        if err != nil {
            return err
        }
        defer file.Close()
	}
	if err != nil {
		slog.Error("error while opening file", "filename", filePath, "error", err)
		return err
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		slog.Error("failed to write JSON to file", "filename", filePath, "error", err)
		return err
	}

	return nil
}

func GetListOfFilesInDirectory(relativePath string) []string {
	var files []string

	fileInfo, err := os.ReadDir(GetFullPath(relativePath))
	if err != nil {
		slog.Error("error while reading directory", "error", err)
	}

	for _, file := range fileInfo {
		if !file.IsDir() {
			slog.Debug("file found", "name", file.Name())
			files = append(files, file.Name())
		}
	}

	return files
}

func CheckPathExists(path string) bool {
	_, err := os.Stat(GetFullPath(path))
	return err == nil || !os.IsNotExist(err)
}

func GetFullPath(relativePath string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		slog.Error("Unable to retrieve caller information")
		return ""
	}
	srcDir := filepath.Dir((filepath.Dir(filename)))
	// slog.Info("Detail", "srcDir", srcDir)
	return filepath.Join(srcDir, relativePath)
}
