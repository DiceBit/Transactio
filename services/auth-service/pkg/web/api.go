package web

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	pgx2 "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
	"user-service/pkg/db/pgx"
	"user-service/pkg/models"
	"user-service/pkg/utils"
)

type API struct {
	router *mux.Router
	db     *pgxpool.Pool
}

func New() *API {
	db, err := pgx.New()
	if err != nil {
		log.Println("Error with DB", err)
	}

	api := API{
		router: mux.NewRouter(),
		db:     db,
	}
	api.Endpoints()
	return &api
}

func (api *API) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, api.router))
}

func (api *API) Endpoints() {
	//api.router.Use(JwtMiddleware)

	api.router.HandleFunc("/test", api.test).Methods(http.MethodGet)
	api.router.HandleFunc("/login", api.login).Methods(http.MethodPost)
	api.router.HandleFunc("/signup", api.registration).Methods(http.MethodPost)

	protected := api.router.PathPrefix("/protected").Subrouter()
	protected.Use(JwtMiddleware)
	protected.Use(RoleMiddleware("admin"))
	protected.HandleFunc("/admin", api.admTest).Methods(http.MethodGet)

	usrProt := api.router.PathPrefix("/usrProt").Subrouter()
	usrProt.Use(JwtMiddleware)
	usrProt.Use(RoleMiddleware("user"))
	usrProt.HandleFunc("/user", api.usrTest).Methods(http.MethodGet)

}

func (api *API) test(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("test123"))
}

func (api *API) usrTest(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("user test"))
}

func (api *API) admTest(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("admin test"))
}

// вход
func (api *API) login(w http.ResponseWriter, req *http.Request) {
	var reqAuth models.AuthRequest
	err := json.NewDecoder(req.Body).Decode(&reqAuth)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		log.Println("Invalid request:", err)
		return
	}

	usrFromBd := getUsrByEmail(context.Background(), api.db, w, reqAuth.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error hashing password:", err)
		return
	}

	if ok := bcrypt.CompareHashAndPassword([]byte(usrFromBd.Password), []byte(reqAuth.Password)); ok != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		log.Println("Invalid password")
		return
	}

	token, err := utils.GenerateJWT(usrFromBd.Email, usrFromBd.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		log.Println("Error generating token:", err)
		return
	}

	json.NewEncoder(w).Encode(token)
}

func (api *API) registration(w http.ResponseWriter, req *http.Request) {
	var usr models.UserAuthRequest
	err := json.NewDecoder(req.Body).Decode(&usr)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Error decoding request body: ", err)
		return
	}
	defer req.Body.Close()

	if exist := checkIfExistUsr(context.Background(), api.db, w, usr); !exist {
		addUser(context.Background(), api.db, w, usr)
		log.Println("User added")
	} else {
		http.Error(w, "User already exists", http.StatusConflict)
		log.Println("User already exists:", err)
		return
	}

}

func addUser(ctx context.Context, db *pgxpool.Pool, w http.ResponseWriter, usr models.UserAuthRequest) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		http.Error(w, "Error with transaction", http.StatusInternalServerError)
		log.Println("Error with transaction:", err)
		return
	}
	defer tx.Rollback(ctx)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error hashing password:", err)
		return
	}
	usr.Password = string(hashedPassword)

	var batch = &pgx2.Batch{}
	batch.Queue(`insert into Users(username, email, password, balance, role, createat) values ($1,$2,$3,$4,$5, $6)`,
		usr.Username,
		usr.Email,
		usr.Password,
		0,
		usr.Role,
		time.Now())

	res := tx.SendBatch(ctx, batch)

	err = res.Close()
	if err != nil {
		http.Error(w, "Error with transaction", http.StatusInternalServerError)
		log.Println("Error with transaction:", err)
		return
	}

	tx.Commit(ctx)
}

func getUsrByEmail(ctx context.Context, db *pgxpool.Pool, w http.ResponseWriter, email string) *models.User {
	tx, err := db.Begin(ctx)
	if err != nil {
		http.Error(w, "Error with transaction", http.StatusInternalServerError)
		log.Println("Error with transaction:", err)
		return nil
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `select email, password, role from users where email=$1`, email)
	if err != nil {
		http.Error(w, "Error with query", http.StatusInternalServerError)
		log.Println("Error with query:", err)
		return nil
	}
	defer rows.Close()

	var usr models.User
	if !rows.Next() {
		http.Error(w, "User not found", http.StatusNotFound)
		log.Println("User not found:", err)
		return nil
	}

	err = rows.Scan(
		&usr.Email,
		&usr.Password,
		&usr.Role,
	)
	if err != nil {
		http.Error(w, "Error with scanning query", http.StatusInternalServerError)
		log.Println("Error with  scanning query:", err)
		return nil
	}

	tx.Commit(ctx)
	return &usr
}

func checkIfExistUsr(ctx context.Context, db *pgxpool.Pool, w http.ResponseWriter, usr models.UserAuthRequest) bool {
	tx, err := db.Begin(ctx)
	if err != nil {
		http.Error(w, "Error with transaction", http.StatusInternalServerError)
		log.Println("Error with transaction:", err)
		return false
	}
	defer tx.Rollback(ctx)

	_, err = db.Query(ctx, `select * from users where email=$1 or username=$2`, usr.Email, usr.Username)
	if err != nil {
		tx.Commit(ctx)
		return false
	} else {
		tx.Commit(ctx)
		return true
	}
}
