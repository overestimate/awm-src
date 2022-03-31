package main

import (
	"errors"
	"io"
	"os"

	"github.com/urfave/cli/v2"
)

func import_cfg(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return errors.New("no path specified")
	}
	f, err := os.Open(c.Args().First())
	if err != nil {
		return errors.New("failed to open file. is the path valid?")
	}

	content := ""
	b := make([]byte, 8192) //8kb buffer
	for {
		r, err := f.Read(b)
		if err != nil {
			if err != io.EOF {
				return errors.New("non-eof error occured while reading file")
			}
			break
		}
		content += string(b[:r])
	}
	f.Close()
	cfg_path := cfg_get_path() + "/config.json"
	err = os.WriteFile(cfg_path, []byte(content), 0755)
	if err != nil {
		return errors.New("could not write to config file. permission denied or file locked")
	}
	return nil
}

func export_cfg(c *cli.Context) error {
	cfg_path := cfg_get_path() + "/config.json"
	if c.Args().Len() < 1 {
		return errors.New("no path specified")
	}
	if !FsExists(cfg_path) {
		return errors.New("config file does not exist, cannot export")
	}
	f, err := os.Open(cfg_path)
	if err != nil {
		return errors.New("failed to open config file. permission denied or file locked")
	}
	defer f.Close()

	content := ""
	b := make([]byte, 8192) //8kb buffer

	for {
		r, err := f.Read(b)
		if err != nil {
			if err != io.EOF {
				return errors.New("non-eof error occured while reading config file")
			}
			break
		}
		content += string(b[:r])
	}

	err = os.WriteFile(c.Args().First(), []byte(content), 0755)
	if err != nil {
		return errors.New("could not write exported config. permission denied or file locked")
	}
	return nil
}

func modify_opt(c *cli.Context) error {
	if c.Args().Len() != 2 {
		return errors.New("exactly 2 arguments must be passed. use double quotes for grouping if needed")
	}
	return nil
}

func get_cd_cmd(c *cli.Context) error {
	return nil
}
