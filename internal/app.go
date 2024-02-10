package internal

import (
	"fmt"
	"gitlab.com/back1ng1/prerender-warming/internal/sitemapper"
	"gitlab.com/back1ng1/prerender-warming/internal/warmer"
	"gitlab.com/back1ng1/prerender-warming/pkg/conf"
	"time"
)

func Run() {
	configuration := conf.New()

	warm := warmer.New(time.Hour, configuration.NumProcs)

	for {
		sitemapParser := sitemapper.New()
		sitemap := sitemapParser.Get(configuration.Url)

		for _, url := range sitemap.URL {
			warm.Add(url.Loc)
		}

		warm.Refresh()
		fmt.Println("Waiting 1 hour for refreshing...")
		time.Sleep(time.Hour)
	}
}
