package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

func loadConfig() {
	// Load from .env into os.Environ()
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, trying OS environment")
	}

	// This allows Viper to read from os.Getenv
	viper.AutomaticEnv()
}

func createTwitterClient() *http.Client {
	config := oauth1.NewConfig(viper.GetString("TWITTER_API_KEY"), viper.GetString("TWITTER_API_SECRET"))
	token := oauth1.NewToken(viper.GetString("TWITTER_ACCESS_TOKEN"), viper.GetString("TWITTER_ACCESS_SECRET"))
	fmt.Println(viper.GetString("TWITTER_API_KEY"))
	return config.Client(oauth1.NoContext, token)
}

func postTweet(client *http.Client, message string) {
	url := "https://api.twitter.com/2/tweets"
	// JSON body
	body := map[string]string{"text": message}
	jsonBody, _ := json.Marshal(body)

	req, err := client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Error sending tweet: %v\n", err)
		return
	}
	defer req.Body.Close()

	respBody, _ := io.ReadAll(req.Body)
	fmt.Println("Status:", req.Status)
	fmt.Println("Response:", string(respBody))
}

func main() {
	loadConfig()
	client := createTwitterClient()
	c := cron.New()

	schedule := viper.GetString("TWITTER_CRON_SCHEDULE")
	if schedule == "" {
		log.Fatal("TWITTER_CRON_SCHEDULE is not set in environment variables")
	}

	println("Cron execution schedule:", viper.GetString("TWITTER_CRON_SCHEDULE"))

	c.AddFunc(viper.GetString("TWITTER_CRON_SCHEDULE"), func() {
		txt := viper.GetString("TWITTER_TWEET_TEXT")
		tweetBase64 := viper.GetString("TWITTER_TWEET_BASE64")
		if tweetBase64 != "" {
			tweetBytes, err := base64.StdEncoding.DecodeString(tweetBase64)
			if err != nil {
				log.Fatal("Error decoding base64 string: ", err)
			}
			var tweet []string
			err = json.Unmarshal(tweetBytes, &tweet)
			if err != nil {
				log.Fatal("Error unmarshalling json: ", err)
				txt = viper.GetString("TWITTER_TWEET_TEXT")
			}
			txt = tweet[rand.Intn(len(tweet))]
		}
		println(txt)
		postTweet(client, txt)
	})

	// Start scheduler
	c.Start()

	fmt.Println("Twitter bot started with cron ⏰")

	// Keep the program alive
	select {}
}
