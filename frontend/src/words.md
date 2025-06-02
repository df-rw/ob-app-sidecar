---
title: words
---

# Word searching

This is an active word search using [htmx](https://htmx.org) and
[hyperscript](https://hyperscript.org).

<div hx-get="/api/words"
     hx-trigger="load"
     hx-swap="outerHTML"
     _="on htmx:responseError set my.innerHTML to event.detail.error"
     ></div>

