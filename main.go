package main

import (
	"context"
	"fmt"
	"saidvandeklundert/agrippa/agrippalogger"
	"saidvandeklundert/agrippa/repository"
	"saidvandeklundert/agrippa/systeminteraction"
	"time"
)

func main() {
	fmt.Println("Hello world!")

	log := agrippalogger.GetLogger()
	log.Infow("agrippa starts")

	commandOutput, err := systeminteraction.RunCommand("ls", "-ltr")
	if err != nil {
		log.Error(err)
	} else {
		log.Info(commandOutput)
	}
	response, err := systeminteraction.RunCommandWithTimeout(2*time.Second, "ls", "-ltr", "/tmp")
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Error("Command timed out")
		} else {
			log.Error("Command failed:", err)
		}
		log.Info(response)
	}
	repository.GetRepository()
}
