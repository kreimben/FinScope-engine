package models

type PCEPI struct {
	ObservationStart string        `json:"observation_start"`
	ObservationEnd   string        `json:"observation_end"`
	Units            string        `json:"units"`
	OutputType       int           `json:"output_type"`
	Observations     []Observation `json:"observations"`
}
