package beubo

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func errHandler(err error) {
	if err != nil {
		log.Print(err)
	}
}

// SetFlash sets a cookie which expires after the next page load
func SetFlash(w http.ResponseWriter, name string, value []byte) {
	c := &http.Cookie{Name: name, Value: encode(value)}
	http.SetCookie(w, c)
}

// GetFlash gets the cookie value set by SetFlash and removes the cookie
func GetFlash(w http.ResponseWriter, r *http.Request, name string) ([]byte, error) {
	c, err := r.Cookie(name)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return nil, nil
		default:
			return nil, err
		}
	}
	value, err := decode(c.Value)
	if err != nil {
		return nil, err
	}
	// Deletes the cookie
	dc := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(1, 0)}
	http.SetCookie(w, dc)
	return value, nil
}

// encode encodes a byte array into an urlencoded string
func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

// decode decodes an url encoded string into a byte array
func decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}

// generateToken generates a random string of len length
func generateToken(len int) (string, error) {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
