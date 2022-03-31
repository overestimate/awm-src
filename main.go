package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var app *cli.App

func main() {
	if !runOnceCheck() {
		firstTimeInit()
		welcome_msg()
		return
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	app = &cli.App{
		Commands: []*cli.Command{
			{
				Name: "cfg_manual_recreate",
				Action: func(c *cli.Context) error {
					firstTimeInit()
					return nil
				},
				Hidden: true,
			},
			{
				Name:    "account",
				Aliases: []string{"a", "acc", "acct"},
				Usage:   "account management commands. use subcommands as listed",
				Action: func(c *cli.Context) error {
					fmt.Println("no subcommand specified. see 'awm account -h' for subcommands")
					os.Exit(0)
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:    "add",
						Aliases: []string{"a", "new", "n"},
						Usage:   "adds an account. syntax: 'add --type=\"type\" --label=\"label\"'",
						Action:  stub_fnc,
					},
					{
						Name:    "edit",
						Aliases: []string{"e", "ed"},
						Usage:   "edits an account. syntax: 'edit <uuid> <form> <value>' to set <form> to value. leave <value> unspecified to remove a form",
						Action:  stub_fnc,
					},
					{
						Name:    "list",
						Aliases: []string{"l", "ls"},
						Usage:   "lists accounts.",
						Action:  stub_fnc,
					},
					{
						Name:    "remove",
						Aliases: []string{"r", "rm", "rem", "d", "del", "delete"},
						Usage:   "removed an account. syntax: 'remove <uuid>', 'remove <label>'",
					},
				},
			},
			{
				Name:    "configure",
				Aliases: []string{"c", "cfg", "conf", "config"},
				Usage:   "configuration commands. use subcommands as listed",
				Action: func(c *cli.Context) error {
					fmt.Println("no subcommand specified. see 'awm account -h' for subcommands")
					os.Exit(0)
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:    "import",
						Aliases: []string{"i", "im", "imp"},
						Usage:   "import <path>. imports config from exported file",
						Action:  import_cfg,
					},
					{
						Name:    "export",
						Aliases: []string{"e", "ex", "exp"},
						Usage:   "export <path>. exports config to file",
						Action:  export_cfg,
					},
					{
						Name:    "modify",
						Aliases: []string{"m", "mod", "s", "set"},
						Usage:   "modify <config.option> <value>. modifies config to set key <config.option> to <value>",
						Action:  export_cfg,
					},
				},
			},
			{
				Name:    "snipe_beta",
				Aliases: []string{"snipe", "s"},
				Usage:   "snipe_beta <name>. will do msa auth.",
				Action:  BetaSnipe,
			},
			{
				Name:    "logs",
				Aliases: []string{"log", "l"},
				Usage:   "",
				Action: func(c *cli.Context) error {
					fmt.Println("no subcommand specified. see 'awm logs -h' for subcommands")
					os.Exit(0)
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:    "clear",
						Aliases: []string{"cls", "c", "new", "n"},
						Usage:   "clears current logs and backs up to config dir.",
						Action:  mklog,
					},
				},
			},
		},
		Name:  "awm",
		Usage: "block game sniper with cli and tui for simple usage\n   specify a command as seen below.",
		Action: func(c *cli.Context) error {
			fmt.Println("no command specified. see 'awm -h' for commands")
			os.Exit(0)
			return nil
		},
	}
}

func stub_fnc(c *cli.Context) error {
	fmt.Println("stub function, currently not implemented")
	return nil
}
