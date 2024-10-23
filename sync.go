package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"cgt.name/pkg/go-mwclient"
	"cgt.name/pkg/go-mwclient/params"
	"github.com/samber/lo"
	"golang.org/x/text/unicode/norm"
)

var (
	character     = flag.Int("character", 131438467, "First character to sync")
	campaign      = flag.Bool("campaign", false, "Also sync all other characters in the campaign")
	cobaltSession = flag.String("cobalt-session", "", "Cobalt session ID")
	wikiUser      = flag.String("wiki-user", "Spelbot", "Wiki username")
	wikiPass      = flag.String("wiki-pass", "", "Wiki password")

	characterTemplateRe = regexp.MustCompile(`(?ms)\n*^{{Character.*?}}\n*`)
)

func main() {
	flag.Parse()
	w := lo.Must(mwclient.New("https://spelslot.nl/codex/api.php", "Spelbot"))
	lo.Must0(w.Login(*wikiUser, *wikiPass))

	accessToken := lo.Must(getDDBAccessToken(*cobaltSession))

	mainCharacter := lo.Must(fetchCharacter(accessToken, *character))

	if err := syncCharacter(w, mainCharacter); err != nil {
		log.Printf("Failed to sync main character: %v", err)
		os.Exit(1)
	}

	if *campaign {
		ok := true
		for _, c := range mainCharacter.Data.Campaign.Characters {
			if c.CharacterID == mainCharacter.Data.ID {
				continue
			}
			ch, err := fetchCharacter(accessToken, c.CharacterID)
			if err != nil {
				log.Printf("Failed to fetch %d (%s by %s): %v", c.CharacterID, c.CharacterName, c.Username, err)
				ok = false
				continue
			}
			if err := syncCharacter(w, ch); err != nil {
				log.Printf("Failed to sync %d (%s by %s): %v", c.CharacterID, c.CharacterName, c.Username, err)
				ok = false
			}
		}
		if !ok {
			os.Exit(1)
		}
	}
}

func syncCharacter(w *mwclient.Client, ch CharacterResponse) error {
	text := characterToWikiPage(ch)

	if err := w.Edit(params.Values{
		"title":   "Character:" + norm.NFC.String(ch.Data.Name),
		"text":    norm.NFC.String(text),
		"summary": "Sync character from D&D Beyond",
		"bot":     "",
	}); err != nil && err != mwclient.ErrEditNoChange {
		return err
	}

	csrfToken, err := w.GetToken(mwclient.CSRFToken)
	if err != nil {
		panic(fmt.Errorf("unable to obtain csrf token: %w", err))
	}

	if _, err := w.Post(params.Values{
		"action":      "protect",
		"title":       "Character:" + norm.NFC.String(ch.Data.Name),
		"protections": "edit=sysop|move=sysop",
		"reason":      "Any changes would be overwritten by Spelbot",
		"token":       csrfToken,
	}); err != nil {
		return err
	}
	return nil
}
