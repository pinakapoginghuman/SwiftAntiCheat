package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type uploadPayload struct {
	Results  ScanReport `json:"results"`
	HWIDHash string     `json:"hwidHash"`
}

func UploadResults(apiURL string, report ScanReport) error {
	payload := uploadPayload{
		Results:  report,
		HWIDHash: report.HWIDHash,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	url := fmt.Sprintf("%s/api/scans/%s/results", apiURL, report.ScanID)
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send results: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	return nil
}
