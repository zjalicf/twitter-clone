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

	// Connect to default keyspace
	cluster := gocql.NewCluster(db)
	cluster.Keyspace = "system"
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	// Create 'student' keyspace
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

	// Connect to student keyspace
	cluster.Keyspace = DATABASE
	cluster.Consistency = gocql.One
	session, err = cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	// Return repository with logger and DB session
	return &TweetRepo{
		session: session,
		logger:  logger,
	}, nil
}

// Disconnect from database
func (sr *TweetRepo) CloseSession() {
	sr.session.Close()
}

// Field picture is missing
func (sr *TweetRepo) CreateTables() {
	//err := sr.session.Query(
	//	fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
	//				(tweet_id UUID, text text, created_at bigint, favorite boolean, favorite_count int, retweeted boolean,
	//				retweet_count, user_id UUID,
	//				PRIMARY KEY ((tweet_id))
	//				WITH CLUSTERING ORDER BY (tweet_id DESC)`,
	//		COLLECTION)).Exec()

	err := sr.session.Query(
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (tweet_id UUID, text text, PRIMARY KEY ((tweet_id)))",
			COLLECTION)).Exec()
	//CREATE TABLE IF NOT EXISTS %s  ( predmet text, PRIMARY KEY ((student_id), ispit_id, ocena)) WITH CLUSTERING ORDER BY (ispit_id ASC, ocena DESC), "ispiti_by_student")).Exec()

	if err != nil {
		sr.logger.Println(err)
	}
}

func (sr *TweetRepo) GetAll() ([]domain.Tweet, error) {
	//scanner := sr.session.Query(`SELECT tweet_id, text, created_at, favorite, favorite_count, retweeted, retweet_count,
	//   user_id FROM tweet`).Iter().Scanner()
	scanner := sr.session.Query(`SELECT tweet_id, text FROM tweet`).Iter().Scanner()

	var tweets []domain.Tweet
	for scanner.Next() {
		var tweet domain.Tweet
		err := scanner.Scan(&tweet.ID, &tweet.Text)
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

//// Dopunjena metoda tako da podrzava dobavljanje za predmet i smer
//func (sr *TweetRepo) GetIspitiByPredmetAndSmer(predmetId string, smerId string) (IspitiByPredmetAndSmer, error) {
//	scanner := sr.session.Query(`SELECT predmet_id, smer_id, predmet_naziv, smer_naziv, indeks, ime, prezime, ocena, ispit_id, datum
//								FROM ispiti_by_predmet_and_smer WHERE predmet_id = ? and smer_id = ?`,
//		predmetId, smerId).Iter().Scanner()
//
//	var ocene IspitiByPredmetAndSmer
//	for scanner.Next() {
//		var pred IspitByPredmetAndSmer
//		err := scanner.Scan(&pred.PredmetId, &pred.SmerId, &pred.PredmetNaziv, &pred.SmerNaziv, &pred.Indeks, &pred.Ime, &pred.Prezime, &pred.Ocena, &pred.IspitId, &pred.Datum)
//		if err != nil {
//			sr.logger.Println(err)
//			return nil, err
//		}
//		ocene = append(ocene, &pred)
//	}
//	if err := scanner.Err(); err != nil {
//		sr.logger.Println(err)
//		return nil, err
//	}
//	return ocene, nil
//}
//
//func (sr *TweetRepo) GetStudentiBySmer(smerId string) (StudentiBySmer, error) {
//	scanner := sr.session.Query(`SELECT smer_id, student_id, indeks, ime, prezime, smer_naziv, stepeni_studija
//								FROM studenti_by_smer WHERE smer_id = ?`,
//		smerId).Iter().Scanner()
//
//	var studenti StudentiBySmer
//	for scanner.Next() {
//		var stud StudentBySmer
//		err := scanner.Scan(&stud.SmerId, &stud.StudentId, &stud.Indeks, &stud.Ime, &stud.Prezime, &stud.SmerNaziv, &stud.StepeniStudija)
//		if err != nil {
//			sr.logger.Println(err)
//			return nil, err
//		}
//		studenti = append(studenti, &stud)
//	}
//	if err := scanner.Err(); err != nil {
//		sr.logger.Println(err)
//		return nil, err
//	}
//	return studenti, nil
//}
//
//func (sr *TweetRepo) InsertIspitByStudent(studentIspit *IspitByStudent) error {
//	ispitId := gocql.TimeUUID()
//	err := sr.session.Query(
//		`INSERT INTO ispiti_by_student (student_id, ime, prezime, indeks, ocena, ispit_id, predmet)
//		VALUES (?, ?, ?, ?, ?, ?, ?)`,
//		studentIspit.StudentId, studentIspit.Ime, studentIspit.Prezime, studentIspit.Indeks,
//		studentIspit.Ocena, ispitId, studentIspit.Predmet).Exec()
//	if err != nil {
//		sr.logger.Println(err)
//		return err
//	}
//	return nil
//}
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
