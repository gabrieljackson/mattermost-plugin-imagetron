package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/fogleman/primitive/primitive"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

const (
	makeImageUsage = `/imagetron primitive [url] [optional flags]
  Generate a primitive image from the provided URL
`
)

type shapeConfig struct {
	Count  int
	Mode   int
	Alpha  int
	Repeat int
}

func (s *shapeConfig) isValid() error {
	if s.Count > 150 {
		return errors.New("count value cannot be over 100")
	}
	if s.Mode > 8 || s.Mode < 0 {
		return errors.New("valid shape options are 0-8")
	}

	return nil
}

func getMakeImageFlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("make", pflag.ContinueOnError)
	flagSet.Int("shape", 4, "The base shape used: 0=combo 1=triangle 2=rect 3=ellipse 4=circle 5=rotatedrect 6=beziers 7=rotatedellipse 8=polygon")
	flagSet.Int("count", 30, "The number of shapes to generate in the output image")

	return flagSet
}

func parseMakeImageArgs(args []string) (*shapeConfig, error) {
	flagSet := getMakeImageFlagSet()
	err := flagSet.Parse(args)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse the command args")
	}

	count, err := flagSet.GetInt("count")
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse count value")
	}
	mode, err := flagSet.GetInt("shape")
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse shape value")
	}

	return &shapeConfig{
		Count:  count,
		Mode:   mode,
		Alpha:  128,
		Repeat: 0,
	}, nil
}

func (p *Plugin) runMakePrimitiveImageCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	if len(args) < 1 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, makeImageUsage), true, nil
	}
	url := args[0]

	config, err := parseMakeImageArgs(args)
	if err != nil {
		return nil, false, err
	}
	err = config.isValid()
	if err != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, err.Error()), true, nil
	}

	// Grab the base image.
	response, err := http.Get(url)
	if err != nil {
		return nil, false, err
	}
	defer response.Body.Close()

	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to create temporary image directory")
	}
	defer os.RemoveAll(tempDir)

	tmpFilename := tempDir + "image.png"
	file, err := os.Create(tmpFilename)
	if err != nil {
		return nil, false, err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return nil, false, err
	}

	baseImage, err := p.uploadLocalImage(tmpFilename, model.NewId()+".png", extra.ChannelId)
	if err != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, err.Error()), true, err
	}

	_, aerr := p.API.CreatePost(&model.Post{
		UserId:    p.BotUserID,
		ChannelId: extra.ChannelId,
		Message:   "Generating primitive image(s) from this base image. Results will be posted here when finished.",
		FileIds:   []string{baseImage.Id},
	})
	if aerr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, aerr.Error()), true, aerr
	}

	// Load in the base image and begin generating the primitave image.
	input, err := primitive.LoadImage(tmpFilename)
	if err != nil {
		return nil, false, err
	}

	bg := primitive.MakeColor(primitive.AverageImageColor(input))
	pmodel := primitive.NewModel(input, bg, 1024, 3)
	output := tempDir + "output.png"

	start := time.Now()
	frame := 0
	for i := 0; i < config.Count; i++ {
		frame++

		t := time.Now()
		n := pmodel.Step(primitive.ShapeType(config.Mode), config.Alpha, config.Repeat)
		nps := primitive.NumberString(float64(n) / time.Since(t).Seconds())
		elapsed := time.Since(start).Seconds()
		p.API.LogInfo(fmt.Sprintf("Creating primitive image - %d: t=%.3f, score=%.6f, n=%d, n/s=%s", frame, elapsed, pmodel.Score, n, nps))

		last := i == config.Count-1
		if last {
			err = primitive.SavePNG(output, pmodel.Context.Image())
			if err != nil {
				return nil, false, err
			}
		}
	}

	newImage, err := p.uploadLocalImage(output, model.NewId()+".png", extra.ChannelId)
	if err != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, aerr.Error()), true, aerr
	}

	_, aerr = p.API.CreatePost(&model.Post{
		UserId:    p.BotUserID,
		ChannelId: extra.ChannelId,
		Message:   "New primitive image generated",
		FileIds:   []string{newImage.Id},
	})
	if aerr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, aerr.Error()), true, aerr
	}

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Your new primitive image is ready!"), false, nil
}

func (p *Plugin) uploadLocalImage(fileLocation, filename, channelID string) (*model.FileInfo, error) {
	data, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		return nil, err
	}

	newImage, aerr := p.API.UploadFile(data, channelID, filename)
	if aerr != nil {
		return nil, aerr
	}

	return newImage, nil
}
