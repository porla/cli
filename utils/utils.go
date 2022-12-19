package utils

import (
	"log"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Contains(s []int, test int) bool {
	for _, v := range s {
		if v == test {
			return true
		}
	}

	return false
}
