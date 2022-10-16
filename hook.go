package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	JSONFormatter = &log.JSONFormatter{}
)

type Hooker struct {
	Token     string
	formatter log.Formatter
}

func NewHook(Token string) (*Hooker, error) {
	if Token == "" {
		return nil, errors.New("enter bearer token logtail")
	}
	return &Hooker{
		Token:     Token,
		formatter: JSONFormatter,
	}, nil
}

func (h *Hooker) Fire(entry *log.Entry) error {

	buf, _ := h.formatter.Format(entry)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	body := bytes.NewReader(buf)

	req, err := http.NewRequest("POST", "https://in.logtail.com", body)
	if err != nil {
		log.Error("Error NewRequest:", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.Token)

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error Do:", err)
	}
	defer resp.Body.Close()

	return nil
}

func (h *Hooker) Levels() []log.Level {
	return log.AllLevels
}
