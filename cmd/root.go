package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/zrl/gifmoji/internal"
	"github.com/zrl/gifmoji/internal/infinite"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gifmoji",
		Short: "Make an image into a gif",
		Long:  ``,
	}

	infiniteCmd = &cobra.Command{
		Use:   "infinite filename",
		Short: "Make an infinitely scrolling gif",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return transform(&infinite.Transformer{}, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(infiniteCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func transform(t internal.Transformer, args []string) error {
	return t.Transform(args)
}
