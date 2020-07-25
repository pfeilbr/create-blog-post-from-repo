+++
author = "Brian Pfeil"
categories = ["{{ .Repo.Language }}", "playground"]
date = {{ .Repo.CreatedAt.Format "2006-01-02" }}
description = ""
draft = false
slug = "{{ .Repo.Name }}"
tags = ["playground"]
title = "{{ .Title }}"

+++

{{ .MarkdownBody }}
