package cmd

import (
	"fmt"
	"github.com/geniusmonkey.com/jump/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
)

var connectCmd = &cobra.Command{
	Use:   "connect [remote] <tunnel>",
	Short: "connects to a remote and open a tunnel",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Read(config.Location)
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

func init() {
	rootCmd.AddCommand(connectCmd)
}
