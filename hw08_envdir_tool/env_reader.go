package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrDirectoryNotExists = errors.New("directory not exists")
	ErrPathIsNotDir       = errors.New("path is not a directory")
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrDirectoryNotExists
		}
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	if !info.IsDir() {
		return nil, ErrPathIsNotDir
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	env := make(Environment)

	for _, entry := range entries {
		envVal, err := makeEnvValue(dir, entry)
		if err != nil {
			return nil, fmt.Errorf("failed to handle entry %s: %w", entry.Name(), err)
		}
		if envVal == nil {
			continue
		}

		env[entry.Name()] = *envVal
	}

	return env, nil
}

func makeEnvValue(dir string, entry fs.DirEntry) (*EnvValue, error) {
	if entry.IsDir() {
		return nil, nil
	}

	if strings.Contains(entry.Name(), "=") {
		return nil, fmt.Errorf("invalid file name %s", entry.Name())
	}

	filePath := filepath.Join(dir, entry.Name())

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	//если файл полностью пустой (длина - 0 байт), то envdir удаляет переменную окружения с именем S.
	if len(content) == 0 {
		return &EnvValue{Value: "", NeedRemove: true}, nil
	}

	return &EnvValue{Value: trimContent(content), NeedRemove: false}, nil
}

// Пробелы и табуляция в конце T удаляются; терминальные нули (0x00) заменяются на перевод строки (\n);
func trimContent(content []byte) string {
	trimmedContent := bytes.TrimRight(content, " \t")

	if newlineIndex := bytes.Index(trimmedContent, []byte("\n")); newlineIndex != -1 {
		trimmedContent = trimmedContent[:newlineIndex]
	}

	finalContent := bytes.ReplaceAll(trimmedContent, []byte{0x00}, []byte("\n"))

	// Проверяем случай, когда остается только одиночный пробел
	if len(trimmedContent) == 1 && trimmedContent[0] == ' ' {
		return ""
	}

	return string(finalContent)
}
