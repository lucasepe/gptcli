package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/lucasepe/gptcli/internal/app"
	"github.com/lucasepe/gptcli/internal/completer"
	"github.com/lucasepe/gptcli/internal/executor"
	"github.com/lucasepe/gptcli/internal/shortcuts"
)

const (
	banner = `╔═╗╔═╗╔╦╗
║ ╦╠═╝ ║ 
╚═╝╩   ╩ CLI {{VER}} (rev: {{REV}})`
)

var (
	Version string
	Build   string
)

func main() {
	cfg, err := app.ConfigPath()
	exitOnErr(err)

	shc, err := shortcuts.FromFile(cfg)
	exitOnErr(err)

	header := strings.Replace(banner, "{{VER}}", Version, 1)
	header = strings.Replace(header, "{{REV}}", Build, 1)

	fmt.Println(header)
	fmt.Println()
	fmt.Println("An interactive ChatGPT client featuring shortcuts and auto-complete.")
	fmt.Println()

	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
	defer fmt.Println("Bye!")
	p := prompt.New(
		executor.GPT(shc),
		completer.FromShortcuts(shc),
		prompt.OptionTitle("gptcli: interactive ChatGPT client"),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.Red),
	)
	p.Run()
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
