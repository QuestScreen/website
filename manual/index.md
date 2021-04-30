---
layout: with-nav-footer
title: User Manual
weight: 2
permalink: /manual/
next: /manual/startup/
---

QuestScreen is an app that renders information on a screen during pen & paper sessions.
Within this manual, the term *session* always refers to a pen & paper session.

You can obtain QuestScreen either as binary release for Windows (TODO: link) or as source code release for any operating system supporting Go, OpenGL (ES) and SDL.
QuestScreen supports plugins.
Any plugin you want to use must be compiled into your binary.
The official binary release does not contain any plugins.
The [building instructions](/building/) describe how you build the app from source and add plugins.

QuestScreen is able to run on small boards like a Raspberry Pi (even without a window manager), but can also run on most desktop environments like Windows, macOS and the various Linux desktop environments.
The hardware requirements depend mainly on the size of your screen (if you have a 4K screen, you obviously need a board or PC that supports it).
Chances are that if your hardware can generally use the screen, it is good enough to run QuestScreen.
