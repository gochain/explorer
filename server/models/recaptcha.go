package models

type ReCaptchaRequest struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
}

type ReCaptchaResponse struct {
	Success    bool     `json:"success"`
	HostName   string   `json:"hostname"`
	Score      float32  `json:"score"`
	Action     string   `json:"action"`
	ErrorCodes []string `json:"error-codes"`
	//challenge_ts time
}
