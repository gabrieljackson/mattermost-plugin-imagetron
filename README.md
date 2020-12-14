<img src="https://github.com/gabrieljackson/mattermost-plugin-imagetron/blob/master/assets/profile.png?raw=true" width="100" height="100" alt="imagetron">

# Mattermost ImageTron Plugin

A fun image-generating [Mattermost](https://mattermost.com) plugin.

## About

ImageTron can put your Mattermost server to work by turning it into an image generator.

![sample](https://user-images.githubusercontent.com/3694686/102115080-0a197380-3e09-11eb-9b89-3bd11d7706cb.png)

## Image Types

Currently, ImageTron supports image generation from the very awesome [primitive](https://github.com/fogleman/primitive) tool. Even if you have no interest in using this plugin, I highly recommend that you try out primitive locally to see what you think.

In the future, more base image types will be supported.

## Commands

#### /imagetron primitive [url] [optional flags]

Produce a [primitive](https://github.com/fogleman/primitive) image from the image at the provided URL.

Flags:
 - count: the number of shapes used in image generation (default 30, max 150)
 - shape: the base shape used in image generation (default 4, options: 0=combo 1=triangle 2=rect 3=ellipse 4=circle 5=rotatedrect 6=beziers 7=rotatedellipse 8=polygon)

## Warnings

While generating images, this plugin can consume many server resources which could possibly affect other processes.

Attempts to address this will be made in the future.

## Ideas for improvement

 - New base types of images
 - Expose more options for image generation
 - Process batches of images
