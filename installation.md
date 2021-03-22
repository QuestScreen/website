---
layout: default
title: Installation
weight: 2
permalink: /installation/
---

## Installation

QuestScreen is written in Go.
Since Go's support for plugins is rather lackluster (for starters, they don't work on Windows), QuestScreen offers an installer that helps you compiling it from source and enables you to add any plugins you want to use directly into the main application.

### Prerequisites

To compile QuestScreen, you need a Go compiler.
There are binary packages available on [golang.org](https://golang.org), which is the preferred way to acquire it for Windows users.
Most Linux distributions have the Go compiler available in their package repositories.
Mac users can use [HomeBrew](https://brew.sh).
You will need at least Go 1.12.

Open a shell / terminal window and ensure that `go version` is available and displays something not older than 1.12.
Once this works, you are ready to install QuestScreen.

### Getting the Installer

QuestScreens installer is a command-line tool.
This is because QuestScreen can be used on systems that do not necessarily have a graphical desktop environment, like for example a Raspberry Pi.
You will need an internet connection.

To get the installer, open a shell / terminal window and enter

    go get github.com/QuestScreen/qs-build

Once that is done, you should be able to call the installer:

    qs-build

What the installer does is this:

 * It searches for the QuestScreen source code in a directory named `QuestScreen` in the current directory.
   If such a directory is not found, it checks whether you are currently inside a directory containing the QuestScreen source code.
   If that also isn't the case, it will ask you whether you want to download the source code into a directory `QuestScreen` inside the current directory.
 * If the previous step was successful, the installer now has access to the QuestScreen source code.
   It will search for several dependencies and automatically download and install them if not found.
   These are all small Go tools that are used for compiling QuestScreen.
 * The installer will now search for plugins inside the `plugins` directory.
   QuestScreen contains a plugin `base` that already resides in this directory.
   The installer will ask you whether you want to add additional plugins.
   If you choose to do this, the installer will quit and you need to copy the desired plugins into the `plugins` directory.
   You can then re-run the installer which will do nothing in the first two steps because the source code was already acquired and the dependencies installed in the previous run.
 * When you tell the installer that you are content with your plugins, you can tell it to finalize the installation.
   This will generate several files in the source directory and then compile the application.
   Depending on your machine, this will take some time and does not have a progress bar.
 * After compilation, the resulting executable file, named `questscreen` (or `questscreen.exe` on Windows) will be moved into your user's `bin` folder for Go tools.

This sums up the installation.
You now have an executable that you can start from the terminal or via double-click.

### Advanced Usage

This section is only relevant for developers.

`qs-build` is not only the installer, but also the build system of QuestScreen.
You can call it without any argument within the QuestScreen source code directory to build QuestScreen.

The build process has several phases:

 * `plugins`: Discover all plugins and generate Go code that loads the plugins.
   This phase also analyses the plugins and detects errors in the module layout of a plugin.
 * `webui`: Compile the part of the code that implements the web frontend.
   This employs GopherJS to compile Go code to JavaScript.
   It uses the tool [askew](https://github.com/flyx/askew) to generate the code managing the HTML UI components defined in the `.askew` files.
 * `assets`: Collect all assets, including the generated JavaScript code, CSS and so on, and write a Go file that includes these files as data.
   This way, the resulting application is completely standalone and does not need any supporting files.
 * `compile`: Calls the Go compiler on the main application.

Each phase depends on the result of the previous phase.
Since the results are always generated files, you can only call the necessary phases when working on the code.
To only execute specific phases, give the names of those phases as command line arguments.
