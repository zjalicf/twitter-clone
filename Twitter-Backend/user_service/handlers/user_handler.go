package handlers

import (
	"encoding/json"
	"github.com/casbin/casbin"
	"github.com/cristalhq/jwt/v4"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"os"
	"strings"
	"user_service/application"
	"user_service/authorization"
	"user_service/domain"
	"user_service/errors"
)

var (
	jwtKey      = []byte(os.Getenv("SECRET_KEY"))
	verifier, _ = jwt.NewVerifierHS(jwt.HS256, jwtKey)
)

type UserHandler struct {
	service *application.UserService
	tracer  trace.Tracer
}

func NewUserHandler(service *application.UserService, tracer trace.Tracer) *UserHandler {
	return &UserHandler{
		service: service,
		tracer:  tracer,
	}
}

func (handler *UserHandler) Init(router *mux.Router) {
	authEnforcer, err := casbin.NewEnforcerSafe("./auth_model.conf", "./policy.csv")
	log.Println("successful init of enforcer")
	if err != nil {
		log.Fatal(err)
	}

	router.Use(ExtractTraceInfoMiddleware)
	router.HandleFunc("/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/", handler.Register).Methods("POST")
	router.HandleFunc("/", handler.GetAll).Methods("GET")
	router.HandleFunc("/getOne/{username}", handler.GetOne).Methods("GET")
	router.HandleFunc("/getMe/", handler.GetMe).Methods("GET")
	router.HandleFunc("/mailExist/{mail}", handler.MailExist).Methods("GET")
	router.HandleFunc("/visibility", handler.ChangeVisibility).Methods("PUT")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8002", authorization.Authorizer(authEnforcer)(router)))
}

func (handler *UserHandler) Register(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "UserHandler.Register")
	defer span.End()

	var user domain.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	saved, err := handler.service.Register(ctx, &user)
	if err != nil {
		if err.Error() == errors.DatabaseError {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		}
		return
	}

	jsonResponse(saved, writer)
}

func (handler *UserHandler) GetAll(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "UserHandler.GetAll")
	defer span.End()

	users, err := handler.service.GetAll(ctx)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse(users, writer)
}

func (handler *UserHandler) Get(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "UserHandler.Get")
	defer span.End()

	vars := mux.Vars(req)
	id, ok := vars["id"]
	if !ok {
		log.Println("id get err")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("primitive get err")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := handler.service.Get(ctx, objectID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	//otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	jsonResponse(user, writer)
}

func (handler *UserHandler) MailExist(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "UserHandler.MailExist")
	defer span.End()

	vars := mux.Vars(req)
	mail, ok := vars["mail"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := handler.service.DoesEmailExist(ctx, mail)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = writer.Write([]byte(id))
	if err != nil {
		log.Println("error in response user service")
		log.Println(err.Error())
		return
	}
}

func (handler *UserHandler) ChangeVisibility(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "UserHandler.ChangeVisibility")
	defer span.End()

	bearer := req.Header.Get("Authorization")
	bearerToken := strings.Split(bearer, "Bearer ")
	tokenString := bearerToken[1]
	parsedToken, err := jwt.Parse([]byte(tokenString), verifier)
	if err != nil {
		log.Println(errors.InvalidTokenError)
		http.Error(writer, errors.InvalidTokenError, http.StatusNotAcceptable)
		return
	}

	claims := parsedToken.Claims()
	var claimsMap map[string]string
	err = json.Unmarshal(claims, &claimsMap)
	if err != nil {
		log.Printf("Unmarshal claims error occured: %s", err.Error())
		http.Error(writer, errors.InvalidTokenError, http.StatusNotAcceptable)
		return
	}

	err = handler.service.ChangeUserVisibility(ctx, claimsMap["user_id"])
	if err != nil {
		log.Printf("Error occured in change user visibility: %s", err.Error())
		if err.Error() == errors.UserNotFound {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (handler *UserHandler) GetOne(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "UserHandler.GetOne")
	defer span.End()

	vars := mux.Vars(req)
	username := vars["username"]

	user, err := handler.service.GetOneUser(ctx, username)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusNotFound)
	}
	jsonResponse(user, writer)
}

func (handler *UserHandler) GetMe(writer http.ResponseWriter, req *http.Request) {
	ctx, span := handler.tracer.Start(req.Context(), "UserHandler.GetMe")
	defer span.End()

	bearer := req.Header.Get("Authorization")
	bearerToken := strings.Split(bearer, "Bearer ")
	tokenString := bearerToken[1]

	token, err := jwt.Parse([]byte(tokenString), verifier)
	if err != nil {
		log.Println(err)
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}

	claims := authorization.GetMapClaims(token.Bytes())
	username := claims["username"]

	user, err := handler.service.GetOneUser(ctx, username)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusNotFound)
	}
	jsonResponse(user, writer)
}

func ExtractTraceInfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
