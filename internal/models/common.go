// Start of Selection
package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Observation struct {
	Date  CustomDate `json:"date"`
	Value float64    `json:"value"`
}

// CustomDate is a custom type for parsing date strings in the format "2006-01-02"
type CustomDate struct {
	time.Time
}

// UnmarshalJSON parses the date string in the format "2006-01-02"
func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	var dateStr string
	if err := json.Unmarshal(b, &dateStr); err != nil {
		return fmt.Errorf("CustomDate UnmarshalJSON: %v", err)
	}
	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("CustomDate UnmarshalJSON: %v", err)
	}
	cd.Time = parsedTime
	return nil
}

// Start of Selection
// UnmarshalJSON for Observation to handle string to float64 conversion
func (o *Observation) UnmarshalJSON(b []byte) error {
	type Alias Observation
	aux := &struct {
		Value interface{} `json:"value"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(b, &aux); err != nil {
		return fmt.Errorf("Observation UnmarshalJSON: %v", err)
	}

	switch v := aux.Value.(type) {
	case float64:
		o.Value = v
	case string:
		if v == "." {
			o.Value = 0.0
		} else {
			parsedValue, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return fmt.Errorf("Observation UnmarshalJSON: cannot parse value: %v", err)
			}
			o.Value = parsedValue
		}
	default:
		return fmt.Errorf("Observation UnmarshalJSON: unexpected type for value")
	}

	return nil
}
