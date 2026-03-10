package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	GoogleAPIKey      string
	GeminiModel       string
	TodayNotesPath    string
	ReportOutputPath  string
	MonthlyNotesPath  string
	MonthlyReportPath string
	SMTPHost          string
	SMTPPort          int
	SMTPUser          string
	SMTPPass          string
	EmailFrom         string
	EmailTo           string
	EmailSubject      string
	StateFilePath     string
}

func Load() (*Config, error) {
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

	monthlyNotesPath := os.Getenv("MONTHLY_NOTES_PATH")
	if monthlyNotesPath == "" {
		monthlyNotesPath = filepath.Join(home, "Documents/personal-notes/Trabalho/Gestão/Notes")
	}
	monthlyNotesPath = expandPath(monthlyNotesPath, home)

	monthlyReportPath := os.Getenv("MONTHLY_REPORT_OUTPUT_PATH")
	if monthlyReportPath == "" {
		monthlyReportPath = filepath.Join(home, "Documents/personal-notes/Trabalho/Gestão/Notes/jarvis-report.md")
	}
	monthlyReportPath = expandPath(monthlyReportPath, home)

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

	stateFilePath := os.Getenv("STATE_FILE_PATH")
	if stateFilePath == "" {
		stateFilePath = filepath.Join(home, ".jarvis-agent", "memory-ticker.json")
	}

	return &Config{
		GoogleAPIKey:      apiKey,
		GeminiModel:       geminiModel,
		TodayNotesPath:    notesPath,
		ReportOutputPath:  reportPath,
		MonthlyNotesPath:  monthlyNotesPath,
		MonthlyReportPath: monthlyReportPath,
		SMTPHost:          smtpHost,
		SMTPPort:          smtpPort,
		SMTPUser:          smtpUser,
		SMTPPass:          smtpPass,
		EmailFrom:         emailFrom,
		EmailTo:           emailTo,
		EmailSubject:      emailSubject,
		StateFilePath:     stateFilePath,
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

func (c *Config) GetMonthlyNotesPath(year, month string) string {
	return filepath.Join(c.MonthlyNotesPath, year, month, "notes.md")
}

func (c *Config) GetMonthlyReportPath(year, month string) string {
	dir := filepath.Dir(c.MonthlyReportPath)
	ext := filepath.Ext(c.MonthlyReportPath)
	base := strings.TrimSuffix(filepath.Base(c.MonthlyReportPath), ext)
	return filepath.Join(dir, fmt.Sprintf("%s-%s-%s%s", base, year, month, ext))
}
