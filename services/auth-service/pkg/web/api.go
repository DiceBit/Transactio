package web

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"user-service/pkg/db/pgx"
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
	//api.router.HandleFunc("/signup", api.registration).Methods(http.MethodPost)
	//api.router.HandleFunc("/login", api.login).Methods(http.MethodPost)

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
/*func (api *API) login(w http.ResponseWriter, req *http.Request) {
	var reqAuth models.AuthRequest
	err := json.NewDecoder(req.Body).Decode(&reqAuth)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		log.Println("Invalid request:", err)
		return
	}

	usrFromBd := userUtils.UsrByEmail(context.Background(), api.db, w, reqAuth.Email)
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

	if exist := userUtils.CheckIfExistUsr(context.Background(), api.db, w, usr); !exist {
		userUtils.AddUser(context.Background(), api.db, w, usr)
		log.Println("User added")
	} else {
		http.Error(w, "User already exists", http.StatusConflict)
		log.Println("User already exists:", err)
		return
	}

}*/
