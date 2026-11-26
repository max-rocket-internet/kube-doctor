// Command kube-doctor gives your Kubernetes cluster a health checkup.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/urfave/cli/v3"

	"github.com/max-rocket-internet/kube-doctor/internal/doctor"
)

// version is set at build time with -ldflags "-X main.version=...".
var version string

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "💣", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cmd := &cli.Command{
		Name:    "kube-doctor",
		Usage:   "Gives your Kubernetes cluster a health checkup",
		Suggest: true,
		Version: buildVersion(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "namespace",
				Usage: "A namespace name to check. e.g. 'kube-system'",
			},
			&cli.StringFlag{
				Name:  "label-selector",
				Usage: "A label selector to check. e.g. 'app.kubernetes.io/name=prometheus'",
			},
			&cli.BoolFlag{
				Name:  "non-namespaced-resources",
				Usage: "Whether to check non-namespaced resources like Nodes, PersistentVolumes, apiserver health etc",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug logging",
			},
			&cli.BoolFlag{
				Name:  "warning-symptoms",
				Usage: "Whether to show warning symptoms (otherwise only critical are shown)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return doctor.Run(ctx, doctor.Options{
				Namespace:       cmd.String("namespace"),
				LabelSelector:   cmd.String("label-selector"),
				NonNamespaced:   cmd.Bool("non-namespaced-resources"),
				Debug:           cmd.Bool("debug"),
				WarningSymptoms: cmd.Bool("warning-symptoms"),
			})
		},
	}

	return cmd.Run(ctx, os.Args)
}

// buildVersion returns the version injected via ldflags, falling back to the
// module version recorded by `go install`.
func buildVersion() string {
	if version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "dev"
}
