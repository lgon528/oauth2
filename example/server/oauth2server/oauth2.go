package oauth2server

import (
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-session/session"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"log"
	"net/http"
	"sync"
)

var (
	svr     *server.Server
	svrOnce sync.Once
)

func GetOauth2Server() *server.Server {
	if svr == nil {
		Init()
	}

	return svr
}

func Init() {
	svrOnce.Do(func() {
		manager := manage.NewDefaultManager()
		manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

		// token store
		manager.MustTokenStorage(store.NewMemoryTokenStore())

		// user store
		userStore := store.NewMemoryUserStore()
		users := []models.User{
			models.User{
				ID:       "forrest",
				Password: "12345",
			},
			models.User{
				ID:       "jason",
				Password: "12345",
			},
			models.User{
				ID:       "eric",
				Password: "12345",
			},
		}

		for _, user := range users {
			err := userStore.SetUser(&user)
			if err != nil {
				// todo what to do with err
			}
		}
		manager.MapUserStorage(userStore)

		// client store
		clientStore := store.NewClientStore()

		client := models.Client{
			ID:     "10000",
			Secret: "10000",
			Domain: "http://localhost:9094",
		}

		for _, user := range users {
			err := client.SetOpenUser(&models.OpenUser{
				User:     user,
				OpenID:   generateOpenID(client.GetID(), user.GetID()),
				ClientID: client.GetID(),
			})

			if err != nil {
				// todo what to do with err
			}
		}
		err := clientStore.Set(client.GetID(), &client)
		if err != nil {
			// todo what to do with err
		}

		manager.MapClientStorage(clientStore)

		// generate jwt access token
		manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte("00000000"), jwt.SigningMethodHS512))

		svr = server.NewServer(server.NewConfig(), manager)
		svr.SetClientInfoHandler(server.ClientFormHandler)

		svr.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
			if username == "test" && password == "test" {
				userID = "test"
			}
			return
		})

		svr.SetUserAuthorizationHandler(userAuthorizeHandler)

		svr.SetInternalErrorHandler(func(err error) (re *errors.Response) {
			log.Println("Internal Error:", err.Error())
			return
		})

		svr.SetResponseErrorHandler(func(re *errors.Response) {
			log.Println("Response Error:", re.Error.Error())
		})
	})
}

func generateOpenID(clientID, userid string) string {
	src := clientID + "#" + userid

	icode := make([]byte, len(src))
	for i, c := range []byte(src) {
		icode[i] = c + 250
	}

	return base64.StdEncoding.EncodeToString(icode)
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {

	if platform := r.FormValue("platform"); platform == "" || platform == "web" {
		store, e := session.Start(nil, w, r)
		if e != nil {
			err = e
			return
		}

		_, ok := store.Get(r.FormValue("userid"))
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

		store.Delete(r.FormValue("userid"))
		store.Save()
	}

	clientID := r.FormValue("client_id")
	cli, err := svr.Manager.GetClient(clientID)
	if err != nil {
		log.Printf("client_id %s not registered", clientID)
		return
	}

	userid := r.FormValue("userid")
	user, err := svr.Manager.GetUser(userid)
	if err != nil {
		log.Printf("user %s not existed", userid)
		return
	}

	openUser, ok := cli.GetOpenUser(userid)
	if !ok {
		openUser = &models.OpenUser{
			User: models.User{
				ID:       user.GetID(),
				Password: user.GetPassword(),
			},
			OpenID:   generateOpenID(clientID, userid),
			ClientID: clientID,
		}
	}
	err = cli.SetOpenUser(openUser)
	if err != nil {
		log.Printf("add open user info")
		return
	}

	userID = openUser.GetOpenID()

	return
}
