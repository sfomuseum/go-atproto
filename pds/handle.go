package pds

type Handle struct {
	DID          string `json:"did"`
	Name         string `json:"name"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}
