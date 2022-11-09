package domain

import (
	"encoding/json"
	"io"
)

type Tweet struct {
	ID        int    `json:"id" validate:"required, unique"`
	Text      string `json:"text"`
	Image     string `json:"image"`
	CreatedOn string `json:"createdOn"`
}

type Tweets []*Tweet

func (tweets *Tweets) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(tweets)
}

func (tweet *Tweet) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(tweet)
}

func (tweet *Tweet) FromJSON(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(tweet)
}
