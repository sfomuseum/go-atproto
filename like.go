package pds

type Like struct {
	DID          string `json:"did"`
	RecordCID    string `json:"record_cid"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}
