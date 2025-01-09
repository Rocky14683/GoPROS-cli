package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func compile(directory string, buildArgs []string) int {
	env := os.Environ()
	prosToolchain := os.Getenv("PROS_TOOLCHAIN")
	if prosToolchain != "" {
		env = append([]string{fmt.Sprintf("PATH=%s%s%s",
			filepath.Join(prosToolchain, "bin"), string(os.PathListSeparator),
			os.Getenv("PATH"))}, env...)
	}

	var makeCmd string
	if os.Getenv("PROS_TOOLCHAIN") != "" && os.PathSeparator == '\\' {
		makeCmd = filepath.Join(prosToolchain, "bin", "make.exe")
	} else {
		makeCmd = "make"
	}

	fmt.Println("Running make command: ", makeCmd, buildArgs)

	cmd := exec.Command(makeCmd, buildArgs...)
	cmd.Dir = directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("could not run command: ", err)
		return 1
	}

	return cmd.ProcessState.ExitCode()
}

func main() {
	app := &cli.App{

		Name:                 "gopros",
		Version:              "0.0.1",
		Usage:                "CLI for PROS",
		Authors:              []*cli.Author{{Name: "Purdue ACM SIGBots", Email: "sig.robotics.purdue@gmail.com"}},
		UsageText:            "gopros [OPTIONS] COMMAND [ARGS]...",
		EnableBashCompletion: true,
		Suggest:              true,

		Commands: []*cli.Command{
			{
				Name:      "make",
				Aliases:   []string{"m", "build"},
				Usage:     "Build current PROS project or cwd",
				Args:      true,
				ArgsUsage: "Compile arguments",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "Multithreading",
						Aliases: []string{"j"},
						Usage:   "Enable multithreading",
					},
				},
				Action: func(cCtx *cli.Context) error {
					cwd, err := os.Getwd()
					if err != nil {
						return fmt.Errorf("could not get current working directory: %v", err)
					}

					buildArgs := cCtx.Args().Slice()
					if cCtx.Bool("Multithreading") {
						buildArgs = append(buildArgs, "-j")
					}

					exitCode := compile(cwd, buildArgs)
					os.Exit(exitCode)
					return nil
				},
			}, {
				Name:    "conductor",
				Aliases: []string{"c"},
				Usage:   "Run PROS Conductor",
				Subcommands: []*cli.Command{
					{
						Name:    "info-project",
						Aliases: []string{"info", "ip"},
						Usage:   "Get information about the current project",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("pros c info-project is runned")
							return nil
						},
					},
					{
						Name:    "query-template",
						Aliases: []string{"query-temp", "qt"},
						Usage:   "Show all the templates available in the conductor",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("pros c query-template is runned")
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
