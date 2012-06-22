package main

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"net/url"
	//"net/http/httputil"
)

const clientLoginUrl = "https://www.google.com/accounts/ClientLogin"

func ClientLogin(client *http.Client, email, password, service string) (string, error) {
	resp, err := client.PostForm(clientLoginUrl, url.Values{
		"accountType": {"HOSTED_OR_GOOGLE"},
		"Email":       {email},
		"Passwd":      {password},
		"service":     {"lh2"},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// check for 200

	r := bufio.NewReader(resp.Body)
	var auth string
	for {
		line := MustReadLine(r)
		if line == nil {
			break
		}

		parts := bytes.SplitN(line, []byte("="), 2)
		if len(parts) != 2 {
			log.Panicf("bad line %q", line)
		}
		if string(parts[0]) == "Auth" {
			auth = string(parts[1])
		}
	}

	return auth, nil
}
