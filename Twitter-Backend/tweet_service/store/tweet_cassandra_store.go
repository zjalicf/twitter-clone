package store

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"os"
	"tweet_service/domain"
)

const (
	DATABASE            = "tweet"
	COLLECTION          = "tweet"
	COLLECTION_BY_USER  = "tweets_by_user"
	COLLECTION_FAVORITE = "favorite"
)

type TweetRepo struct {
	session *gocql.Session
	logger  *log.Logger
}

func New(logger *log.Logger) (*TweetRepo, error) {
	db := os.Getenv("TWEET_DB")

	cluster := gocql.NewCluster(db)
	cluster.Keyspace = "system"
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	err = session.Query(
		fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
					WITH replication = {
						'class' : 'SimpleStrategy',
						'replication_factor' : %d
					}`, DATABASE, 1)).Exec()
	if err != nil {
		logger.Println(err)
	}
	session.Close()

	cluster.Keyspace = DATABASE
	cluster.Consistency = gocql.One
	session, err = cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	return &TweetRepo{
		session: session,
		logger:  logger,
	}, nil
}

func (sr *TweetRepo) CloseSession() {
	sr.session.Close()
}

// Field picture is missing
func (sr *TweetRepo) CreateTables() {
	err := sr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
					(id UUID, text text, created_at time, favorited boolean, favorite_count int, retweeted boolean,
					retweet_count int, username text,
					PRIMARY KEY ((id)))`, //for now there is no clustering order!!
			COLLECTION)).Exec()

	err = sr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
					(id UUID, text text, created_at time, favorited boolean, favorite_count int, retweeted boolean,
					retweet_count int, username text,
					PRIMARY KEY ((username), created_at))
					WITH CLUSTERING ORDER BY (created_at DESC)`, //clustering key by creating date and pk for tweet id and user_id
			COLLECTION_BY_USER)).Exec()

	err = sr.session.Query(
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (tweet_id UUID, username text, PRIMARY KEY ((tweet_id)))",
			COLLECTION_FAVORITE)).Exec()

	if err != nil {
		sr.logger.Printf("CASSANDRA CREATE TABLE ERR: %s", err.Error())
	}
}

//insert into tweet (tweet_id, created_at, favorite_count, favorited, retweet_count, retweeted, text, user_id) values
//(60089906-68d2-11ed-9022-0242ac120002, 1641540002, 0, false, 0, false, 'cao', dae71a94-68d2-11ed-9022-0242ac120002) ;

func (sr *TweetRepo) GetAll() ([]domain.Tweet, error) {
	scanner := sr.session.Query(`SELECT * FROM tweet`).Iter().Scanner()

	var tweets []domain.Tweet
	for scanner.Next() {
		var tweet domain.Tweet
		err := scanner.Scan(&tweet.ID, &tweet.CreatedAt, &tweet.FavoriteCount, &tweet.Favorited, &tweet.RetweetCount,
			&tweet.Retweeted, &tweet.Text, &tweet.Username)
		if err != nil {
			sr.logger.Println(err)
			return nil, err
		}

		tweets = append(tweets, tweet)
	}

	if err := scanner.Err(); err != nil {
		sr.logger.Println(err)
		return nil, err
	}
	return tweets, nil
}

func (sr *TweetRepo) GetTweetsByUser(username string) ([]*domain.Tweet, error) {
	query := fmt.Sprintf(`SELECT * FROM tweets_by_user WHERE username = '%s'`, username)
	fmt.Println(query)
	scanner := sr.session.Query(query).Iter().Scanner()

	var tweets []*domain.Tweet
	for scanner.Next() {
		var tweet domain.Tweet
		err := scanner.Scan(&tweet.Username, &tweet.CreatedAt, &tweet.FavoriteCount, &tweet.Favorited, &tweet.ID,
			&tweet.RetweetCount, &tweet.Retweeted, &tweet.Text)
		if err != nil {
			sr.logger.Println(err)
			return nil, err
		}

		tweets = append(tweets, &tweet)
	}

	if err := scanner.Err(); err != nil {
		sr.logger.Println(err)
		return nil, err
	}
	return tweets, nil
}

func (sr *TweetRepo) Post(tweet *domain.Tweet) (*domain.Tweet, error) {
	insertGeneral := fmt.Sprintf("INSERT INTO %s "+
		"(id, created_at, favorite_count, favorited, retweet_count, retweeted, text, username) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)", COLLECTION)

	insertByUser := fmt.Sprintf("INSERT INTO %s "+
		"(id, created_at, favorite_count, favorited, retweet_count, retweeted, text, username) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)", COLLECTION_BY_USER)

	err := sr.session.Query(
		insertGeneral, tweet.ID, tweet.CreatedAt, tweet.FavoriteCount, tweet.Favorited,
		tweet.RetweetCount, tweet.Retweeted, tweet.Text, tweet.Username).Exec()

	err = sr.session.Query(
		insertByUser, tweet.ID, tweet.CreatedAt, tweet.FavoriteCount, tweet.Favorited,
		tweet.RetweetCount, tweet.Retweeted, tweet.Text, tweet.Username).Exec()

	if err != nil {
		sr.logger.Println(err)
		return nil, err
	}
	return tweet, nil
}

func (sr *TweetRepo) Favorite(tweet_id *gocql.UUID, username *string) (int, error) {

	query_favorite := fmt.Sprintf(`SELECT count(*) FROM favorite WHERE tweet_id = '%s' AND username = '%s'`,
		tweet_id.String(), username)
	log.Printf(query_favorite)
	log.Printf("tu sam")
	return 1, nil
	//var favorite domain.Favorite
	//err := sr.session.Query(query_favorite).Scan(&favorite)
	//
	//if err != nil {
	//	sr.logger.Println(err)
	//	return 502, err
	//}
	//
	//
	//
	//query := fmt.Sprintf(`SELECT * FROM tweet WHERE tweet_id = '%s'`, tweet_id.String())
	//var tweet domain.Tweet
	//err = sr.session.Query(query).Scan(&tweet)
	//
	//if err != nil {
	//	sr.logger.Println(err)
	//	return 502, err
	//}
	//
	//username := tweet.Username
	//createdAt := tweet.CreatedAt
	//favoriteCount := tweet.FavoriteCount + 1
	//
	//err = sr.session.Query(
	//	`UPDATE tweet SET favorited=true, favorite_count=? where tweet_id=?`, favoriteCount, tweet_id.String()).Exec()
	//
	//if err != nil {
	//	sr.logger.Println(err)
	//	return 502, err
	//}
	//
	//err = sr.session.Query(
	//	`UPDATE tweet_by_user SET favorited=true, favorite_count=? where username=? and created_at=?`,
	//	favoriteCount, username, createdAt).Exec()
	//
	//if err != nil {
	//	sr.logger.Println(err)
	//	return 502, err
	//}
	//
	//insert := fmt.Sprintf("INSERT INTO %s "+"(id, username) "+"VALUES (?, ?)", COLLECTION_FAVORITE)
	//
	//err = sr.session.Query(
	//	insert, tweet.ID, username).Exec()
	//
	//if err != nil {
	//	sr.logger.Println(err)
	//	return 502, err
	//}
	//
	//return 200, nil
}
