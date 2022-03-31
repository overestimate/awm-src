package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

func snipe_listcmd(c *cli.Context) error {
	if c.Args().Len() == 0 {
		fmt.Println("[inf] no names specified, exiting")
		return nil
	}

	listener := make(chan string)
	for _, name := range c.Args().Slice() {
		// call concurrent snipe setup twice each
		dtres := FetchDropTime(name)
		if !dtres.IsDropping {
			fmt.Printf("[inf] %v not dropping, omitting.\n", name)
			continue
		}
		// TODO: remove hardcoded 2 request/account/name cap

	}
	skipLog := false
	f, err := os.OpenFile(cfg_get_path()+"/awm.log", os.O_APPEND, os.ModeAppend)
	if err != nil {
		skipLog = true
	}
	for range c.Args().Slice() {
		r1 := <-listener
		if !skipLog {
			f.Write([]byte(r1))
		}
		fmt.Print(r1)
		r2 := <-listener
		if !skipLog {
			f.Write([]byte(r2))
		}
		fmt.Print(r2)
	}
	return nil
}

func concurrent_snipe(a MCaccount, name string, when time.Time, recvr chan string) {
	time.Sleep(time.Until(when.Add(-30 * time.Second)))
	r := a.Namechange(name, when)
	if !r.IsSuccess {
		recvr <- fmt.Sprintf("[inf] failed to snipe %v\n  - response time = %v\n  - status code err = %v\n", name, r.ResponseTimestamp.Local().Format("15:04:05.0000"), (strconv.Itoa(r.Status) + r.Error.Error()))
	} else {
		recvr <- fmt.Sprintf("[inf] sniped %v\n  - response time = %v\n  - status code + err = %v\n", name, r.ResponseTimestamp.Local().Format("15:04:05.0000"), (strconv.Itoa(r.Status) + r.Error.Error()))
	}
}
