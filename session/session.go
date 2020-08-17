package session

// Session token is mongo
type Session struct {
	ID        string `json:"-" bson:"_id"`
	IsServer  bool   `json:"isServer" bson:"isServer"`
	ServerID  string `json:"-" bson:"serverid,omitempty"`
	UserID    string `json:"userid,omitempty" bson:"userid,omitempty"`
	ExpiresAt int64  `json:"-" bson:"expiresAt"`
	IssuedAt  int64  `json:"-" bson:"issuedAt"`
	Extended  bool   `json:"-" bson:"extended"`
}
