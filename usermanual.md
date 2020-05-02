---
layout: default
title: User Manual
weight: 2
permalink: /usermanual/
---
## Introduction

<section class="highlighted"><i class="fas fa-info-circle"></i>
<div>
This manual covers the operation of a Quest Screen installation.
For building and installing the app, please see the Readme file on GitHub.
</div>
</section>

Quest Screen is an app that renders information on a screen during pen & paper sessions.
This manual uses the term *session* to refer to a pen & paper session.

Quest Screen is a modular display, meaning that the image rendered on the screen is assembled from a number of independent *modules*.
For example, the *background* module renders the current background image, and the *herolist* module renders the list of heroes on top of that.

The set of modules that can contribute to the current image make up the current *scene*.
For example, the default scene consists of the modules *background*, *herolist*, *title* and *overlays*.
You can use different scenes during a session, for example if you have a module displaying some kind of map which needs too much place to fit in with other modules.
Additional modules can be added via plugins.

The information shown on the display is the current *state*.
Each scene has a state and that state is retained when you swich between sessions, and also from one session to the next.
For example, the *title* module's state is the text it currently displays.
Changing the state will be animated.

Additionally, you can *configure* some modules.
Such configuration typically contains stuff like background color or text font.
Configuration should be customized between sessions and changes to configuration are not animated (they are applied immediately).

Quest Screen allows you to manage any number of *groups* with each group having its own set of scenes, configuration and state.
You can use *systems* (meaning *pen & paper systems*) to manage configuration and resources (e.g. images) available to all groups using that system.

Generally, Quest Screen is designed to be controlled by the game master via web interface (which runs on any modern browser, be it on a smartphone, tablet or laptop).
There is no interface for players.

## Setup

You need the `questscreen` executable.
Refer to the **Readme** file on GitHub to read up how to compile it.
We'll start it once to initialize its default data directory.

On first startup, it will display an empty screen with the Quest Screen logo.
This is because initially, it does not have any fonts available and thus cannot render text.
In a desktop environment, it will also start windowed.
You can quit it by pressing *Escape* two times or closing the window.

After the first run, it will have created its data directory in `~/.local/share/questscreen`.
Now to actually use it, you need to give it at least one font to render text.
Fonts go in the `fonts` subdirectory; you can put TTF and OTF files there.
A nice free font for fantasy-themed groups would be [Garamond](https://garamond.org/).
You can also put symlinks to fonts installed on your system into the directory.

Another important subdirectory is `textures` where you can put grayscale images.
These images will be used to texture backgrounds with two colors, where white is filled with the first color, black with the second, and gray shades with the corrensponding mixture of the two colors.
The textures should be repeatable in both directions.
Quest Screen does not require any textures to run.

Both fonts and textures are loaded at startup; if you add some, you'll need to restart Quest Screen.

Let's have a quick look at the default configuration that has been created in `config.yaml` inside the data directory:

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

`width` and `height` define the window's size if it is not `fullscreen`.
`port` is the port the web interface will be available at.
`keyActions` is a list of keys you can press to exit the app.

When the app is running and at least one font is available, pressing any key will display a popup which shows all configured keys that can be used to quit the application, along with their description.
The return value is returned by the app.
This is a convenience feature for using Quest Screen on a board like the Raspberry Pi.
You can write a simple script for deciding what to do after the app quits based on which key has been pressed (e.g. shutting down the system, launching some media center, etc).

Apart from quitting the app, you cannot control it locally – everything else is done via the web interface.
When you have at least one font installed, the app will show you the URL you need to enter in your browser to reach the web interface.
This will be `http://<ip>:<port>` with `<ip>` being the public IP address of the current device, and `<port>` being the configured port.

## The Web Interface

The main page of the web interface shows you the list of currently loaded modules.
The menu at the top gives you three options:

 * *Switch Group* lets you switch to a group.
   Switching the group will load the group's recently active scene and show it on the display.
   On the web interface, you will get an interface for manipulating the scene's state.
 * *Configuration* lets you modify the configuration of your modules.
 * *Datasets* lets you manage your groups and systems.
   This is where you can create your first group.

### Managing Datasets

Here you manage your systems, groups, scenes and heroes.
Scenes and heroes always belong to a group, a group may link to a system.
Sometimes plugins require some system to exist, in which case you cannot remove it as long as the plugin is installed.

When you create a new group or scene, you can select between from number of templates.
Templates are defined by plugins and are used to provide default configurations which you can then modify.
For example, a plugin might provide a group template which contains a default scene (with background, list of heroes and so on) and a battle scene which contains the module showing battle stats that is provided by the plugin.

The scene templates contained in a group template are also available for creating single scenes.
The *base* plugin (which provides basic modules that are always available) provides empty templates if you want to start with an empty group or scene.

### Configuring the Modules

By default, modules have a pretty simple look.
For example, the *title* module by default has a white background.
You can change the module's look in the *configuration* section.

A module's configuration has multiple layers.
On each layer, you may define configuration values that override the ones from a lower layer.
If you don't set some configuration value for a layer, the configuration value from the lower layer will be used.
The following layers exist:

 * The **default** configuration is defined for each module and is not editable.
 * The **base** configuration applies to all groups and overrides the default configuration.
 * The **system** configuration applies to all groups linked to that system and overrides the base configuration.
 * The **group** configuration applies to all scenes in that group and overrides the system configuration (or the base configuration if the group is not linked to a system).
 * The **scene** configuration applies to a single scene and overrides the group configuration.

Usually, you want to edit the base or system configuration unless you want to have a different look for each group.
Scene configuration is useful if you want to use a module in multiple scenes but want to have it look different.

### Using Quest Screen during a session

Finally, when you have created a group and configured it to your liking, you can start using Quest Screen during a session.
You start by selecting the active group.
Only one group can be active and you can only modify the state of the active group.

Unlike the configuration and dataset pages (where you can click *reset* or *save* to commit your changes), changes to the state are typically sent immediately.
Each action triggers the corresponding animation.
Which interactions are available depends on the active modules in the current scene.
Changing the scene will change the possible interactions depending on the modules in the new scene.
Module states are local to the scene, so if you have the *title* module active in two scenes, updating the text in one scene will not modify the text displayed in the other scene.

## Providing files to modules

Some modules, like the *background* module, depend on files (in this case, images).
Currently, the web interface does not allow file uploads; you must place them on the host system manually.
The files for each module have to be put in a directory whose name matches the module's ID.
You can look up the ID of a module on the web interface's main page.

You can create such a directory in any of the following places (relative to the root directory `~/.local/share/questscreen`):

 * `base`: Files placed will be available in all groups.
 * `systems/<system-id>`: Files placed here will be available in all groups that are linked to the system identified by `<system-id>`.
 * `groups/<group-id>`: Files placed here will be available in all scenes of the group identified by `<group-id>`.
 * `groups/<group-id>/scenes/<scene-id>`: Files placed here will be available in the scene identified by `<scene-id>`.

So for example, if you have the following files:

    base/background/one.jpg
    groups/firstgroup/background/two.jpg
    groups/othergroup/scenes/primary/background/three.jpg

The background module in any scene of `firstgroup` will have access to `one.jpg` and `two.jpg`, while in the `primary` scene if group `othergroup`, it will have access to `one.jpg` and `three.jpg`.

Modules may define multiple sets of files that need to be placed in separate subdirectories with defined names.
They may also require a specific file with a defined name.
Make sure to read the documentation of the plugin providing the module to figure out what files it expects and where you have to put them!

## Persistence

Whenever you modify the state, the active scene of a group, any configuration or any dataset, the change will be immediately written to the file system.
When you quit Quest Screen and later start it again, any group will still be in the state where you left it.

Configuration will be written to `config.yaml` files in the corresponding base, system, group or scene directory.
State will be written to `state.yaml` in the scene directory.

If you want to backup your state, you can simply backup the whole base directory.
This will save your whole configuration and state, including all additional files you provide to modules.

## Manage Plugins

Plugins are placed in the `plugins` directory inside the root directory.
They are loaded at startup; you can't add a plugin while Quest Screen is running.

Plugins provide you with additional modules, additional systems and/or additional group and scene templates.
Writing plugins is covered by the [plugin tutorial](/plugins/tutorial/).