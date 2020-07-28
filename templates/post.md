+++
author = "Brian Pfeil"
categories = ["{{ .Repo.Language }}", "playground"]
date = {{ .Repo.CreatedAt.Format "2006-01-02" }}
description = ""
summary = " "
draft = false
slug = "{{ .Slug }}"
tags = [{{range $val := .Tags}}"{{$val}}",{{end}}]
title = "{{ .Title }}"
repoFullName = "{{ .Repo.FullName }}"
repoHTMLURL = "{{ .Repo.HTMLURL }}"
truncated = true

+++

<div class="alert alert-info small bg-info" role="alert">
<span class="text-muted">code for article</span>&nbsp;<a href="{{ .Repo.HTMLURL }}" target="_blank"><i class="fab fa-github fa-sm"></i>&nbsp;{{ .Repo.FullName }}</a>
</div>

{{ .MarkdownBody }}


