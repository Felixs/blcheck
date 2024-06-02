package url

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
)

var (
	csvHeader = []string{
		"url", "is_reachable", "status_message", "content_length", "response_time", "num_occured",
	}
)

// Helper construct of UrlReport to customize the JSON conversion
type JsonUrlReport struct {
	ExecutedAt string          `json:"executed_at"`
	Runtime    string          `json:"runtime"`
	UrlStatus  []JsonUrlStatus `json:"url_status"`
}

// Helper construct of UrlStatus to customize the JSON conversion
type JsonUrlStatus struct {
	Url           string `json:"url"`
	IsReachable   bool   `json:"is_reachable"`
	StatusMessage string `json:"status_message"`
	ContentLength int64  `json:"content_length"`
	ResponseTime  string `json:"response_time"`
	NumOccured    int    `json:"num_occured"`
}

// Converts UrlReport to JSON string
func (r *UrlReport) Json() (string, error) {
	jsonReport, err := convertToJsonStruct(*r)
	if err != nil {
		return "", err
	}
	jsonBytes, err := json.Marshal(jsonReport)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), err
}

// Converts UrlTeport to Csv string
func (r *UrlReport) Csv(writeHeader bool) (string, error) {
	buf := bytes.Buffer{}
	w := csv.NewWriter(&buf)

	if writeHeader {
		w.Write(csvHeader)
	}
	for _, status := range r.UrlStatus {
		lineContent := []string{
			status.Url,
			fmt.Sprint(status.IsReachable),
			status.StatusMessage,
			fmt.Sprint(status.ContentLength),
			status.ResponseTime.String(),
			fmt.Sprint(status.NumOccured),
		}
		w.Write(lineContent)
	}
	w.Flush()
	err := w.Error()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Internal conversion, to set time.* values as we want them to be
func convertToJsonStruct(report UrlReport) (JsonUrlReport, error) {
	timeConverted, err := report.ExecutedAt.MarshalText()
	if err != nil {
		return JsonUrlReport{}, err
	}

	return JsonUrlReport{
		ExecutedAt: string(timeConverted),
		Runtime:    report.Runtime.String(),
		UrlStatus:  convertUrlStatusToJsonStuct(report.UrlStatus),
	}, nil
}

// Internal conversion, to set time.* values as we want them to be
func convertUrlStatusToJsonStuct(status []UrlStatus) []JsonUrlStatus {
	jsonUrlStatus := []JsonUrlStatus{}
	for _, u := range status {
		j := JsonUrlStatus{
			Url:           u.Url,
			IsReachable:   u.IsReachable,
			StatusMessage: u.StatusMessage,
			ContentLength: u.ContentLength,
			ResponseTime:  u.ResponseTime.String(),
			NumOccured:    u.NumOccured,
		}
		jsonUrlStatus = append(jsonUrlStatus, j)
	}
	return jsonUrlStatus
}
