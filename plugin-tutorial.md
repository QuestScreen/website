---
layout: default
title: Plugin Tutorial
permanlink: /plugins/tutorial/
weight: 3
---

## Introduction

This is a step-by-step tutorial that shows how to build a QuestScreen plugin.
It deliberately does not describe the design of the plugin API in detail; this is what the [Plugin Documentation](/plugins/documentation/) is for.

This article assumes you are somewhat familiar with Go, JavaScript, HTML and programming in general.
It tries to describe everything good enough so that a beginner should be able to understand it.

As for the build environment, you'll need the [Go compiler][1] (at least Go 1.12), the [SDL2 library][2] including its headers, on a supported operating system (Linux and macOS are both fine, Windows is unsupported).

 [1]: https://golang.org/
 [2]: https://www.libsdl.org/download-2.0.php