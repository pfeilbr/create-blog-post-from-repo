+++
author = "Brian Pfeil"
categories = ["{{ .Repo.Language }}", "playground"]
date = {{ .Repo.CreatedAt.Format "2006-01-02" }}
description = ""
summary = "{{ .Summary }}"
draft = false
slug = "{{ .Repo.Name }}"
tags = [{{range $val := .Tags}}"{{$val}}",{{end}}]
title = "{{ .Title }}"
repoFullName = "{{ .Repo.FullName }}"
repoHTMLURL = "{{ .Repo.HTMLURL }}"
truncated = true

+++

{{ .MarkdownBody }}
