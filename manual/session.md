---
layout: with-nav-footer
title: Running a Session
weight: 7
permalink: /manual/session/
breadcrumb: User Manual
parent: /manual/
previous: /manual/configuration/
---

When you have configured everything to your liking, you can start a pen & paper session.
This is done on the *Home* page by clicking the group with which you want to start the session.
Starting a session throws you on the *Session* page:

![Session page](/media/web-interface-session.png)

Here, each active module provides an interface to control it.
For example, the *Background Image* module has a dropdown menu where you can select the current background image.

Unlike on the other pages, changes on the session page are usually commited immediately and trigger an animation on the screen.
So for example, if you select a different background image, the screen will transition to that image as soon as you select it.

The *Hero List* module lets you hide some of the heroes and there's also a button to hide all (for example if you want to show the whole background image).

The *Overlays* module allows you to show multiple images anchored at the bottom of the scene.
It is designed to show items or persons the group encounters.
If you show many images at once, they may be scaled down to fit on the screen.

The *Scene Title* module shows a title at the top of the scene.
It uses the selected font size but will scale down the text if it doesn't fit.

The state of a session will be persisted per group.
So when you start the next session with the same group, the scene state will be what it was when you ended the session.
You can end the session with the big red *End* button.