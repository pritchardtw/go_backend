package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestHashPassword(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{
			"angryMonkey",
			"ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==",
		},
	}
	for _, c := range cases {
		got := hashPassword(c.in)
		if got != c.want {
			t.Errorf("hashPassword(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestHashRoute(t *testing.T) {
	//Test successful password post
	resp, err := http.PostForm("http://localhost:5000/hash",
		url.Values{"password": {"angryMonkey"}})

	if err != nil {
		t.Errorf("hashRoute test failed. Err: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			t.Errorf("hashRoute test failed. Unable to read response body. Err: %v", err)
		}
		bodyString := string(bodyBytes)
		if len(bodyString) < 1 {
			t.Errorf("hashRoute test failed. Did not return index of password. Err: %v", err)
		}

		i, err := strconv.Atoi(bodyString)
		if err != nil {
			t.Errorf("hashRoute test failed. Unable to parse returned password index. Err: %v", err)
		} else if i < 0 {
			t.Errorf("hashRoute test failed. Returned negative password index. Index: %v", i)
		}

	} else {
		t.Errorf("hashRoute test failed. Expected 200, received %v", resp.StatusCode)
	}
	resp.Body.Close()

	// Wait 5 seconds to continue for password to be saved
	hashDelay := time.NewTimer(HASH_DELAY)
	<-hashDelay.C

	//Test invalid password post
	resp, err = http.PostForm("http://localhost:5000/hash",
		url.Values{"password": {""}})

	if err != nil {
		t.Errorf("hashRoute test failed. Err: %v", err)
	}

	if resp.StatusCode == http.StatusBadRequest {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			t.Errorf("hashRoute test failed. Unable to read response body. Err: %v", err)
		}
		bodyString := string(bodyBytes)
		if bodyString != "Password cannot be blank\n" {
			t.Errorf("hashRoute test failed. Incorrect Error message. Incorrect message: %v", bodyString)
		}
	} else {
		t.Errorf("hashRoute test failed. Expected 400, received %v", resp.StatusCode)
	}
	resp.Body.Close()

	//Test GET password
	resp, err = http.Get("http://localhost:5000/hash/0")
	if err != nil {
		t.Errorf("hashRoute test failed. Err: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			t.Errorf("hashRoute test failed. Unable to read response body. Err: %v", err)
		}
		bodyString := string(bodyBytes)
		want := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
		if bodyString != want {
			t.Errorf("Expected PWD to be %v, got %v", want, bodyString)
		}
	} else {
		t.Errorf("hashRoute test failed. Expected 200, received %v", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestStatsRoute(t *testing.T) {
	resp, err := http.Get("http://localhost:5000/stats")
	if err != nil {
		t.Errorf("statsRoute test failed. Err: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			t.Errorf("statsRoute test failed. Unable to read response body. Err: %v", err)
		}
		var stats Stats
		json.Unmarshal(bodyBytes, &stats)
		if stats.Total < 1 {
			t.Errorf("statsRoute test failed. Inaccurate stats. Err: %v", err)
		}
	} else {
		t.Errorf("statsRoute test failed. Expected 200, received %v", resp.StatusCode)
	}
	resp.Body.Close()
}
