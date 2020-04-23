---
layout: default
title: Home
weight: 1
permalink: /
---
**Quest Screen** is an app that displays images and information on a screen during pen & paper sessions.
It has been designed to be installed on a small board (e.g. Raspberry Pi) but you can also run it on your desktop machine or laptop.

<figure class="video-demo">
  <p>JavaScript is required to show video demos</p>
</figure>
<template id="demo-video">
  <video loop controls>
    <source type="video/mp4"/>
  </video>
  <figcaption>
    <a class="demo-selector left" href="#"><i class="fas fa-caret-left"></i></a>
    <span></span>
    <a class="demo-selector right" href="#"><i class="fas fa-caret-right"></i></a>
  </figcaption>
</template>

The app is controlled via **web interface** which works on any modern browser.
The web interface supports both large (tablet, laptop) and small (smartphone) screens.

Quest Screen allows you to **manage multiple groups** and **persists the state** of each group between sessions.
It allows you to modify presentation (colors, font etc.) per group and system.

A **plugin API** allows you to add modules that display any kind of additional information that is not covered by the modules provided by the core.

Quest Screen is **open source** and released under the terms of the [GNU GPL v3][4].

The app is currently in **alpha stage**; in order to use it, you need to [get the code from GitHub][5] and build it yourself.

---

Credits for the pictures used in the demonstration videos go to [Martin Damboldt][1], [Krivec Ales][2], [Skitterphoto][3], [susannp4][6] and [Posthuman Studios, LLC][7] (the humanoid person). The latter is distributed under the [CC Attribution-NonCommercial-ShareAlike 3.0 Unported][8] license and thus, the videos are also available under the terms of that license.

<script>
  let demoVideos = [];
  const demoVideoTmpl = document.querySelector("#demo-video");

  function loadDemoVideo(index) {
    const node = document.importNode(demoVideoTmpl, true).content;
    node.querySelector("source").setAttribute("src", demoVideos[index].src);
    node.querySelector("span").textContent = demoVideos[index].description;
    const len = demoVideos.length;
    node.querySelector(".left").addEventListener("click", e => {
      loadDemoVideo((index + len - 1) % len);
      e.preventDefault();
    });
    node.querySelector(".right").addEventListener("click", e => {
      loadDemoVideo((index + 1) % len);
      e.preventDefault();
    });
    const cur = document.querySelector(".video-demo");
    const replacement = cur.cloneNode(false);
    replacement.appendChild(node);
    cur.parentNode.replaceChild(replacement, cur);
  }

  window.addEventListener('DOMContentLoaded', async (e) => {
    const response = await fetch("/media/demo.json");
    demoVideos = await response.json();
    loadDemoVideo(0);
    const fig = document.querySelector(".video-demo");
  });
</script>

 [1]: https://www.pexels.com/photo/gray-bridge-and-trees-814499/
 [2]: https://www.pexels.com/photo/adventure-alps-amazing-beautiful-552785/
 [3]: https://www.pexels.com/photo/trees-in-the-middle-of-body-of-water-819699/
 [4]: https://www.gnu.org/licenses/gpl-3.0.en.html
 [5]: https://github.com/QuestScreen/QuestScreen
 [6]: https://pixabay.com/illustrations/clock-time-time-of-fantasy-1516967/
 [7]: https://eclipsephase.com/resources
 [8]: https://creativecommons.org/licenses/by-nc-sa/3.0/