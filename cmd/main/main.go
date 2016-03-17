package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kr/pretty"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	_ "github.com/lib/pq"
)

func main() {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	insertQry, err := db.Prepare("INSERT INTO twitter_urls VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		log.Fatal(err)
	}

	stream := streamTwitter()

	for {
		switch msg := (<-stream).(type) {
		case *twitter.Tweet:
			go handleTweet(msg, insertQry)
		}
	}
}

func handleTweet(tweet *twitter.Tweet, insertQry *sql.Stmt) {

	for _, url := range tweet.Entities.Urls {

		_, err := insertQry.Exec(url.URL[13:], url.ExpandedURL, tweet.ID, tweet.CreatedAt, time.Now())
		if err != nil {
			pretty.Print(err)
		}
		fmt.Println(url.URL[13:])
	}
}

func streamTwitter() <-chan interface{} {

	var (
		config     = oauth1.NewConfig("NL7gSv640hzRMlSJJN4yAv5h2", "zwHbd7s10L9XED3s3rxs2Lgb9gVQPhc7oKi68T649zr17cBWIS")
		token      = oauth1.NewToken("3092527138-IeVGfw0CIx2QAmSvzoHqTxw3HSq4jgeEkqbgBpB", "h7oIpW5BcJ6adDGqWOfwhGColGLu8fQ6GMkTGidzrdiJr")
		httpClient = config.Client(oauth1.NoContext, token) // OAuth1 http.Client will automatically authorize Requests
		client     = twitter.NewClient(httpClient)          // Twitter Client
	)

	stream, err := client.Streams.User(& twitter.StreamUserParams{
		StallWarnings: twitter.Bool(true),
		With:          "followings",
		Language:      []string{"en"}})
	if err != nil {
		log.Fatal(err)
	}

	return stream.Messages
}
