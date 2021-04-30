---
layout: with-nav-footer
title: The Datasets Page
weight: 5
permalink: /manual/datasets/
breadcrumb: User Manual
parent: /manual/
previous: /manual/files/
next: /manual/configuration/
---

The *Datasets* page looks like this:

![Datasets page](/media/web-interface-datasets.png)

The navigation on the left lets you view the following tabs:

The **Base** tab lets you add and remove systems and groups.
Plugins may *require* a system to exist in which case it cannot be deleted.
The idea is that a plugin can cater for a specific system, in which case it defines the default configuration for that system which will be used to create the system if it doesn't exist.

The **System** tabs let you rename each system.
Currently, there is no other data for systems.

The **Group** tabs let you rename each group and define the system it uses.
If you define that a group uses a certain system, you will have access to the files in the module directories of that system during a session with that group.
Systems can also define a basic look that is used for all groups using that system.
Such looks are defined in the *Configuration* page.

In a group tab, you can also define *scenes* for the group.
A scene is a composition of modules that will be used together on the screen.
For example, a group might have a default scene that displays the current location as background and the name of the current location as scene title.
Another scene could be used during battles and can contain an initiative table and maybe even the health of the PCs.
Mind that modules rendering such information are specific to a system and are therefore not part of QuestScreen's basic modules.
They can be added as plugin.

You can also manage the members of the group in the group tab.
These will currently only be used by the module *Base > Hero List* that display the names and descriptions of the group's members.

Each scene of a group has a **Scene** tab.
In that tab you can configure the scene's name and the modules that should be active in that scene.