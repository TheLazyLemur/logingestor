package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type option struct {
	vus      int
	duration time.Duration
}

var ch = make(chan struct{}, 100)
var results = []struct{}{}

func main() {
	go func() {
		for {
			res, ok := <-ch
			if !ok {
				fmt.Println("UHM")
			}

			results = append(results, res)
		}
	}()

	opts := option{
		vus:      1,
		duration: time.Second * 30,
	}
	_ = opts

	pl := `
	{
		"level": "dude",
		"message": "Failed to connect to DB",
		"resourceId": "server-1234",
		"timestamp": "2023-09-15T08:00:00Z",
		"traceId": "abc-xyz-123",
		"spanId": "span-456",
		"commit": "5e5342f",
		"metadata": {
			"parentResourceId": "server-0987"
		}
	}
	`

	wg := sync.WaitGroup{}
	wg.Add(opts.vus)

	for i := 0; i < opts.vus; i++ {
		go func() {
			defer wg.Done()
			do([]byte(pl), opts.duration)
		}()
	}

	wg.Wait()
	fmt.Println(len(results))
}

func do(pl []byte, duration time.Duration) {
	stop := time.After(duration)
	for {
		select {
		case <-stop:
			fmt.Println("Loop stopped after", duration)
			return
		default:
			req, err := http.NewRequest("POST", "http://localhost:3000", bytes.NewBuffer(pl))
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			defer resp.Body.Close()

			fmt.Println("Response Status:", resp.Status)
			fmt.Println("Response Body:", resp.Body)

			if resp.StatusCode != 202 {
				panic("fucked")
			}

			ch <- struct{}{}

			time.Sleep(1 * time.Second)
		}
	}

}
