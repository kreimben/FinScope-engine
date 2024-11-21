package economic_indicators

type FSQuery struct {
	// FinScope Economic Indicator Query URL Struct
	url string
}

func NewFSQuery(url string) *FSQuery {
	return (&FSQuery{url: url + "?"})
}

func (q *FSQuery) Add(key string, value string) *FSQuery {
	q.url += key + "=" + value
	return q
}

func (q *FSQuery) And() *FSQuery {
	q.url += "&"
	return q
}

func (q *FSQuery) Build() string {
	q.And().Add("file_type", "json")
	return q.url
}
