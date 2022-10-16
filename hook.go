package logrus2logtail

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
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

	var p map[string]interface{} // Decode into Struct

	err := json.Unmarshal(buf, &p)
	if err != nil {
		log.Error("Unmarshall:")
	}

	b, err := msgpack.Marshal(p)
	if err != nil {
		log.Error("MsgPack:", err)
	}

	req, err := http.NewRequest("POST", "https://in.logtail.com", bytes.NewReader(b))
	if err != nil {
		log.Error("Error NewRequest:", err)
	}
	req.Header.Set("Content-Type", "application/msgpack")
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
