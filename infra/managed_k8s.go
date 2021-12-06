package infra

import (
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var region = "lon1"

// UpManagedK8s creates a dedicated project and a VPC for three-node Kubernetes cluster.
func UpManagedK8s(ctx *pulumi.Context) (map[string]interface{}, error) {
	cluster := make(map[string]interface{})

	project, err := digitalocean.NewProject(ctx, "main", &digitalocean.ProjectArgs{
		Name:        pulumi.String("k8s-challenge"),
		Description: pulumi.String("A project for https://www.digitalocean.com/community/pages/kubernetes-challenge"),
		Environment: pulumi.String("Development"),
		Purpose:     pulumi.String("Learning"),
	})
	if err != nil {
		return nil, err
	}

	vpc, err := digitalocean.NewVpc(ctx, "main", &digitalocean.VpcArgs{
		Name:    pulumi.String("k8s-challenge"),
		IpRange: pulumi.String("10.40.10.0/24"),
		Region:  pulumi.String(region),
	})
	if err != nil {
		return nil, err
	}

	kube, err := digitalocean.NewKubernetesCluster(ctx, "main", &digitalocean.KubernetesClusterArgs{
		Region:  pulumi.String(region),
		Version: pulumi.String("1.21.5-do.0"),
		Ha:      pulumi.Bool(false),
		VpcUuid: vpc.ID(),
		NodePool: &digitalocean.KubernetesClusterNodePoolArgs{
			Name:      pulumi.String("k8s-challenge"),
			Size:      pulumi.String("s-2vcpu-2gb"),
			NodeCount: pulumi.Int(3),
			AutoScale: pulumi.Bool(false),
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = digitalocean.NewProjectResources(ctx, "main", &digitalocean.ProjectResourcesArgs{
		Project: project.ID(),
		Resources: pulumi.StringArray{
			kube.ClusterUrn,
		},
	})
	if err != nil {
		return nil, err
	}

	cluster["kubeconfig"] = kube.KubeConfigs.Index(pulumi.Int(0)).RawConfig().Elem()

	return cluster, nil
}
