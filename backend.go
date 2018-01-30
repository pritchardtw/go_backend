package main

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

const HASH_DELAY = 5 * time.Second // A 5 second response delay.
const PORT = ":5000"

var shutDown bool = false
var processes = 0
var srv *http.Server
var passwords []string
var passwordsSize int = 0
var totalRequests int = 0
var totalTime time.Duration = 0

func hashPassword(pwd string) string {
	pwdBytes := []byte(pwd)                                      // Convert the string to a byte array to be hashed.
	hashPwd := sha512.Sum512(pwdBytes)                           // Hash the password
	encodedHash := base64.StdEncoding.EncodeToString(hashPwd[:]) // Convert hash to base64 encoded string.
	return encodedHash
}

// Should be used async with go shutdownServer
func shutdownServer() {
	srv.Shutdown(context.Background())
}

func shutdown(w http.ResponseWriter, r *http.Request) {
	shutDown = true
	if processes == 0 {
		go shutdownServer()
	}
}

// If the server is shutting down, do not process any more requests.
func checkShutdownState(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if shutDown {
			http.Error(w, "Server Shutting Down", http.StatusServiceUnavailable)
			return
		}

		processes++
		next.ServeHTTP(w, r)
		processes--
		if shutDown && processes == 0 {
			go shutdownServer()
		}
	})
}

func recordStats(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		totalRequests++
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		totalTime += elapsed
	})
}

func delayHash(pwd string) {
	hashDelay := time.NewTimer(HASH_DELAY)
	<-hashDelay.C
	encodedHash := hashPassword(pwd)           // Hash & encode the password
	passwords = append(passwords, encodedHash) // Respond with the password after 5 seconds
}

func hashRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if r.Method == "GET" {
		idString := r.URL.Path[len("/hash/"):]
		if len(idString) < 1 {
			http.Error(w, "Must provide an index", http.StatusBadRequest)
			return
		}
		id, _ := strconv.Atoi(idString)
		if id < len(passwords) {
			w.Write([]byte(passwords[id]))
		} else {
			http.Error(w, "Invalid index provided", http.StatusBadRequest)
			return
		}
	} else { // Must be a post
		r.ParseForm() // Parse form to retrieve values.
		pwd := r.FormValue("password")
		if len(pwd) == 0 {
			http.Error(w, "Password cannot be blank", http.StatusBadRequest)
			return
		}
		go delayHash(pwd) //Hash after 5 seconds on asynch thread.
		pwdIndex := passwordsSize
		passwordsSize++
		w.Write([]byte(strconv.Itoa(pwdIndex))) // Respond with password index.
	}
}

type Stats struct {
	Total   int
	Average float64
}

func statsRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Total time is time.Duration which is in nS.
	var totalNano int64 = totalTime.Nanoseconds()
	totalMilli := float64(totalNano) / float64(1000000)
	var averageTime float64
	if totalRequests > 0 {
		averageTime = totalMilli / float64(totalRequests)
	} else {
		averageTime = 0
	}
	stats := Stats{totalRequests, averageTime}

	js, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	// Create custom server
	srv = &http.Server{
		Addr: PORT,
	}
	hashHandler := http.HandlerFunc(hashRoute)
	statsHandler := http.HandlerFunc(statsRoute)
	http.Handle("/hash", checkShutdownState(recordStats(hashHandler)))
	http.Handle("/hash/", checkShutdownState(recordStats(hashHandler)))
	http.Handle("/stats", checkShutdownState(statsHandler))
	http.HandleFunc("/shutdown", shutdown)
	srv.ListenAndServe()
}
