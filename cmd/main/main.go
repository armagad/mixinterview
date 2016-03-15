package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kr/pretty"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	_ "github.com/lib/pq"
)

func main() {

	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)

	consumerKey := flags.String("consumer-key", "NL7gSv640hzRMlSJJN4yAv5h2", "")
	consumerSecret := flags.String("consumer-secret", "zwHbd7s10L9XED3s3rxs2Lgb9gVQPhc7oKi68T649zr17cBWIS", "")
	accessToken := flags.String("access-token", "3092527138-IeVGfw0CIx2QAmSvzoHqTxw3HSq4jgeEkqbgBpB", "")
	accessSecret := flags.String("access-secret", "h7oIpW5BcJ6adDGqWOfwhGColGLu8fQ6GMkTGidzrdiJr", "")
	flags.Parse(os.Args[1:])

	fmt.Println(os.Getenv("DATABASE_URL"))
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// DROP TABLE twitter_urls;
	// CREATE TABLE twitter_urls (url varchar, expanded_url varchar, created_at timestamp, id bigint, ts timestamp default current_timestamp, _ID serial PRIMARY KEY);
	insertQry, err := db.Prepare("INSERT INTO twitter_urls VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		pretty.Print(err)
		return
	}

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {

		for _, url := range tweet.Entities.Urls {

			_, err := insertQry.Exec(url.URL[13:], url.ExpandedURL, tweet.ID, tweet.CreatedAt, time.Now())
			if err != nil {
				pretty.Print(err)
			}
			fmt.Println(url.URL[13:])
		}
	}

	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		fmt.Printf("Event: %#v\n", event)
	}

	// USER (quick test: auth'd user likes a tweet -> event)
	stream, err := client.Streams.User(&twitter.StreamUserParams{
		StallWarnings: twitter.Bool(true),
		With:          "followings",
		Language:      []string{"en"},
	},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Stop()

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
}
