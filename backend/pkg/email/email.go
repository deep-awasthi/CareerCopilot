package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"github.com/deepawasthi/careercopilot/pkg/config"
	"github.com/deepawasthi/careercopilot/pkg/logger"
	"go.uber.org/zap"
)

type Client struct {
	cfg *config.SMTPConfig
}

func NewClient(cfg *config.SMTPConfig) *Client {
	return &Client{cfg: cfg}
}

type EmailMessage struct {
	To      string
	Subject string
	Body    string
	IsHTML  bool
}

func (c *Client) Send(msg *EmailMessage) error {
	if c.cfg.Username == "" {
		logger.Warn("SMTP not configured, skipping email", zap.String("to", msg.To))
		return nil
	}

	auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)

	from := fmt.Sprintf("%s <%s>", c.cfg.FromName, c.cfg.FromEmail)
	contentType := "text/plain"
	if msg.IsHTML {
		contentType = "text/html"
	}

	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: %s; charset=UTF-8\r\n\r\n",
		from, msg.To, msg.Subject, contentType)

	body := headers + msg.Body
	addr := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)

	err := smtp.SendMail(addr, auth, c.cfg.FromEmail, []string{msg.To}, []byte(body))
	if err != nil {
		logger.Error("failed to send email", zap.String("to", msg.To), zap.Error(err))
		return fmt.Errorf("smtp send error: %w", err)
	}

	logger.Info("email sent", zap.String("to", msg.To), zap.String("subject", msg.Subject))
	return nil
}

// DailyDigestData holds data for the daily digest email
type DailyDigestData struct {
	UserName          string
	NewJobsCount      int
	KeywordMatches    []KeywordMatch
	CompanyOpenings   []CompanyOpening
	ReferralCount     int
	UpcomingInterviews []UpcomingInterview
	SavedJobsCount    int
}

type KeywordMatch struct {
	Keyword string
	Count   int
}

type CompanyOpening struct {
	Company  string
	JobCount int
}

type UpcomingInterview struct {
	Company     string
	Role        string
	Stage       string
	ScheduledAt string
}

var dailyDigestTemplate = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
  body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #f8fafc; margin: 0; padding: 20px; }
  .container { max-width: 600px; margin: 0 auto; background: white; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 20px rgba(0,0,0,0.08); }
  .header { background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%); padding: 32px; color: white; text-align: center; }
  .header h1 { margin: 0; font-size: 28px; font-weight: 700; }
  .header p { margin: 8px 0 0; opacity: 0.85; font-size: 14px; }
  .content { padding: 32px; }
  .stat-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 16px; margin-bottom: 32px; }
  .stat-card { background: #f8fafc; border-radius: 8px; padding: 16px; text-align: center; border: 1px solid #e2e8f0; }
  .stat-card .number { font-size: 32px; font-weight: 700; color: #6366f1; }
  .stat-card .label { font-size: 12px; color: #64748b; margin-top: 4px; text-transform: uppercase; letter-spacing: 0.05em; }
  .section { margin-bottom: 24px; }
  .section h3 { font-size: 16px; color: #1e293b; font-weight: 600; margin-bottom: 12px; padding-bottom: 8px; border-bottom: 2px solid #f1f5f9; }
  .item { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #f1f5f9; font-size: 14px; color: #475569; }
  .item strong { color: #1e293b; }
  .badge { background: #ede9fe; color: #7c3aed; padding: 2px 8px; border-radius: 12px; font-size: 12px; font-weight: 600; }
  .footer { background: #f8fafc; padding: 20px 32px; text-align: center; font-size: 12px; color: #94a3b8; border-top: 1px solid #e2e8f0; }
  .cta { display: block; text-align: center; margin: 24px 0; background: #6366f1; color: white; padding: 14px 28px; border-radius: 8px; text-decoration: none; font-weight: 600; }
</style>
</head>
<body>
<div class="container">
  <div class="header">
    <h1>🚀 CareerCopilot Daily Report</h1>
    <p>Good morning, {{.UserName}}! Here's your career update.</p>
  </div>
  <div class="content">
    <div class="stat-grid">
      <div class="stat-card">
        <div class="number">{{.NewJobsCount}}</div>
        <div class="label">New Jobs Today</div>
      </div>
      <div class="stat-card">
        <div class="number">{{.ReferralCount}}</div>
        <div class="label">Referral Opportunities</div>
      </div>
      <div class="stat-card">
        <div class="number">{{len .UpcomingInterviews}}</div>
        <div class="label">Upcoming Interviews</div>
      </div>
      <div class="stat-card">
        <div class="number">{{.SavedJobsCount}}</div>
        <div class="label">Saved — Apply Now</div>
      </div>
    </div>

    {{if .KeywordMatches}}
    <div class="section">
      <h3>🔔 Keyword Matches</h3>
      {{range .KeywordMatches}}
      <div class="item">
        <strong>{{.Keyword}}</strong>
        <span class="badge">{{.Count}} jobs</span>
      </div>
      {{end}}
    </div>
    {{end}}

    {{if .CompanyOpenings}}
    <div class="section">
      <h3>🏢 Company Openings</h3>
      {{range .CompanyOpenings}}
      <div class="item">
        <strong>{{.Company}}</strong>
        <span class="badge">{{.JobCount}} new</span>
      </div>
      {{end}}
    </div>
    {{end}}

    {{if .UpcomingInterviews}}
    <div class="section">
      <h3>📅 Upcoming Interviews</h3>
      {{range .UpcomingInterviews}}
      <div class="item">
        <div><strong>{{.Company}}</strong> — {{.Role}}</div>
        <div><span class="badge">{{.Stage}}</span> {{.ScheduledAt}}</div>
      </div>
      {{end}}
    </div>
    {{end}}

    <a href="https://careercopilot.io" class="cta">Open CareerCopilot Dashboard →</a>
  </div>
  <div class="footer">
    CareerCopilot · Track jobs. Find referrals. Grow your career.<br>
    <a href="https://careercopilot.io/settings/notifications" style="color: #6366f1;">Manage notifications</a>
  </div>
</div>
</body>
</html>`

func (c *Client) SendDailyDigest(to string, data *DailyDigestData) error {
	tmpl, err := template.New("digest").Parse(dailyDigestTemplate)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}

	kwParts := make([]string, len(data.KeywordMatches))
	for i, km := range data.KeywordMatches {
		kwParts[i] = fmt.Sprintf("%d %s", km.Count, km.Keyword)
	}

	subject := fmt.Sprintf("CareerCopilot Daily Report · %d New Jobs", data.NewJobsCount)
	if len(kwParts) > 0 {
		subject += " · " + strings.Join(kwParts[:min(2, len(kwParts))], ", ")
	}

	return c.Send(&EmailMessage{
		To:      to,
		Subject: subject,
		Body:    buf.String(),
		IsHTML:  true,
	})
}

func (c *Client) SendKeywordAlert(to, keyword string, jobCount int) error {
	body := fmt.Sprintf(`<h2>🔔 %d new jobs match your keyword: <strong>%s</strong></h2>
<p>Visit <a href="https://careercopilot.io/jobs?q=%s">CareerCopilot</a> to view these jobs.</p>`, jobCount, keyword, keyword)

	return c.Send(&EmailMessage{
		To:      to,
		Subject: fmt.Sprintf("CareerCopilot: %d new jobs match '%s'", jobCount, keyword),
		Body:    body,
		IsHTML:  true,
	})
}

var passwordResetTemplate = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
  body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #f8fafc; margin: 0; padding: 20px; }
  .container { max-width: 576px; margin: 40px auto; background: white; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 20px rgba(0,0,0,0.08); border: 1px solid #e2e8f0; }
  .header { background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%); padding: 32px; color: white; text-align: center; }
  .header h1 { margin: 0; font-size: 24px; font-weight: 700; }
  .content { padding: 32px; text-align: center; }
  .content p { font-size: 16px; color: #475569; line-height: 1.5; margin-bottom: 24px; }
  .cta { display: inline-block; background: #6366f1; color: white !important; padding: 14px 28px; border-radius: 8px; text-decoration: none; font-weight: 600; font-size: 16px; }
  .footer { background: #f8fafc; padding: 20px 32px; text-align: center; font-size: 12px; color: #94a3b8; border-top: 1px solid #e2e8f0; }
</style>
</head>
<body>
<div class="container">
  <div class="header">
    <h1>🔑 Reset Your Password</h1>
  </div>
  <div class="content">
    <p>We received a request to reset your password for your CareerCopilot account. Click the button below to proceed.</p>
    <a href="{{.ResetLink}}" class="cta">Reset Password</a>
    <p style="margin-top: 24px; font-size: 12px; color: #94a3b8;">If you did not request this reset, you can safely ignore this email. This link will expire in 1 hour.</p>
  </div>
  <div class="footer">
    CareerCopilot · Track jobs. Find referrals. Grow your career.
  </div>
</div>
</body>
</html>`

func (c *Client) SendPasswordReset(to, resetLink string) error {
	tmpl, err := template.New("reset").Parse(passwordResetTemplate)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	var buf bytes.Buffer
	data := struct{ ResetLink string }{ResetLink: resetLink}
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}

	return c.Send(&EmailMessage{
		To:      to,
		Subject: "CareerCopilot: Reset Your Password",
		Body:    buf.String(),
		IsHTML:  true,
	})
}

func min(a, b int) int {
	if a < b { return a }
	return b
}

