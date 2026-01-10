package data

type VotingObject interface {
	GetDescription(int) string
	GetButtonTitle() string
	GetCallbackData(int) string
}
