package warmer

import (
	"fmt"
	"github.com/gosuri/uilive"
	"net/http"
	"sync"
	"time"
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

		latestResponse = time.Duration(time.Now().Sub(startReq).Seconds())
	}
	// TODO: check list of urls, if has same - check TTL

	// TODO: load url until he wont be low latency

	// TODO: after this, save and set the TTL before needed warming
}

func (w *warmer) Refresh(urlStream <-chan string) {
	counter := 1

	for url := range urlStream {
		fmt.Fprintf(w.writer, "Warming up. Left process: %d\n", len(urlStream))
		w.Process(url)
		counter++
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
