package k8s

import (
	"do-k8s-challenge-kubegres/k8s/kubegres"

	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Cluster struct {
	ctx      *pulumi.Context
	provider *kubernetes.Provider
}

// Init creates a custom pulumi provider for k8s.
func Init(ctx *pulumi.Context, kubeconfig pulumi.StringOutput) (*Cluster, error) {
	provider, err := kubernetes.NewProvider(ctx, "main", &kubernetes.ProviderArgs{
		Kubeconfig: kubeconfig,
	})
	if err != nil {
		return &Cluster{}, err
	}

	return &Cluster{
		ctx:      ctx,
		provider: provider,
	}, nil
}

func (c *Cluster) DeployKubegres(manifestsPath string) *kubegres.Deployed {
	return kubegres.Deploy(c.ctx, c.provider, manifestsPath)
}
