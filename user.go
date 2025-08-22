package pds

type User struct {
	DID          string   `json:"did"`
	PublicKey    string   `json:"public_key"`
	Handle       *Handle  `json:"handle"`
	Aliases      []*Alias `json:"aliases"`
	Created      int64    `json:"created"`
	LastModified int64    `json:"lastmodified"`
}
