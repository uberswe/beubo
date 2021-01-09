package middleware

import (
	"github.com/uberswe/beubo/pkg/structs"
	"gorm.io/gorm"
	"net/http"
)

// CanAccess checks if a user is allowed to access a specified feature
func CanAccess(db *gorm.DB, FeatureKey string, r *http.Request) bool {
	self := r.Context().Value(UserContextKey)
	if self != nil && self.(structs.User).ID > 0 {
		if self.(structs.User).CanAccess(db, FeatureKey) {
			return true
		}
	}
	return false
}
