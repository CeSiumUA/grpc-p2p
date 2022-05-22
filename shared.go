package main

import "log"

func HandleError(err error) {
	if err != nil {
		log.Println("error in application occurred:", err)
		panic(err)
	}
}
