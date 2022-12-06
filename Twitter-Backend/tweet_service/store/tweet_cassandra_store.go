package store

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"os"
	"tweet_service/domain"
)

const (
	DATABASE   = "tweet"
	COLLECTION = "tweet"
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
					retweet_count int, user_id text,
					PRIMARY KEY ((id)))`, //for now there is no clustering order!!
			COLLECTION)).Exec()

	//err := sr.session.Query(
	//	fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (tweet_id UUID, text text, PRIMARY KEY ((tweet_id)))",
	//		COLLECTION)).Exec()

	if err != nil {
		sr.logger.Println(err)
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
			&tweet.Retweeted, &tweet.Text, &tweet.UserID)
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

func (sr *TweetRepo) Post(tweet *domain.Tweet) (*domain.Tweet, error) {
	err := sr.session.Query(
		`INSERT INTO tweet (id, created_at, favorite_count, favorited, retweet_count, retweeted, text, user_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, tweet.ID, tweet.CreatedAt, tweet.FavoriteCount, tweet.Favorited,
		tweet.RetweetCount, tweet.Retweeted, tweet.Text, tweet.UserID).Exec()
	if err != nil {
		sr.logger.Println(err)
		return nil, err
	}
	return tweet, nil
}

//
//func (sr *TweetRepo) InsertIspitByPredmetAndSmer(predmetSmerIspit *IspitByPredmetAndSmer) error {
//	ispitId, _ := gocql.RandomUUID()
//	err := sr.session.Query(
//		`INSERT INTO ispiti_by_predmet_and_smer (predmet_id, smer_id, predmet_naziv, smer_naziv, indeks, ocena, ispit_id, datum, ime, prezime)
//		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
//		predmetSmerIspit.PredmetId, predmetSmerIspit.SmerId, predmetSmerIspit.PredmetNaziv, predmetSmerIspit.SmerNaziv,
//		predmetSmerIspit.Indeks, predmetSmerIspit.Ocena, ispitId, predmetSmerIspit.Datum, predmetSmerIspit.Ime, predmetSmerIspit.Prezime).Exec()
//	if err != nil {
//		sr.logger.Println(err)
//		return err
//	}
//	return nil
//}
//
//// Zadatak 1
//func (sr *TweetRepo) InsertStudentBySmer(studentSmer *StudentBySmer) error {
//	studentId, _ := gocql.RandomUUID()
//	err := sr.session.Query(
//		`INSERT INTO studenti_by_smer (smer_id, student_id, indeks, ime, prezime, smer_naziv, stepeni_studija)
//		VALUES (?, ?, ?, ?, ?, ?, ?)`,
//		studentSmer.SmerId, studentId, studentSmer.Indeks, studentSmer.Ime, studentSmer.Prezime, studentSmer.SmerNaziv,
//		studentSmer.StepeniStudija).Exec()
//	if err != nil {
//		sr.logger.Println(err)
//		return err
//	}
//	return nil
//}
//
//// Zadatak 4: dodavanje informacije o zavrsenom stepenu studija studenta
//func (sr *TweetRepo) UpdateIspitByPredmetAddStepenStudija(smerId string, studentId string, indeks string, stepenStudija string) error {
//	// za Update je neophodno da pronadjemo vrednost po PRIMARNOM KLJUCU = PK + CK (ukljucuje sve kljuceve particije i klastera)
//	// u ovom slucaju: PK = smerId, CK = student_id, indeks
//	err := sr.session.Query(
//		`UPDATE studenti_by_smer SET stepeni_studija=stepeni_studija+? where smer_id = ? and student_id = ? and indeks = ?`,
//		[]string{stepenStudija}, smerId, studentId, indeks).Exec()
//	if err != nil {
//		sr.logger.Println(err)
//		return err
//	}
//	return nil
//}
//
//// NoSQL: Performance issue, we never want to fetch all the data
//// (In order to get all student ids we need to contact every partition which are usually located on different servers!)
//// Here we are doing it for demonstration purposes (so we can see all student/predmet ids)
//func (sr *TweetRepo) GetDistinctIds(idColumnName string, tableName string) ([]string, error) {
//	scanner := sr.session.Query(
//		fmt.Sprintf(`SELECT DISTINCT %s FROM %s`, idColumnName, tableName)).
//		Iter().Scanner()
//	var ids []string
//	for scanner.Next() {
//		var id string
//		err := scanner.Scan(&id)
//		if err != nil {
//			sr.logger.Println(err)
//			return nil, err
//		}
//		ids = append(ids, id)
//	}
//	if err := scanner.Err(); err != nil {
//		sr.logger.Println(err)
//		return nil, err
//	}
//	return ids, nil
//}
