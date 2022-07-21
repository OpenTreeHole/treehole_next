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
	message := fmt.Sprintf("%v, %vID: %v, UserID: %v", action, model, objectID, userID)
	if len(msg) > 0 {
		message += msg[0]
	}
	logger.Println(message)
}
