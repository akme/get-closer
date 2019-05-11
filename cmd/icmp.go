// Copyright © 2019 Vyacheslav Akhmetov <iam@itaddict.ru>

package cmd

import (
	"github.com/spf13/cobra"
)

// icmpCmd represents the icmp command
var icmpCmd = &cobra.Command{
	Use:     "icmp",
	Aliases: []string{"icmp", "ping"},
	Short:   "Measuring RTT for hosts from list.",
	Long:    `Measuring RTT for hosts from list.`,
	/* 	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("icmp called")
	}, */

	Run: startMeasurements,
}

func init() {
	rootCmd.AddCommand(icmpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// icmpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// icmpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
