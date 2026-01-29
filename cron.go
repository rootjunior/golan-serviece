package main

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"
)

func StartCron(ctx context.Context) {
	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("*/10 * * * * *", func() {
		log.Println("⏰ Cron: every 10 seconds")
	})
	if err != nil {
		return
	}

	_, err = c.AddFunc("*/7 * * * * *", func() {
		log.Println("⏰ Cron: fetch posts")

		client := PostClient{
			URL: "https://jsonplaceholder.typicode.com/posts",
		}

		posts, err := client.GetPosts()
		if err != nil {
			log.Printf("❌ Cron error: %v", err)
			return
		}
		for _, post := range posts {
			log.Printf("❌ Post data: %v", post)

		}
		log.Printf("📬 Got %d posts", len(posts))
	})
	if err != nil {
		return
	}

	c.Start()
	log.Println("✅ Cron scheduler started")

	<-ctx.Done()
	log.Println("🛑 Cron shutting down")
	c.Stop()
}
