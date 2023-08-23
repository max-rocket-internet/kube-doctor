package main

import (
	"log"
	"os"

	"github.com/max-rocket-internet/kube-doctor/pkg/doctor"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "kube-doctor",
		Usage:   "Gives your Kubnernetes cluster a health checkup",
		Suggest: true,
		Version: "0.1",
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
				Value: false,
				Usage: "Whether to check non-namespaced resources like Nodes, PersistentVolumes, apiserver health etc",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "Enable debug logging",
			},
			&cli.BoolFlag{
				Name:  "warning-symptoms",
				Value: false,
				Usage: "Whether to show warning symptoms (otherwise only critical are shown)",
			},
		},
		Action: func(cCtx *cli.Context) error {
			doctor.DoCheckUp(cCtx)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
