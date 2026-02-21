package main

import (
	"embed"
	"errors"
	"os"

	"gitea.roboalch.com/roboalchemist/tfc/cmd"
	"gitea.roboalch.com/roboalchemist/tfc/pkg/output"
)

// version is set via ldflags at build time: -X main.version=x.y.z
var version = "dev"

//go:embed skill/SKILL.md
var skillMD string

//go:embed skill/reference/commands.md
var commandsRef string

//go:embed skill
var skillFS embed.FS

func main() {
	cmd.SetVersion(version)
	cmd.SetSkillData(skillMD, commandsRef, skillFS)
	if err := cmd.Execute(); err != nil {
		var se *output.StructuredError
		if errors.As(err, &se) {
			os.Exit(se.ExitCode)
		}
		os.Exit(1)
	}
}
