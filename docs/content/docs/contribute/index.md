---
title: 'Contribute'
date: 2024-08-22T16:35:17-06:00
draft: false
weight: 11
---

This page is for describing how to add, improve, and update ACDC documentation.

### Installation
The documentation is generated using the [Hugo](https://gohugo.io) documentation generator. We recommend running below commands for each machine.

#### macOS/Linux

```bash
$ brew install hugo
```

#### Windows

```bash
$ choco install hugo-extended
```

Note that Hugo with Extended edition is required for building the documentation with [current theme](https://themes.gohugo.io/themes/hugo-whisper-theme/). Make sure you see **/extended** after the version number, once you run `hugo version`.


### Getting started

#### Add new page
To add a new page, `cd docs` and run:

```bash
$ hugo new content content/docs/[new_page_name]/index.md
```

Open the file with your editor and add contents.

#### View site
To test your site, run:

```bash
$ hugo server
```

Now enter `localhost:1313` in the address bar of your browser. The server rebuilds your site and refreshes your browser whenever it detects the change.