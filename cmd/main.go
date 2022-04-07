// The program supports scanning Istio configuration and analzye CVE.
package main

// Plan of the record.
// Basic parser for Istio authoriation policy.
// - [ ] Check reports the total number of the policy scanned. (58 security policies, 32 networking policies.)
// - [ ] Consider to add more accurate error message tied to the field.
// Analyze the Istio cluster.
// scanner --kube <kubeconfig-path>

import (
	"github.com/tetratelabs/istio-security-scanner/pkg/k8s"

	"github.com/spf13/cobra"
	"istio.io/pkg/log"
)

var (
	kubeConfigPath = "."
	runOnce        = false
	loggingOptions log.Options

	scannerCmd = &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			RunAll(&analyerOptions{
				KubeConfig: kubeConfigPath,
			})
		},
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			// TODO(incfly): this seems returning an error if print out. Check why.
			_ = log.Configure(&loggingOptions)
			return nil
		},
	}
)

type analyerOptions struct {
	Dir             string
	KubeConfig      string
	CVEDatabaseURL  string
	CVEDatabasePath string
}

func RunAll(options *analyerOptions) {
	c, err := k8s.NewClient(options.KubeConfig, runOnce)
	if err != nil {
		log.Fatalf("error %v", err)
	}
	stopCh := make(chan struct{})
	c.Run(stopCh)
}

func init() {
	flags := scannerCmd.Flags()
	flags.StringVarP(&kubeConfigPath, "config", "c", "~/.kube/config", "The path to the kubeconfig of a cluster to be analyzed.")
	flags.BoolVar(&runOnce, "once", true, "Whether running the scanning only one shot. If false, will continue in a loop")
	// By default if `--log_output_level` is not specified by users, we supress the output to make report clean.
	// Setting "default" is needed, otherwise setting "kube" alone does not work, due to issue possible
	// log package itself. That make log output level field as ",kube:none", no effect.
	loggingOptions.SetOutputLevel("default", log.InfoLevel)
	loggingOptions.SetOutputLevel("kube", log.NoneLevel)
	loggingOptions.AttachCobraFlags(scannerCmd)
}

func main() {
	if err := scannerCmd.Execute(); err != nil {
		log.Fatalf("failed to run scanner: %v", err)
	}
}
