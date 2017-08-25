package main

import (
	"net/http"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	logChan := make(chan string, 1)
	logger(logChan, "test")

	output := <-logChan
	if strings.HasSuffix(output, ": test") == false {
		t.Errorf("logger did not end with :test: \"%s\"", output)
	}
}

func TestAuthenticateInvalidUser(t *testing.T) {
	r, err := http.NewRequest("GET", "/test", strings.NewReader(""))
	if err != nil {
		t.Error("Error", err)
	}
	r.SetBasicAuth("baduser", "goodpassword")
	if authenticate(r, "gooduser", "goodpassword") == true {
		t.Errorf("authenticate returned true on baduser")
	}
}

func TestAuthenticateInvalidPass(t *testing.T) {
	r, err := http.NewRequest("GET", "/test", strings.NewReader(""))
	if err != nil {
		t.Error("Error", err)
	}
	r.SetBasicAuth("gooduser", "badpassword")
	if authenticate(r, "gooduser", "goodpassword") == true {
		t.Errorf("authenticate returned true on badpass")
	}
}

func TestAuthenticateValidUserPass(t *testing.T) {
	r, err := http.NewRequest("GET", "/test", strings.NewReader(""))
	if err != nil {
		t.Error("Error", err)
	}
	r.SetBasicAuth("gooduser", "goodpassword")
	if authenticate(r, "gooduser", "goodpassword") != true {
		t.Errorf("authenticate returned false on valid credentials")
	}
}
