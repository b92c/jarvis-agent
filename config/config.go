package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	GoogleAPIKey     string
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

	notesPath := os.Getenv("TODAY_NOTES_PATH")
	if notesPath == "" {
		notesPath = defaultNotesPath
	}

	reportPath := os.Getenv("REPORT_OUTPUT_PATH")
	if reportPath == "" {
		reportPath = defaultReportPath
	}

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
