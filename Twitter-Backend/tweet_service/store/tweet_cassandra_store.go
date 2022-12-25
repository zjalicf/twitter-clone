package store

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"os"
	"tweet_service/domain"
)

const (
	DATABASE                = "tweet"
	COLLECTION              = "tweet"
	COLLECTION_BY_USER      = "tweets_by_user"
	COLLECTION_FAVORITE     = "favorite"
	COLLECTION_FEED_BY_USER = "feed_by_user"
)

type TweetRepo struct {
	session *gocql.Session
	logger  *log.Logger
}

func New(logger *log.Logger) (*TweetRepo, error) {
	db := os.Getenv("TWEET_DB")

	cluster := gocql.NewCluster(db)
	cluster.Keyspace = "system"
	cluster.PageSize = 65535
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

	log.Printf("tweet ERROR IN CREATE TABLE EXECUTION : %s", err)

	err = sr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
					(id UUID, text text, created_at time, favorited boolean, favorite_count int, retweeted boolean,
					retweet_count int, username text,
					PRIMARY KEY ((username), created_at))
					WITH CLUSTERING ORDER BY (created_at DESC)`, //clustering key by creating date and pk for tweet id and user_id
			COLLECTION_BY_USER)).Exec()

	err = sr.session.Query(
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id UUID, tweet_id UUID, username text, PRIMARY KEY ((tweet_id), username))",
			COLLECTION_FAVORITE)).Exec()

	//err = sr.session.Query(
	//	fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id UUID, tweet_id UUID, username text, PRIMARY KEY ((tweet_id)))",
	//		COLLECTION_FAVORITE_BY_TWEET)).Exec()
	log.Printf("tweets_by_user ERROR IN CREATE TABLE EXECUTION : %s", err)

	err = sr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
					(id UUID, text text, created_at time, favorited boolean, favorite_count int, retweeted boolean,
					retweet_count int, username text,
					PRIMARY KEY (username, created_at))
					WITH CLUSTERING ORDER BY (created_at DESC)`,
			COLLECTION_FEED_BY_USER)).Exec()
	if err != nil {
		log.Printf("feed_by_user ERROR IN CREATE TABLE EXECUTION : %s", err)
	}

	//err := sr.session.Query(
	//	fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (tweet_id UUID, text text, PRIMARY KEY ((tweet_id)))",
	//		COLLECTION)).Exec()

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

func (sr *TweetRepo) GetFeedByUser(followings []string) ([]*domain.Tweet, error) {

	fmt.Println(followings)

	query := sr.session.Query(`SELECT * FROM feed_by_user WHERE username IN ? ORDER BY created_at DESC`, followings)
	query.PageSize(0)
	scanner := query.Iter().Scanner()

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
	insert := fmt.Sprintf("INSERT INTO %s "+
		"(id, created_at, favorite_count, favorited, retweet_count, retweeted, text, username) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)", COLLECTION)

	insertByUser := fmt.Sprintf("INSERT INTO %s "+
		"(id, created_at, favorite_count, favorited, retweet_count, retweeted, text, username) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)", COLLECTION_BY_USER)

	insertFeedByUser := fmt.Sprintf("INSERT INTO %s "+
		"(id, created_at, favorite_count, favorited, retweet_count, retweeted, text, username) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)", COLLECTION_FEED_BY_USER)

	err := sr.session.Query(
		insert, tweet.ID, tweet.CreatedAt, tweet.FavoriteCount, tweet.Favorited,
		tweet.RetweetCount, tweet.Retweeted, tweet.Text, tweet.Username).Exec()

	err = sr.session.Query(
		insertByUser, tweet.ID, tweet.CreatedAt, tweet.FavoriteCount, tweet.Favorited,
		tweet.RetweetCount, tweet.Retweeted, tweet.Text, tweet.Username).Exec()

	err = sr.session.Query(
		insertFeedByUser, tweet.ID, tweet.CreatedAt, tweet.FavoriteCount, tweet.Favorited,
		tweet.RetweetCount, tweet.Retweeted, tweet.Text, tweet.Username).Exec()

	if err != nil {
		sr.logger.Println(err)
		return nil, err
	}
	return tweet, nil
}

func (sr *TweetRepo) Favorite(tweetID string, username string) (int, error) {
	id, err := gocql.ParseUUID(tweetID)
	if err != nil {
		return -1, nil
	}

	query := fmt.Sprintf(`SELECT * FROM favorite WHERE tweet_id = %s AND username = '%s'`, id.String(), username)
	scanner := sr.session.Query(query).Iter().Scanner()

	var favorites []*domain.Favorite

	for scanner.Next() {
		var favorite domain.Favorite
		err = scanner.Scan(&favorite.TweetID, &favorite.Username, &favorite.ID)
		if err != nil {
			sr.logger.Println(err)
			return 502, err
		}
		favorites = append(favorites, &favorite)
	}

	//Nekako proslediti idList u funkciju GetLikesByTweet
	//var idList []string
	//for favorite in favorites
	//id, err := gocql.ParseUUID(tweetID)
	//if err != nil {
	//	return nil, err
	//}
	//idList[] = append(idList, favorite.id)

	scanner = sr.session.Query(`SELECT * FROM tweet WHERE id = ?`, id.String()).Iter().Scanner()

	var tweetUsername string
	var tweets []*domain.Tweet
	for scanner.Next() {
		var tweet domain.Tweet
		err := scanner.Scan(&tweet.ID, &tweet.CreatedAt, &tweet.FavoriteCount, &tweet.Favorited, &tweet.RetweetCount,
			&tweet.Retweeted, &tweet.Text, &tweet.Username)
		tweetUsername = tweet.Username
		if err != nil {
			sr.logger.Println(err)
			return 500, err
		}

		tweets = append(tweets, &tweet)
	}

	if len(tweets) == 0 {
		sr.logger.Println("No such tweet")
		return 500, nil
	}

	favorited := false
	create := false
	favoriteCount := 0
	createdAt := tweets[0].CreatedAt

	if len(favorites) != 0 {
		log.Println(214)
		favoriteCount = tweets[0].FavoriteCount - 1
		favorited = false
		create = false
	} else {
		log.Println(218)
		favoriteCount = tweets[0].FavoriteCount + 1
		favorited = true
		create = true
	}

	isDeleted := false

	if create {
		insert := fmt.Sprintf("INSERT INTO %s "+"(id, tweet_id, username) "+"VALUES (?, ?, ?)", COLLECTION_FAVORITE)
		idFav, _ := gocql.RandomUUID()

		err = sr.session.Query(
			insert, idFav.String(), tweets[0].ID.String(), username).Exec()
		if err != nil {
			sr.logger.Println(err)
			return 502, err
		}

		//insert = fmt.Sprintf("INSERT INTO %s "+"(id, tweet_id, username) "+"VALUES (?, ?, ?)", COLLECTION_FAVORITE_BY_TWEET)
		//idFav, _ = gocql.RandomUUID()
		//
		//err = sr.session.Query(
		//	insert, idFav.String(), tweets[0].ID.String(), username).Exec()
		//if err != nil {
		//	sr.logger.Println(err)
		//	return 502, err
		//}
	} else {
		log.Println("u delete sam")
		log.Println(username)
		deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE tweet_id=%s AND username='%s'", COLLECTION_FAVORITE, id, username)

		err = sr.session.Query(deleteQuery).Exec()

		if err != nil {
			sr.logger.Println(err)
			return 502, err
		}

		//deleteQuery = fmt.Sprintf("DELETE FROM %s WHERE tweet_id=%s AND username='%s'", COLLECTION_FAVORITE_BY_TWEET, id, username)
		//
		//err = sr.session.Query(deleteQuery).Exec()
		//
		//if err != nil {
		//	sr.logger.Println(err)
		//	return 502, err
		//}

		isDeleted = true
	}

	err = sr.session.Query(
		`UPDATE tweet SET favorited=?, favorite_count=? where id=?`, favorited, favoriteCount, id.String()).Exec()

	if err != nil {
		sr.logger.Println(err)
		return 502, err
	}

	err = sr.session.Query(
		`UPDATE tweets_by_user SET favorited=?, favorite_count=? where username=? and created_at=?`,
		favorited, favoriteCount, tweetUsername, createdAt).Exec()

	if err != nil {
		sr.logger.Println(err)
		return 502, err
	}

	if isDeleted {
		return 200, nil
	}

	return 201, nil
}

func (sr *TweetRepo) GetLikesByTweet(tweetID string) ([]*domain.Favorite, error) {
	id, err := gocql.ParseUUID(tweetID)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT * FROM favorite WHERE tweet_id = %s`, id.String())
	scanner := sr.session.Query(query).Iter().Scanner()

	var favorites []*domain.Favorite
	for scanner.Next() {
		var favorite domain.Favorite
		err := scanner.Scan(&favorite.TweetID, &favorite.Username, &favorite.ID)
		if err != nil {
			sr.logger.Println(err)
			return nil, err
		}
		favorites = append(favorites, &favorite)
	}

	if err := scanner.Err(); err != nil {
		sr.logger.Println(err)
		return nil, err
	}

	return favorites, nil
}
