package utils

import (
	"log"
)

// test function for local package importing
func ConsoleLog(value ...any) {
	log.Default().Println(value...)
}
