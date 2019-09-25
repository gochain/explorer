package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/gochain-io/explorer/server/models"
	"go.uber.org/zap"
)

const RECAPTCHA_URL = "https://www.google.com/recaptcha/api/siteverify"

func verifyReCaptcha(token, reCaptchaSecret, action, remoteIP string) error {
	if reCaptchaSecret == "" {
		return nil
	}
	params := url.Values{}
	params.Add("secret", reCaptchaSecret)
	params.Add("response", token)
	if remoteIP != "" {
		params.Add("remoteip", remoteIP)
	}
	resp, err := http.PostForm(RECAPTCHA_URL, params)
	if err != nil {
		logger.Warn("error occurred during making recaptcha request", zap.Error(err))
		err := errors.New("error occurred during processing your request. please try again")
		return err
	}
	var result *models.ReCaptchaResponse
	json.NewDecoder(resp.Body).Decode(&result)
	// resp.Body.Close()
	if result.Success == false {
		err := errors.New("error occurred during anti-bot checking. please try again")
		return err
	}
	if result.Score < 0.5 {
		err := errors.New("not handling bot request")
		return err
	}
	return nil
}
