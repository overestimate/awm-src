package main

// parts of this code licensed under MIT license.
// created by contributors to Kqzz/mcgo and myself.

import (
	"time"
)

// xbl-related types here

type xBLSignInBody struct {
	Properties struct {
		Authmethod string `json:"AuthMethod"`
		Sitename   string `json:"SiteName"`
		Rpsticket  string `json:"RpsTicket"`
	} `json:"Properties"`
	Relyingparty string `json:"RelyingParty"`
	Tokentype    string `json:"TokenType"`
}

type xBLSignInResp struct {
	Issueinstant  time.Time `json:"IssueInstant"`
	Notafter      time.Time `json:"NotAfter"`
	Token         string    `json:"Token"`
	Displayclaims struct {
		Xui []struct {
			Uhs string `json:"uhs"`
		} `json:"xui"`
	} `json:"DisplayClaims"`
}

type xSTSPostBody struct {
	Properties struct {
		Sandboxid  string   `json:"SandboxId"`
		Usertokens []string `json:"UserTokens"`
	} `json:"Properties"`
	Relyingparty string `json:"RelyingParty"`
	Tokentype    string `json:"TokenType"`
}

type xSTSAuthorizeResponse struct {
	Issueinstant  time.Time `json:"IssueInstant"`
	Notafter      time.Time `json:"NotAfter"`
	Token         string    `json:"Token"`
	Displayclaims struct {
		Xui []struct {
			Uhs string `json:"uhs"`
		} `json:"xui"`
	} `json:"DisplayClaims"`
}

type xSTSAuthorizeResponseFail struct {
	Identity string `json:"Identity"`
	Xerr     int64  `json:"XErr"`
	Message  string `json:"Message"`
	Redirect string `json:"Redirect"`
}

type msGetMojangbearerBody struct {
	Identitytoken       string `json:"identityToken"`
	Ensurelegacyenabled bool   `json:"ensureLegacyEnabled"`
}

type msGetMojangBearerResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
	Foci         string `json:"foci"`
}

type msDeviceInitResponse struct {
	Message    string `json:"message"`
	Interval   int    `json:"interval"`
	DeviceCode string `json:"device_code"`
}

type msErrorPollResponse struct {
	Error string `json:"error"`
}

type msSuccessPollResponse struct {
	AccessToken string `json:"access_token"`
}

// account-related structs

// This type holds a Minecraft account.
type MCaccount struct {
	Bearer   *string // The account's bearer token. If nil, has not been authenticated.
	ShouldNc bool    // A boolean that determines if the account should namechange or not.
	Uuid     string  // UUID for internal usage. Not related to the UUID from MC's API.
	Label    *string // A label for the account. If omitted, automatically generates from type and UUID in the format "AUTO-UUID.TYPE"
}

func (a MCaccount) GetType() string {
	if a.ShouldNc {
		return "mspr"
	} else {
		return "msnc"
	}
}

// configuration structs

// ConfigOptions is an object that represents config.json. All options are nullable unless noted otherwise.
type ConfigOptions struct {
}
