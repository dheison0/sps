package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net"
	"sps/pkg"
)

const DefaultPort = 8888

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Long: `Start the server using the default configuration if a config
file or flags are not especified`,
	Run: start,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func start(cmd *cobra.Command, args []string) {
	server, err := net.ListenTCP("tcp", &net.TCPAddr{Port: DefaultPort})
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	fmt.Printf("Server started at port %d!\n", DefaultPort)
	for {
		client, err := server.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		go pkg.ProccessRequest(client)
	}
}
