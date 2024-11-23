package economic_indicators

import "os"

const FRED_OBSERVATIONS_BASE_URL = "https://api.stlouisfed.org/fred/series/observations"

func FRED_API_KEY() string {
	return os.Getenv("FRED_API_KEY")
}

func getFREDQuery(seriesID string, observationStartDate string, frequency string, units string) string {
	query := NewFSQuery(FRED_OBSERVATIONS_BASE_URL)
	query.Add("api_key", FRED_API_KEY()).And()
	query.Add("series_id", seriesID).And()
	query.Add("frequency", frequency).And()
	query.Add("observation_start", observationStartDate)

	if units != "" {
		query.And().Add("units", units)
	}
	return query.Build()
}
