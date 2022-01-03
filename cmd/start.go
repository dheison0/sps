package cmd

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"log"
	"net"
	"sps/pkg"
	"sps/pkg/forwards"
	"sps/util"
	"sps/types"
)

var config = types.Config{}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Long: `Start the server using the default configuration if a config
file or flags are not especified`,
	Run: start,
}

func init() {
	startCmd.Flags().IntVarP(
		&config.Main.Port,
		"port",
		"p",
		8888,
		"The port to listen the server",
	)
	startCmd.Flags().StringVarP(
		&config.Filter.File,
		"filter",
		"f",
		"",
		"A simple text file containing the filters to match",
	)
	startCmd.Flags().StringVarP(
		&config.ConfigFile,
		"config",
		"c",
		"",
		"A TOML file containing the configuration",
	)
	rootCmd.AddCommand(startCmd)
}

func start(cmd *cobra.Command, args []string) {
	if config.ConfigFile != "" {
		file, err := util.ReadFile(config.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		err = toml.Unmarshal(file, &config)
		if err != nil {
			log.Fatal(err)
		}
	}
	forwards.SetConfigAndParse(config.Filter)
	server, err := net.ListenTCP("tcp", &net.TCPAddr{Port: config.Main.Port})
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	fmt.Printf("Server started at port %d!\n", config.Main.Port)
	for {
		client, err := server.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		go pkg.ProccessRequest(client, config.Filter.EnableRegex)
	}
}
