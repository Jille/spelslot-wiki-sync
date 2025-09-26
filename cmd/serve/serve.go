package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/Jille/convreq"
	"github.com/Jille/convreq/respond"
	charsync "github.com/Jille/spelslot-wiki-sync"
)

var (
	datafile = flag.String("datafile", "characters.json", "Path to the JSON file")
	port     = flag.Int("port", 8080, "HTTP port number")
)

func main() {
	flag.Parse()
	http.Handle("/full", convreq.Wrap(handleFull))
	http.Handle("/summary", convreq.Wrap(handleSummary))
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func handleFull(ctx context.Context, r *http.Request) convreq.HttpResponse {
	return respond.ServeFile(*datafile)
}

func handleSummary(ctx context.Context, r *http.Request) convreq.HttpResponse {
	b, err := os.ReadFile(*datafile)
	if err != nil {
		return respond.Error(err)
	}
	var data map[int]charsync.CharacterInfo
	if err := json.Unmarshal(b, &data); err != nil {
		return respond.Error(err)
	}

	ret := map[int]Summary{}
	for _, ch := range data {
		ret[ch.ID] = characterSummary(ch)
	}
	return respond.JSON(ret)
}

type Summary struct {
	ID               int            `json:"id"`
	Name             string         `json:"name"`
	DNDBeyondAccount string         `json:"dndbeyond_account"`
	Campaign         string         `json:"campaign"`
	CharacterSheet   string         `json:"character_sheet"`
	Avatar           string         `json:"avatar"`
	Race             string         `json:"race"`
	BaseRace         string         `json:"base_race"`
	ClassDescription string         `json:"class_description"`
	Classes          map[string]int `json:"classes"`
	Subclasses       []string       `json:"subclasses"`
	Level            int            `json:"level"`
	// Spellslots          map[int]int    `json:"spellslots"`
	Alignment           string `json:"alignment"`
	Age                 string `json:"age"`
	Hair                string `json:"hair"`
	Eyes                string `json:"eyes"`
	Skin                string `json:"skin"`
	Height              string `json:"height"`
	PersonalityTraits   string `json:"personality_traits"`
	Ideals              string `json:"ideals"`
	Bonds               string `json:"bonds"`
	Flaws               string `json:"flaws"`
	Appearance          string `json:"appearance"`
	PersonalPossessions string `json:"personal_possessions"`
	Organizations       string `json:"organizations"`
	Allies              string `json:"allies"`
	Enemies             string `json:"enemies"`
	Backstory           string `json:"backstory"`
	OtherNotes          string `json:"other_notes"`
}

func characterSummary(ch charsync.CharacterInfo) Summary {
	ret := Summary{
		ID:               ch.ID,
		Name:             ch.Name,
		DNDBeyondAccount: ch.Username,
		Campaign:         ch.Campaign.Name,
		CharacterSheet:   ch.ReadonlyURL,
		Avatar:           ch.Decorations.AvatarURL,
		Race:             ch.Race.FullName,
		BaseRace:         ch.Race.BaseRaceName,
		Classes:          map[string]int{},
		// Spellslots:          map[int]int{},
		Alignment:           translateAlignment(ch.AlignmentID),
		Hair:                ch.Hair,
		Eyes:                ch.Eyes,
		Skin:                ch.Skin,
		Height:              ch.Height,
		PersonalityTraits:   ch.Traits.PersonalityTraits,
		Ideals:              ch.Traits.Ideals,
		Bonds:               ch.Traits.Bonds,
		Flaws:               ch.Traits.Flaws,
		Appearance:          ch.Traits.Appearance,
		PersonalPossessions: ch.Notes.PersonalPossessions,
		Organizations:       ch.Notes.Organizations,
		Allies:              ch.Notes.Allies,
		Enemies:             ch.Notes.Enemies,
		Backstory:           ch.Notes.Backstory,
		OtherNotes:          ch.Notes.OtherNotes,
	}
	sort.Slice(ch.Classes, func(i, j int) bool {
		if ch.Classes[i].Level != ch.Classes[j].Level {
			return ch.Classes[i].Level > ch.Classes[j].Level
		}
		return ch.Classes[i].Definition.Name < ch.Classes[j].Definition.Name
	})
	for _, c := range ch.Classes {
		f := c.Definition.Name
		if len(ch.Classes) > 1 {
			if c.Level == 1 {
				f += " (1 level)"
			} else {
				f += fmt.Sprintf(" (%d levels)", c.Level)
			}
		}
		ret.Level += c.Level
		ret.Classes[c.Definition.Name] = c.Level
		if ret.ClassDescription != "" {
			ret.ClassDescription += ", "
		}
		ret.ClassDescription += f

		if c.SubclassDefinition.Name != "" {
			ret.Subclasses = append(ret.Subclasses, c.SubclassDefinition.Name)
		}
	}
	/* Is always empty?
	for _, s := range ch.SpellSlots {
		if s.Available > 0 {
			ret.Spellslots[s.Level] = s.Available
		}
	}
	*/
	if ch.Age > 0 {
		ret.Age = strconv.Itoa(ch.Age)
	}
	return ret
}

func translateAlignment(n int) string {
	switch n {
	case 1:
		return "Lawful Good"
	case 2:
		return "Neutral Good"
	case 3:
		return "Chaotic Good"
	case 4:
		return "Lawful Neutral"
	case 5:
		return "True Neutral"
	case 6:
		return "Chaotic Neutral"
	case 7:
		return "Lawful Evil"
	case 8:
		return "Neutral Evil"
	case 9:
		return "Chaotic Evil"
	default:
		return ""
	}
}
