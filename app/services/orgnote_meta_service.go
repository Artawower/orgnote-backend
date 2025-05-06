package services

import (
	"fmt"
	"orgnote/app/configs"
	"orgnote/app/models"
	"orgnote/app/tools"
	"regexp"

	"github.com/google/go-github/github"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"golang.org/x/mod/semver"
)

type OrgNoteMetaConfig struct {
	ClientRepoName   string
	ClientRepoOwner  string
	DisableScheduler *bool
}

type OrgNoteMetaService struct {
	repoConfig       OrgNoteMetaConfig
	config           configs.Config
	cachedClientInfo *github.RepositoryRelease
	queue            *cron.Cron
}

func NewOrgNoteMetaService(repoConfig OrgNoteMetaConfig, config configs.Config) *OrgNoteMetaService {
	metaService := &OrgNoteMetaService{repoConfig, config, nil, nil}
	if repoConfig.DisableScheduler == nil || !*repoConfig.DisableScheduler {
		metaService.RunScheduler()
	}
	return metaService
}

func (o *OrgNoteMetaService) LoadClientMeta() error {
	client := github.NewClient(nil)
	ctx, _ := tools.DefaultContextTimeout()
	release, _, err := client.Repositories.GetLatestRelease(ctx, o.repoConfig.ClientRepoOwner, o.repoConfig.ClientRepoName)

	if err != nil {
		return fmt.Errorf("orgnote meta: load client version: %w", err)
	}

	o.cachedClientInfo = release
	return nil
}

func (o *OrgNoteMetaService) LoadReleasesChanges() error {
	return fmt.Errorf("orgnote meta: load releases changes: method unimplemented yet")
}

func (o *OrgNoteMetaService) GetChangesFrom(version string) *models.OrgNoteClientUpdateInfo {
	if o.cachedClientInfo == nil || o.cachedClientInfo.TagName == nil {
		return nil
	}

	needUpdate := semver.Compare(tools.NormalizeVersion(version), tools.NormalizeVersion(*o.cachedClientInfo.TagName)) == -1

	if !needUpdate {
		return nil
	}

	return &models.OrgNoteClientUpdateInfo{
		Version:   *o.cachedClientInfo.TagName,
		Url:       o.cachedClientInfo.GetHTMLURL(),
		ChangeLog: o.formatChangeLog(o.cachedClientInfo.Body),
	}
}

func (o *OrgNoteMetaService) formatChangeLog(changeLog *string) string {
	if changeLog == nil {
		return ""
	}

	r := regexp.MustCompile(`(.*: )`)
	return r.ReplaceAllString(*changeLog, "")

}

func (o *OrgNoteMetaService) RunScheduler() {
	if o.queue != nil {
		return
	}
	err := o.LoadClientMeta()
	if err != nil {
		log.Error().Msgf("orgnote meta: run scheduler: %s", err)
	}

	o.queue = cron.New()
	o.queue.AddFunc("@every 30m", func() {
		err := o.LoadClientMeta()
		if err != nil {
			log.Error().Msgf("orgnote meta: run scheduler: %s", err)
		}
	})

	o.queue.Start()
}

func (o *OrgNoteMetaService) GetEnvironmentInfo() models.EnvironmentInfo {
	return models.EnvironmentInfo{
		SelfHosted: tools.IsEmpty(o.config.AccessCheckerURL) || tools.IsEmpty(o.config.AccessCheckToken),
	}
}
