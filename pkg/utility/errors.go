package utility

import "log"

// ErrorHandler is a utility function to handle errors so we don't have to write
// if err != nil {
//     // Do something here
// }
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
