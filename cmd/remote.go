package cmd

import (
	"fmt"
	"github.com/geniusmonkey/jump/config"
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)

var rmtFlags = remoteFlags{}

type remoteFlags struct {
	ident string
	force bool
	user  string
	port  int16
}

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "manage remote ssh connections",
}

var remoteShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "manage remote ssh connections",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Read(config.Location)
		if err != nil {
			log.Fatal(err)
		}

		if len(args) == 1 {
			remote := cfg.GetRemote(args[0])
			if remote == nil {
				log.Fatalf("remote with name %s does not exist", args[0])
			}
			printRemote(args[0], *remote)
		} else {
			for name, r := range cfg.Remotes {
				printRemote(name, r)
			}
		}
	},
}

var remoteAddCmd = &cobra.Command{
	Use:   "add [name] [remoteAddr]",
	Short: "add a new remote",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Flags().PrintDefaults()
		cfg, err := config.Read(config.Location)
		if err != nil {
			cfg = config.Config{}
		}

		name := args[0]
		addr := args[1]

		if r := cfg.GetRemote(name); r != nil && !rmtFlags.force {
			log.Fatalf("a remote with name %s already exist", name)
		}

		remote := config.Remote{
			Addr: addr,
			Port: rmtFlags.port,
		}

		if rmtFlags.ident != "" {
			remote.IdentFile = &rmtFlags.ident
		}

		if rmtFlags.user != "" {
			remote.User = &rmtFlags.user
		}

		cfg.AddRemote(name, remote)
		if _, err := config.Write(config.Location, cfg); err != nil {
			log.Fatalf("failed to save remote to config, %v", err)
		}
		log.Printf("added new remote %s", name)
	},
}

var remoteRmCmd = &cobra.Command{
	Use:   "rm [name]",
	Short: "remove remote",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Read(config.Location)
		if err != nil {
			log.Fatal(err)
		}

		name := args[0]
		if _, exists := cfg.Remotes[name]; !exists {
			log.Fatalf("remote with name %v does not exist", name)
		}

		delete(cfg.Remotes, name)

		if _, err := config.Write(config.Location, cfg); err != nil {
			log.Fatalf("failed to write config file, %v", err)
		} else {
			log.Infof("removed remote %v", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(remoteCmd)

	remoteCmd.AddCommand(remoteAddCmd)
	remoteAddCmd.Flags().StringVarP(&rmtFlags.ident, "identity", "i", "", "identity file")
	remoteAddCmd.Flags().StringVarP(&rmtFlags.user, "user", "u", "", "")
	remoteAddCmd.Flags().Int16VarP(&rmtFlags.port, "port", "p", 22, "")
	remoteAddCmd.Flags().BoolVarP(&rmtFlags.force, "force", "f", false, "force adding a remote if it already exists")

	remoteCmd.AddCommand(remoteShowCmd)
	remoteCmd.AddCommand(remoteRmCmd)
}

func printRemote(name string, remote config.Remote) {
	fmt.Println(name)
	fmt.Printf("  Address: %s\n", remote.Addr)
	fmt.Printf("  Port: %v\n", remote.Port)

	if remote.User != nil {
		fmt.Printf("  User: %v\n", *remote.User)
	}
	if remote.IdentFile != nil {
		fmt.Printf("  IdentityFile: %v\n", *remote.IdentFile)
	}
}
