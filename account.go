package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Takes an account and puts it in the base64 format.
func (account MCaccount) ToBase64() string {
	strToEncode := ""
	strToEncode += (map[bool]string{true: "type=\"msnc\";", false: "type=\"mspr\";"})[account.ShouldNc]
	lbl := fmt.Sprintf("AUTO-%v.%v", account.Uuid, account.GetType())
	if account.Label == nil {
		account.Label = &lbl
	}
	strToEncode += fmt.Sprintf("lbl=\"%v\";", account.Label)
	strToEncode += fmt.Sprintf("uuid=\"%v\"", account.Uuid)
	return base64.StdEncoding.EncodeToString([]byte(strToEncode))
}

// Takes a base64 string and makes an account object from it.
func (account *MCaccount) FromBase64(str string) error {
	var b []byte
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return errors.New("could not decode base64 string")
	}
	acctStr := string(b)
	parts := strings.Split(acctStr, ";")
	for _, part := range parts {
		kvp := strings.Split(part, "=")
		key := kvp[0]
		val := kvp[1]
		val = strings.ReplaceAll(val, "\"", "")
		switch key {
		case "type":
			switch val {
			case "mcnc":
				account.ShouldNc = true
			default:
				account.ShouldNc = false
			}
		case "lbl":
			account.Label = &val

		case "uuid":
			account.Uuid = val
		}
	}
	return nil
}

// Writes an account to disk. Returns true on success and false on failure.
func (a MCaccount) WriteToDisk(path string) error {
	e := ReadArrayFromDisk(path)
	for i, ea := range e {
		if a.Uuid == ea.Uuid {
			e[i] = a
		}
	}
	return WriteArrayToDisk(path, e)
}

// Writes an array of accounts to disk. Returns nil if successful.
func WriteArrayToDisk(path string, e []MCaccount) error {
	strToWrite := ""
	for _, a := range e {
		strToWrite += a.ToBase64() + "\n"
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	file.Write([]byte(strToWrite))
	return nil
}

// Reads an array of accounts from disk. Returns the array read.
func ReadArrayFromDisk(path string) []MCaccount {
	f := TextToSliceStr(path)
	e := []MCaccount{}
	for _, s := range f {
		var a *MCaccount
		if err := a.FromBase64(s); err != nil {
			continue
		}
		e = append(e, *a)
	}
	return e
}

func NewMCaInstance(IsNC bool, Label ...string) {
	if len(Label) == 0 {
		
	}
}
