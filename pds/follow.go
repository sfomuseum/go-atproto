package pds

type Follow struct {
	FollowerDID  string `json:"follower_did"`
	FollowingDID string `json:"following_did"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}
