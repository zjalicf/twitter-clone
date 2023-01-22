package domain

type TweetCache interface {
	PostCacheData(key string, value *[]byte) error
	GetCachedValue(key string) (*[]byte, error)
}
