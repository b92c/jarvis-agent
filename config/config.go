package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	GoogleAPIKey     string
	GeminiModel      string
	TodayNotesPath   string
	ReportOutputPath string
	SMTPHost         string
	SMTPPort         int
	SMTPUser         string
	SMTPPass         string
	EmailFrom        string
	EmailTo          string
	EmailSubject     string
}

func Load() (*Config, error) {
	// Carrega variáveis do arquivo .env (ignora erro caso não exista)
	_ = godotenv.Load()

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter home directory: %w", err)
	}

	defaultNotesPath := filepath.Join(home, "Documents/personal-notes/Trabalho/Gestão/Notes/today.md")
	defaultReportPath := filepath.Join(home, "Documents/personal-notes/Trabalho/Gestão/Notes/jarvis-report.md")

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_API_KEY não configurada")
	}

	geminiModel := os.Getenv("GEMINI_MODEL")
	if geminiModel == "" {
		geminiModel = "gemini-2.5-flash"
	}

	notesPath := os.Getenv("TODAY_NOTES_PATH")
	if notesPath == "" {
		notesPath = defaultNotesPath
	}
	notesPath = expandPath(notesPath, home)

	reportPath := os.Getenv("REPORT_OUTPUT_PATH")
	if reportPath == "" {
		reportPath = defaultReportPath
	}
	reportPath = expandPath(reportPath, home)

	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com"
	}

	smtpPort := 587
	if port := os.Getenv("SMTP_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &smtpPort)
	}

	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	emailFrom := os.Getenv("EMAIL_FROM")
	emailTo := os.Getenv("EMAIL_TO")
	emailSubject := os.Getenv("EMAIL_SUBJECT")

	if emailSubject == "" {
		emailSubject = "[JARVIS] Seu report está pronto senhor"
	}

	return &Config{
		GoogleAPIKey:     apiKey,
		GeminiModel:      geminiModel,
		TodayNotesPath:   notesPath,
		ReportOutputPath: reportPath,
		SMTPHost:         smtpHost,
		SMTPPort:         smtpPort,
		SMTPUser:         smtpUser,
		SMTPPass:         smtpPass,
		EmailFrom:        emailFrom,
		EmailTo:          emailTo,
		EmailSubject:     emailSubject,
	}, nil
}

// expandPath resolve caminhos com ~ e caminhos relativos
func expandPath(path string, home string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}
	if !filepath.IsAbs(path) {
		return filepath.Join(home, path)
	}
	return path
}
