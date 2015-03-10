package controllers

// SessionResponse ...
type SessionResponse struct {
	Token     string `json:"token,omitempty"`
	UserID    string `json:"uid"`
	Expire    int64  `json:"expire"`
	Remaining int    `json:"remaining"`
	Domain    string `json:"domain,omitempty"`
}
