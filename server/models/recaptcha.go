package models

import "time"

type ReCaptchaRequest struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
	RemoteIp string `json:"remoteip"`
}

type ReCaptchaResponse struct {
	Success     bool      `json:"success"`
	HostName    string    `json:"hostname"`
	Score       float32   `json:"score"`
	Action      string    `json:"action"`
	ErrorCodes  []string  `json:"error-codes"`
	ChallengeTs time.Time `json:"challenge_ts"`
}
