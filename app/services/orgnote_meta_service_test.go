package services

import (
	"orgnote/app/configs"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/google/go-github/github"
)

func TestGetLatestChangesReturnsCachedRelease(t *testing.T) {
	disableScheduler := true
	repoConfig := OrgNoteMetaConfig{
		DisableScheduler: &disableScheduler,
	}

	config := configs.Config{}
	orgNoteMetaService := NewOrgNoteMetaService(repoConfig, config)
	tag := "0.42.0"
	body := "release: 0.42.0\n- feat: latest changes"
	htmlURL := "https://example.com/release/0.42.0"
	orgNoteMetaService.cachedClientInfo = &github.RepositoryRelease{
		TagName: &tag,
		Body:    &body,
		HTMLURL: &htmlURL,
	}

	changes := orgNoteMetaService.GetLatestChange()

	if changes == nil {
		t.Fatal("expected latest changes")
	}
	if changes.Version != tag {
		t.Fatalf("expected version %s, got %s", tag, changes.Version)
	}
	if changes.Url != htmlURL {
		t.Fatalf("expected url %s, got %s", htmlURL, changes.Url)
	}
	if changes.ChangeLog != "0.42.0\nlatest changes" {
		t.Fatalf("expected changelog to be formatted, got %q", changes.ChangeLog)
	}
}

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
	repoConfig := OrgNoteMetaConfig{
		DisableScheduler: &disableScheduler,
	}

	config := configs.Config{}

	orgNoteMetaService := NewOrgNoteMetaService(repoConfig, config)

	formatted := orgNoteMetaService.formatChangeLog(&link)

	snaps.MatchSnapshot(t, formatted)

}
