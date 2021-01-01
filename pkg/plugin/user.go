package plugin

import (
	"gorm.io/gorm"
)

func (p Handler) UserBeforeCreate(tx *gorm.DB) (err error) {
	for _, plug := range p.Plugins {
		symbol, err := plug.Plugin.Lookup("UserBeforeCreate")
		if err != nil {
			continue
		}
		UserBeforeCreateFunc, ok := symbol.(func(*gorm.DB) error)
		if !ok {
			continue
		}
		return UserBeforeCreateFunc(tx)
	}
	return nil
}
