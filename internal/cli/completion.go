// Package cli provides the command-line interface for LazyFocus.
package cli

import (
	"github.com/spf13/cobra"
)

// NewCompletionCommand creates the completion command for shell completion scripts
func NewCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for lazyfocus.

To load completions:

Bash:
  $ source <(lazyfocus completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ lazyfocus completion bash > /etc/bash_completion.d/lazyfocus
  # macOS:
  $ lazyfocus completion bash > $(brew --prefix)/etc/bash_completion.d/lazyfocus

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ lazyfocus completion zsh > "${fpath[1]}/_lazyfocus"

  # You may need to start a new shell for this setup to take effect.

Fish:
  $ lazyfocus completion fish | source
  # To load completions for each session, execute once:
  $ lazyfocus completion fish > ~/.config/fish/completions/lazyfocus.fish

PowerShell:
  PS> lazyfocus completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> lazyfocus completion powershell > lazyfocus.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Annotations: map[string]string{
			"skipServiceSetup": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.Root().OutOrStdout()
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(out)
			case "zsh":
				return cmd.Root().GenZshCompletion(out)
			case "fish":
				return cmd.Root().GenFishCompletion(out, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(out)
			}
			return nil
		},
	}

	return cmd
}
