package main

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/gonum/stat/distuv"
)

type hostpool struct {
	Addrs []string
	alpha map[string]float64
	beta  map[string]float64
}

var counter = make(map[string]int)

func (hp *hostpool) Do(req *http.Request) (*http.Response, error) {
	start := time.Now()
	host := hp.selectHost()
	counter[host]++
	req.URL.Host = host
	req.URL.Scheme = "http"

	client := &http.Client{}
	resp, err := client.Do(req)

	var reward float64
	if resp.StatusCode == http.StatusOK {
		reward = 1.0 - float64(time.Since(start).Milliseconds())/1000.0
		reward = math.Max(0, math.Min(1, reward)) // Clamp reward between 0 and 1
	} else {
		reward = 0.0 // No reward for non-OK responses
	}
	hp.update(host, reward)
	return resp, err
}

func newHostPool(addrs []string) *hostpool {
	hp := &hostpool{
		Addrs: addrs,
		alpha: make(map[string]float64),
		beta:  make(map[string]float64),
	}

	for _, addr := range addrs {
		hp.alpha[addr] = 1.0
		hp.beta[addr] = 1.0
	}

	return hp
}

func (hp *hostpool) update(host string, reward float64) {
	hp.alpha[host] += reward
	hp.beta[host] += 1 - reward
}

func (hp *hostpool) selectHost() string {
	var selected string
	minScore := -1.0 // Initialize to a large number

	for _, addr := range hp.Addrs {
		dist := distuv.Beta{Alpha: hp.alpha[addr], Beta: hp.beta[addr]}
		score := dist.Rand()
		if score > minScore {
			minScore = score
			selected = addr
		}
	}
	return selected
}

func startServer(host string, respTime time.Duration) {
	fmt.Println("Starting server on", host)
	h := http.NewServeMux()
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(respTime) // Simulate response time
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, host)
	})
	if err := http.ListenAndServe(host, h); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func main() {
	hosts := []string{":8081", ":8084", ":8082", ":8083"}
	hp := newHostPool(hosts)
	for i, host := range hosts {
		go startServer(host, time.Duration(i*100)*time.Millisecond)
	}
	time.Sleep(5 * time.Second) // Give servers time to start
	for range 1000 {
		resp, err := hp.Do(&http.Request{
			Method: "GET",
			URL:    &url.URL{Path: "/"},
		})
		if err != nil {
			fmt.Println("Error fetching from host:", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Received non-OK response:", resp.StatusCode)
			continue
		}
	}
	fmt.Println("Final counts:", counter)
}
