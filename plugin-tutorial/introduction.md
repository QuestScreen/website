## Before you start

Ensure you have a working [Go compiler][1] (at least Go 1.12) as well as the [SDL2 library][2] including its headers installed on a supported operating system (currently Linux and macOS).
Windows is currently unsupported, mostly because [-builtmode=plugin][3] does not work on it.

This tutorial assumes you have a GitHub account so you can use the plugin template repository.
If you don't, that's not a problem.
Simply download the repository's content as zip file and upload it to your favorite source code hosting platform.
While code hosting is not strictly required, we'll use it to imitate the process of creating an actual plugin as close as possible.

## What we're gonna do

In this tutorial, we'll be writing a module which displays an in-universe calendar.
As lots of settings have custom in-universe calendars, this seems to be a worthy endeavor which can be a prototype for other calendars.
The calendar we will implement is the one of [Discworld][4] because it is rather simple.

The display should show the current day, month and year.
We'll be counting years according to the University Calendar, which counts *common years* having 400 days each (while a full rotation of the disc takes 800 days, in which each month occurs twice).

The user should be able to skip through the calendar easily.
Therefore, we want to have **+1**, **+3**, **+10**, **-1**, **-3** and **-10** buttons for days, months and years each.
An additional date chooser is left as an exercise for the reader.

The complete result code of this tutorial is available in [this repository][5] and the tutorial is actually autogenerated from that code.

## Creating the skeleton

We'll be starting from the [plugin template][6] repository.

| Click on the *Use this template* button to create a repository for our plugin.

This tutorial will assume that your repo is named `discworld` and refer to it as such.

| Use this command to clone the repo into your GOPATH:
| 
|     go get -d <url-to-your-repo>
| 
| Navigate to the clone in your GOPATH (`$GOPATH/src/<url>`).

For the rest of this tutorial, each page will be linked to one file in this directory, whose relative path is given below the heading.
Before we start, let's look at the plugin's structure:

 * `go.mod` links to the plugin API, including its dependencies, `go-yaml` and `go-sdl2`.
   The former is used for storing configuration and state in the file system, while the latter is used for rendering our module.
 * `moduletemplate` contains the template for a module.
   A plugin may contain multiple templates in which case you'd want to copy it for each additional module.
   We'll just implement one module (the calendar).
 * `web` contains HTML and JS files that will be used with the Web Client.
   It can potentially also contain CSS files.
 * We have a `Makefile` which is unusual for a Go project.
   It is mainly used to encode all files in `web` as Go source file so that they can be part of our plugin's binary file.
   Gophers like to say you should use `go generate` and commit the generated file into the repository, but I really don't fancy cluttering my git commit history with large chunks of encoded data.
   YMMV.
 * `plugin.go` is the main file of our project.
   It describes the plugin's metadata and lists its modules and web sources.

| Rename the directory `moduletemplate` to `calendar`, since it will become our calendar module.

This tutorial will now show and discuss the code you should put into each source file.
Each code snippet shown should be added to the current file, potentially replacing existing skeleton code from the plugin template.
The code snippets always cover the whole content of a file.
This means that you do not really need the skeleton code from the template to follow the tutorial.
It is still useful since you see what parts are automatically provided and what parts you need to write yourself.

 [1]: https://golang.org/
 [2]: https://www.libsdl.org/download-2.0.php
 [3]: https://golang.org/pkg/plugin/
 [4]: https://wiki.lspace.org/mediawiki/Discworld_calendar
 [5]: https://github.com/QuestScreen/plugin-tutorial
 [6]: https://github.com/QuestScreen/PluginTemplate