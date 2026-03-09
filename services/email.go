package services

import (
	"fmt"

	"jarvis-agent/config"

	"github.com/go-mail/mail/v2"
)

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{cfg: cfg}
}

func (s *EmailService) SendReport(reportContent string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.cfg.EmailFrom)
	m.SetHeader("To", s.cfg.EmailTo)
	m.SetHeader("Subject", s.cfg.EmailSubject)
	m.SetBody("text/html", formatEmailBody(reportContent))

	d := mail.NewDialer(
		s.cfg.SMTPHost,
		s.cfg.SMTPPort,
		s.cfg.SMTPUser,
		s.cfg.SMTPPass,
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("erro ao enviar email: %w", err)
	}

	return nil
}

func formatEmailBody(content string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 800px; margin: 0 auto; padding: 20px; }
        .header { background: #1a73e8; color: white; padding: 20px; border-radius: 8px 8px 0 0; }
        .content { background: #f5f5f5; padding: 20px; border-radius: 0 0 8px 8px; }
        pre { white-space: pre-wrap; word-wrap: break-word; background: white; padding: 15px; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>🤖 Jarvis Report</h2>
            <p>Seu relatório diário está pronto!</p>
        </div>
        <div class="content">
            <pre>%s</pre>
        </div>
    </div>
</body>
</html>
`, content)
}
