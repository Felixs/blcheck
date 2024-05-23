package report

import (
	"encoding/json"
	"fmt"

	"github.com/Felixs/blcheck/pkg/url"
)

// Helper construct of UrlReport to customize the JSON conversion
type JsonUrlReport struct {
	ExecutedAt string          `json:"executed_at"`
	Runtime    string          `json:"runtime"`
	UrlStatus  []url.UrlStatus `json:"url_status"`
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
		UrlStatus:  report.UrlStatus,
	}
}
