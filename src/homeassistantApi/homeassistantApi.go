package homeassistantApi

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"solar-monitor/variables"
	"strings"
)

func CallHomeAssisatantApiSet(apiCall string, data *strings.Reader) error {
	var errorReturn error = nil
	client := &http.Client{}
	url := variables.UrlHomeAssistant
	token := variables.TokenHomeAssistant

	req, err := http.NewRequest("POST", url+apiCall, data)
	if err != nil {
		errorReturn = err
		slog.Error("Failed to create new http request", "ERROR", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		errorReturn = err
		slog.Error("API request failed", "ERROR", err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	_ = bodyText
	if err != nil {
		errorReturn = err
		slog.Error("Failed to get API response", "ERROR", err)
	} else {
		if resp.StatusCode != 200 {
			errorReturn = errors.New("API-Request failed: " + resp.Status)
			slog.Error("API-Request failed", "Status", resp.Status, "Request-URL", req.URL, "Request-Method", req.Method)
		} else {
			if variables.HomeAssistantLogApiCallResult200 {
				slog.Info("API-Request successful", "Status", resp.Status, "Request-URL", req.URL, "Request-Method", req.Method)
			}
		}
	}
	//fmt.Printf("%s\n", bodyText)
	return errorReturn
}
