package main

import (
	"context"
    "database/sql"
	"strings"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ryansurv/go-blog-aggregator/internal/database"
)

func scrapeFeeds(s *state) error {
	// Get next feed to fetch
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	// Mark that feed as fetched
	if err := s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID: feed.ID,
		LastFetchedAt: sql.NullTime{
            Time:  time.Now(),
            Valid: true,
        },
	}); err != nil {
		return err
	}

	// Fetch the feed
	fmt.Printf("Fetching from %s..\n", feed.Url)
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	// Loop through items and save to database
	for _, rssItem := range rssFeed.Channel.Item {
		pubDate, err := parsePubDate(rssItem.PubDate)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if _, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			Title: rssItem.Title,
			Url: rssItem.Link,
			Description: rssItem.Description,
			PublishedAt: pubDate,
			FeedID: feed.ID,
		}); err != nil {
			if !strings.Contains(err.Error(), "unique constraint") {
				fmt.Println(err)
    		}
		}
	}

	fmt.Printf("Fetched %d posts from %s..\n", len(rssFeed.Channel.Item), feed.Url)
	return nil
}

var layouts = []string{
    time.RFC1123Z,              // "Wed, 03 Jul 2019 00:00:00 +0000"
    time.RFC1123,               // "Wed, 03 Jul 2019 00:00:00 UTC"
    time.RFC3339,               // "2019-07-03T00:00:00Z"
    time.RFC822Z,               // "03 Jul 19 00:00 +0000"
}

func parsePubDate(s string) (time.Time, error) {
    var lastErr error
    for _, layout := range layouts {
        t, err := time.Parse(layout, s)
        if err == nil {
            return t, nil
        }
        lastErr = err
    }
    return time.Time{}, lastErr
}