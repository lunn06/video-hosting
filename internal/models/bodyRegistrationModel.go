package models

type BodyForRegistration struct {
	Email       string
	ChannelName string
	Password    string
}

type BodyForAuthorization struct {
	Email    string
	Password string
}
