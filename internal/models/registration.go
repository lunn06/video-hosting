package models

// RegisterRequest swagger::model
type RegisterRequest struct {
	Email       string `json:"email"`
	ChannelName string `json:"channel_name"`
	Password    string `json:"password"`
}
