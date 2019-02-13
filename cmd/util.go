package cmd

import "log"

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
