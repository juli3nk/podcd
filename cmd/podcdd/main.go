package main

import (
	"github.com/juli3nk/go-utils"
)

func main() {
	cmd := newCommand()

	cmd.SetHelpTemplate(helpTemplate)
	cmd.SetUsageTemplate(usageTemplate)

	if err := cmd.Execute(); err != nil {
		utils.Exit(err)
	}
}

var usageTemplate = `{{ .Short | trim }}

Usage:{{ if .Runnable }}
  {{ if .HasAvailableFlags }}{{ appendIfNotPresent .UseLine "[flags]" }}{{ else }}{{ .UseLine }}{{ end }}{{ end }}{{ if .HasAvailableSubCommands }}
  {{ .CommandPath }} [command]{{ end }}{{ if gt .Aliases 0 }}

Aliases:
  {{ .NameAndAliases }}{{ end }}{{ if .HasExample }}

Examples:
  {{ .Example }}{{ end }}{{ if .HasAvailableSubCommands }}

Available Commands:{{ range .Commands }}{{ if .IsAvailableCommand }}
  {{ rpad .Name .NamePadding }} {{ .Short }}{{ end }}{{ end }}{{ end }}{{ if .HasAvailableLocalFlags }}

Flags:
{{ .LocalFlags.FlagUsages | trimRightSpace }}{{ end }}{{ if .HasAvailableInheritedFlags }}

Global Flags:
  {{ .InheritedFlags.FlagUsages | trimRightSpace }}{{ end }}{{ if .HasHelpSubCommands }}

Additional help topics:{{ range .Commands }}{{ if .IsHelpCommand }}
  {{ rpad .CommandPath .CommandPathPadding }} {{ .Short }}{{ end }}{{ end }}{{ end }}{{ if .HasAvailableSubCommands }}

Use "{{ .CommandPath }} [command] --help" for more information about a command.{{ end }}
`

var helpTemplate = `
{{ if or .Runnable .HasSubCommands }}{{ .UsageString }}{{ end }}`
