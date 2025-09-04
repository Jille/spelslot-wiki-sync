package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"

	charsync "github.com/Jille/spelslot-wiki-sync"
	"github.com/samber/lo"
)

var (
	characters    = flag.String("characters", "131438467", "Comma separated list of characters to sync")
	campaign      = flag.Bool("campaign", false, "Also sync all other characters in their campaigns")
	cobaltSession = flag.String("cobalt-session", "", "Cobalt session ID")
)

func main() {
	flag.Parse()

	accessToken := lo.Must(getDDBAccessToken(*cobaltSession))

	var todo []int
	for _, n := range strings.Split(*characters, ",") {
		id, err := strconv.Atoi(n)
		if err != nil {
			log.Fatalf("--characters got invalid %q: %v", n, err)
		}
		todo = append(todo, id)
	}

	ret := map[int]charsync.CharacterInfo{}

	ok := true
	for len(todo) > 0 {
		id := todo[len(todo)-1]
		todo = todo[:len(todo)-1]

		if _, seen := ret[id]; seen {
			continue
		}

		ch, err := fetchCharacter(accessToken, id)
		if err != nil {
			log.Printf("Failed to fetch %d: %v", id, err)
			ok = false
			continue
		}

		if *campaign {
			if len(ch.Data.Campaign.Characters) == 0 {
				log.Printf("Warning: %s (%d) does not seem to be in a campaign", ch.Data.Name, id)
			}
			for _, c := range ch.Data.Campaign.Characters {
				todo = append(todo, c.CharacterID)
			}
		}

		ret[id] = ch.Data
	}

	b, err := json.Marshal(ret)
	if err != nil {
		log.Fatalf("json.Marshal: %v", err)
	}
	if err := os.WriteFile("characters.json", b, 0644); err != nil {
		log.Fatalf("os.WriteFile: %v", err)
	}

	if !ok {
		os.Exit(1)
	}
}
