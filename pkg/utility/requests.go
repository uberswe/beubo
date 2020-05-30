package utility

import (
	"encoding/base64"
	"net/http"
	"time"
)

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
	value, err := Decode(c.Value)
	if err != nil {
		return nil, err
	}
	// Deletes the cookie
	dc := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(1, 0)}
	http.SetCookie(w, dc)
	return value, nil
}

// SetFlash sets a cookie which expires after the next page load
func SetFlash(w http.ResponseWriter, name string, value []byte) {
	c := &http.Cookie{Name: name, Value: Encode(value)}
	http.SetCookie(w, c)
}

// Encode encodes a byte array into an urlencoded string
func Encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

// Decode decodes an url encoded string into a byte array
func Decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}
