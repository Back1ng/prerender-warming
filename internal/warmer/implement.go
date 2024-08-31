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

	writer   *uilive.Writer
	messages chan string
}

func New() Warmer {
	// todo move writer to another package
	writer := uilive.New()
	writer.Start()

	messages := make(chan string)
	go func() {
		for m := range messages {
			fmt.Fprint(writer, m)
			<-time.After(time.Millisecond * 10)
		}
	}()

	return &warmer{
		client:   http.Client{},
		mu:       &sync.Mutex{},
		writer:   writer,
		messages: messages,
	}
}

func (w *warmer) Print(message string) {
	w.messages <- message
}

func (w *warmer) ResetWriter() {
	w.writer.Stop()
}

func (w *warmer) StartWriter() {
	w.writer.Start()
}

// Process Perform check on low latency
func (w *warmer) Process(url string) {
	latencyPeek := 2
	latestResponse := 20

	for latestResponse > latencyPeek {
		req := prepareUrl(url)

		startReq := time.Now()
		resp, err := w.client.Do(req)
		if err != nil {
			fmt.Println(err)
			continue
		}

		resp.Body.Close()

		if resp.StatusCode != 200 {
			<-time.After(time.Second * 20)
			continue
		}

		if int(time.Now().Sub(startReq).Seconds()) > 15 {
			<-time.After(time.Second * 20)
		}

		latestResponse = int(time.Now().Sub(startReq).Seconds())
	}
}

func (w *warmer) Refresh(urlStream <-chan string, countLinks *int) {
	for url := range urlStream {
		w.mu.Lock()
		*countLinks--
		w.mu.Unlock()

		w.messages <- fmt.Sprintf("Warming up. Left process: %d\n", *countLinks)
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
