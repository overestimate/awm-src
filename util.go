package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// DroptimeResp represents a response from http://api.star.shopping/droptime/:name as a struct.
type DroptimeResp struct {
	UnixSeconds *int64     `json:"unix,omitempty"` // The UNIX timestamp for when the name is dropping, in milliseconds.
	Droptime    *time.Time // The droptime as a time.Time pointer.
	RfcTime     *string    `json:"droptime,omitempty"` // The droptime, in ISO 8601 format.
	Username    string     `json:"username,omitempty"` // The username requested.
	IsDropping  bool       // False if Error is not nil.
	Error       *string    `json:"error,omitempty"`
}

// FsExists returns true if ``path`` exists, and false otherwise.
func FsExists(path string) bool {
	// check for existing file taken from https://stackoverflow.com/a/12518877/10334831. same idea as runOnceCheck(), but with any path.
	_, err := os.Stat(path)
	return err == nil
}

func FetchDropTime(uname string) DroptimeResp {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://api.star.shopping/droptime/%s", uname), nil)
	if err != nil {
		fmt.Printf("[err] msg => %v\n", err.Error())
		return DroptimeResp{IsDropping: false}
	}
	req.Header.Set("user-agent", "Sniper")
	r, err := http.DefaultClient.Do(req)
	if err != nil || r == nil {
		fmt.Printf("[err] msg => %v\n", err.Error())
		return DroptimeResp{IsDropping: false}
	}
	defer r.Body.Close()
	var dtRes DroptimeResp
	rB, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[err] msg => %v\n", err.Error())
		e := "error with client-side content read"
		return DroptimeResp{IsDropping: false, Username: uname, Error: &e}
	}
	json.Unmarshal(rB, &dtRes)
	dtRes.IsDropping = (dtRes.Error == nil)
	if dtRes.IsDropping {
		dt := time.Unix(*dtRes.UnixSeconds, 0)
		dtRes.Droptime = &dt
	}

	return dtRes
}

// This function will read a file at path and return a []string of all the lines.
func TextToSliceStr(path string) []string {
	file, err := os.Open(path)
	i := 0
	if err == nil {
		var txtSlice []string
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "\n") {
				line = line[:len(line)-1]
			}
			if strings.Contains(line, "\r") {
				line = line[:len(line)-1]
			}
			txtSlice = append(txtSlice, scanner.Text())
			i++
		}
		return txtSlice
	}
	return make([]string, 0)
} // stolen from gosnipe LMAO

// This function removes any account lines starting with '#' and empty lines and returns the rest.
func RemoveAccountComments(acct []string) []string {
	acctNew := []string{}
	for _, l := range acct {
		if len(l) == 0 { // initial blank line check
			continue
		}
		if l[0] == '#' { // remove comments
			continue
		} else { // blank and comments gone
			acctNew = append(acctNew, l)
		}
	}
	return acctNew
}
