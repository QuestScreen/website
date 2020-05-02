# QuestScreen website

This repo contains the QuestScreen website.
It uses [Git Large File Storage][2] to hold the demo videos, so you need to install it before cloning the repo.
You need to have a clone of the [plugin tutorial][1] available in the parent folder of this clone.
The easiest way to do this is to clone the website via `go` even though it is not really a Go package (it does contain some Go code):

    go get -d github.com/QuestScreen/website
    go get -d github.com/QuestScreen/plugin-tutorial

For building the website, you need Go and Ruby.
Install the dependencies for building the website with these commands:

    go get -u github.com/tdewolff/parse/...
    bundle install

Then, you can build the website with `make`.

## How it works

The basic stuff is Markdown processed with Jekyll.

The plugin tutorial is extracted from the source file of its repository.
This makes it easy to check whether the code shown in the tutorial actually compiles and works.
Text enclosed in comments (`/* … */` in Go and JS, `<!-- … -->` in HTML) is treated as Markdown source, everything else is transformed into a syntax-highlighted Markdown code block.

All processed source files of the plugin tutorial are put into the transient `generated` directory as Markdown files to be processed with Jekyll.

 [1]: https://github.com/QuestScreen/plugin-tutorial
 [2]: https://git-lfs.github.com/