package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

type NameResponse struct {
	Name              string
	Status            int
	Contents          *string
	IsSuccess         bool
	Error             error
	ResponseTimestamp *time.Time
}

// Runs a namechange request and returns a NameResponse object.
func (a MCaccount) Namechange(name string, when time.Time) NameResponse {
	errObj := NameResponse{Name: name, Status: 0, Contents: nil, IsSuccess: false, Error: errors.New("generic failure"), ResponseTimestamp: nil}
	successObj := NameResponse{Name: name, Status: 200, Contents: nil, IsSuccess: true, Error: nil, ResponseTimestamp: nil}
	payload := fmt.Sprintf("PUT /minecraft/profile/name/%v HTTP/1.1\r\nUser-Agent: awm/0.1\r\nHost: api.minecraftservices.com\r\nAuthorization: Bearer %v\r\n\r\n", name, a.Bearer)
	pb := []byte(payload)
	resp := make([]byte, 8192)
	testconn, err := tls.Dial("tcp", "api.minecraftservices.com:443", nil)
	if err != nil {
		errObj.Error = errors.New("failed to establish connection to minecraftservices api while testing")
		fmt.Printf("dbg:: %v\n", errObj)
		return errObj
	}
	testconn.Close()
	time.Sleep(time.Until(when) - (time.Second * 15))
	conn, err := tls.Dial("tcp", "api.minecraftservices.com:443", nil)
	if err != nil {
		errObj.Error = errors.New("failed to establish connection to minecraftservices api while testing")
		fmt.Printf("dbg:: %v\n", errObj)
		return errObj
	}
	time.Sleep(time.Until(when))
	conn.Write(pb)
	ts := time.Now()
	nr, err := conn.Read(resp)
	if err != nil {
		fmt.Printf("[!] err dbg | %v\n", err)
		errObj.Error = errors.New("socket read failed")
		fmt.Printf("dbg:: nr=%v  %v\n", nr, errObj)
		return errObj
	}
	fmt.Printf("dbg:: nr=%v  %v\n", nr, resp[:nr])
	status_code, err := strconv.Atoi(string(resp[9:12]))
	if err != nil {
		fmt.Printf("[!] err dbg | %v\n", err)
		errObj.Error = errors.New("invalid status code for protocol http/1.1")
		fmt.Printf("dbg:: %v\n", errObj)
		return errObj
	}
	errObj.Status = status_code
	if status_code != 200 {
		errObj.Error = errors.New("namechange failed. likely status 400")
		fmt.Printf("dbg:: %v\n", errObj)
		return errObj
	}
	successObj.ResponseTimestamp = &ts
	respStr := string(resp)
	contents := strings.SplitN(respStr, "\r\n\r\n", 2)[1]
	successObj.Contents = &contents
	fmt.Printf("dbg:: %v\n", successObj)
	return successObj
}

// Runs a profile creation request and returns a NameResponse object.
func (a MCaccount) CreateProfile(name string, when time.Time) NameResponse {
	errObj := NameResponse{Name: name, Status: 0, Contents: nil, IsSuccess: false, Error: errors.New("generic failure"), ResponseTimestamp: nil}
	successObj := NameResponse{Name: name, Status: 200, Contents: nil, IsSuccess: true, Error: nil, ResponseTimestamp: nil}
	payload := fmt.Sprintf("POST /minecraft/profile HTTP/1.1\r\nHost: api.minecraftservices.com\r\nAuthorization: Bearer %v\r\n\r\n{\"profileName\": \"%v\"}", a.Bearer, name)
	pb := []byte(payload)
	resp := make([]byte, 8192)
	conn, err := tls.Dial("tcp", "api.minecraftservices.com:443", nil)
	if err != nil {
		errObj.Error = errors.New("failed to establish connection to minecraftservices api")
		return errObj
	}
	time.Sleep(time.Until(when))
	conn.Write(pb)
	ts := time.Now()
	conn.Read(resp)
	status_code, err := strconv.Atoi(string(resp[9:12]))
	if err != nil {
		errObj.Error = errors.New("invalid status code for protocol http/1.1")
		return errObj
	}
	errObj.Status = status_code
	if status_code != 200 {
		errObj.Error = errors.New("profile creation failed. likely status 400")
		return errObj
	}
	successObj.ResponseTimestamp = &ts
	respStr := string(resp)
	contents := strings.SplitN(respStr, "\r\n\r\n", 2)[1]
	successObj.Contents = &contents
	return successObj
}

func Snipe(a MCaccount, name string, when time.Time) {
	r := a.Namechange(name, when)
	if !r.IsSuccess {
		fmt.Printf("failed to snipe %v.\nstatus=%v\nerr=%v\n", name, r.Status, r.Error)
	} else {
		fmt.Printf("ez. got %v.\nsunglasses emoji?\n", name)
	}
}

func BetaSnipe(c *cli.Context) error {
	fmt.Println("snipe test run.")
	name := c.Args().First()
	fmt.Printf("name=%v\n", name)
	lbl := "beta-acct"
	account := MCaccount{Label: &lbl, Uuid: "none", ShouldNc: true}
	fmt.Println("getting droptime...")

	dtr := FetchDropTime(name)
	if !dtr.IsDropping {
		return errors.New("name not dropping")
	}
	when := *dtr.Droptime
	fmt.Println("done. authing...")
	err := account.InitAuthFlow()
	if err != nil {
		fmt.Printf("auth failed, err=%v\n", err)
	}
	fmt.Println("authed. waiting...")
	fmt.Printf("dbg:: a.bearer=%v\n", account.Bearer)
	go Snipe(account, name, when.Add(15*time.Millisecond))
	go Snipe(account, name, when.Add(24*time.Millisecond))
	IdleTimer(when)
	//account.FakeNC(name)
	return nil
}

func IdleTimer(when time.Time) {
	fmt.Printf("waiting %.0f seconds                               \r", time.Until(when).Seconds())
	for {
		time.Sleep(2 * time.Second)
		fmt.Printf("waiting %.0f seconds                               \r", time.Until(when).Seconds())
		if time.Until(when) < (10 * time.Second) {
			break
		}
	}
	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("waiting %.0f seconds                               \r", time.Until(when).Seconds())
		if time.Until(when) < (0 * time.Second) {
			break
		}
	}
	fmt.Println("waiting done.                               ")
	time.Sleep(3 * time.Second)
}

// Run a Namechange request without authorization.
func (a MCaccount) FakeNC(name string) {
	when := time.Now()
	errObj := NameResponse{Name: name, Status: 0, Contents: nil, IsSuccess: false, Error: errors.New("generic failure"), ResponseTimestamp: nil}
	successObj := NameResponse{Name: name, Status: 200, Contents: nil, IsSuccess: true, Error: nil, ResponseTimestamp: nil}
	payload := fmt.Sprintf("PUT /minecraft/profile/name/%v HTTP/1.1\r\nHost: api.minecraftservices.com\r\nAuthorization: Bearer %v\r\n\r\n", name, a.Bearer)
	pb := []byte(payload)
	resp := make([]byte, 8192)
	conn, err := tls.Dial("tcp", "api.minecraftservices.com:443", nil)
	if err != nil {
		errObj.Error = errors.New("failed to establish connection to minecraftservices api")
		fmt.Printf("dbg:: %v\n", errObj)
		return
	}
	time.Sleep(time.Until(when))
	conn.Write(pb)
	ts := time.Now()
	conn.Read(resp)
	status_code, err := strconv.Atoi(string(resp[9:12]))
	if err != nil {
		fmt.Printf("[!] err dbg | %v\n", err)
		errObj.Error = errors.New("invalid status code for protocol http/1.1")
		fmt.Printf("dbg:: %v\n", errObj)
		return
	}
	errObj.Status = status_code
	if status_code != 200 {
		errObj.Error = errors.New("namechange failed. likely status 400")
		fmt.Printf("dbg:: %v\n", errObj)
		return
	}
	successObj.ResponseTimestamp = &ts
	respStr := string(resp)
	contents := strings.SplitN(respStr, "\r\n\r\n", 2)[1]
	successObj.Contents = &contents
	fmt.Printf("dbg:: %v\n", successObj)
	return
}
