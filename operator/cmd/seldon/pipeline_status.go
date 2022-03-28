package main

import (
	"os"

	"github.com/seldonio/seldon-core/operatorv2/pkg/cli"
	"github.com/spf13/cobra"
)

func createPipelineStatus() *cobra.Command {
	cmdPipelineUnload := &cobra.Command{
		Use:   "status",
		Short: "status of a pipeline",
		Long:  `status of a pipeline`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			schedulerHost, err := cmd.Flags().GetString(schedulerHostFlag)
			if err != nil {
				return err
			}
			schedulerPort, err := cmd.Flags().GetInt(schedulerPortFlag)
			if err != nil {
				return err
			}
			pipelineName, err := cmd.Flags().GetString(pipelineNameFlag)
			if err != nil {
				return err
			}
			verbose, err := cmd.Flags().GetBool(verboseFlag)
			if err != nil {
				return err
			}
			schedulerClient := cli.NewSchedulerClient(schedulerHost, schedulerPort)
			err = schedulerClient.PipelineStatus(pipelineName, verbose)
			return err
		},
	}
	cmdPipelineUnload.Flags().StringP(pipelineNameFlag, "p", "", "pipeline name for status")
	if err := cmdPipelineUnload.MarkFlagRequired(pipelineNameFlag); err != nil {
		os.Exit(-1)
	}
	cmdPipelineUnload.Flags().String(schedulerHostFlag, "0.0.0.0", "seldon scheduler host")
	cmdPipelineUnload.Flags().Int(schedulerPortFlag, 9004, "seldon scheduler port")
	return cmdPipelineUnload
}