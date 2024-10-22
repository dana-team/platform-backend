package mocks

import (
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/dana-team/platform-backend/internal/utils"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	clusterRoleKey = "ClusterRole"
	cappUserPrefix = "capp-user-"
)

// PrepareRoleBinding returns a mock RoleBinding object.
func PrepareRoleBinding(name, namespace, role string) rbacv1.RoleBinding {
	return rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    utils.AddManagedLabel(map[string]string{}),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:     rbacv1.UserKind,
				Name:     name,
				APIGroup: rbacv1.GroupName,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind: clusterRoleKey,
			Name: cappUserPrefix + role,
		},
	}
}

// PrepareUserType returns a mock User type object.
func PrepareUserType(name, role string) types.User {
	return types.User{
		Name: name,
		Role: role,
	}
}

// PrepareUpdateUserDataType returns a mock User type object.
func PrepareUpdateUserDataType(role string) types.UpdateUserData {
	return types.UpdateUserData{
		Role: role,
	}
}
