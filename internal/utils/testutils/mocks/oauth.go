package mocks

import (
	"fmt"
	configv1 "github.com/openshift/api/config/v1"
	userv1 "github.com/openshift/api/user/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PrepareHTPasswdProvider returns a mock HTPasswordIdentityProvider object.
func PrepareHTPasswdProvider(secretName string) configv1.HTPasswdIdentityProvider {
	return configv1.HTPasswdIdentityProvider{
		FileData: configv1.SecretNameReference{
			Name: secretName,
		},
	}
}

// PrepareUser returns a mock User object.
func PrepareUser(username string) userv1.User {
	return userv1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: username,
		},
	}
}

// PrepareIdentity returns a mock Identity object.
func PrepareIdentity(providerName, username string) userv1.Identity {
	name := fmt.Sprintf("%s:%s", providerName, username)

	return userv1.Identity{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

// PrepareHTPasswdIdentityProvider returns a mock HTPasswordIde
func PrepareHTPasswdIdentityProvider(providerName string, providerType configv1.IdentityProviderType, htpasswdProvider configv1.HTPasswdIdentityProvider) configv1.IdentityProvider {
	return configv1.IdentityProvider{
		Name:          providerName,
		MappingMethod: configv1.MappingMethodClaim,
		IdentityProviderConfig: configv1.IdentityProviderConfig{
			Type:     providerType,
			HTPasswd: &htpasswdProvider,
		},
	}
}

func PrepareClusterRoleBinding(username, clusterRoleName string) rbacv1.ClusterRoleBinding {
	return rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", clusterRoleName, username),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: rbacv1.UserKind,
				Name: username,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Name: clusterRoleName,
			Kind: clusterRoleKey,
		},
	}
}
