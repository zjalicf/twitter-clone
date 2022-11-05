package data

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type TweetRepoPostgreSql struct {
	log      *log.Logger
	database *gorm.DB
}

func NewPostgreSql(log *log.Logger) (TweetRepoPostgreSql, error) {
	username := os.Getenv("db_username")
	host := os.Getenv("db_host")
	password := os.Getenv("db_password")
	name := os.Getenv("db_name")
	port := os.Getenv("db_port")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s", host, username, name, password, port)

	db, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})

	if err != nil {
		return TweetRepoPostgreSql{}, err
	}

	setup(db)
	return TweetRepoPostgreSql{log, db}, nil
}

func setup(database *gorm.DB) {
	database.AutoMigrate(&Tweet{})
}

func (tweetRepo *TweetRepoPostgreSql) GetAll() Tweets {
	tweetRepo.log.Println("{TweetRepoPostgreSql} - getting all Tweets")
	var tweets []*Tweet

	tweetRepo.database.Find(&tweets)

	return tweets
}

func (tweetRepo *TweetRepoPostgreSql) PostTweet(tweet *Tweet) Tweet {
	tweetRepo.log.Println("{TweetRepoPostgreSql} - adding tweet")

	tweet.CreatedOn = time.Now().UTC().String()

	tweetRepo.database.Create(tweet)
	return *tweet
}

func (tweetRepo *TweetRepoPostgreSql) PutTweet(tweet *Tweet, id int) error {
	tweetRepo.log.Println("{TweetRepoPostgreSql} - Updating Tweet")

	var foundTweet Tweet

	tweetRepo.database.Where("id = ?", id).Find(&foundTweet)

	if foundTweet.ID == 0 {
		return errors.New(fmt.Sprintf("Product with id %d not found", id))
	}

	foundTweet.Text = tweet.Text
	foundTweet.Image = tweet.Image

	tweetRepo.database.Save(&foundTweet)

	*tweet = foundTweet
	return nil
}

func (tweetRepo *TweetRepoPostgreSql) DeleteTweet(id int) error {
	tweetRepo.log.Println("{TweetRepoPostgreSql} - Deleting Tweet")

	var tweet Tweet

	tweetRepo.database.Where("id = ?", id).Find(&tweet)

	if tweet.ID == 0 {
		return errors.New(fmt.Sprintf("Tweet with id %d not found", id))
	}

	tweetRepo.database.Save(&tweet)
	return nil
}
