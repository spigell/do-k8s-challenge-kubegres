package main

import (
	"do-k8s-challenge-kubegres/infra"
	"do-k8s-challenge-kubegres/k8s"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// This is a main pulumi program.
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		clusterInfo, err := infra.UpManagedK8s(ctx)
		if err != nil {
			return err
		}

		kubeconfig := clusterInfo["kubeconfig"].(pulumi.StringOutput)
		ctx.Export("kubeconfig", kubeconfig)

		cfg := config.New(ctx, "")
		manifests := cfg.Require("manifestPath")

		cluster, err := k8s.Init(ctx, kubeconfig)
		if err != nil {
			return err
		}

		return cluster.DeployKubegres(manifests).WithDefaultPGCluster()
	})
}
