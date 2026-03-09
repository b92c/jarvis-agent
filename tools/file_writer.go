package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"jarvis-agent/config"
)

type FileWriter struct {
	cfg *config.Config
}

func NewFileWriter(cfg *config.Config) *FileWriter {
	return &FileWriter{cfg: cfg}
}

func (w *FileWriter) WriteReport(content string) error {
	dir := filepath.Dir(w.cfg.ReportOutputPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório %s: %w", dir, err)
	}

	if err := os.WriteFile(w.cfg.ReportOutputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("erro ao escrever relatório %s: %w", w.cfg.ReportOutputPath, err)
	}

	return nil
}

func (w *FileWriter) GetReportPath() string {
	return w.cfg.ReportOutputPath
}
