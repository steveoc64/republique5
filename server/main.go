package main

import "github.com/sirupsen/logrus"

const version = "0.1"

func main() {
	log := logrus.New()
	log.WithField("version", version).Println("Republique 5.0 Server")
	log.SetFormatter(&logrus.JSONFormatter{})
	log.Println("Starting")
}
