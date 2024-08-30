package utils

import "encoding/json"

func JsinStatus(message string) []byte {
	m, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	return m
}