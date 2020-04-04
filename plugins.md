---
layout: default
title: Plugin Documentation
permanlink: /plugins/documentation/
weight: 4
---
## Introduction

This is the high-level plugin API documentation.
It describes the facilities defined in the plugin API and how they interact with the QuestScreen core.

## QuestScreen's Architecture

QuestScreen consists of three components, each running their own thread:
**Render**, **Storage** and the **Web Client**.
*Render* and *Storage* are implemented in Go as part of QuestScreen's executable.
The *Web Client* is implemented in HTML/CSS/JS, its files are included as data in the executable and will be delivered to the user's browser.

<figure>
  <svg viewBox="-1 -1 402 152" xmlns="http://www.w3.org/2000/svg">
    <defs>
      <marker id="head" orient="auto" markerWidth="4" markerHeight="8"
              refX="4" refY="4">
        <path d="M0,0 V8 L4,4 Z" fill="black" />
      </marker>
    </defs>
    <g class="component">
      <rect x="0" y="0" width="250" height="150" />
      <text x="16" y="12">App</text>
      <polyline points="0,17 32,17 32,0" />
      <g class="thread">
        <rect x="15" y="25" width="90" height="17" />
        <text x="60" y="36">Render Thread</text>
        <line x1="60" y1="42" x2="60" y2="145" />
      </g>
      <text x="115" y="11" style="text-anchor: end;" class="part">Render</text>
      <text x="135" y="11" style="text-anchor: start;" class="part">Storage</text>
      <line x1="125" y1="0" x2="125" y2="150" stroke-dasharray="7 4" />
      <g class="thread">
        <rect x="145" y="25" width="90" height="17" />
        <text x="190" y="36">Data Thread</text>
        <line x1="190" y1="42" x2="190" y2="145" />
      </g>
    </g>
    <g class="component">
      <rect x="280" y="0" width="120" height="150" />
      <text x="306" y="12">Browser</text>
      <polyline points="280,17 332,17 332,0" />
      <g class="thread">
        <rect x="295" y="25" width="90" height="17" />
        <text x="340" y="36">Web Client</text>
        <line x1="340" y1="42" x2="340" y2="145" />
      </g>
    </g>
    <g class="communication">
      <rect x="235" y="60" width="60" height="34" />
      <text x="265" y="72">HTTP</text>
      <text x="265" y="88">REST</text>
      <line x1="340" y1="70" x2="295" y2="70" />
      <line x1="235" y1="70" x2="190" y2="70" marker-end="url(#head)" />
      <line x1="190" y1="84" x2="235" y2="84" />
      <line x1="295" y1="84" x2="340" y2="84" marker-end="url(#head)" />
    </g>
    <g class="communication">
      <rect x="95" y="100" width="60" height="34" />
      <text x="125" y="112">inter-</text>
      <text x="125" y="128">thread</text>
      <line x1="190" y1="117" x2="155" y2="117" />
      <line x1="95" y1="117" x2="60" y2="117" marker-end="url(#head)" />
    </g>
  </svg>
  <figcaption>Communication channels of the QuestScreen app</figcaption>
</figure>

*Render* is a pure visualization component.
It uses SDL with an OpenGL (OpenGL ES for small boards) backend to render the current state to the screen.
Render has only access to data that has been sent by Storage (it can load files, such as images, from the file system, though).

*Storage* is the central component which handles incoming HTTP requests from the *Web Client* and notifies *Render* of any updates to the state or config.

The *Web Client* runs on the user's browser and communicates with Storage via a HTTP REST API.
Since the Web Client is the sole source for any changes, Storage does not need to actively push changes to the Web Client.
It only answers to requests from the Web Client.

## Modules

Your plugin will provide a number of modules.
Each module will extend each of QuestScreen's components with additional code.

<figure>
  <svg viewBox="-1 -1 382 182" xmlns="http://www.w3.org/2000/svg">
    <rect x="0" y="30" width="380" height="130" stroke="black" fill="lightgray" />
    <g class="components">
      <rect x="20" y="0" width="100" height="180" />
      <text x="70" y="11">Render</text>
      <rect x="140" y="0" width="100" height="180" />
      <text x="190" y="11">Storage</text>
      <rect x="260" y="0" width="100" height="180" />
      <text x="310" y="11">Web Client</text>
    </g>
    <g class="implementations">
      <rect x="30" y="40" width="80" height="20" />
      <text x="70" y="53">Renderer</text>
      <rect x="150" y="40" width="80" height="20" />
      <text x="190" y="53">Data</text>
      <rect x="270" y="40" width="80" height="20" />
      <text x="310" y="53">Controller</text>
      <rect x="95" y="70" width="70" height="20" />
      <text x="130" y="83">renderData</text>
      <rect x="95" y="100" width="70" height="20" />
      <text x="130" y="113">Configuration</text>
      <rect x="190" y="70" width="60" height="20" />
      <text x="220" y="82">Endpoints</text>
      <rect x="270" y="100" width="80" height="20" />
      <text x="310" y="113">HTML / CSS</text>
      <rect x="30" y="130" width="200" height="20" />
      <text x="130" y="143">Descriptor</text>
    </g>
  </svg>
  <figcaption>Components of a module and where they are used</figcaption>
</figure>

A module's **Renderer**, implementing `api.ModuleRenderer`, provides all functionality to render a module's current state to the screen.
To achieve this, it receives *renderData* from the Storage.
*renderData* can be objects of different types that do not need to implement any interface.
While a module's **State** must always be able to create a *renderData* object that contains all relevant data (used for initialization or whenever a scene or group is changed), the endpoints typically create smaller *renderData* objects which only initiate animated transitions.

A module must provide a type for its **Configuration**, which is then used inside Storage to store and load the module's configuration items.
Remember that configuration is handled on multiple levels (default, base, system, group, scene) so while most other objects are singletons, **Configuration** will instantiated multiple times.
Storage merges the different levels and sends the resulting Configuration object to the renderer whenever necessary.
Unlike *renderData*, configuration data is always send completely and does not trigger an animation.

The module's **State**, which must implement `api.ModuleState`, takes care of loading and storing the module's state relative to the current scene.
Changing the scene (and, consequently, selecting a different group) will trigger loading the state of that scene into State.
Currently, it is not possible to share state values between scenes (unlike configuration, which is shared between all group scenes if configured at group level or further up in the hierarchy).

To communicate with the Web Client, the module must define one or more **Endpoints**.
These must implement either `api.ModulePureEndpoint` or `api.ModuleIDEndpoint`.
Each Endpoint binds to a URI path, with the URI including a variable *id** in the case of `api.ModuleIDEndpoint`.
The endpoints receive HTTP requests from the client and produce both the answer to the client and the data that should be passed on to the Renderer.
Endpoints modify the **State**; you cannot have them modify the **Configuration** or anything else (since the pre-defined processing of configuration does not allow custom modification).
Any State modification done by an Endpoint must be consistent with the *renderData* it sends to the Renderer – meaning that the resulting state of the Renderer must be equivalent to that which would result from the State sending a complete *renderData* to the Renderer.

The **Controller** implemented in JavaScript provides the function rendering the user interface for the module's current state, as well as handlers for user interactions that should result in state change.
As it will typically render some HTML, it is usually accompanied with a set of HTML templates.
The whole client is pure vanilla ECMAScript 2018 and includes basic utilities for rendering HTML based on `<template>` elements.
Therefore, you can send along **HTML** (and also **CSS**) files.

Finally, the **Descriptor** contains the necessary metadata for setting up the module.
It defines how to create the described objects in the app, e.g. how many endpoints the module has and which paths they bind to, how to create the State and so on.
It also describes which *resources* the module uses, i.e. files on the file system.
These are selected via parent directory and merged in the same way as configuration objects are merged.

The HTML/CSS/JS code that implements the web client is supplied by the plugin descriptor containing the module.

### The Render Loop

The most complex interface to implement is `api.ModuleRenderer`.
A module's renderer takes care of displaying the module on the screen, including its animations.
Its functions are called in a specific order.

<figure>
  <svg viewBox="24 19 337 102" xmlns="http://www.w3.org/2000/svg">
    <defs>
      <rect id="callr" width="80" height="20" fill="none" stroke="black" />
    </defs>
    <g class="flow">
      <circle cx="60" cy="70" r="10" />
      <text x="100" y="73">Start</text>
      <path d="m50,70 h-10 a5,5 0 0 1 -5,-5 v-20 a5,5 0 0 1 5,-5 h10" marker-end="url(#head)" />
      <g class="call">
        <use xlink:href="#callr" x="50" y="30" />
        <text x="90" y="43">CreateRenderer</text>
      </g>
      <line x1="130" y1="40" x2="160" y2="40" marker-end="url(#head)" />
      <g class="call">
        <use xlink:href="#callr" x="160" y="30" />
        <text x="200" y="43">Rebuild</text>
      </g>
      <line x1="240" y1="40" x2="270" y2="40" marker-end="url(#head)" />
      <g class="call">
        <use xlink:href="#callr" x="270" y="30" />
        <text x="310" y="43">Render</text>
      </g>
      <path d="m350,40 h5 a5,5 0 0 0 5,-5 v-10 a5,5 0 0 0 -5,-5 h-205 a5,5 0 0 0 -5,5 v10 a5,5 0 0 0 5,5" />
      <path d="m260,20 a5,5 0 0 0 -5,5 v10 a5,5 0 0 0 5,5" />
      <path d="m355,40 a5,5 0 0 1 5,5 v20 a5,5 0 0 1 -5,5 h-5" marker-end="url(#head)" />
      <g class="call">
        <use xlink:href="#callr" x="270" y="60" />
        <text x="310" y="73">InitTransition</text>
      </g>
      <path d="m270,70 h-5 a5,5 0 0 0 -5,5 v10 a5,5 0 0 1 -5,5 h-105 a5,5 0 0 0 -5,5 v10 a 5,5 0 0 0 5,5 h10" marker-end="url(#head)" />
      <path d="m265,70 h-5 a5,5 0 0 1 -5,-5" />
      <g class="call">
        <use xlink:href="#callr" x="160" y="100" />
        <text x="200" y="113">TransitionStep</text>
      </g>
      <line x1="240" y1="110" x2="270" y2="110" marker-end="url(#head)" />
      <g class="call">
        <use xlink:href="#callr" x="270" y="100" />
        <text x="310" y="113">Render</text>
      </g>
      <path d="m350,110 h5 a5,5 0 0 0 5,-5 v-10 a5,5 0 0 0 -5,-5 h-100" />
      <path d="m150,90 a5,5 0 0 1 -5,-5 v-10 a5,5 0 0 1 5,-5 h10" marker-end="url(#head)" />
      <g class="call">
        <use xlink:href="#callr" x="160" y="60" />
        <text x="200" y="73">FinishTransition</text>
      </g>
      <path d="m240,70 h10 a5,5 0 0 0 5,-5 v-20 a5,5 0 0 1 5,-5" />
    </g>
  </svg>
  <figcaption>Sequences in which ModuleRenderer's functions get called</figcaption>
</figure>

After creating the renderer via `api.Module`, the first function that will be called is `Rebuild`.
This happens when a scene is activated which uses the module.
It receives a merge config and the state data generated with ModuleState's `CreateRendererData`.
`Rebuild` should set up everything so that the module can be rendered with the given config and state.
It has access to the `RenderContext` so that it can pre-render stuff to textures that can later be rendered via `Render`.

`Render` should render the module's visualization to the given context.
This should be as fast as possible; try to move every expensive operation to `Rebuild` and the animation functions if it doesn't need to be recalculated every time.
This should usually be possible because whenever the data to be rendered changes, it will either happen via `Rebuild` or via `InitTransition`.
The displayed picture is not refreshed without user input, so you can't implement continuous animation in `Render`!

Whenever the current scene changes, `Rebuild` will receive the merged configuration for the new scene (even if it is identical to the current configuration) and new data from the ModuleState.
`Rebuild` will also be called when the user changes part of the configuration.

`InitTransition` will be called to receive data created by the module's endpoint(s).
It gets the data created by the endpoint that initialized the transition.
It returns a time duration which is the duration of the following animation.
You may return -1 to tell the *Render* component that nothing changed and no animation will occur, which skips the animation loop and `FinishTransition` altogether.
Returning 0 will lead to an immediate call to `FinishTransition`.

If you returned a positive number, the animation loop – consisting of a call to `TransitionStep` followed by a call to `Render` – will start and last for the specified duration.
During this loop, you can animate the transition..
All other module's `Render` functions will be called together with the current module's `Render` function.

`FinishTransition` should put the ModuleRenderer in a stable state, i.e. the resulting state should be identical to the one the module renderer would be in when the current `ModuleState` would have been sent via `Rebuild`.

## Plugins

A plugin contains a list of modules; this list may be empty.
The web (HTML/CSS/JS) code for the modules on the client side are provided directly by the plugin.
Each is provided as a single byte-array that is expected to contain UTF-8 encoded content.
You can use [go-bindata][1] to encode the web files into go byte arrays.

Apart from that, a plugin may provide a list of system templates.
For each system template, if a system with that ID doesn't exist, it gets automatically created on app launch from the template.
Systems whose ID equals a system template in any loaded plugin cannot be deleted since they are seen as *required* by the plugin.

You typically provide system templates for plugins that support specific systems.
For example, a plugin for a specific role playing system may provide modules for showing combat information specific to that system's ruleset, or calendar information specific to that system's setting.

System templates provide values for the system's configuration.
Those are not restricted to the configuration of the plugins provided by the plugin; you can also give configuration for base modules or modules from other plugins.
This allows you to provide a standard look for the given system, e.g. by setting background color and font for the *Title* or *Herolist* plugin.

Besides system templates, you can also provide group and scene templates.
Group templates may link to systems whose ID belongs with a system template of that module.
Since system templates are automatically created, this ensures that the linked system exists.
Group template can be selected by the user when creating a new group; likewise, scene templates can be selected when creating a new scene.

Since a group always contains at least one scene, a group template must link to at least one of the scene templates in the same plugin.
When creating a group from a group template, for each linked scene template, a scene will be created in the new group.
Scene templates can also be used manually by the user when creating a scene for an existing group.

## Building a plugin

TBD

 [1]: https://github.com/go-bindata/go-bindata