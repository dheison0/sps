package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net"
	"sps/pkg"
)

var Port uint16 = 8888

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Long: `Start the server using the default configuration if a config
file or flags are not especified`,
	Run: start,
}

func init() {
	startCmd.Flags().Uint16VarP(&Port, "port", "p", 8888, "The port to listen the server")
	rootCmd.AddCommand(startCmd)
}

func start(cmd *cobra.Command, args []string) {
	server, err := net.ListenTCP("tcp", &net.TCPAddr{Port: int(Port)})
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	fmt.Printf("Server started at port %d!\n", Port)
	for {
		client, err := server.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		go pkg.ProccessRequest(client)
	}
}
