package cmd

import (
	"github.com/geniusmonkey.com/jump/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
)

var forceTunnel = false

var tunnelCmd = &cobra.Command{
	Use:   "tunnel",
	Short: "manage port forwarding tunnel+",
}

var tunnelShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "manage remote ssh connections",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Read(config.Location)
		if err != nil {
			log.Fatal(err)
		}

		if len(args) == 1 {
			tunnel := cfg.GetTunnel(args[0])
			if tunnel == nil {
				log.Fatalf("tunnel with name %s does not exist", args[0])
			}
			printTunnel(args[0], *tunnel)
		} else {
			for name, t := range cfg.Tunnels {
				printTunnel(name, t)
			}
		}
	},
}

var tunnelAddCmd = &cobra.Command{
	Use:   "add [name] [localPort] [remoteAddr] [remotePort]",
	Short: "add a new remote",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Read(config.Location)
		if err != nil {
			log.Fatal(err)
		}

		name := args[0]
		if t := cfg.GetTunnel(name); t != nil && !forceTunnel {
			log.Fatalf("a tunnel with the name %v already exists")
		}

		locPort, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("invalid localPort, must be a number")
		}
		rmtPort, err := strconv.Atoi(args[3])
		if err != nil {
			log.Fatalf("invalid remotePort, must be a number")
		}

		tunnel := config.Tunnel{
			Addr:       args[2],
			LocalPort:  int16(locPort),
			RemotePort: int16(rmtPort),
		}

		cfg.AddTunnel(name, tunnel)
		if _, err := config.Write(config.Location, cfg); err != nil {
			log.Fatalf("failed to save config, %v", err)
		} else {
			log.Printf("added new tunnel %v to config", name)
		}
	},
}

var tunnelRemoveCmd = &cobra.Command{
	Use:   "rm [name]",
	Short: "remove tunnel",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Read(config.Location)
		if err != nil {
			log.Fatal(err)
		}

		name := args[0]
		if _, exists := cfg.Tunnels[name]; !exists {
			log.Fatalf("tunnel with name %v does not exist", name)
		}

		delete(cfg.Tunnels, name)

		if _, err := config.Write(config.Location, cfg); err != nil {
			log.Fatalf("failed to write config file, %v", err)
		} else {
			log.Infof("removed tunnel %v", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(tunnelCmd)

	tunnelCmd.AddCommand(tunnelAddCmd)
	tunnelAddCmd.Flags().BoolVarP(&forceTunnel, "force", "f", false, "force the addition of a tunnel even if it exists")

	tunnelCmd.AddCommand(tunnelShowCmd)
	tunnelCmd.AddCommand(tunnelRemoveCmd)
}

func printTunnel(name string, tunnel config.Tunnel) {
	log.Infof("%s: forwarding port %v -> %v:%v", name, tunnel.LocalPort, tunnel.Addr, tunnel.RemotePort)
}
