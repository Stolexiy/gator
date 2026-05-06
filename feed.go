package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"strings"
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
	if len(cmd.arg) < 1 {
		return fmt.Errorf("not enought arguments, expecting time between requests (duration string, like 1s, 1m, 1h, etc)")
	}

	dur, err := time.ParseDuration(cmd.arg[0])
	if err != nil {
		return err
	}

	fmt.Printf("collectiong feeds every %s\n", cmd, cmd.arg[0])

	tick := time.NewTicker(dur)
	defer tick.Stop()
	for ; ; <-tick.C {
		fmt.Printf("fetching feeds at %s\n", time.Now())
		err = scrapeFeeds(st)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleAddfeed(st *state, cmd command, user database.User) error {
	if len(cmd.arg) < 2 {
		return fmt.Errorf("not enought arguments, expecting feed name and url")
	}

	now := time.Now()
	args := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.arg[0],
		Url:       cmd.arg[1],
		UserID:    user.ID,
	}
	f, err := st.db.CreateFeed(context.Background(), args)
	if err != nil {
		return err
	}

	fmt.Println(f)

	err = followFeed(st, f.Url, user)
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

func handleFollow(st *state, cmd command, user database.User) error {
	if len(cmd.arg) < 1 {
		return fmt.Errorf("not enought arguments, expecting feed url")
	}

	err := followFeed(st, cmd.arg[0], user)
	if err != nil {
		return err
	}

	return nil
}

func handleFollowing(st *state, cmd command, user database.User) error {
	ff, err := st.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}

	for _, f := range ff {
		fmt.Printf("- %s\n", f.FeedName)
	}
	return nil
}

func handleUnfollow(st *state, cmd command, user database.User) error {
	if len(cmd.arg) < 1 {
		return fmt.Errorf("not enought arguments, expecting feed url")
	}

	feed, err := st.db.GetFeedByUrl(context.Background(), cmd.arg[0])
	if err != nil {
		return err
	}

	st.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	fmt.Printf("user %s unfollowed the \"%s\" feed\n", user.Name, feed.Name)

	return nil
}

func handleBrowse(st *state, cmd command) error {
	limit := 2
	if len(cmd.arg) >= 1 {
		val, err := strconv.Atoi(cmd.arg[0])
		if err == nil {
			limit = val
		} else {
			return err
		}
	}

	posts, err := st.db.GetPosts(context.Background(), int32(limit))
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("- %s\n", post.Title)
	}

	return nil
}

func followFeed(st *state, feed_url string, user database.User) error {
	f, err := st.db.GetFeedByUrl(context.Background(), feed_url)
	if err != nil {
		return err
	}

	now := time.Now()
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    f.ID,
	}
	ff, err := st.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("User %s is following %s feed now\n", ff.UserName, ff.FeedName)
	return nil
}

func scrapeFeeds(st *state) error {
	feed, err := st.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = st.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, item := range feedData.Channel.Item {
		pubTime, err := time.Parse("Mon Jan 2 15:04:05 MST 2006", item.PubDate)
		if err != nil {
			pubTime, err = time.Parse("Mon, 02 Jan 2006 15:04:05 +0000", item.PubDate)
		}
		if err != nil {
			return err
		}

		args := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: pubTime,
			FeedID:      feed.ID,
		}
		post, err := st.db.CreatePost(context.Background(), args)
		if err != nil && !strings.Contains(err.Error(), "duplicate key value") {
			return err
		}
		if err == nil {
			fmt.Printf("created new post - %s", post.Title)
		}
	}

	return nil
}
