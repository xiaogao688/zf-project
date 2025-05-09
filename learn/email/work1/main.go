package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
按照磁盘梯度发送邮件告警
*/

// SMTPConfig holds SMTP server settings and email recipients.
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

// Send sends an email with the provided subject and body.
func (c *SMTPConfig) Send(subject, body string) error {
	auth := smtp.PlainAuth("", c.Username, c.Password, c.Host)
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	// Build the message
	header := make(map[string]string)
	header["From"] = c.From
	header["To"] = strings.Join(c.To, ",")
	header["Subject"] = subject

	var msg strings.Builder
	for k, v := range header {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n" + body)

	return smtp.SendMail(addr, auth, c.From, c.To, []byte(msg.String()))
}

// DiskAlert handles threshold-based email alerts for disk usage.
type DiskAlert struct {
	thresholds     []int
	lastAlertTimes map[int]time.Time
	mu             sync.Mutex
	smtpCfg        *SMTPConfig
}

// NewDiskAlert initializes a DiskAlert with the given SMTP settings.
func NewDiskAlert(cfg *SMTPConfig) *DiskAlert {
	// Define thresholds: 80%, then every 5% up to 100%
	ths := []int{}
	for t := 80; t <= 100; t += 5 {
		ths = append(ths, t)
	}
	return &DiskAlert{
		thresholds:     ths,
		lastAlertTimes: make(map[int]time.Time),
		smtpCfg:        cfg,
	}
}

// Check evaluates the current usage against thresholds and sends alerts.
func (da *DiskAlert) Check(usage float64) error {
	da.mu.Lock()
	defer da.mu.Unlock()
	now := time.Now()
	for _, t := range da.thresholds {
		if usage >= float64(t) {
			last := da.lastAlertTimes[t]
			// Alert if never sent or more than an hour has passed
			if last.IsZero() || now.Sub(last) >= time.Hour {
				subj := fmt.Sprintf("Disk usage warning: reached %d%%", t)
				body := fmt.Sprintf("Disk usage is at %.2f%%, exceeding the %d%% threshold.", usage, t)
				if err := da.smtpCfg.Send(subj, body); err != nil {
					return err
				}
				da.lastAlertTimes[t] = now
			}
		}
	}
	return nil
}

// requestPayload is the expected JSON structure for incoming requests.
type requestPayload struct {
	Usage float64 `json:"usage"`
}

// startServer launches an HTTP API at the given address.
func startServer(addr string, da *DiskAlert) {
	http.HandleFunc("/usage", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		var req requestPayload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request: invalid JSON", http.StatusBadRequest)
			return
		}
		go func(u float64) {
			if err := da.Check(u); err != nil {
				log.Printf("Error sending alert: %v", err)
			}
		}(req.Usage)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func main() {
	// Load SMTP settings from environment variables
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USERNAME")
	pass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")
	toList := os.Getenv("SMTP_TO") // comma-separated

	if host == "" || portStr == "" || user == "" || pass == "" || from == "" || toList == "" {
		log.Fatal("Missing one or more SMTP configuration environment variables")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT: %v", err)
	}

	cfg := &SMTPConfig{
		Host:     host,
		Port:     port,
		Username: user,
		Password: pass,
		From:     from,
		To:       strings.Split(toList, ","),
	}

	da := NewDiskAlert(cfg)
	// Start HTTP server on port 8080 (or set via env)
	addr := ":8080"
	if a := os.Getenv("ALERT_SERVER_ADDR"); a != "" {
		addr = a
	}
	startServer(addr, da)
}
