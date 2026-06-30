package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type uploadPayload struct {
	ReportCode string     `json:"reportCode"`
	Results    ScanReport `json:"results"`
	HWIDHash   string     `json:"hwidHash"`
}

func UploadResults(apiURL string, reportCode string, report ScanReport) error {
	payload := uploadPayload{
		ReportCode: reportCode,
		Results:    report,
		HWIDHash:   report.HWIDHash,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	url := fmt.Sprintf("%s/api/reports/upload", apiURL)
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
