package main

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

func cfg_get_path() string {
	critical_config_err := func() {
		criticalErrCol := color.New(color.BgRed, color.FgWhite).Render
		fgRed := color.FgRed.Render
		fmt.Printf(criticalErrCol("            [ !! CRITICAL ERROR !! ]             \nUnable to determine user configuartion directory.")+
			color.OpReset.Render("\n")+criticalErrCol("This usually means $HOME or %%APPDATA%% isn't set. ")+
			color.OpReset.Render("\n")+"[%s] awm cannot start. Exiting, code=1.\n\n", fgRed("error"))
		os.Exit(1)
	}
	path, err := os.UserConfigDir()
	if err != nil {
		critical_config_err()
	}
	return path + "/awm"
}

func firstTimeInit() {
	init_critical_err := func() {
		criticalErrCol := color.New(color.BgRed, color.FgWhite).Render
		fgRed := color.FgRed.Render
		fmt.Printf(criticalErrCol("         [ !! CRITICAL ERROR !! ]             \nUnable to create required files for first run.")+
			color.OpReset.Render("\n")+criticalErrCol("This is likely an issue with your filesystem. ")+
			color.OpReset.Render("\n")+"[%s] awm cannot start. Exiting, code=1.\n\n", fgRed("error"))
		os.Exit(1)
	}

	acctsformattingdemo := `# accounts format template for awm.
	#
	# unless you are using a custom interface, this is probably useless. use 'awm account' if you need to edit files (see 'awm account -h')
	#
	# all lines beginning with # are not interpreted as an account.
	#
	# accounts are done as a base64-encoded string, format defined below.
	# form_name_in_snake_case="value";some_number=123;do_not_leave_a="trailing ";but_use_between_forms="okay?"
	#
	# the forms supported are:
	# type (type of account, see below) uuid (unique id for account) label (label for account)
	#
	# the supported account types are:
	# msnc (microsoft oauth login, name change) mspr (microsoft oauth login, no current username)
	#
	# manual bearers are unsupported. if you need to use manual bearers, this program is not for you.
	#
	# the only required fields are uuid and type.
	# if the type field is set to moj, also requires moj_cred
	#
	# the following would be an account of type "msnc" with a label of "personal acct" and a uuid of 01234567-89ab-cdef-0123-456789abcdef, prior to encoding
	# type="msa";label="personal acct";uuid="01234567-89ab-cdef-0123-456789abcdef"
	#
	# encoded to a proper string, this is the following line
	dHlwZT0ibXNhIjtsYWJlbD0icGVyc29uYWwgYWNjdCI7dXVpZD0iMDEyMzQ1NjctODlhYi1jZGVmLTAxMjMtNDU2Nzg5YWJjZGVmIg==
`
	cfg := cfg_get_path()

	if !FsExists(cfg_get_path()) {
		os.Mkdir(cfg_get_path(), 0755)
	}

	_, err := os.Create(cfg + "/.runonce")
	if err != nil {
		init_critical_err()
		return
	}

	f, err := os.Create(cfg + "/formatting_guidelines.txt")
	if err != nil {
		init_critical_err()
		return
	}
	defer f.Close()
	f.WriteString(acctsformattingdemo)

	_, err = os.Create(cfg + "/accounts.txt")
	if err != nil {
		init_critical_err()
		return
	}
}

func runOnceCheck() bool {
	cfg := cfg_get_path()
	// check for existing file taken from https://stackoverflow.com/a/12518877/10334831
	_, err := os.Stat(cfg)
	if err != nil {
		return false
	}
	_, err = os.Stat(cfg + "/.runonce")
	return err == nil
}

func logo_print() {
	fmt.Println(`
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃       ___        ____  __    ┃                v0.1.0 beta                  ┃
┃      / \ \      / |  \/  |   ┃              written by emma.               ┃
┃     / _ \ \ /\ / /| |\/| |   ┃                                             ┃
┃    / ___ \ V  V / | |  | |   ┃                                             ┃
┃   /_/   \_\_/\_/  |_|  |_|   ┃                                             ┃
┃                              ┃                                             ┃                   
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛`)
}
func welcome_msg() {
	logo_print()
	fmt.Println(`
inf > welcome to AWM. you are seeing this message because this is your first
      time running the program. you should make sure to add accounts using
	  'awm account add --label="CUSTOM_TEXT"' and configure your offset.
	  if you need help with offset, view 'awm whatis offset'. set your
	  offset with 'awm conf mod defaults.offset <offset>' and use
	  'awm util offset_base' to get a starter and set it. if you need more 
	  help with usage, see 'awm whatis usage'. good luck, have fun.`)
}
