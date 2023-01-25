package application

import (
	"fmt"
	"follow_service/domain"
	"follow_service/errors"
	"github.com/google/uuid"
	"github.com/zjalicf/twitter-clone-common/common/saga/create_user"
	"log"
)

type FollowService struct {
	store domain.FollowRequestStore
}

func NewFollowService(store domain.FollowRequestStore) *FollowService {
	return &FollowService{
		store: store,
	}
}

func (service *FollowService) GetAll() ([]*domain.FollowRequest, error) {
	return service.store.GetAll()
}

func (service *FollowService) FollowExist(followRequest *domain.FollowRequest) (bool, error) {
	return service.store.FollowExist(followRequest)
}

func (service *FollowService) GetFollowingsOfUser(username string) ([]*string, error) {
	followings, err := service.store.GetFollowingsOfUser(username)
	if err != nil {
		return nil, err
	}

	var usernameList []*string
	for i := 0; i < len(followings); i++ {
		usernameList = append(usernameList, &followings[i].Username)
	}
	log.Println("LIST OF USERNAMES FOLLOW SERVICE: ")
	log.Println(usernameList)
	usernameList = append(usernameList, &username)
	return usernameList, nil
}

func (service *FollowService) GetRequestsForUser(username string) ([]*domain.FollowRequest, error) {
	return service.store.GetRequestsForUser(username)
}

func (service *FollowService) CreateRequest(request *domain.FollowRequest, username string, visibility bool) error {

	request.ID = uuid.New().String()
	request.Requester = username

	isExist, err := service.FollowExist(request)
	if err != nil {
		return err
	}

	if isExist {
		return fmt.Errorf("You already follow this user!")
	}

	if visibility {
		existing, err := service.store.GetRequestByRequesterReceiver(&request.Requester, &request.Receiver)
		if err != nil {
			if err.Error() == errors.ErrorRequestNotExists {
				request.Status = 1
				err = service.store.SaveRequest(request)
				if err != nil {
					log.Println(err)
					return fmt.Errorf("Request not inserted in db")
				}
				return nil
			} else {
				return err
			}
		}

		existing.Status = 1
		err = service.store.UpdateRequest(existing)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("Request not inserted in db")
		}

	} else {
		err := service.store.SaveFollow(request)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *FollowService) CreateUser(user *domain.User) error {
	err := service.store.SaveUser(user)
	if err != nil {
		log.Printf("Error with saving user node in neo4j: %s", err.Error())
		return err
	}

	return nil
}

func (service *FollowService) AcceptRequest(id *string) error {
	request, err := service.store.AcceptRequest(id)
	if err != nil {
		return fmt.Errorf(errors.ErrorInAcceptRequest)
	}

	err = service.store.SaveFollow(request)
	if err != nil {
		return fmt.Errorf(errors.ErrorInSaveFollow)
	}

	return nil

}

func (service *FollowService) DeleteUser(id *string) error {
	return service.store.DeleteUser(id)
}

func (service *FollowService) DeclineRequest(id *string) error {
	return service.store.DeclineRequest(id)
}

func (service *FollowService) HandleRequest(followRequest *domain.FollowRequest) error {
	return service.store.SaveRequest(followRequest)
}

func (service *FollowService) SaveAd(ad *domain.Ad) error {
	return service.store.SaveAd(ad)
}

func (service *FollowService) GetRecommendationsByUsername(username string) ([]string, error) {

	countFollowings, err := service.store.CountFollowings(username)
	if err != nil {
		log.Println("Error in getting count of followings")
		return nil, err
	}

	if countFollowings == 0 {
		recommendations, err := service.store.RecommendationWithoutFollowings(username, []string{})
		if err != nil {
			log.Println("Error in getting similar recommendations without followings.")
			return nil, err
		}
		return recommendations, nil
	} else {
		var allRecommendations []string
		recommendations, err := service.store.RecommendWithFollowings(username)
		if err != nil {
			log.Println("Error in getting recommendations full with followings.")
			return nil, err
		}

		recommendations2, err := service.store.RecommendationWithoutFollowings(username, recommendations)
		if err != nil {
			log.Println("Error in getting similar recommendations with followings.")
			return nil, err
		}

		allRecommendations = append(recommendations, recommendations2...)
		return allRecommendations, nil
	}

}

func (service *FollowService) UserToDomain(userIn create_user.User) domain.User {
	var user domain.User
	user.ID = userIn.ID.Hex()
	user.Age = userIn.Age
	user.Residence = userIn.Residence
	user.Username = userIn.Username

	return user
}
