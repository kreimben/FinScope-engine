package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

/*
	{
	    "realtime_start": "2024-11-21",
	    "realtime_end": "2024-11-21",
	    "observation_start": "2000-01-01",
	    "observation_end": "2002-01-01",
	    "units": "lin",
	    "output_type": 1,
	    "file_type": "json",
	    "order_by": "observation_date",
	    "sort_order": "asc",
	    "count": 3,
	    "offset": 0,
	    "limit": 100000,
	    "observations": [
	        {
	            "realtime_start": "2024-11-21",
	            "realtime_end": "2024-11-21",
	            "date": "2000-01-01",
	            "value": "10250.952"
	        },
	        {
	            "realtime_start": "2024-11-21",
	            "realtime_end": "2024-11-21",
	            "date": "2001-01-01",
	            "value": "10581.929"
	        },
	        {
	            "realtime_start": "2024-11-21",
	            "realtime_end": "2024-11-21",
	            "date": "2002-01-01",
	            "value": "10929.108"
	        }
	    ]
	}
*/
type GDP struct {
	ObservationStart string        `json:"observation_start"`
	ObservationEnd   string        `json:"observation_end"`
	Units            string        `json:"units"` // in this case, it's always "billions of US dollars"
	OutputType       int           `json:"output_type"`
	Observations     []Observation `json:"observations"`
}

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
	// Remove quotes from the JSON string
	str := string(b)
	str = str[1 : len(str)-1]

	// Parse the date string
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return fmt.Errorf("error parsing date: %v", err)
	}

	cd.Time = t
	return nil
}

// UnmarshalJSON for Observation to handle string to float64 conversion
func (o *Observation) UnmarshalJSON(b []byte) error {
	var temp struct {
		Date  CustomDate `json:"date"`
		Value string     `json:"value"`
	}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	// Convert the value from string to float64
	value, err := strconv.ParseFloat(temp.Value, 64)
	if err != nil {
		return fmt.Errorf("error parsing value: %v", err)
	}

	o.Date = temp.Date
	o.Value = value
	return nil
}
