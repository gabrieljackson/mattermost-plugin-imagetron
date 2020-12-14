package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const helpText = `ImageTron Plugin - Slash Command Help

/imagetron primitive [url] [optional flags]
/imagetron info
`

func getHelp() string {
	return fmt.Sprintf("```%s```", helpText)
}

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "imagetron",
		DisplayName:      "ImageTron",
		Description:      "Generate awesome images",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: primitive, info",
		AutoCompleteHint: "[command]",
		AutocompleteData: getAutocompleteData(),
	}
}

func getCommandResponse(responseType, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Username:     "imagetron",
		IconURL:      fmt.Sprintf("/plugins/%s/profile.png", manifest.Id),
	}
}

// ExecuteCommand executes a given command and returns a command response.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	stringArgs := strings.Split(args.Command, " ")

	if len(stringArgs) < 2 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getHelp()), nil
	}

	command := stringArgs[1]

	var handler func([]string, *model.CommandArgs) (*model.CommandResponse, bool, error)

	switch command {
	case "primitive":
		if len(stringArgs) < 2 {
			break
		}

		handler = p.runMakePrimitiveImageCommand
		stringArgs = stringArgs[2:]
	case "info":
		handler = p.runInfoCommand
		stringArgs = stringArgs[2:]
	}

	if handler == nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getHelp()), nil
	}

	resp, userError, err := handler(stringArgs, args)

	if err != nil {
		p.API.LogError(err.Error())
		if userError {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("__Error: %s__\n\nRun `/imagetron help` for usage instructions.", err.Error())), nil
		}

		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "An unknown error occurred. Please talk to your administrator for help."), nil
	}

	return resp, nil
}

func (p *Plugin) runInfoCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	resp := fmt.Sprintf("imagetron plugin version: %s, "+
		"[%s](https://github.com/gabrieljackson/mattermost-plugin-imagetron/commit/%s), built %s\n\n",
		manifest.Version, BuildHashShort, BuildHash, BuildDate)

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, resp), false, nil
}

func getAutocompleteData() *model.AutocompleteData {
	command := model.NewAutocompleteData("imagetron", "[command]", "Available commands: primitive, info, help")

	primitive := model.NewAutocompleteData("primitive", "[url]", "Generate a primitive image from the provided url")
	command.AddCommand(primitive)

	info := model.NewAutocompleteData("info", "", "Shows plugin information")
	command.AddCommand(info)

	help := model.NewAutocompleteData("help", "", "Shows detailed help information")
	command.AddCommand(help)

	return command
}
