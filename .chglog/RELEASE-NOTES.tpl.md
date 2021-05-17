{{ if .Unreleased.CommitGroups -}}
<a name="unreleased"></a>
## Unreleased

{{ range .Unreleased.CommitGroups -}}
### {{ .Title }}
{{ range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }} ({{ .Hash.Long }}){{ if .Refs }}{{ range .Refs }} #{{ .Ref }}{{ end }}{{ end }}
{{ end }}
{{ end -}}

{{- if .Unreleased.RevertCommits -}}
### Reverts
{{ range .Unreleased.RevertCommits -}}
- {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if .Unreleased.MergeCommits -}}
### Pull Requests
{{ range .Unreleased.MergeCommits -}}
- {{ .Header }}
{{ end }}
{{ end -}}

{{- if .Unreleased.NoteGroups -}}
{{ range .Unreleased.NoteGroups -}}
### {{ .Title }}
{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}

{{ else -}}
<a name="changelog"></a>
## Changelog
{{ range .Versions }}
{{ range .CommitGroups -}}
### {{ .Title }}
{{ range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }} ({{ .Hash.Long }}){{ if .Refs }}{{ range .Refs }} #{{ .Ref }}{{ end }}{{ end }}
{{ end }}
{{ end -}}

{{- if .RevertCommits -}}
### Reverts
{{ range .RevertCommits -}}
- {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if .MergeCommits -}}
### Pull Requests
{{ range .MergeCommits -}}
- {{ .Header }}
{{ end }}
{{ end -}}

{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
### {{ .Title }}
{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}
{{ end -}}
