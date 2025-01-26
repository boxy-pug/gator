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

func HandlerAddFeed(s *config.State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("expecting two args, name and url")
	}

	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	currentUser := s.Config.CurrentUserName

	user, err := s.Db.GetUser(context.Background(), currentUser)
	if err != nil {
		return fmt.Errorf("error retrieving current user: %v", err)
	}

	feedID := uuid.New()
	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:     feedID,
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
		Name:   feedName,
		Url:    sql.NullString{String: feedUrl, Valid: true},
	})
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
