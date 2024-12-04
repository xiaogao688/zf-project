package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use: "print [OPTIONS] [COMMANDS]",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run print...")
		fmt.Printf("printFlag: %v\n", printFlag)
		fmt.Printf("LocalLogo: %v\n", LocalLogo)
	},
}

var printFlag string

func init() {
	rootCmd.AddCommand(printCmd)

	// 本地标志
	printCmd.Flags().StringVarP(&printFlag, "flag", "f", "", "print flag for local")
}
