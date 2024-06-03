package warmer

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gosuri/uilive"
)

type warmer struct {
	client http.Client
	mu     *sync.Mutex
	writer *uilive.Writer
}

func New() Warmer {
	writer := uilive.New()
	writer.Start()

	return &warmer{
		client: http.Client{},
		mu:     &sync.Mutex{},
		writer: writer,
	}
}

func (w *warmer) ResetWriter() {
	w.writer.Stop()
}

func (w *warmer) StartWriter() {
	w.writer.Start()
}

// Process Perform check on low latency
func (w *warmer) Process(url string) {
	latencyPeek := time.Second * 10
	latestResponse := time.Second * 20

	for latestResponse > latencyPeek {
		req := prepareUrl(url)

		startReq := time.Now()
		resp, err := w.client.Do(req)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		latestResponse = time.Duration(time.Since(startReq).Seconds())
	}
}

func (w *warmer) Refresh(urlStream <-chan string, countLinks *int) {
	for url := range urlStream {
		w.mu.Lock()
		*countLinks--
		w.mu.Unlock()

		fmt.Fprintf(w.writer, "Warming up. Left process: %d\n", *countLinks)
		w.Process(url)

		<-time.After(time.Millisecond * 10)
	}
}

func prepareUrl(url string) *http.Request {
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "googlebot")
	req.Header.Set("x-prerender-warmer", "googlebot")

	return req
}
