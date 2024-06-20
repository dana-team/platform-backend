package controllers

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/src/types"
	"go.uber.org/zap"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

const (
	AdminClusterRole        = "capp-user-admin"
	ContributorClusterRole  = "capp-user-contributor"
	ViewerClusterRole       = "capp-user-viewer"
	AdminPlatformRole       = "admin"
	ContributorPlatformRole = "contributor"
	ViewerPlatformRole      = "viewer"
)

type UserController interface {
	// GetUsers get users from specified namespace and returns them as users.
	GetUsers(namespace string) (types.UsersOutput, error)

	// GetUser gets a specific user from specified namespace and returns it as user.
	GetUser(userIdentifier types.UserIdentifier) (types.User, error)

	// AddUser creates a new roleBinding in the specified namespace.
	AddUser(user types.UserInput) (types.User, error)

	// DeleteUser deletes user in the specified namespace.
	DeleteUser(userIdentifier types.UserIdentifier) (types.DeleteUserResponse, error)

	// UpdateUser updates user in the specified namespace.
	UpdateUser(user types.UserInput) (types.User, error)
}

type userController struct {
	client kubernetes.Interface
	ctx    context.Context
	logger *zap.Logger
}

func NewUserController(client kubernetes.Interface, context context.Context, logger *zap.Logger) UserController {
	return &userController{
		logger: logger,
		client: client,
		ctx:    context,
	}
}

func (u *userController) GetUsers(namespace string) (types.UsersOutput, error) {
	u.logger.Debug(fmt.Sprintf("Trying to get all rolebindings in %q namespace", namespace))

	roleBindings, err := u.client.RbacV1().RoleBindings(namespace).List(u.ctx, metav1.ListOptions{})
	userOutputs := types.UsersOutput{}
	if err != nil {
		u.logger.Error(fmt.Sprintf("Could not get rolebindings with error: %v", err.Error()))
		return userOutputs, err
	}
	u.logger.Debug("Fetched all rolebindings successfully")

	for _, roleBinding := range roleBindings.Items {
		userOutputs.Users = append(userOutputs.Users, types.User{Name: roleBinding.Name, Role: convertToPlatformRole(roleBinding.RoleRef.Name)})
	}
	userOutputs.Count = len(roleBindings.Items)

	return userOutputs, nil
}

func (u *userController) GetUser(userIdentifier types.UserIdentifier) (types.User, error) {
	userOutput := types.User{}
	u.logger.Debug(fmt.Sprintf("Trying to fetch rolebinding %q in %q namespace", userIdentifier.UserName, userIdentifier.NamespaceName))

	roleBinding, err := u.client.RbacV1().RoleBindings(userIdentifier.NamespaceName).Get(u.ctx, userIdentifier.UserName, metav1.GetOptions{})
	if err != nil {
		u.logger.Error(fmt.Sprintf("Could fetch rolebinding %q with error: %v", userIdentifier.UserName, err.Error()))
		return userOutput, err
	}
	u.logger.Debug(fmt.Sprintf("fetched roleBinding %q successfully", roleBinding.Name))

	userOutput.Name = roleBinding.Name
	userOutput.Role = convertToPlatformRole(roleBinding.RoleRef.Name)

	return userOutput, nil
}

func (u *userController) AddUser(user types.UserInput) (types.User, error) {
	userOutput := types.User{}
	u.logger.Debug(fmt.Sprintf("Trying to create rolebinding %q in %q namespace", user.Name, user.Namespace))

	roleBinding, err := u.client.RbacV1().RoleBindings(user.Namespace).Create(u.ctx,
		prepareRoleBinding(user.Name, user.Role), metav1.CreateOptions{})
	if err != nil {
		u.logger.Error(fmt.Sprintf("Could create rolebinding %q with error: %v", user.Name, err.Error()))
		return userOutput, err
	}
	u.logger.Debug(fmt.Sprintf("created roleBinding %q successfully", roleBinding.Name))
	userOutput.Name = roleBinding.Name
	userOutput.Role = user.Role

	return userOutput, nil
}

func (u *userController) UpdateUser(user types.UserInput) (types.User, error) {
	userOutput := types.User{}

	u.logger.Debug(fmt.Sprintf("Trying to update rolebinding %q in %q namespace", user.Name, user.Namespace))

	// K8s does not allow to update role ref of roleBinding. So we need to delete the old roleBinding and create the desired one
	_, err := u.DeleteUser(types.UserIdentifier{UserName: user.Name, NamespaceName: user.Namespace})
	if err != nil {
		return userOutput, err
	}

	// The retry is needed because K8s is eventually consistent the old roleBinding takes time to delete.
	err = retry.OnError(retry.DefaultRetry, func(err error) bool {
		return err != nil

	}, func() error {
		_, err := u.AddUser(user)
		return err
	})

	if err != nil {
		u.logger.Error(fmt.Sprintf("Could not update rolebinding %q with error: %v", user.Name, err.Error()))
		return userOutput, err
	}

	u.logger.Debug(fmt.Sprintf("updated roleBinding %q successfully", user.Name))
	userOutput.Name = user.Name
	userOutput.Role = user.Role

	return userOutput, nil
}

func (u *userController) DeleteUser(userIdentifier types.UserIdentifier) (types.DeleteUserResponse, error) {
	u.logger.Debug(fmt.Sprintf("Trying to delete rolebinding %q in %q namespace", userIdentifier.UserName, userIdentifier.NamespaceName))

	if err := u.client.RbacV1().RoleBindings(userIdentifier.NamespaceName).Delete(u.ctx, userIdentifier.UserName, metav1.DeleteOptions{}); err != nil {
		u.logger.Error(fmt.Sprintf("Could note delete rolebinding %q with error: %v",
			userIdentifier.UserName, err.Error()))
		return types.DeleteUserResponse{Message: fmt.Sprintf("Could note delete rolebinding %q with error: %s",
			userIdentifier.UserName, err.Error())}, err
	}

	u.logger.Debug(fmt.Sprintf("deleted roleBinding %q successfully", userIdentifier.UserName))
	return types.DeleteUserResponse{Message: fmt.Sprintf("deleted roleBinding %q successfully", userIdentifier.UserName)}, nil
}

func prepareRoleBinding(roleBindingName string, role string) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: roleBindingName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:     rbacv1.UserKind,
				Name:     roleBindingName,
				APIGroup: rbacv1.GroupName,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     convertToK8sRoles(role),
			APIGroup: rbacv1.GroupName,
		},
	}
}

// convertToK8sRoles converts user given roles to the role which exists in the cluster
func convertToK8sRoles(requestRole string) string {
	rolesConvertor := map[string]string{AdminPlatformRole: AdminClusterRole, ContributorPlatformRole: ContributorClusterRole, ViewerPlatformRole: ViewerClusterRole}
	return rolesConvertor[requestRole]
}

// convertToPlatformRoles converts from the cluster role to platform roles.
func convertToPlatformRole(k8sRole string) string {
	rolesConvertor := map[string]string{AdminClusterRole: AdminPlatformRole, ContributorClusterRole: ContributorPlatformRole, ViewerClusterRole: ViewerClusterRole}
	return rolesConvertor[k8sRole]
}
