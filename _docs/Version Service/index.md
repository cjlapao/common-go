---
layout: default
title: Version Service
nav_order: 2
search_enabled: true
has_children: true
has_toc: false
---
# Version Service

Designed to help generate versioning in a easy package to either display in a command line type application or as a maintainer

# Usage

To create a new version service you can use the function GET()

```go
ver := version.Get()
ver.Author = "Carlos Lapao"
```
