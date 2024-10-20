package main

import (
	"fmt"
	"saidvandeklundert/agrippa/agrippalogger"
	"saidvandeklundert/agrippa/repository"
)

func main() {
	fmt.Println("Hello world!")

	logger := agrippalogger.GetLogger()
	logger.Infow("agrippa starts")
	repository.GetRepository()

}
