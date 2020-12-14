package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
)

// PostToChannelByIDAsBot posts a message to the provided channel.
func (p *Plugin) PostToChannelByIDAsBot(channelID, message string) error {
	_, appError := p.API.CreatePost(&model.Post{
		UserId:    p.BotUserID,
		ChannelId: channelID,
		Message:   message,
	})
	if appError != nil {
		return appError
	}

	return nil
}
