package services

import (
	"fmt"

	"jarvis-agent/config"

	"github.com/go-mail/mail/v2"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
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
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	opts := html.RendererOptions{
		Flags: html.CommonFlags | html.HrefTargetBlank,
	}
	renderer := html.NewRenderer(opts)

	htmlContent := markdown.ToHTML([]byte(content), p, renderer)

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: 'Segoe UI', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            margin: 0;
            padding: 0;
            background-color: #eef2f7;
        }
        .container {
            max-width: 800px;
            margin: 20px auto;
            background: #ffffff;
            border-radius: 10px;
            overflow: hidden;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        .header {
            background: linear-gradient(135deg, #1a73e8, #0d47a1);
            color: white;
            padding: 24px 30px;
        }
        .header h2 {
            margin: 0 0 4px 0;
            font-size: 22px;
        }
        .header p {
            margin: 0;
            opacity: 0.9;
            font-size: 14px;
        }
        .content {
            padding: 24px 30px;
        }
        .content h1 {
            color: #1a73e8;
            font-size: 22px;
            border-bottom: 2px solid #e0e0e0;
            padding-bottom: 8px;
            margin-top: 0;
        }
        .content h2 {
            color: #333;
            font-size: 18px;
            margin-top: 24px;
            margin-bottom: 12px;
            border-bottom: 1px solid #eee;
            padding-bottom: 6px;
        }
        .content ul {
            padding-left: 20px;
        }
        .content li {
            margin-bottom: 6px;
        }
        .content p {
            margin: 8px 0;
        }
        .content strong {
            color: #1a73e8;
        }
        .content code {
            background: #f0f0f0;
            padding: 2px 6px;
            border-radius: 3px;
            font-size: 13px;
        }
        .content pre {
            background: #f8f8f8;
            padding: 14px;
            border-radius: 6px;
            overflow-x: auto;
            border: 1px solid #e0e0e0;
        }
        .footer {
            background: #f5f5f5;
            padding: 14px 30px;
            text-align: center;
            font-size: 12px;
            color: #999;
            border-top: 1px solid #e0e0e0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>🤖 Jarvis Report</h2>
            <p>Seu relatório diário está pronto!</p>
        </div>
        <div class="content">
            %s
        </div>
        <div class="footer">
            Gerado automaticamente por Jarvis Agent
        </div>
    </div>
</body>
</html>
`, string(htmlContent))
}
