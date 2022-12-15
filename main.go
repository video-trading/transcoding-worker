package main

import (
	"github.com/spf13/viper"
	"log"
	"video_transcoding_worker/worker"
)

func Init() {
	err := viper.BindEnv("endpoint", "endpoint")
	if err != nil {
		return
	}
	err = viper.BindEnv("jwt_token", "jwt_token")
	if err != nil {
		return
	}
	err = viper.BindEnv("message_queue", "message_queue")
	if err != nil {
		return
	}
	viper.AutomaticEnv()

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file, %s", err)
	}
}

func main() {
	Init()
	endpoint := viper.GetString("endpoint")
	jwtToken := viper.GetString("jwt_token")
	messageQueue := viper.GetString("message_queue")

	if len(endpoint) == 0 || len(jwtToken) == 0 || len(messageQueue) == 0 {
		log.Fatal("Please provide all required environment variables")
	}

	worker.Setup(endpoint, jwtToken, messageQueue)
}
