package cmd

import (
	"errors"

	ia "github.com/iter8-tools/iter8/action"
	"github.com/iter8-tools/iter8/base/log"
	"github.com/iter8-tools/iter8/driver"
	"github.com/spf13/cobra"
)

// kAssertDesc is the description of the k assert cmd
const kAssertDesc = `
Assert if the result of a Kubernetes experiment satisfies the specified conditions. If all conditions are satisfied, the command exits with code 0. Else, the command exits with code 1. 

Assertions are especially useful for automation inside CI/CD/GitOps pipelines.

Supported conditions are 'completed', 'nofailure', 'slos', which indicate that the experiment has completed, none of the tasks have failed, and the SLOs are satisfied.

	iter8 k assert -c completed -c nofailure -c slos
	# same as iter8 k assert -c completed,nofailure,slos

You can optionally specify a timeout, which is the maximum amount of time to wait for the conditions to be satisfied:

	iter8 k assert -c completed,nofailure,slos -t 5s
`

// newAssertCmd creates the Kubernetes assert command
func newKAssertCmd(kd *driver.KubeDriver) *cobra.Command {
	actor := ia.NewAssertOpts(kd)

	cmd := &cobra.Command{
		Use:          "assert",
		Short:        "Assert if Kubernetes experiment result satisfies conditions",
		Long:         kAssertDesc,
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			allGood, err := actor.KubeRun()
			if err != nil {
				return err
			}
			if !allGood {
				e := errors.New("assert conditions failed")
				log.Logger.Error(e)
				return e
			}
			return nil
		},
	}
	// options specific to k assert
	addExperimentGroupFlag(cmd, &actor.Group)
	actor.EnvSettings = settings

	// options shared with assert
	addConditionFlag(cmd, &actor.Conditions)
	addTimeoutFlag(cmd, &actor.Timeout)
	return cmd
}

// initialize with k assert cmd
func init() {
	kCmd.AddCommand(newKAssertCmd(kd))
}
