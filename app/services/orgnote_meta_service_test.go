package services

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestChangelogShouldBeFormatted(t *testing.T) {
	link := `- 648596d release: 0.17.0
- 7796632 feat: ability to upload GPG keys from files
- 6b1970e feat: include encryption info to debug page
- 27f0c5e fix: actions block alignment for raw editor block widgets
- e63d2d3 fix: line class decorations for complex blocks inside quote block
- 40bba17 feat: encryption using orgnote api
- ed6286e feat: GPG encryption
- 7aa69fb feat: build main action toolbar from commands (#22)`

	disableScheduler := true
	config := OrgNoteMetaConfig{
		DisableScheduler: &disableScheduler,
	}
	orgNoteMetaService := NewOrgNoteMetaService(config)

	formatted := orgNoteMetaService.formatChangeLog(&link)

	snaps.MatchSnapshot(t, formatted)

}
