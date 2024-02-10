package warmer

import (
	"fmt"
	"github.com/gosuri/uilive"
	"net/http"
	"sync"
	"time"
)

type warmer struct {
	Urls map[string]int64
	mu   *sync.Mutex

	// ttl Set up the needed TTL from current time in seconds
	ttl time.Duration

	numProcs int
}

var client http.Client

func New(ttl time.Duration, numProcs int) Warmer {
	client = http.Client{}

	return &warmer{
		Urls:     make(map[string]int64),
		mu:       &sync.Mutex{},
		ttl:      ttl,
		numProcs: numProcs,
	}
}

// Process Perform check on low latency
func (w *warmer) Process(urls <-chan string, res chan<- struct{}) {
	for url := range urls {
		latencyPeek := time.Second * 10
		latestResponse := time.Second * 20
		for latestResponse > latencyPeek {
			req := prepareUrl(url)

			startReq := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if resp.StatusCode != 200 {
				resp.Body.Close()
				continue
			}
			resp.Body.Close()

			latestResponse = time.Duration(time.Now().Sub(startReq).Seconds())
		}

		w.mu.Lock()
		w.Urls[url] = time.Now().Unix() + int64(w.ttl)
		w.mu.Unlock()
		res <- struct{}{}
		// TODO: check list of urls, if has same - check TTL

		// TODO: load url until he wont be low latency

		// TODO: after this, save and set the TTL before needed warming
	}
}

func (w *warmer) Add(url string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, ok := w.Urls[url]
	if !ok {
		w.Urls[url] = time.Now().Unix() + int64(w.ttl)
	}
}

func (w *warmer) Refresh() {
	writer := uilive.New()
	writer.Start()

	counter := 0
	urls := make(chan string, len(w.Urls))
	res := make(chan struct{})

	for i := 0; i < w.numProcs; i++ {
		go w.Process(urls, res)
	}

	for url, _ := range w.Urls {
		urls <- url
	}

	for range res {
		counter++
		fmt.Fprintf(writer, "Warming up [%d/%d]\n", counter, len(w.Urls))

		if counter == len(w.Urls) {
			close(urls)
			close(res)
		}
	}
}

func prepareUrl(url string) *http.Request {
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "googlebot")

	return req
}
