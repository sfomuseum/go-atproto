package pds

import (
	"fmt"
)

type Record struct {
	CID          string `json:"cid"`
	DID          string `json:"did"`
	Collection   string `json:"collection"`
	RKey         string `json:"rkey"`
	Value        string `json:"value"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}

func (r *Record) BlockURI() string {
	return fmt.Sprintf("repo:%s/%s/%s", r.DID, r.Collection, r.RKey)
}
