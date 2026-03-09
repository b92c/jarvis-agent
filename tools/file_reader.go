package tools

import (
	"fmt"
	"os"

	"jarvis-agent/config"
)

type FileReader struct {
	cfg *config.Config
}

func NewFileReader(cfg *config.Config) *FileReader {
	return &FileReader{cfg: cfg}
}

func (r *FileReader) ReadTodayNotes() (string, error) {
	content, err := os.ReadFile(r.cfg.TodayNotesPath)
	if err != nil {
		return "", fmt.Errorf("erro ao ler arquivo %s: %w", r.cfg.TodayNotesPath, err)
	}
	return string(content), nil
}

func (r *FileReader) ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("erro ao ler arquivo %s: %w", path, err)
	}
	return string(content), nil
}

func (r *FileReader) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
