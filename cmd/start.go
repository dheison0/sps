package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net"
	"sps/pkg"
	"sps/util"
	"sps/pkg/forwards"
)

var port uint16
var filterFile string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Long: `Start the server using the default configuration if a config
file or flags are not especified`,
	Run: start,
}

func init() {
	startCmd.Flags().Uint16VarP(
		&port,
		"port",
		"p",
		8888,
		"The port to listen the server",
	)
	startCmd.Flags().StringVarP(
		&filterFile,
		"filter",
		"f",
		"",
		"A text file containing the filters to match",
	)
	rootCmd.AddCommand(startCmd)
}

func start(cmd *cobra.Command, args []string) {
	if filterFile != "" {
		fmt.Println("Processing filter file...")
		file, err := util.ReadFile(filterFile)
		if err != nil {
			log.Fatal(err)
		}
		lines := util.ReadLinesFromBytes(file)
		for _, l := range lines {
			forwards.AddFilter(string(l))
		}
	}
	server, err := net.ListenTCP("tcp", &net.TCPAddr{Port: int(port)})
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	fmt.Printf("Server started at port %d!\n", port)
	for {
		client, err := server.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		go pkg.ProccessRequest(client)
	}
}
