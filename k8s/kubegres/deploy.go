package kubegres

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Deployed struct {
	ctx      *pulumi.Context
	provider *kubernetes.Provider
	error    error
}

// Deploy deploys kubegres controller with HA settings.
func Deploy(ctx *pulumi.Context, provider *kubernetes.Provider, manifestsPath string) *Deployed {
	_, err := yaml.NewConfigFile(ctx, "kubegres", &yaml.ConfigFileArgs{
		File: manifestsPath,
		Transformations: []yaml.Transformation{
			// Increase replicas to 2 and set podAntiAffinity to deployment
			func(state map[string]interface{}, opts ...pulumi.ResourceOption) {
				if state["kind"] == "Deployment" {
					spec := state["spec"].(map[string]interface{})
					spec["replicas"] = 2
					template := spec["template"].(map[string]interface{})
					tSpec := template["spec"].(map[string]interface{})
					tSpec["affinity"] = map[string]interface{}{
						"podAntiAffinity": map[string]interface{}{
							"preferredDuringSchedulingIgnoredDuringExecution": []map[string]interface{}{
								{
									"podAffinityTerm": map[string]interface{}{
										"labelSelector": map[string]interface{}{
											"matchExpressions": []map[string]interface{}{
												{
													"key":      "control-plane",
													"operator": "In",
													"values": []string{
														"controller-manager",
													},
												},
											},
										},
										"topologyKey": "kubernetes.io/hostname",
									},
									"weight": 100,
								},
							},
						},
					}
				}
			},
		},
	}, pulumi.Provider(provider))

	return &Deployed{
		ctx:      ctx,
		error:    err,
		provider: provider,
	}
}
