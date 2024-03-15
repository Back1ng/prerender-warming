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
	sitemapLinksStream := make(chan string, 1)
	warm := warmer.New()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			sitemap := sitemapParser.Get(*url)

			warm.ResetWriter()
			log.Printf("Waiting 1 hour for refreshing...\n\n")
			warm.StartWriter()

			for _, url := range sitemap.URL {
				sitemapLinksStream <- url.Loc
			}

			<-time.After(sleeping)
		}
	}()

	for i := 0; i < *threads; i++ {
		go warm.Refresh(sitemapLinksStream)
	}

	<-done
	fmt.Println("Gracefully shutdown..")
	close(sitemapLinksStream)
}
