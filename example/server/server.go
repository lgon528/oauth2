package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-session/session"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

var srv *server.Server

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// user store
	userStore := store.NewUserStore()
	userStore.SetUser("10000", &models.User{
		ClientID: "10000",
		ID:       "forrest",
		OpenID:   generateOpenID("10000#forrest"),
	})
	userStore.SetUser("10000", &models.User{
		ClientID: "10000",
		ID:       "jason",
		OpenID:   generateOpenID("10000#jason"),
	})
	userStore.SetUser("10000", &models.User{
		ClientID: "10000",
		ID:       "eric",
		OpenID:   generateOpenID("10000#eric"),
	})
	manager.MapUserStorage(userStore)

	// generate jwt access token
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte("00000000"), jwt.SigningMethodHS512))

	clientStore := store.NewClientStore()
	clientStore.Set("10000", &models.Client{
		ID:     "10000",
		Secret: "10000",
		Domain: "http://localhost:9094",
	})
	manager.MapClientStorage(clientStore)

	srv = server.NewServer(server.NewConfig(), manager)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		if username == "test" && password == "test" {
			userID = "test"
		}
		return
	})

	srv.SetUserAuthorizationHandler(userAuthorizeHandler2)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/oauth2/test/login", loginHandler)
	http.HandleFunc("/oauth2/test/auth", authHandler)

	http.HandleFunc("/oauth2/test/authorize", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("we're here, %v", r)
		store, err := session.Start(nil, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var form url.Values
		if v, ok := store.Get("ReturnUri"); ok {
			form = v.(url.Values)
		}
		r.Form = form

		store.Delete("ReturnUri")
		store.Save()

		err = srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/oauth2/test/access_token", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("we're here %v", r)
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			log.Printf("we're here, err %v", err)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error_code": -1,
				"error_msg":  err.Error,
			})
			// http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/oauth2/test/refresh_token", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("we're here, %v", r)

		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/oauth2/test/test", func(w http.ResponseWriter, r *http.Request) {
		token, err := srv.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := map[string]interface{}{
			"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
			"client_id":  token.GetClientID(),
			"user_id":    token.GetUserID(),
		}
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(data)
	})

	log.Println("Server is running at 9096 port.")
	log.Fatal(http.ListenAndServe(":9096", nil))
}

func generateOpenID(src string) string {
	icode := make([]byte, len(src))

	for i, c := range []byte(src) {
		icode[i] = c + 250
	}

	return base64.StdEncoding.EncodeToString(icode)
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/oauth2/test/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}

func userAuthorizeHandler2(w http.ResponseWriter, r *http.Request) (userID string, err error) {

	if platform := r.FormValue("platform"); platform == "" || platform == "web" {
		store, e := session.Start(nil, w, r)
		if e != nil {
			err = e
			return
		}

		_, ok := store.Get("LoggedInUserID")
		if !ok {
			if r.Form == nil {
				r.ParseForm()
			}

			store.Set("ReturnUri", r.Form)
			store.Save()

			w.Header().Set("Location", "/oauth2/test/login")
			w.WriteHeader(http.StatusFound)
			return
		}

		store.Delete("LoggedInUserID")
		store.Save()
	}

	user, err := srv.Manager.GetUser(r.FormValue("client_id"), r.FormValue("userid"))
	if err != nil {
		// json.NewEncoder(w).Encode(map[string]interface{}{
		// 	"error_code": -1,
		// 	"error_msg":  "user not logged in",
		// }
		return
	}

	userID = user.GetOpenID()

	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("we're here, loginHandler")
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		store.Set("LoggedInUserID", "000000")
		store.Save()

		w.Header().Set("Location", "/oauth2/test/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "static/login.html")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		w.Header().Set("Location", "/oauth2/test/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	outputHTML(w, r, "static/auth.html")
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}
