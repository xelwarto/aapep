package cmd

import (
	"aapep/util"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

type rootCmdFlags struct {
	ClientCount int
	Interval    int
	Timer       string
}

var wg sync.WaitGroup
var rootFlags rootCmdFlags

func init() {
	rootCmd.PersistentFlags().IntVarP(&rootFlags.ClientCount, "clients", "c", 10, "Number of clients to start.")
	rootCmd.PersistentFlags().IntVarP(&rootFlags.Interval, "interval", "i", 1000, "Client interval time in milliseconds.")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.Timer, "timer", "d", "10s", "Time duration for running the test.")
}

var rootCmd = &cobra.Command{
	Use:   "aapep",
	Short: "Aapep - Network Stress Tester",
	Long:  "aapep - A simple tool used for stress (load) testing systems on a network.",
}

func Execute(v string) {
	fmt.Println(util.BannerHelp(v))
	rootCmd.Version = v
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func showHelp(cmd *cobra.Command, err string) {
	cmd.Help()
	if err != "" {
		fmt.Printf("\n%v\n", util.ErrorMsgStyle.Render(err))
	}
	os.Exit(1)
}
