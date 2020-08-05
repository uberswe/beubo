package middleware

import "github.com/jinzhu/gorm"

type BeuboMiddleware struct {
	DB *gorm.DB
}
