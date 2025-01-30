package commands

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/boxy-pug/gator/internal/config"
	"github.com/boxy-pug/gator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func UnescapeHTML(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read resp body: %w", err)
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling xml: %w", err)
	}

	UnescapeHTML(&feed)

	return &feed, nil

}

func HandlerAddFeed(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("expecting two args, name and url")
	}

	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	feedID := uuid.New()
	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:     feedID,
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
		Name:   feedName,
		Url:    sql.NullString{String: feedUrl, Valid: true},
	})
	if err != nil {
		fmt.Errorf("error creating feed: %w", err)
	}

	// Call HandlerFollow to follow the newly added feed
	followCmd := Command{Name: "follow", Args: []string{feedUrl}}
	err = HandlerFollow(s, followCmd, user)
	if err != nil {
		return fmt.Errorf("error following feed: %w", err)
	}

	// Print the details of the new feed
	fmt.Printf("Feed added successfully:\nID: %s\nName: %s\nURL: %s\n", feed.ID, feed.Name, feed.Url)
	log.Printf("Feed added: ID=%s, Name=%s, URL=%s, UserID=%s\n", feed.ID, feed.Name, feed.Url, user.ID)

	return nil

}
func HandlerFeeds(s *config.State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds: %w", err)
	}

	for _, feed := range feeds {
		userName, err := s.Db.GetUserFromId(context.Background(), feed.UserID.UUID)
		if err != nil {
			return fmt.Errorf("error fetching username from uuid: %w", err)
		}
		fmt.Printf("%s\n", feed.Name)
		fmt.Printf("%v\n", feed.Url.String)
		fmt.Printf("%s\n", userName)
	}
	return nil
}

func HandlerFollow(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("expecting url argument")
	}

	feedUrl := cmd.Args[0]

	feedId, err := s.Db.GetFeedByUrl(context.Background(), sql.NullString{String: feedUrl, Valid: true})
	if err != nil {
		return fmt.Errorf("error getting feed by url: %w", err)
	}

	feedFollowId := uuid.New()
	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        feedFollowId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		FeedID:    uuid.NullUUID{UUID: feedId, Valid: true},
	})
	if err != nil {
		fmt.Errorf("error creating feed follow: %w", err)
	}

	fmt.Printf("Successfully followed feed:\nUser: %s\nFeed: %s\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

//Add a follow command.
// It takes a single url argument and creates a new feed follow record for the current user.
//It should print the name of the feed and the current user once the record is created
//(which the query we just made should support). You'll need a query to look up feeds by URL.

func HandlerFollowing(s *config.State, cmd Command, user database.User) error {

	feeds, err := s.Db.GetFeedFollowsForUser(context.Background(), uuid.NullUUID{UUID: user.ID, Valid: true})
	if err != nil {
		fmt.Errorf("error retreiving user feeds: %w", err)
	}

	for _, feed := range feeds {
		fmt.Printf("%s\n", feed.FeedName)
	}

	return nil

}

func HandlerDeleteFeed(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		fmt.Errorf("please provide url for unfollowing")
	}
	url := cmd.Args[0]

	err := s.Db.DeleteFollowFeed(context.Background(), database.DeleteFollowFeedParams{
		Url:    sql.NullString{String: url, Valid: true},
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("could not delete follow feed: %w", err)
	}
	fmt.Printf("Successfully deleted feed %s from follow list\n", url)

	return nil
}

/*
Write an aggregation function, I called mine scrapeFeeds. It should:

    Get the next feed to fetch from the DB.
    Mark it as fetched.
    Fetch the feed using the URL (we already wrote this function)
    Iterate over the items in the feed and print their titles to the console.
*/

func ScrapeFeeds(s *config.State) error {
	nextFeed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("could not fetch next feed; %w", err)
	}
	nextFeed.LastFetchedAt = sql.NullTime{Time: time.Now(), Valid: true}

	fetchedFeed, err := FetchFeed(context.Background(), nextFeed.Url.String)
	if err != nil {
		return fmt.Errorf("could not fetch feed content:%w", err)
	}

	for _, item := range fetchedFeed.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("error fetching published date for %s, using current time instead", item.Title)
			publishedAt = time.Now()
		}
		postID := uuid.New()
		post := database.CreatePostParams{
			ID:          postID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: sql.NullTime{Time: publishedAt, Valid: true},
			FeedID:      uuid.NullUUID{UUID: nextFeed.ID, Valid: true},
		}

		_, err = s.Db.CreatePost(context.Background(), post)
		if err != nil {
			log.Printf("error saving post %v", item.Title)
		}
	}
	return nil
}

func HandlerBrowse(s *config.State, cmd Command, user database.User) error {
	limit := 2
	if len(cmd.Args) > 0 {
		l, err := strconv.Atoi(cmd.Args[0])
		if err == nil {
			limit = l
		}

	}
	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error retrieving posts for user: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("Post Title: %s, URL: %s, Published At: %v\n", post.Title, post.Url, post.PublishedAt)
	}

	return nil
}
