package middleware

import (
	"fmt"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/jinzhu/gorm"
	"github.com/markustenghamn/beubo/pkg/structs"
	"log"
	"net/http"
)

type BeuboMiddleware struct {
	DB *gorm.DB
}

func (bmw *BeuboMiddleware) Auth(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session := sessions.GetSession(r)
	token := session.Get("SES_ID")

	user := structs.FetchUserFromSession(bmw.DB, fmt.Sprintf("%v", token))

	if user.ID == 0 {
		log.Println("user is not logged in")
		http.Redirect(rw, r, "/login", 302)
		return
	}

	log.Printf("Auth middleware hit: {user:%d}\n", user.ID)

	next(rw, r)
}
