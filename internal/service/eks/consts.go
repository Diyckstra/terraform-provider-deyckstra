package eks

import "time"

const (
	ErrCodeClusterNotFound   = "InvalidKubernetesCluster.NotFound"
	ErrCodeNodegroupNotFound = "InvalidKubernetesNodegroup.NotFound"
)

const (
	IdentityProviderConfigTypeOIDC = "oidc"
)

const (
	ResourcesSecrets = "secrets"
)

func Resources_Values() []string {
	return []string{
		ResourcesSecrets,
	}
}

const (
	propagationTimeout = 2 * time.Minute
)
