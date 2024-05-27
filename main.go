package main

import (
	"crypto-lottery-srv/pkg/worker"
	"fmt"
	"os"
	"time"
)

func main() {
	timeString := fmt.Sprintf("%s-%d-%d", time.Now().Format(time.DateOnly), time.Now().Hour(), time.Now().Minute())
	logPath := "log/" + timeString + "/"
	err := os.MkdirAll(logPath, os.ModePerm)
	recordF, err := os.Create(fmt.Sprintf("./%s/record.log", logPath))
	if err != nil {
		panic(err)
	}
	defer recordF.Close()
	goodF, err := os.Create(fmt.Sprintf("./%s/omg-hope-you-are-not-empty.log", logPath))
	if err != nil {
		panic(err)
	}
	defer recordF.Close()
	worker.InitApp(recordF, goodF)
	worker.Lottery()
}
