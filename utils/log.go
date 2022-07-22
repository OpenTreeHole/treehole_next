package utils

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(
	os.Stdout,
	"",
	log.Lmicroseconds|log.Ldate,
)

func MyLog(model string, action string, objectID, userID int, msg ...string) {
	logger.SetPrefix(fmt.Sprintf("[%v] ", model))
	message := fmt.Sprintf("%vID: %v, UserID: %v. %v, ", model, objectID, userID, action)
	for _, v := range msg {
		message += v
	}
	logger.Println(message)
}
