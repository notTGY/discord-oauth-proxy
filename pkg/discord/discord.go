package discord

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Auth(
	code, client_id, client_secret, redirect_uri string,
) (AuthResponse, error) {
	var authResponse AuthResponse
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirect_uri)

	req_body := strings.NewReader(data.Encode())

	creds := fmt.Sprintf("%s:%s", client_id, client_secret)
	creds = b64.StdEncoding.EncodeToString([]byte(creds))
	creds = fmt.Sprintf("Basic %s", creds)

	req, err := http.NewRequest(
		"POST",
		"https://discord.com/api/v10/oauth2/token",
		req_body,
	)
	if err != nil {
		return authResponse, err
	}
	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)
	req.Header.Set("Authorization", creds)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return authResponse, err
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		err = errors.New(fmt.Sprintf(
			"Code %d: %s", resp.StatusCode, string(body),
		))
		return authResponse, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return authResponse, err
	}

	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return authResponse, err
	}
	return authResponse, nil
}

func Me(
	access_token, token_type string,
) (DiscordUser, error) {
	var discordUser DiscordUser
	creds := fmt.Sprintf("%s %s", token_type, access_token)

	req, err := http.NewRequest(
		"GET",
		"https://discord.com/api/v10/users/@me",
		nil,
	)
	if err != nil {
		return discordUser, err
	}
	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)
	req.Header.Set("Authorization", creds)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return discordUser, err
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		err = errors.New(fmt.Sprintf(
			"Code %d: %s", resp.StatusCode, string(body),
		))
		return discordUser, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return discordUser, err
	}

	err = json.Unmarshal(body, &discordUser)
	if err != nil {
		return discordUser, err
	}
	return discordUser, nil
}
