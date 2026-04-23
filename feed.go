package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/stolexiy/gator/internal/database"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "gator")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	feed := &RSSFeed{}
	err = xml.Unmarshal(data, feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	return feed, nil
}

func handleAgg(st *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func handleAddfeed(st *state, cmd command) error {
	if len(cmd.arg) < 2 {
		return fmt.Errorf("not enought arguments, expecting feed name and url")
	}

	now := time.Now()

	u, err := st.db.GetUser(context.Background(), st.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	args := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.arg[0],
		Url:       cmd.arg[1],
		UserID:    u.ID,
	}
	f, err := st.db.CreateFeed(context.Background(), args)
	if err != nil {
		return err
	}

	fmt.Println(f)

	err = followFeed(st, f.Url)
	if err != nil {
		return err
	}

	return nil
}

func handleFeeds(st *state, cmd command) error {
	feeds, err := st.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, f := range feeds {
		u, err := st.db.GetUserById(context.Background(), f.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("- %s, URL: %s, Created by: %s\n", f.Name, f.Url, u.Name)
	}
	return nil
}

func handleFollow(st *state, cmd command) error {
	if len(cmd.arg) < 1 {
		return fmt.Errorf("not enought arguments, expecting feed url")
	}

	err := followFeed(st, cmd.arg[0])
	if err != nil {
		return err
	}

	return nil
}

func handleFollowing(st *state, cmd command) error {
	ff, err := st.db.GetFeedFollowsForUser(context.Background(), st.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	for _, f := range ff {
		fmt.Printf("- %s\n", f.FeedName)
	}
	return nil
}

func followFeed(st *state, feed_url string) error {
	u, err := st.db.GetUser(context.Background(), st.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	f, err := st.db.GetFeedByUrl(context.Background(), feed_url)
	if err != nil {
		return err
	}

	now := time.Now()
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    u.ID,
		FeedID:    f.ID,
	}
	ff, err := st.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("User %s is following %s feed now\n", ff.UserName, ff.FeedName)
	return nil
}
