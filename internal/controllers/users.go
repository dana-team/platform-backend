package controllers

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/dana-team/platform-backend/internal/utils"
	"github.com/dana-team/platform-backend/internal/utils/pagination"
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

const (
	ErrCouldNotListUsers         = "Could not list users"
	ErrCouldNotGetRoleBinding    = "Could not get rolebinding %q"
	ErrCouldNotCreateRolebinding = "Could not create rolebinding %q"
	ErrCouldNotUpdateRolebinding = "Could not update rolebinding %q"
	ErrCouldNotGetRolebindings   = "Could not get rolebindings"
	ErrCouldNotDeleteRolebinding = "Could not delete rolebinding %q"
)

type UserController interface {
	// GetUsers get users from specified namespace and returns them as users.
	GetUsers(namespace string, limit, page int) (types.UsersOutput, error)

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

// UserPaginator paginates through secrets in a specified namespace.
type UserPaginator struct {
	pagination.GenericPaginator
	client    kubernetes.Interface
	namespace string
}

func NewUserController(client kubernetes.Interface, context context.Context, logger *zap.Logger) UserController {
	return &userController{
		logger: logger,
		client: client,
		ctx:    context,
	}
}

func (u *userController) GetUsers(namespace string, limit, page int) (types.UsersOutput, error) {
	userOutputs := types.UsersOutput{}
	u.logger.Debug(fmt.Sprintf("Trying to get all rolebindings in %q namespace", namespace))

	userPaginator := &UserPaginator{
		GenericPaginator: pagination.CreatePaginator(u.ctx, u.logger),
		namespace:        namespace,
		client:           u.client,
	}

	roleBindings, err := pagination.FetchPage[rbacv1.RoleBinding](limit, page, userPaginator)
	if err != nil {
		u.logger.Error(fmt.Sprintf("%v with error: %v", ErrCouldNotListUsers, err))
		return types.UsersOutput{}, customerrors.NewAPIError(ErrCouldNotListUsers, err)
	}
	for _, roleBinding := range roleBindings {
		userOutputs.Users = append(userOutputs.Users, types.User{Name: roleBinding.Name, Role: convertToPlatformRole(roleBinding.RoleRef.Name)})
	}
	userOutputs.Count = len(roleBindings)

	return userOutputs, nil
}

func (u *userController) GetUser(userIdentifier types.UserIdentifier) (types.User, error) {
	userOutput := types.User{}
	u.logger.Debug(fmt.Sprintf("Trying to fetch rolebinding %q in %q namespace", userIdentifier.UserName, userIdentifier.NamespaceName))

	roleBinding, err := u.client.RbacV1().RoleBindings(userIdentifier.NamespaceName).Get(u.ctx, userIdentifier.UserName, metav1.GetOptions{})
	if err != nil {
		u.logger.Error(fmt.Sprintf("%v with error: %v", fmt.Sprintf(ErrCouldNotGetRoleBinding, userIdentifier.UserName), err.Error()))
		return userOutput, customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotGetRoleBinding, userIdentifier.UserName), err)
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
		u.logger.Error(fmt.Sprintf("%v with error: %v", fmt.Sprintf(ErrCouldNotCreateRolebinding, user.Name), err.Error()))
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
		u.logger.Error(fmt.Sprintf("%v with error: %v", fmt.Sprintf(ErrCouldNotUpdateRolebinding, user.Name), err.Error()))
		return userOutput, customerrors.NewAPIError(fmt.Sprintf(ErrCouldNotUpdateRolebinding, user.Name), err)
	}

	u.logger.Debug(fmt.Sprintf("updated roleBinding %q successfully", user.Name))
	userOutput.Name = user.Name
	userOutput.Role = user.Role

	return userOutput, nil
}

func (u *userController) DeleteUser(userIdentifier types.UserIdentifier) (types.DeleteUserResponse, error) {
	u.logger.Debug(fmt.Sprintf("Trying to delete rolebinding %q in namespace %q", userIdentifier.UserName, userIdentifier.NamespaceName))

	if err := u.client.RbacV1().RoleBindings(userIdentifier.NamespaceName).Delete(u.ctx, userIdentifier.UserName, metav1.DeleteOptions{}); err != nil {
		message := fmt.Sprintf(ErrCouldNotDeleteRolebinding, userIdentifier.UserName)
		u.logger.Error(fmt.Sprintf("%v with error: %v", message, err.Error()))
		return types.DeleteUserResponse{Message: fmt.Sprintf("%v with error: %v", message, err.Error())}, customerrors.NewAPIError(message, err)
	}

	u.logger.Debug(fmt.Sprintf("Deleted roleBinding %q in namespace %q successfully", userIdentifier.UserName, userIdentifier.NamespaceName))
	return types.DeleteUserResponse{Message: fmt.Sprintf("Deleted roleBinding %q in namespace %q successfully", userIdentifier.UserName, userIdentifier.NamespaceName)}, nil
}

// FetchList retrieves a list of secrets from the specified namespace with given options.
func (p *UserPaginator) FetchList(listOptions metav1.ListOptions) (*types.List[rbacv1.RoleBinding], error) {
	roleBindings, err := p.client.RbacV1().RoleBindings(p.namespace).List(p.Ctx, listOptions)
	if err != nil {
		p.Logger.Error(fmt.Sprintf("%v with error: %v", ErrCouldNotGetRolebindings, err.Error()))
		return nil, customerrors.NewAPIError(ErrCouldNotGetRolebindings, err)
	}

	p.Logger.Debug("Fetched all rolebindings successfully")
	return (*types.List[rbacv1.RoleBinding])(roleBindings), nil
}

func prepareRoleBinding(roleBindingName string, role string) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   roleBindingName,
			Labels: utils.AddManagedLabel(map[string]string{}),
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
