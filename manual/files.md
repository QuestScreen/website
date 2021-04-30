---
layout: with-nav-footer
title: Managing Data Files
weight: 4
permalink: /manual/files/
breadcrumb: User Manual
parent: /manual/
previous: /manual/startup/
next: /manual/datasets/
---

All of QuestScreen's data is kept in its *data directory*.
By default, the data directory is located at `~/.local/share/questscreen` (`~` being your user directory).
You can tell QuestScreen to use the data directory at *path* instead via the command line option `-d `*`path`*.
If the data directory does not exist when QuestScreen starts, it will generate the directory with all needed subdirectories.

Some modules might need images or other files to function.
For example, the *Base > Background Image* module can render a background image over the boring white background, but you need to supply it with images it can display.

Whenever a module wants you to supply files, you put them in a *Module directory*.
The name of this directory is given on the *Home* page per module.
There can be multiple module directories per module, located at the following places (rooted in the data directory):

 * `base`: The module directories here contain files that should always be available.
 * `systems/`*`id`*: The module directories here contain files that should only be available when using the system with the ID *id*.
 * `groups/`*`id`*: The module directories here contain files that should only be available for the group with the ID *id*.
 * `groups/`*`gid`*`/scenes/`*`sid`*: The module directories here contain files that should only be available for the group with the ID *gid* while displaying the scene with the ID *sid*.

This setup enables you to restrict QuestScreen to only present you with relevant files e.g. when selecting a background – for example, general scenery images that could always be useful would go in `base`, official images from a system's artwork would go in the system's directory, and pictures of the PCs would go in the group's directory.

The module directories will not automatically be created by QuestScreen; to add files for a module, just create a directory for the module with the name specified in the web interface's *Home* page.
The IDs of systems, groups and scenes can be inspected on the *Datasets* page.

QuestScreen scans the data directory once on startup.
This means that you currently cannot access files that have been added after QuestScreen was started – you need to restart it to see the new files.
This might change in the future.

## The startup configuration file

While we're at it, let's take a quick look at the file `config.yaml` in our data directory's root:

```yaml
fullscreen: false
width: 800
height: 600
port: 8080
keyActions:
  - key: Escape
    returnValue: 0
    description: Exit
```

This file will be created if it doesn't exist.
Here you can modify the defaults for all command line arguments, and you can configure additional key actions on the quit screen.
For example, to enable fullscreen by default, change its value to `true`.
`keyActions` takes a sequence of keys, in each sequence item you define the key to enter, the description that should be displayed, and the return value that should be returned by the application when it is exited via that key.