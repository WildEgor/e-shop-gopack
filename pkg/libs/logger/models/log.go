package models

type LogEntry struct {
	Organization string                 `json:"organization"`
	App          string                 `json:"app"`
	Label        string                 `json:"label"`
	Err          any                    `json:"error"`
	Props        map[string]interface{} `json:"props"`
}
