package internal

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/back1ng1/prerender-warming/internal/sitemapper"
	"gitlab.com/back1ng1/prerender-warming/internal/warmer"
)

var sleeping time.Duration

func Run() {
	var url = flag.String("url", "https://example.com/sitemap.xml", "Sitemap that will be parsed.")
	var threads = flag.Int("threads", 2, "Count of threads to warm prerender.")
	flag.Parse()

	if *threads < 1 {
		log.Fatal("Count of threads cannot be less then 1.")
	}

	sleeping = time.Hour * 1

	sitemapParser := sitemapper.New()
	sitemapLinks := make(chan string)
	warm := warmer.New()
	countLinks := 0

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func(countLikes *int) {
		for {
			sitemap := sitemapParser.Get(*url)
			sitemapLinksStream := make(chan string, len(sitemap.URL))
			countLinks = len(sitemap.URL)

			for _, url := range sitemap.URL {
				sitemapLinksStream <- url.Loc
			}

			close(sitemapLinksStream)

			for v := range sitemapLinksStream {
				select {
				case sitemapLinks <- v:
				case <-done:
					return
				}
			}

			warm.ResetWriter()
			log.Printf("Waiting 1 hour for refreshing...\n\n")
			warm.StartWriter()
			<-time.After(time.Millisecond * 100)
			<-time.After(sleeping)
		}
	}(&countLinks)

	for i := 0; i < *threads; i++ {
		go warm.Refresh(sitemapLinks, &countLinks)
	}

	<-done
	fmt.Println("Gracefully shutdown..")
}
