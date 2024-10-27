package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/TejasThombare20/weather-engine/repositories"
	"gopkg.in/gomail.v2"
)

type AlertService struct {
	emailDialer *gomail.Dialer
	emailFrom   string
	alertRepo   *repositories.AlertRepository
}

func NewAlertService(alertRepo *repositories.AlertRepository, smtpHost string, smtpPort int, smtpUser, smtpPass, emailFrom string) *AlertService {
	return &AlertService{
		emailDialer: gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass),
		emailFrom:   emailFrom,
		alertRepo:   alertRepo,
	}
}

func (s *AlertService) ProcessAlerts() error {

	rows, err := s.alertRepo.QyeryBreaches_view()

	if err != nil {
		log.Println("error getting rows from query breatches view ")
	}

	for rows.Next() {
		var alert struct {
			ID          int
			UserID      int
			CityName    string
			Temperature float64
			Threshold   float64
			Email       string
		}

		if err := rows.Scan(
			&alert.ID, &alert.UserID, &alert.CityName,
			&alert.Temperature, &alert.Threshold, &alert.Email,
		); err != nil {
			log.Printf("Error scanning alert: %v", err)
			continue
		}

		// Send email alert
		if err := s.sendAlertEmail(alert); err != nil {
			log.Printf("Error sending alert email: %v", err)
			continue
		}

		// Update alert as sent
		err = s.alertRepo.UpdateTemeratureAlert(alert.ID)

		if err != nil {
			log.Printf("Error updating alert status: %v", err)
		}
	}

	return nil
}

func (s *AlertService) sendAlertEmail(alert struct {
	ID          int
	UserID      int
	CityName    string
	Temperature float64
	Threshold   float64
	Email       string
}) error {
	const emailTemplate = `
        <h2>Weather Alert for {{.CityName}}</h2>
        <p>The temperature has exceeded your configured threshold:</p>
        <ul>
            <li>Current Temperature: {{printf "%.1f" .Temperature}}°C</li>
            <li>Your Threshold: {{printf "%.1f" .Threshold}}°C</li>
        </ul>
        <p>This condition has been observed in consecutive measurements.</p>
    `

	tmpl, err := template.New("alert").Parse(emailTemplate)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, alert); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.emailFrom)
	m.SetHeader("To", alert.Email)
	m.SetHeader("Subject", fmt.Sprintf("Temperature Alert for %s", alert.CityName))
	m.SetBody("text/html", body.String())

	return s.emailDialer.DialAndSend(m)
}

func (s *AlertService) StartAlertProcessing() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			if err := s.ProcessAlerts(); err != nil {
				log.Printf("Error processing alerts: %v", err)
			}
		}
	}()
}
