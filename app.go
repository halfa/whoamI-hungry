package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	// "github.com/pkg/profile"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

var port string
var running int

func init() {
	flag.StringVar(&port, "port", "80", "give me a port number")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	// defer profile.Start().Stop()
	flag.Parse()
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/bench", benchHandler)
	http.HandleFunc("/burn", burnHandler)
	http.HandleFunc("/", whoamI)
	http.HandleFunc("/api", api)
	fmt.Println("Starting up on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func printBinary(s []byte) {
	fmt.Printf("Received b:")
	for n := 0; n < len(s); n++ {
		fmt.Printf("%d,", s[n])
	}
	fmt.Printf("\n")
}
func benchHandler(w http.ResponseWriter, r *http.Request) {
	// body := "Hello World\n"
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/plain")
	// fmt.Fprint(w, body)
}
func echoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		printBinary(p)
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			return
		}
	}
}

func burnHandler(w http.ResponseWriter, req *http.Request) {
	go pi(100000) // TODO iterations as GET parameter
	whoamI(w, req)
}

// Pi launches n goroutines to compute an
// approximation of pi.
func pi(n int) float64 {
	ch := make(chan float64)
	for k := 0; k <= n; k++ {
		go term(ch, float64(k))
	}
	f := 0.0
	for k := 0; k <= n; k++ {
		f += <-ch
	}
	return f
}

func term(ch chan float64, k float64) {
	ch <- 4 * math.Pow(-1, k) / (2*k + 1)
}

func whoamI(w http.ResponseWriter, req *http.Request) {
	u, _ := url.Parse(req.URL.String())
	queryParams := u.Query()
	wait := queryParams.Get("wait")
	if len(wait) > 0 {
		duration, err := time.ParseDuration(wait)
		if err == nil {
			time.Sleep(duration)
		}
	}
	hostname, _ := os.Hostname()
	fmt.Fprintln(w, "Hostname:", hostname)
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			fmt.Fprintln(w, "IP:", ip)
		}
	}
	req.Write(w)
}

func api(w http.ResponseWriter, req *http.Request) {
	hostname, _ := os.Hostname()
	data := struct {
		Hostname string      `json:"hostname,omitempty"`
		IP       []string    `json:"ip,omitempty"`
		Headers  http.Header `json:"headers,omitempty"`
	}{
		hostname,
		[]string{},
		req.Header,
	}

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			data.IP = append(data.IP, ip.String())
		}
	}
	json.NewEncoder(w).Encode(data)
}
