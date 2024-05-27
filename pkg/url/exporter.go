package url

import (
	"encoding/json"
	"fmt"
)

// Helper construct of UrlReport to customize the JSON conversion
type JsonUrlReport struct {
	ExecutedAt string          `json:"executed_at"`
	Runtime    string          `json:"runtime"`
	UrlStatus  []JsonUrlStatus `json:"url_status"`
}

type JsonUrlStatus struct {
	Url           string `json:"url"`
	IsReachable   bool   `json:"is_reachable"`
	StatusMessage string `json:"status_message"`
	ContentLength int64  `json:"content_length"`
	ResponseTime  string `json:"response_time"`
	NumOccured    int    `json:"num_occured"`
}

// Converts UrlReport to JSON string
func (r *UrlReport) Json() string {
	jsonReport := convertToJsonStruct(*r)
	jsonBytes, err := json.Marshal(jsonReport)
	if err != nil {
		fmt.Printf("failed to marshal UrlReport to json. %v\n", r)
	}
	return string(jsonBytes)
}

// Internal conversion, to set time.* values as we want them to be
func convertToJsonStruct(report UrlReport) JsonUrlReport {
	timeConverted, err := report.ExecutedAt.MarshalText()
	if err != nil {
		fmt.Printf("Failed to convert time into json string, %v", report.ExecutedAt)
	}

	return JsonUrlReport{
		ExecutedAt: string(timeConverted),
		Runtime:    report.Runtime.String(),
		UrlStatus:  convertUrlStatusToJsonStuct(report.UrlStatus),
	}
}

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
