// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package codespace

import (
	"context"
	"sync"

	"github.com/cli/cli/v2/internal/codespaces/api"
)

// apiClientMock is a mock implementation of apiClient.
//
// 	func TestSomethingThatUsesapiClient(t *testing.T) {
//
// 		// make and configure a mocked apiClient
// 		mockedapiClient := &apiClientMock{
// 			AuthorizedKeysFunc: func(ctx context.Context, user string) ([]byte, error) {
// 				panic("mock out the AuthorizedKeys method")
// 			},
// 			CreateCodespaceFunc: func(ctx context.Context, params *api.CreateCodespaceParams) (*api.Codespace, error) {
// 				panic("mock out the CreateCodespace method")
// 			},
// 			DeleteCodespaceFunc: func(ctx context.Context, name string) error {
// 				panic("mock out the DeleteCodespace method")
// 			},
// 			GetCodespaceFunc: func(ctx context.Context, token string, user string, name string) (*api.Codespace, error) {
// 				panic("mock out the GetCodespace method")
// 			},
// 			GetCodespaceRegionLocationFunc: func(ctx context.Context) (string, error) {
// 				panic("mock out the GetCodespaceRegionLocation method")
// 			},
// 			GetCodespaceRepositoryContentsFunc: func(ctx context.Context, codespace *api.Codespace, path string) ([]byte, error) {
// 				panic("mock out the GetCodespaceRepositoryContents method")
// 			},
// 			GetCodespaceTokenFunc: func(ctx context.Context, user string, name string) (string, error) {
// 				panic("mock out the GetCodespaceToken method")
// 			},
// 			GetCodespacesSKUsFunc: func(ctx context.Context, user *api.User, repository *api.Repository, branch string, location string) ([]*api.SKU, error) {
// 				panic("mock out the GetCodespacesSKUs method")
// 			},
// 			GetRepositoryFunc: func(ctx context.Context, nwo string) (*api.Repository, error) {
// 				panic("mock out the GetRepository method")
// 			},
// 			GetUserFunc: func(ctx context.Context) (*api.User, error) {
// 				panic("mock out the GetUser method")
// 			},
// 			ListCodespacesFunc: func(ctx context.Context) ([]*api.Codespace, error) {
// 				panic("mock out the ListCodespaces method")
// 			},
// 			StartCodespaceFunc: func(ctx context.Context, name string) error {
// 				panic("mock out the StartCodespace method")
// 			},
// 		}
//
// 		// use mockedapiClient in code that requires apiClient
// 		// and then make assertions.
//
// 	}
type apiClientMock struct {
	// AuthorizedKeysFunc mocks the AuthorizedKeys method.
	AuthorizedKeysFunc func(ctx context.Context, user string) ([]byte, error)

	// CreateCodespaceFunc mocks the CreateCodespace method.
	CreateCodespaceFunc func(ctx context.Context, params *api.CreateCodespaceParams) (*api.Codespace, error)

	// DeleteCodespaceFunc mocks the DeleteCodespace method.
	DeleteCodespaceFunc func(ctx context.Context, name string) error

	// GetCodespaceFunc mocks the GetCodespace method.
	GetCodespaceFunc func(ctx context.Context, token string, user string, name string) (*api.Codespace, error)

	// GetCodespaceRegionLocationFunc mocks the GetCodespaceRegionLocation method.
	GetCodespaceRegionLocationFunc func(ctx context.Context) (string, error)

	// GetCodespaceRepositoryContentsFunc mocks the GetCodespaceRepositoryContents method.
	GetCodespaceRepositoryContentsFunc func(ctx context.Context, codespace *api.Codespace, path string) ([]byte, error)

	// GetCodespaceTokenFunc mocks the GetCodespaceToken method.
	GetCodespaceTokenFunc func(ctx context.Context, user string, name string) (string, error)

	// GetCodespacesSKUsFunc mocks the GetCodespacesSKUs method.
	GetCodespacesSKUsFunc func(ctx context.Context, user *api.User, repository *api.Repository, branch string, location string) ([]*api.SKU, error)

	// GetRepositoryFunc mocks the GetRepository method.
	GetRepositoryFunc func(ctx context.Context, nwo string) (*api.Repository, error)

	// GetUserFunc mocks the GetUser method.
	GetUserFunc func(ctx context.Context) (*api.User, error)

	// ListCodespacesFunc mocks the ListCodespaces method.
	ListCodespacesFunc func(ctx context.Context) ([]*api.Codespace, error)

	// StartCodespaceFunc mocks the StartCodespace method.
	StartCodespaceFunc func(ctx context.Context, name string) error

	// calls tracks calls to the methods.
	calls struct {
		// AuthorizedKeys holds details about calls to the AuthorizedKeys method.
		AuthorizedKeys []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// User is the user argument value.
			User string
		}
		// CreateCodespace holds details about calls to the CreateCodespace method.
		CreateCodespace []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Params is the params argument value.
			Params *api.CreateCodespaceParams
		}
		// DeleteCodespace holds details about calls to the DeleteCodespace method.
		DeleteCodespace []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
		}
		// GetCodespace holds details about calls to the GetCodespace method.
		GetCodespace []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Token is the token argument value.
			Token string
			// User is the user argument value.
			User string
			// Name is the name argument value.
			Name string
		}
		// GetCodespaceRegionLocation holds details about calls to the GetCodespaceRegionLocation method.
		GetCodespaceRegionLocation []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// GetCodespaceRepositoryContents holds details about calls to the GetCodespaceRepositoryContents method.
		GetCodespaceRepositoryContents []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Codespace is the codespace argument value.
			Codespace *api.Codespace
			// Path is the path argument value.
			Path string
		}
		// GetCodespaceToken holds details about calls to the GetCodespaceToken method.
		GetCodespaceToken []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// User is the user argument value.
			User string
			// Name is the name argument value.
			Name string
		}
		// GetCodespacesSKUs holds details about calls to the GetCodespacesSKUs method.
		GetCodespacesSKUs []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// User is the user argument value.
			User *api.User
			// Repository is the repository argument value.
			Repository *api.Repository
			// Branch is the branch argument value.
			Branch string
			// Location is the location argument value.
			Location string
		}
		// GetRepository holds details about calls to the GetRepository method.
		GetRepository []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Nwo is the nwo argument value.
			Nwo string
		}
		// GetUser holds details about calls to the GetUser method.
		GetUser []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// ListCodespaces holds details about calls to the ListCodespaces method.
		ListCodespaces []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// StartCodespace holds details about calls to the StartCodespace method.
		StartCodespace []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
		}
	}
	lockAuthorizedKeys                 sync.RWMutex
	lockCreateCodespace                sync.RWMutex
	lockDeleteCodespace                sync.RWMutex
	lockGetCodespace                   sync.RWMutex
	lockGetCodespaceRegionLocation     sync.RWMutex
	lockGetCodespaceRepositoryContents sync.RWMutex
	lockGetCodespaceToken              sync.RWMutex
	lockGetCodespacesSKUs              sync.RWMutex
	lockGetRepository                  sync.RWMutex
	lockGetUser                        sync.RWMutex
	lockListCodespaces                 sync.RWMutex
	lockStartCodespace                 sync.RWMutex
}

// AuthorizedKeys calls AuthorizedKeysFunc.
func (mock *apiClientMock) AuthorizedKeys(ctx context.Context, user string) ([]byte, error) {
	if mock.AuthorizedKeysFunc == nil {
		panic("apiClientMock.AuthorizedKeysFunc: method is nil but apiClient.AuthorizedKeys was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		User string
	}{
		Ctx:  ctx,
		User: user,
	}
	mock.lockAuthorizedKeys.Lock()
	mock.calls.AuthorizedKeys = append(mock.calls.AuthorizedKeys, callInfo)
	mock.lockAuthorizedKeys.Unlock()
	return mock.AuthorizedKeysFunc(ctx, user)
}

// AuthorizedKeysCalls gets all the calls that were made to AuthorizedKeys.
// Check the length with:
//     len(mockedapiClient.AuthorizedKeysCalls())
func (mock *apiClientMock) AuthorizedKeysCalls() []struct {
	Ctx  context.Context
	User string
} {
	var calls []struct {
		Ctx  context.Context
		User string
	}
	mock.lockAuthorizedKeys.RLock()
	calls = mock.calls.AuthorizedKeys
	mock.lockAuthorizedKeys.RUnlock()
	return calls
}

// CreateCodespace calls CreateCodespaceFunc.
func (mock *apiClientMock) CreateCodespace(ctx context.Context, params *api.CreateCodespaceParams) (*api.Codespace, error) {
	if mock.CreateCodespaceFunc == nil {
		panic("apiClientMock.CreateCodespaceFunc: method is nil but apiClient.CreateCodespace was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		Params *api.CreateCodespaceParams
	}{
		Ctx:    ctx,
		Params: params,
	}
	mock.lockCreateCodespace.Lock()
	mock.calls.CreateCodespace = append(mock.calls.CreateCodespace, callInfo)
	mock.lockCreateCodespace.Unlock()
	return mock.CreateCodespaceFunc(ctx, params)
}

// CreateCodespaceCalls gets all the calls that were made to CreateCodespace.
// Check the length with:
//     len(mockedapiClient.CreateCodespaceCalls())
func (mock *apiClientMock) CreateCodespaceCalls() []struct {
	Ctx    context.Context
	Params *api.CreateCodespaceParams
} {
	var calls []struct {
		Ctx    context.Context
		Params *api.CreateCodespaceParams
	}
	mock.lockCreateCodespace.RLock()
	calls = mock.calls.CreateCodespace
	mock.lockCreateCodespace.RUnlock()
	return calls
}

// DeleteCodespace calls DeleteCodespaceFunc.
func (mock *apiClientMock) DeleteCodespace(ctx context.Context, name string) error {
	if mock.DeleteCodespaceFunc == nil {
		panic("apiClientMock.DeleteCodespaceFunc: method is nil but apiClient.DeleteCodespace was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Name string
	}{
		Ctx:  ctx,
		Name: name,
	}
	mock.lockDeleteCodespace.Lock()
	mock.calls.DeleteCodespace = append(mock.calls.DeleteCodespace, callInfo)
	mock.lockDeleteCodespace.Unlock()
	return mock.DeleteCodespaceFunc(ctx, name)
}

// DeleteCodespaceCalls gets all the calls that were made to DeleteCodespace.
// Check the length with:
//     len(mockedapiClient.DeleteCodespaceCalls())
func (mock *apiClientMock) DeleteCodespaceCalls() []struct {
	Ctx  context.Context
	Name string
} {
	var calls []struct {
		Ctx  context.Context
		Name string
	}
	mock.lockDeleteCodespace.RLock()
	calls = mock.calls.DeleteCodespace
	mock.lockDeleteCodespace.RUnlock()
	return calls
}

// GetCodespace calls GetCodespaceFunc.
func (mock *apiClientMock) GetCodespace(ctx context.Context, token string, user string, name string) (*api.Codespace, error) {
	if mock.GetCodespaceFunc == nil {
		panic("apiClientMock.GetCodespaceFunc: method is nil but apiClient.GetCodespace was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Token string
		User  string
		Name  string
	}{
		Ctx:   ctx,
		Token: token,
		User:  user,
		Name:  name,
	}
	mock.lockGetCodespace.Lock()
	mock.calls.GetCodespace = append(mock.calls.GetCodespace, callInfo)
	mock.lockGetCodespace.Unlock()
	return mock.GetCodespaceFunc(ctx, token, user, name)
}

// GetCodespaceCalls gets all the calls that were made to GetCodespace.
// Check the length with:
//     len(mockedapiClient.GetCodespaceCalls())
func (mock *apiClientMock) GetCodespaceCalls() []struct {
	Ctx   context.Context
	Token string
	User  string
	Name  string
} {
	var calls []struct {
		Ctx   context.Context
		Token string
		User  string
		Name  string
	}
	mock.lockGetCodespace.RLock()
	calls = mock.calls.GetCodespace
	mock.lockGetCodespace.RUnlock()
	return calls
}

// GetCodespaceRegionLocation calls GetCodespaceRegionLocationFunc.
func (mock *apiClientMock) GetCodespaceRegionLocation(ctx context.Context) (string, error) {
	if mock.GetCodespaceRegionLocationFunc == nil {
		panic("apiClientMock.GetCodespaceRegionLocationFunc: method is nil but apiClient.GetCodespaceRegionLocation was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockGetCodespaceRegionLocation.Lock()
	mock.calls.GetCodespaceRegionLocation = append(mock.calls.GetCodespaceRegionLocation, callInfo)
	mock.lockGetCodespaceRegionLocation.Unlock()
	return mock.GetCodespaceRegionLocationFunc(ctx)
}

// GetCodespaceRegionLocationCalls gets all the calls that were made to GetCodespaceRegionLocation.
// Check the length with:
//     len(mockedapiClient.GetCodespaceRegionLocationCalls())
func (mock *apiClientMock) GetCodespaceRegionLocationCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockGetCodespaceRegionLocation.RLock()
	calls = mock.calls.GetCodespaceRegionLocation
	mock.lockGetCodespaceRegionLocation.RUnlock()
	return calls
}

// GetCodespaceRepositoryContents calls GetCodespaceRepositoryContentsFunc.
func (mock *apiClientMock) GetCodespaceRepositoryContents(ctx context.Context, codespace *api.Codespace, path string) ([]byte, error) {
	if mock.GetCodespaceRepositoryContentsFunc == nil {
		panic("apiClientMock.GetCodespaceRepositoryContentsFunc: method is nil but apiClient.GetCodespaceRepositoryContents was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		Codespace *api.Codespace
		Path      string
	}{
		Ctx:       ctx,
		Codespace: codespace,
		Path:      path,
	}
	mock.lockGetCodespaceRepositoryContents.Lock()
	mock.calls.GetCodespaceRepositoryContents = append(mock.calls.GetCodespaceRepositoryContents, callInfo)
	mock.lockGetCodespaceRepositoryContents.Unlock()
	return mock.GetCodespaceRepositoryContentsFunc(ctx, codespace, path)
}

// GetCodespaceRepositoryContentsCalls gets all the calls that were made to GetCodespaceRepositoryContents.
// Check the length with:
//     len(mockedapiClient.GetCodespaceRepositoryContentsCalls())
func (mock *apiClientMock) GetCodespaceRepositoryContentsCalls() []struct {
	Ctx       context.Context
	Codespace *api.Codespace
	Path      string
} {
	var calls []struct {
		Ctx       context.Context
		Codespace *api.Codespace
		Path      string
	}
	mock.lockGetCodespaceRepositoryContents.RLock()
	calls = mock.calls.GetCodespaceRepositoryContents
	mock.lockGetCodespaceRepositoryContents.RUnlock()
	return calls
}

// GetCodespaceToken calls GetCodespaceTokenFunc.
func (mock *apiClientMock) GetCodespaceToken(ctx context.Context, user string, name string) (string, error) {
	if mock.GetCodespaceTokenFunc == nil {
		panic("apiClientMock.GetCodespaceTokenFunc: method is nil but apiClient.GetCodespaceToken was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		User string
		Name string
	}{
		Ctx:  ctx,
		User: user,
		Name: name,
	}
	mock.lockGetCodespaceToken.Lock()
	mock.calls.GetCodespaceToken = append(mock.calls.GetCodespaceToken, callInfo)
	mock.lockGetCodespaceToken.Unlock()
	return mock.GetCodespaceTokenFunc(ctx, user, name)
}

// GetCodespaceTokenCalls gets all the calls that were made to GetCodespaceToken.
// Check the length with:
//     len(mockedapiClient.GetCodespaceTokenCalls())
func (mock *apiClientMock) GetCodespaceTokenCalls() []struct {
	Ctx  context.Context
	User string
	Name string
} {
	var calls []struct {
		Ctx  context.Context
		User string
		Name string
	}
	mock.lockGetCodespaceToken.RLock()
	calls = mock.calls.GetCodespaceToken
	mock.lockGetCodespaceToken.RUnlock()
	return calls
}

// GetCodespacesSKUs calls GetCodespacesSKUsFunc.
func (mock *apiClientMock) GetCodespacesSKUs(ctx context.Context, user *api.User, repository *api.Repository, branch string, location string) ([]*api.SKU, error) {
	if mock.GetCodespacesSKUsFunc == nil {
		panic("apiClientMock.GetCodespacesSKUsFunc: method is nil but apiClient.GetCodespacesSKUs was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		User       *api.User
		Repository *api.Repository
		Branch     string
		Location   string
	}{
		Ctx:        ctx,
		User:       user,
		Repository: repository,
		Branch:     branch,
		Location:   location,
	}
	mock.lockGetCodespacesSKUs.Lock()
	mock.calls.GetCodespacesSKUs = append(mock.calls.GetCodespacesSKUs, callInfo)
	mock.lockGetCodespacesSKUs.Unlock()
	return mock.GetCodespacesSKUsFunc(ctx, user, repository, branch, location)
}

// GetCodespacesSKUsCalls gets all the calls that were made to GetCodespacesSKUs.
// Check the length with:
//     len(mockedapiClient.GetCodespacesSKUsCalls())
func (mock *apiClientMock) GetCodespacesSKUsCalls() []struct {
	Ctx        context.Context
	User       *api.User
	Repository *api.Repository
	Branch     string
	Location   string
} {
	var calls []struct {
		Ctx        context.Context
		User       *api.User
		Repository *api.Repository
		Branch     string
		Location   string
	}
	mock.lockGetCodespacesSKUs.RLock()
	calls = mock.calls.GetCodespacesSKUs
	mock.lockGetCodespacesSKUs.RUnlock()
	return calls
}

// GetRepository calls GetRepositoryFunc.
func (mock *apiClientMock) GetRepository(ctx context.Context, nwo string) (*api.Repository, error) {
	if mock.GetRepositoryFunc == nil {
		panic("apiClientMock.GetRepositoryFunc: method is nil but apiClient.GetRepository was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Nwo string
	}{
		Ctx: ctx,
		Nwo: nwo,
	}
	mock.lockGetRepository.Lock()
	mock.calls.GetRepository = append(mock.calls.GetRepository, callInfo)
	mock.lockGetRepository.Unlock()
	return mock.GetRepositoryFunc(ctx, nwo)
}

// GetRepositoryCalls gets all the calls that were made to GetRepository.
// Check the length with:
//     len(mockedapiClient.GetRepositoryCalls())
func (mock *apiClientMock) GetRepositoryCalls() []struct {
	Ctx context.Context
	Nwo string
} {
	var calls []struct {
		Ctx context.Context
		Nwo string
	}
	mock.lockGetRepository.RLock()
	calls = mock.calls.GetRepository
	mock.lockGetRepository.RUnlock()
	return calls
}

// GetUser calls GetUserFunc.
func (mock *apiClientMock) GetUser(ctx context.Context) (*api.User, error) {
	if mock.GetUserFunc == nil {
		panic("apiClientMock.GetUserFunc: method is nil but apiClient.GetUser was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockGetUser.Lock()
	mock.calls.GetUser = append(mock.calls.GetUser, callInfo)
	mock.lockGetUser.Unlock()
	return mock.GetUserFunc(ctx)
}

// GetUserCalls gets all the calls that were made to GetUser.
// Check the length with:
//     len(mockedapiClient.GetUserCalls())
func (mock *apiClientMock) GetUserCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockGetUser.RLock()
	calls = mock.calls.GetUser
	mock.lockGetUser.RUnlock()
	return calls
}

// ListCodespaces calls ListCodespacesFunc.
func (mock *apiClientMock) ListCodespaces(ctx context.Context) ([]*api.Codespace, error) {
	if mock.ListCodespacesFunc == nil {
		panic("apiClientMock.ListCodespacesFunc: method is nil but apiClient.ListCodespaces was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockListCodespaces.Lock()
	mock.calls.ListCodespaces = append(mock.calls.ListCodespaces, callInfo)
	mock.lockListCodespaces.Unlock()
	return mock.ListCodespacesFunc(ctx)
}

// ListCodespacesCalls gets all the calls that were made to ListCodespaces.
// Check the length with:
//     len(mockedapiClient.ListCodespacesCalls())
func (mock *apiClientMock) ListCodespacesCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockListCodespaces.RLock()
	calls = mock.calls.ListCodespaces
	mock.lockListCodespaces.RUnlock()
	return calls
}

// StartCodespace calls StartCodespaceFunc.
func (mock *apiClientMock) StartCodespace(ctx context.Context, name string) error {
	if mock.StartCodespaceFunc == nil {
		panic("apiClientMock.StartCodespaceFunc: method is nil but apiClient.StartCodespace was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Name string
	}{
		Ctx:  ctx,
		Name: name,
	}
	mock.lockStartCodespace.Lock()
	mock.calls.StartCodespace = append(mock.calls.StartCodespace, callInfo)
	mock.lockStartCodespace.Unlock()
	return mock.StartCodespaceFunc(ctx, name)
}

// StartCodespaceCalls gets all the calls that were made to StartCodespace.
// Check the length with:
//     len(mockedapiClient.StartCodespaceCalls())
func (mock *apiClientMock) StartCodespaceCalls() []struct {
	Ctx  context.Context
	Name string
} {
	var calls []struct {
		Ctx  context.Context
		Name string
	}
	mock.lockStartCodespace.RLock()
	calls = mock.calls.StartCodespace
	mock.lockStartCodespace.RUnlock()
	return calls
}
