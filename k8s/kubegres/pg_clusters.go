package kubegres

import (
	kubegresv1 "do-k8s-challenge-kubegres/k8s/kubegres/crd/generated/kubegres/v1"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// WithDefaultPGCluster based on https://www.kubegres.io/doc/getting-started.html
func (k *Deployed) WithDefaultPGCluster() error {
	superUserPassword := "postgresSuperUserPsw"
	replicationUserPassword := "postgresReplicaPsw" //nolint: gosec

	if k.error != nil {
		return k.error
	}

	creds, err := corev1.NewSecret(k.ctx, "default-creds", &corev1.SecretArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Secret"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("pg-default-cluster"),
		},
		Type: pulumi.String("Opaque"),
		StringData: pulumi.StringMap{
			"superUserPassword":       pulumi.String(superUserPassword),
			"replicationUserPassword": pulumi.String(replicationUserPassword),
		},
	}, pulumi.Provider(k.provider))
	if err != nil {
		return err
	}

	backupPvc, err := corev1.NewPersistentVolumeClaim(k.ctx, "pg-default-backup", &corev1.PersistentVolumeClaimArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("PersistentVolumeClaim"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("pg-backup"),
			Annotations: pulumi.StringMap{
				"pulumi.com/skipAwait": pulumi.String("true"),
			},
		},
		Spec: &corev1.PersistentVolumeClaimSpecArgs{
			AccessModes: pulumi.StringArray{
				pulumi.String("ReadWriteOnce"),
			},
			Resources: &corev1.ResourceRequirementsArgs{
				Requests: pulumi.StringMap{
					"storage": pulumi.String("1Gi"),
				},
			},
		},
	}, pulumi.Provider(k.provider))
	if err != nil {
		return err
	}

	_, err = kubegresv1.NewKubegres(k.ctx, "default", &kubegresv1.KubegresArgs{
		ApiVersion: pulumi.String("kubegres.reactive-tech.io/v1"),
		Kind:       pulumi.String("Kubegres"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("pg-default-cluster"),
		},
		Spec: &kubegresv1.KubegresSpecArgs{
			Replicas: pulumi.Int(3),
			Image:    pulumi.String("postgres:14.1"),
			Database: &kubegresv1.KubegresSpecDatabaseArgs{
				Size: pulumi.String("1Gi"),
			},
			Env: &kubegresv1.KubegresSpecEnvArray{
				&kubegresv1.KubegresSpecEnvArgs{
					Name: pulumi.String("POSTGRES_PASSWORD"),
					ValueFrom: &kubegresv1.KubegresSpecEnvValueFromArgs{
						SecretKeyRef: &kubegresv1.KubegresSpecEnvValueFromSecretKeyRefArgs{
							Name: creds.Metadata.Name(),
							Key:  pulumi.String("superUserPassword"),
						},
					},
				},
				&kubegresv1.KubegresSpecEnvArgs{
					Name: pulumi.String("POSTGRES_REPLICATION_PASSWORD"),
					ValueFrom: &kubegresv1.KubegresSpecEnvValueFromArgs{
						SecretKeyRef: &kubegresv1.KubegresSpecEnvValueFromSecretKeyRefArgs{
							Name: creds.Metadata.Name(),
							Key:  pulumi.String("replicationUserPassword"),
						},
					},
				},
			},
			SecurityContext: &kubegresv1.KubegresSpecSecurityContextArgs{
				// Run as postgres user
				RunAsUser:    pulumi.Int(999),
				FsGroup:      pulumi.Int(999),
				RunAsNonRoot: pulumi.Bool(true),
			},
			Backup: &kubegresv1.KubegresSpecBackupArgs{
				Schedule:    pulumi.String("0 0 * * *"),
				PvcName:     backupPvc.Metadata.Name(),
				VolumeMount: pulumi.String("/var/lib/pg-backup"),
			},
		},
	}, pulumi.Provider(k.provider))
	if err != nil {
		return err
	}

	k.ctx.Export("pgpassword", creds.StringData.MapIndex(pulumi.String("superUserPassword")))
	k.ctx.Export("backupPVC", backupPvc.ID())

	return nil
}
