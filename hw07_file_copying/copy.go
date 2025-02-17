package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFromPathEqualToPath   = errors.New("from path equal to path")
)

func isSystemFile(fromPath string, fileInfo os.FileInfo) error {
	// Проверяем, является ли файл устройством (блочным или символьным), FIFO или символической ссылкой
	if fileInfo.Mode()&(os.ModeDevice|os.ModeNamedPipe|os.ModeSymlink) != 0 {
		return ErrUnsupportedFile
	}

	// Проверяем, не находится ли файл в запрещенных системных директориях
	absPath, err := filepath.Abs(fromPath)
	if err != nil {
		return err
	}

	forbiddenPaths := []string{"/dev/", "/proc/", "/sys/", "/run/", "/var/run/"}
	for _, forbidden := range forbiddenPaths {
		if strings.HasPrefix(absPath, forbidden) {
			return ErrUnsupportedFile
		}
	}

	// Проверяем права доступа к файлу
	if fileInfo.Mode().Perm()&(os.ModeSetuid|os.ModeSetgid) != 0 {
		return ErrUnsupportedFile
	}

	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == toPath {
		return ErrFromPathEqualToPath
	}

	file, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	if err = isSystemFile(fromPath, fileInfo); err != nil {
		return fmt.Errorf("%w: %s", ErrUnsupportedFile, fromPath)
	}

	fileSize := fileInfo.Size()
	if fileSize < offset {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || offset+limit > fileSize {
		limit = fileSize - offset
	}

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error setting file offset: %w", err)
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer dstFile.Close()

	progressBar := pb.Start64(limit)
	defer progressBar.Finish()

	reader := progressBar.NewProxyReader(io.LimitReader(file, limit))

	_, err = io.CopyN(dstFile, reader, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("error copying data: %w", err)
	}

	fmt.Printf("Copied %d bytes from %s to %s\n", limit, fromPath, toPath)

	return nil
}
