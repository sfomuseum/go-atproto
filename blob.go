package pds

type Blob struct {
	CID          string `json:"cid"`
	DID          string `json:"did"`
	MediaType    string `json:"media_type"`
	Size         int64  `json:"size"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}
