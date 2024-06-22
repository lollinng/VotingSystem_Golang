package models

type VoteCount struct {
	Option string `json:"option"`
	Count  int    `json:"count"`
}
