package models

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
