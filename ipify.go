package main

import (
	"fmt"
	"io"
	"net/http"
)

func getMyIP(ipifyURL string) (string, error) {
	sugar.Debugw("calling ipify API", "url", ipifyURL)
	res, err := http.Get(ipifyURL)
	if err != nil {
		return "", fmt.Errorf("error calling ipify API: %w", err)
	}
	ip, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading ipify response: %w", err)
	}
	sugar.Debugw("successfully got ip from ipify API", "ip", string(ip))
	return string(ip), nil
}
