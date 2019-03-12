package cmd

import (
	"fmt"
	"github.com/geniusmonkey.com/jump/config"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"sort"
)

var interactive bool
var connectCmd = &cobra.Command{
	Use:   "connect [remote] <tunnel>",
	Short: "connects to a remote and open a tunnel",
	Args: func(cmd *cobra.Command, args []string) error {
		if interactive {
			return nil
		} else {
			return cobra.RangeArgs(1, 2)(cmd, args)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Read(config.Location)
		if interactive {
			remote := getRemoteInteractive(cfg)
			tunnel := getTunnelInteractive(cfg)
			args = []string{remote, tunnel}

		}

		if err != nil {
			log.Fatal(err)
		}

		sshArgs := make([]string, 0)
		remote := cfg.GetRemote(args[0])
		if remote == nil {
			log.Fatalf("remote with name %s does not exist", args[0])
		}
		msg := fmt.Sprintf("Connecting to %s", remote.Addr)

		if len(args) == 2 {
			tunnel := cfg.GetTunnel(args[1])
			if tunnel == nil {
				log.Fatalf("tunnel with name %s does not exist", args[1])
			} else {
				msg = fmt.Sprintf("%s forwarding %v -> %v:%v", msg, tunnel.LocalPort, tunnel.Addr, tunnel.RemotePort)
			}
			sshArgs = append(sshArgs, "-L", fmt.Sprintf("%v:%v:%v", tunnel.LocalPort, tunnel.Addr, tunnel.RemotePort))
		}

		if remote.IdentFile != nil {
			sshArgs = append(sshArgs, "-i", *remote.IdentFile)
		}

		if remote.User != nil {
			sshArgs = append(sshArgs, fmt.Sprintf("%v@%v", *remote.User, remote.Addr))
		} else {
			sshArgs = append(sshArgs, remote.Addr)
		}

		sshArgs = append(sshArgs, "-p", fmt.Sprintf("%v", remote.Port))

		command := exec.Command("ssh", sshArgs...)
		log.Print(msg)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Stdin = os.Stdin
		if err := command.Run(); err != nil {
			log.Printf("failed to connect, %v", err)
		}
	},
}

func getTunnelInteractive(cfg config.Config) string {
	names := make([]string, 0)
	for name := range cfg.Tunnels {
		names = append(names, name)
	}
	sort.Strings(names)

	prompt := promptui.Select{
		Label: "Tunnel",
		Items: names,
	}
	_, tunnel, err := prompt.Run()
	if err != nil {
		log.Fatalf("failed to select tunnel, %v", err)
	}
	return tunnel
}

func getRemoteInteractive(cfg config.Config) string {
	names := make([]string, 0)
	for name := range cfg.Remotes {
		names = append(names, name)
	}

	sort.Strings(names)

	prompt := promptui.Select{
		Label: "Remote",
		Items: names,
	}
	_, remote, err := prompt.Run()
	if err != nil {
		log.Fatalf("failed to select remote, %v", err)
	}
	return remote
}

func init() {
	connectCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "will list options for remotes and tunnels")
	rootCmd.AddCommand(connectCmd)
}
