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
</head>
<body style="margin:0; padding:0; background-color:#eef2f7; font-family:'Segoe UI',Arial,sans-serif; color:#333333;">
    <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color:#eef2f7;">
        <tr>
            <td align="center" style="padding:20px 0;">
                <table width="800" cellpadding="0" cellspacing="0" border="0" style="max-width:800px; border-radius:10px; overflow:hidden;">
                    <!-- HEADER -->
                    <tr>
                        <td bgcolor="#1a73e8" style="background-color:#1a73e8; padding:24px 30px; color:#ffffff;">
                            <h2 style="margin:0 0 4px 0; font-size:22px; color:#ffffff; font-family:'Segoe UI',Arial,sans-serif;">🤖 Jarvis Report</h2>
                            <p style="margin:0; font-size:14px; color:#ffffff; font-family:'Segoe UI',Arial,sans-serif;">Seu relatório diário está pronto!</p>
                        </td>
                    </tr>
                    <!-- CONTENT -->
                    <tr>
                        <td bgcolor="#ffffff" style="background-color:#ffffff; padding:24px 30px; color:#333333; font-family:'Segoe UI',Arial,sans-serif; font-size:15px; line-height:1.6;">
                            %s
                        </td>
                    </tr>
                    <!-- FOOTER -->
                    <tr>
                        <td bgcolor="#f5f5f5" style="background-color:#f5f5f5; padding:14px 30px; text-align:center; font-size:12px; color:#999999; border-top:1px solid #e0e0e0; font-family:'Segoe UI',Arial,sans-serif;">
                            Gerado automaticamente por Jarvis Agent
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, string(htmlContent))
}
