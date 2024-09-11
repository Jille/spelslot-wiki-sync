package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/samber/lo"
)

var knownImageExtensions = regexp.MustCompile(`\.(jpeg?|webp|png|gif)$`)

func characterToWikiPage(ch CharacterResponse) string {
	avatarURL := ch.Data.Decorations.AvatarURL
	if avatarURL != "" && !knownImageExtensions.MatchString(avatarURL) {
		if strings.Contains(avatarURL, "?") {
			avatarURL = strings.TrimSuffix(avatarURL, "&") + "&_mediawiki_hack=.jpg"
		} else {
			avatarURL += "?_mediawiki_hack=.jpg"
		}
	}
	avatarURL = strings.ReplaceAll(strings.ReplaceAll(avatarURL, "=", "{{=}}"), "|", "{{!}}")

	var charParams []string
	var infoTable []string
	var categories []string
	charParams = append(charParams, fmt.Sprintf("|id=%d", ch.Data.ID))
	charParams = append(charParams, "|Name="+ch.Data.Name)
	infoTable = append(infoTable, "Name", ch.Data.Name)
	if ch.Data.Username != "" {
		charParams = append(charParams, "|Player="+translateUsername(ch.Data.Username))
		infoTable = append(infoTable, "Player", translateUsername(ch.Data.Username))
	}
	if ch.Data.ReadonlyURL != "" {
		charParams = append(charParams, "|CharacterSheet="+ch.Data.ReadonlyURL)
		infoTable = append(infoTable, "Character sheet", ch.Data.ReadonlyURL)
	} else {
		infoTable = append(infoTable, "Character sheet", "''not available''")
	}
	if avatarURL != "" {
		charParams = append(charParams, "|Avatar="+avatarURL)
	}
	if ch.Data.Race.FullName != "" {
		charParams = append(charParams, "|Race="+ch.Data.Race.FullName)
		infoTable = append(infoTable, "Race", ch.Data.Race.FullName)
		categories = append(categories, "Races/"+ch.Data.Race.BaseRaceName)
	}
	var classes []string
	var classAndLevels []string
	sort.Slice(ch.Data.Classes, func(i, j int) bool {
		if ch.Data.Classes[i].Level != ch.Data.Classes[j].Level {
			return ch.Data.Classes[i].Level > ch.Data.Classes[j].Level
		}
		return ch.Data.Classes[i].Definition.Name < ch.Data.Classes[j].Definition.Name
	})
	for _, c := range ch.Data.Classes {
		f := c.Definition.Name
		if len(ch.Data.Classes) > 1 {
			if c.Level == 1 {
				f += " (1 level)"
			} else {
				f += fmt.Sprintf(" (%d levels)", c.Level)
			}
		}
		classes = append(classes, c.Definition.Name)
		classAndLevels = append(classAndLevels, f)
		categories = append(categories, "Classes/"+c.Definition.Name)
	}
	if len(classes) > 0 {
		charParams = append(charParams, "|Class="+strings.Join(classes, ", "))
		infoTable = append(infoTable, "Class", strings.Join(classAndLevels, ", "))
	}
	var spellslots []string
	for _, s := range ch.Data.SpellSlots {
		if s.Available > 0 {
			spellslots = append(spellslots, fmt.Sprintf("level %d: %dx", s.Level, s.Available))
		}
	}
	if len(spellslots) > 0 {
		charParams = append(charParams, "|Spellslots="+strings.Join(spellslots, ", "))
		infoTable = append(infoTable, "Spellslots", strings.Join(spellslots, ", "))
	}
	if ch.Data.AlignmentID != 0 {
		charParams = append(charParams, "|Alignment="+translateAlignment(ch.Data.AlignmentID))
		infoTable = append(infoTable, "Alignment", translateAlignment(ch.Data.AlignmentID))
	}
	if ch.Data.Age > 0 {
		charParams = append(charParams, fmt.Sprintf("|Age=%d", ch.Data.Age))
		infoTable = append(infoTable, "Age", fmt.Sprint(ch.Data.Age))
	}
	if ch.Data.Hair != "" {
		charParams = append(charParams, "|Hair="+ch.Data.Hair)
		infoTable = append(infoTable, "Hair", ch.Data.Hair)
	}
	if ch.Data.Eyes != "" {
		charParams = append(charParams, "|Eyes="+ch.Data.Eyes)
		infoTable = append(infoTable, "Eyes", ch.Data.Eyes)
	}
	if ch.Data.Skin != "" {
		charParams = append(charParams, "|Skin="+ch.Data.Skin)
		infoTable = append(infoTable, "Skin", ch.Data.Skin)
	}
	if ch.Data.Height != "" {
		charParams = append(charParams, "|Height="+ch.Data.Height)
		infoTable = append(infoTable, "Height", ch.Data.Height)
	}

	var out strings.Builder
	fmt.Fprintf(&out, "<!-- DO NOT EDIT this page. -->")
	fmt.Fprintf(&out, "<!-- It will be automatically reverted. -->")
	fmt.Fprintf(&out, "<!-- Make your changes in D&D Beyond or talk to Jille. -->\n")
	if avatarURL != "" {
		fmt.Fprintf(&out, "{{Avatar|%s}}\n\n", avatarURL)
	}
	fmt.Fprintf(&out, "{{Character\n")
	for _, p := range charParams {
		fmt.Fprintf(&out, "  %s\n", p)
	}
	fmt.Fprintf(&out, "}}\n")
	fmt.Fprintf(&out, "{| class=wikitable\n")
	fmt.Fprintf(&out, "|-\n")
	fmt.Fprintf(&out, "! colspan=2| ''Character''\n")
	for _, p := range lo.Chunk(infoTable, 2) {
		fmt.Fprintf(&out, "|-\n")
		fmt.Fprintf(&out, "! scope=\"row\"| %s\n", p[0])
		fmt.Fprintf(&out, "| %s\n", p[1])
	}
	fmt.Fprintf(&out, "|}\n")

	if ch.Data.Traits.PersonalityTraits != "" {
		fmt.Fprintf(&out, "== Personality Traits ==\n%s\n\n", ch.Data.Traits.PersonalityTraits)
	}
	if ch.Data.Traits.Ideals != "" {
		fmt.Fprintf(&out, "== Ideals ==\n%s\n\n", ch.Data.Traits.Ideals)
	}
	if ch.Data.Traits.Bonds != "" {
		fmt.Fprintf(&out, "== Bonds ==\n%s\n\n", ch.Data.Traits.Bonds)
	}
	if ch.Data.Traits.Flaws != "" {
		fmt.Fprintf(&out, "== Flaws ==\n%s\n\n", ch.Data.Traits.Flaws)
	}
	if ch.Data.Traits.Appearance != "" {
		fmt.Fprintf(&out, "== Appearance ==\n%s\n\n", ch.Data.Traits.Appearance)
	}

	if ch.Data.Notes.PersonalPossessions != "" {
		fmt.Fprintf(&out, "== Personal Possessions ==\n%s\n\n", ch.Data.Notes.PersonalPossessions)
	}
	if ch.Data.Notes.Organizations != "" {
		fmt.Fprintf(&out, "== Organizations ==\n%s\n\n", ch.Data.Notes.Organizations)
	}
	if ch.Data.Notes.Allies != "" {
		fmt.Fprintf(&out, "== Allies ==\n%s\n\n", ch.Data.Notes.Allies)
	}
	if ch.Data.Notes.Enemies != "" {
		fmt.Fprintf(&out, "== Enemies ==\n%s\n\n", ch.Data.Notes.Enemies)
	}
	if ch.Data.Notes.Backstory != "" {
		fmt.Fprintf(&out, "== Backstory ==\n%s\n\n", ch.Data.Notes.Backstory)
	}
	if ch.Data.Notes.OtherNotes != "" {
		fmt.Fprintf(&out, "== Other ==\n%s\n\n", ch.Data.Notes.OtherNotes)
	}

	fmt.Fprintf(&out, "\n")
	for _, c := range categories {
		fmt.Fprintf(&out, "[[Category:%s]]\n", c)
	}

	return out.String()
}

func translateUsername(u string) string {
	switch u {
	case "Quis__":
		return "Jille"
	default:
		return u
	}
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
