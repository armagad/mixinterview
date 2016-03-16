package main

import (
	"database/sql"
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

var(
	config = oauth1.NewConfig("NL7gSv640hzRMlSJJN4yAv5h2", "zwHbd7s10L9XED3s3rxs2Lgb9gVQPhc7oKi68T649zr17cBWIS")
	token = oauth1.NewToken("3092527138-IeVGfw0CIx2QAmSvzoHqTxw3HSq4jgeEkqbgBpB", "h7oIpW5BcJ6adDGqWOfwhGColGLu8fQ6GMkTGidzrdiJr")
	httpClient = config.Client(oauth1.NoContext, token) // OAuth1 http.Client will automatically authorize Requests
	client = twitter.NewClient(httpClient) // Twitter Client
)

func main() {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	insertQry, err := db.Prepare("INSERT INTO twitter_urls VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		pretty.Print(err)
		return
	}

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

	// USER (quick test: auth'd user likes a tweet -> event)
	stream, err := client.Streams.User(twitter.StreamUserParams{
		StallWarnings: true,
		With:          "followings",
		Language:      []string{"en"}})
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

