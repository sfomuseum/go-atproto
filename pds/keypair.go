package pds

type KeyPair struct {
	DID          string `json:"did"`
	Label        string `json:"label"`
	PublicKey    string `json:"public_key"`
	PrivateKey   string `json:"private_key"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"last_modified"`
}
