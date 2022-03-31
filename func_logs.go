package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

func mklog(c *cli.Context) error {
	oops := func() {
		criticalErrCol := color.New(color.BgRed, color.FgWhite).Render
		fgRed := color.FgRed.Render
		fmt.Printf(criticalErrCol(" !!!ERROR!!! \nUnable to create log file.")+
			color.OpReset.Render("\n")+criticalErrCol("You may need to manually create 'awm.log' in the config folder.")+
			color.OpReset.Render("\n")+"[%s] awm cannot continue operation. Exiting, code=1.\n\n", fgRed("error"))
		os.Exit(1)
	}

	path := cfg_get_path() + "/awm.log"
	if _, err := os.Stat(path); err == nil {
		stamp := time.Now().Format(time.ANSIC)
		os.Rename(path, "awm - "+stamp+".log")
	}
	_, err := os.Create(path)
	if err != nil {
		oops()
	}
	return nil
}

func log_string(to_log string) {
	path := cfg_get_path() + "/awm.log"
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("[err] failed to write to awm.log, make sure it exists.")
		return
	}
	defer f.Close()
	f.Write([]byte(to_log))
}
