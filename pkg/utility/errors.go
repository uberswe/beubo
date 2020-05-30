package utility

import "log"

func ErrorHandler(e error, fatal bool) bool {
	if e != nil {
		if fatal {
			log.Fatal(e)
		} else {
			log.Print(e)
		}
		return true
	}
	return false
}
