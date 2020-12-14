package main

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	BotUserID string

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// BuildHash is the full git hash of the build.
var BuildHash string

// BuildHashShort is the short git hash of the build.
var BuildHashShort string

// BuildDate is the build date of the build.
var BuildDate string

// OnActivate runs when the plugin activates and ensures the plugin is properly
// configured.
func (p *Plugin) OnActivate() error {
	bot := &model.Bot{
		Username:    "imagetron",
		DisplayName: "ImageTron",
		Description: "Created by the ImageTron plugin.",
	}
	options := []plugin.EnsureBotOption{
		plugin.ProfileImagePath("assets/profile.png"),
	}

	botID, err := p.Helpers.EnsureBot(bot, options...)
	if err != nil {
		return errors.Wrap(err, "failed to ensure imagetron bot")
	}
	p.BotUserID = botID

	return p.API.RegisterCommand(getCommand())
}
