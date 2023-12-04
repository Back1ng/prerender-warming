package internal

import (
	"fmt"
	"gitlab.com/back1ng1/prerender-warming/internal/sitemapper"
	"gitlab.com/back1ng1/prerender-warming/internal/warmer"
	"gitlab.com/back1ng1/prerender-warming/pkg/conf"
	"time"
)

var sleeping time.Duration

func Run() {
	sleeping = time.Hour * 1
	warm := warmer.New(sleeping)

	configuration := conf.New()

	for {
		sitemapParser := sitemapper.New()
		sitemap := sitemapParser.Get(configuration.Url)

		for _, url := range sitemap.URL {
			warm.Add(url.Loc)
		}

		warm.Refresh()
		fmt.Println("Waiting 1 hour for refreshing...")
		time.Sleep(sleeping)
	}
}
