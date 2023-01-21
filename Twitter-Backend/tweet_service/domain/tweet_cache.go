package domain

type TweetCache interface {
	PostCacheData(key string, value string) error
	GetCachedValue(key string) (string, error)
}
