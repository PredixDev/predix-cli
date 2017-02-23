package cli

import (
	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.com/urfave/cli"
)

var AppHelpTemplate = h("NAME:") + `
   {{.Name}} - {{.Usage}}

` + h("USAGE:") + `
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}
   {{if .Version}}{{if not .HideVersion}}
` + h("VERSION:") + `
   {{.Version}}
   {{end}}{{end}}{{if len .Authors}}
` + h("AUTHOR(S):") + `
   {{range .Authors}}{{.}}{{end}}
   {{end}}{{if .VisibleCommands}}
` + h("COMMANDS:") + `{{range .VisibleCategories}}{{if .Name}}
   ` + t("{{.Name}}:") + `{{end}}{{range .VisibleCommands}}
     {{.Name}}{{with .ShortName}}, {{.}}{{end}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{end}}
` + h("ENVIRONMENT VARIABLES:") + `
  PREDIX_NO_CACHE=true{{"\t"}}Do not use the cache to lookup apps and services
  PREDIX_NO_CF_BYPASS=true{{"\t"}}Do not try to run an unknown command as a CF CLI command
  {{if .VisibleFlags}}
` + h("GLOBAL OPTIONS:") + `
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright}}
` + h("COPYRIGHT:") + `
COPYRIGHT:
   {{.Copyright}}
   {{end}}

` + h("NOTE: This is a beta release under Predix Labs.") + `
`

var CommandHelpTemplate = h("NAME:") + `
   {{.HelpName}} - {{.Usage}}
` + h("USAGE:") + `
   {{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{if .Category}}
` + h("CATEGORY:") + `
   {{.Category}}{{end}}{{if .Description}}
` + h("DESCRIPTION:") + `
   {{.Description}}{{end}}{{if .VisibleFlags}}
` + h("OPTIONS:") + `
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

var SubcommandHelpTemplate = h("NAME:") + `
   {{.HelpName}} - {{.Usage}}
` + h("USAGE:") + `
   {{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
` + h("COMMANDS:") + `{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{.Name}}{{with .ShortName}}, {{.}}{{end}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{if .VisibleFlags}}
` + h("OPTIONS:") + `
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

func h(s string) string {
	return terminal.HeaderColor(s)
}

func t(s string) string {
	return terminal.TableContentHeaderColor(s)
}

func init() {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate
	cli.SubcommandHelpTemplate = SubcommandHelpTemplate
}
