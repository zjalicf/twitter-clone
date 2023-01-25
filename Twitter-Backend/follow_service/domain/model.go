package domain

type FollowRequest struct {
	ID        string `json:"id,omitempty"`
	Receiver  string `json:"receiver"`
	Requester string `json:"requester"`
	Status    Status `json:"status,omitempty"`
}

type User struct {
	ID        string `json:"id,omitempty"`
	Username  string `json:"username"`
	Age       int    `json:"age"`
	Residence string `json:"residence"`
}

type Ad struct {
	TweetID   string `json:"tweet_id"`
	AgeFrom   int    `json:"age_from"`
	AgeTo     int    `json:"age_to"`
	Gender    string `json:"gender"`
	Residence string `json:"residence"`
}

type Status int

const (
	Pending Status = iota + 1
	Declined
	Accepted
)

type Gender int

const (
	Male Status = iota + 1
	Female
	Both
)

func (status Status) String() string {
	return [...]string{"Pending", "Declined", "Accepted"}[status-1]
}

func (status Status) EnumIndex() int {
	return int(status)
}

func (gender *Gender) String() string {
	return [...]string{"Male", "Female", "Both"}[*gender-1]
}

func (gender *Gender) EnumIndex() int {
	return int(*gender)
}
