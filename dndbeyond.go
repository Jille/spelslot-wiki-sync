package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type authResponse struct {
	Token string `json:"token"`
	TTL   int    `json:"ttl"`
}

func getDDBAccessToken(cobaltSession string) (string, error) {
	req, err := http.NewRequest("POST", "https://auth-service.dndbeyond.com/v1/cobalt-token", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Cookie", "CobaltSession="+cobaltSession)
	req.Header.Set("User-Agent", "Spelslot campaign sync (jille@quis.cx)")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("/cobalt-token returned HTTP %s", resp.Status)
	}
	defer resp.Body.Close()
	var ar authResponse
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		return "", err
	}
	return ar.Token, nil
}

func fetchCharacter(accessToken string, id int) (CharacterResponse, error) {
	var ret CharacterResponse
	req, err := http.NewRequest("GET", fmt.Sprintf("https://character-service.dndbeyond.com/character/v5/character/%d", id), nil)
	if err != nil {
		return ret, err
	}
	req.Header.Set("User-Agent", "Spelslot campaign sync (jille@quis.cx)")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ret, err
	}
	if resp.StatusCode != 200 {
		return ret, fmt.Errorf("character-service returned HTTP %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return ret, err
	}
	_ = os.WriteFile(fmt.Sprintf("testdata/%d.json", id), body, 0644)
	if err := json.Unmarshal(body, &ret); err != nil {
		return ret, err
	}
	return ret, nil
}
