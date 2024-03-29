// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	iam_service "github.com/hashicorp/hcp-sdk-go/clients/cloud-iam/stable/2019-12-10/client/iam_service"
	mock "github.com/stretchr/testify/mock"

	runtime "github.com/go-openapi/runtime"
)

// ClientService is an autogenerated mock type for the ClientService type
type ClientService struct {
	mock.Mock
}

// IamServiceBatchGetPrincipals provides a mock function with given fields: params, authInfo, opts
func (_m *ClientService) IamServiceBatchGetPrincipals(params *iam_service.IamServiceBatchGetPrincipalsParams, authInfo runtime.ClientAuthInfoWriter, opts ...iam_service.ClientOption) (*iam_service.IamServiceBatchGetPrincipalsOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params, authInfo)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *iam_service.IamServiceBatchGetPrincipalsOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceBatchGetPrincipalsParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) (*iam_service.IamServiceBatchGetPrincipalsOK, error)); ok {
		return rf(params, authInfo, opts...)
	}
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceBatchGetPrincipalsParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) *iam_service.IamServiceBatchGetPrincipalsOK); ok {
		r0 = rf(params, authInfo, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam_service.IamServiceBatchGetPrincipalsOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*iam_service.IamServiceBatchGetPrincipalsParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) error); ok {
		r1 = rf(params, authInfo, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IamServiceCreateUserPrincipal provides a mock function with given fields: params, authInfo, opts
func (_m *ClientService) IamServiceCreateUserPrincipal(params *iam_service.IamServiceCreateUserPrincipalParams, authInfo runtime.ClientAuthInfoWriter, opts ...iam_service.ClientOption) (*iam_service.IamServiceCreateUserPrincipalOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params, authInfo)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *iam_service.IamServiceCreateUserPrincipalOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceCreateUserPrincipalParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) (*iam_service.IamServiceCreateUserPrincipalOK, error)); ok {
		return rf(params, authInfo, opts...)
	}
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceCreateUserPrincipalParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) *iam_service.IamServiceCreateUserPrincipalOK); ok {
		r0 = rf(params, authInfo, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam_service.IamServiceCreateUserPrincipalOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*iam_service.IamServiceCreateUserPrincipalParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) error); ok {
		r1 = rf(params, authInfo, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IamServiceDeleteOrganizationMembership provides a mock function with given fields: params, authInfo, opts
func (_m *ClientService) IamServiceDeleteOrganizationMembership(params *iam_service.IamServiceDeleteOrganizationMembershipParams, authInfo runtime.ClientAuthInfoWriter, opts ...iam_service.ClientOption) (*iam_service.IamServiceDeleteOrganizationMembershipOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params, authInfo)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *iam_service.IamServiceDeleteOrganizationMembershipOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceDeleteOrganizationMembershipParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) (*iam_service.IamServiceDeleteOrganizationMembershipOK, error)); ok {
		return rf(params, authInfo, opts...)
	}
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceDeleteOrganizationMembershipParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) *iam_service.IamServiceDeleteOrganizationMembershipOK); ok {
		r0 = rf(params, authInfo, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam_service.IamServiceDeleteOrganizationMembershipOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*iam_service.IamServiceDeleteOrganizationMembershipParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) error); ok {
		r1 = rf(params, authInfo, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IamServiceGetCallerIdentity provides a mock function with given fields: params, authInfo, opts
func (_m *ClientService) IamServiceGetCallerIdentity(params *iam_service.IamServiceGetCallerIdentityParams, authInfo runtime.ClientAuthInfoWriter, opts ...iam_service.ClientOption) (*iam_service.IamServiceGetCallerIdentityOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params, authInfo)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *iam_service.IamServiceGetCallerIdentityOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceGetCallerIdentityParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) (*iam_service.IamServiceGetCallerIdentityOK, error)); ok {
		return rf(params, authInfo, opts...)
	}
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceGetCallerIdentityParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) *iam_service.IamServiceGetCallerIdentityOK); ok {
		r0 = rf(params, authInfo, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam_service.IamServiceGetCallerIdentityOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*iam_service.IamServiceGetCallerIdentityParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) error); ok {
		r1 = rf(params, authInfo, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IamServiceGetCurrentUserPrincipal provides a mock function with given fields: params, authInfo, opts
func (_m *ClientService) IamServiceGetCurrentUserPrincipal(params *iam_service.IamServiceGetCurrentUserPrincipalParams, authInfo runtime.ClientAuthInfoWriter, opts ...iam_service.ClientOption) (*iam_service.IamServiceGetCurrentUserPrincipalOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params, authInfo)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *iam_service.IamServiceGetCurrentUserPrincipalOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceGetCurrentUserPrincipalParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) (*iam_service.IamServiceGetCurrentUserPrincipalOK, error)); ok {
		return rf(params, authInfo, opts...)
	}
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceGetCurrentUserPrincipalParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) *iam_service.IamServiceGetCurrentUserPrincipalOK); ok {
		r0 = rf(params, authInfo, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam_service.IamServiceGetCurrentUserPrincipalOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*iam_service.IamServiceGetCurrentUserPrincipalParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) error); ok {
		r1 = rf(params, authInfo, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IamServiceGetOrganizationAuthMetadata provides a mock function with given fields: params, authInfo, opts
func (_m *ClientService) IamServiceGetOrganizationAuthMetadata(params *iam_service.IamServiceGetOrganizationAuthMetadataParams, authInfo runtime.ClientAuthInfoWriter, opts ...iam_service.ClientOption) (*iam_service.IamServiceGetOrganizationAuthMetadataOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params, authInfo)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *iam_service.IamServiceGetOrganizationAuthMetadataOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceGetOrganizationAuthMetadataParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) (*iam_service.IamServiceGetOrganizationAuthMetadataOK, error)); ok {
		return rf(params, authInfo, opts...)
	}
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceGetOrganizationAuthMetadataParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) *iam_service.IamServiceGetOrganizationAuthMetadataOK); ok {
		r0 = rf(params, authInfo, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam_service.IamServiceGetOrganizationAuthMetadataOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*iam_service.IamServiceGetOrganizationAuthMetadataParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) error); ok {
		r1 = rf(params, authInfo, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IamServiceGetUserPrincipalByIDInOrganization provides a mock function with given fields: params, authInfo, opts
func (_m *ClientService) IamServiceGetUserPrincipalByIDInOrganization(params *iam_service.IamServiceGetUserPrincipalByIDInOrganizationParams, authInfo runtime.ClientAuthInfoWriter, opts ...iam_service.ClientOption) (*iam_service.IamServiceGetUserPrincipalByIDInOrganizationOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params, authInfo)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *iam_service.IamServiceGetUserPrincipalByIDInOrganizationOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceGetUserPrincipalByIDInOrganizationParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) (*iam_service.IamServiceGetUserPrincipalByIDInOrganizationOK, error)); ok {
		return rf(params, authInfo, opts...)
	}
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceGetUserPrincipalByIDInOrganizationParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) *iam_service.IamServiceGetUserPrincipalByIDInOrganizationOK); ok {
		r0 = rf(params, authInfo, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam_service.IamServiceGetUserPrincipalByIDInOrganizationOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*iam_service.IamServiceGetUserPrincipalByIDInOrganizationParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) error); ok {
		r1 = rf(params, authInfo, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IamServiceGetUserPrincipalsByIDsInOrganization provides a mock function with given fields: params, authInfo, opts
func (_m *ClientService) IamServiceGetUserPrincipalsByIDsInOrganization(params *iam_service.IamServiceGetUserPrincipalsByIDsInOrganizationParams, authInfo runtime.ClientAuthInfoWriter, opts ...iam_service.ClientOption) (*iam_service.IamServiceGetUserPrincipalsByIDsInOrganizationOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params, authInfo)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *iam_service.IamServiceGetUserPrincipalsByIDsInOrganizationOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceGetUserPrincipalsByIDsInOrganizationParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) (*iam_service.IamServiceGetUserPrincipalsByIDsInOrganizationOK, error)); ok {
		return rf(params, authInfo, opts...)
	}
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceGetUserPrincipalsByIDsInOrganizationParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) *iam_service.IamServiceGetUserPrincipalsByIDsInOrganizationOK); ok {
		r0 = rf(params, authInfo, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam_service.IamServiceGetUserPrincipalsByIDsInOrganizationOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*iam_service.IamServiceGetUserPrincipalsByIDsInOrganizationParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) error); ok {
		r1 = rf(params, authInfo, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IamServiceListUserPrincipalsByOrganization provides a mock function with given fields: params, authInfo, opts
func (_m *ClientService) IamServiceListUserPrincipalsByOrganization(params *iam_service.IamServiceListUserPrincipalsByOrganizationParams, authInfo runtime.ClientAuthInfoWriter, opts ...iam_service.ClientOption) (*iam_service.IamServiceListUserPrincipalsByOrganizationOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params, authInfo)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *iam_service.IamServiceListUserPrincipalsByOrganizationOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceListUserPrincipalsByOrganizationParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) (*iam_service.IamServiceListUserPrincipalsByOrganizationOK, error)); ok {
		return rf(params, authInfo, opts...)
	}
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceListUserPrincipalsByOrganizationParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) *iam_service.IamServiceListUserPrincipalsByOrganizationOK); ok {
		r0 = rf(params, authInfo, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam_service.IamServiceListUserPrincipalsByOrganizationOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*iam_service.IamServiceListUserPrincipalsByOrganizationParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) error); ok {
		r1 = rf(params, authInfo, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IamServicePing provides a mock function with given fields: params, opts
func (_m *ClientService) IamServicePing(params *iam_service.IamServicePingParams, opts ...iam_service.ClientOption) (*iam_service.IamServicePingOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *iam_service.IamServicePingOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*iam_service.IamServicePingParams, ...iam_service.ClientOption) (*iam_service.IamServicePingOK, error)); ok {
		return rf(params, opts...)
	}
	if rf, ok := ret.Get(0).(func(*iam_service.IamServicePingParams, ...iam_service.ClientOption) *iam_service.IamServicePingOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam_service.IamServicePingOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*iam_service.IamServicePingParams, ...iam_service.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IamServiceSearchPrincipals provides a mock function with given fields: params, authInfo, opts
func (_m *ClientService) IamServiceSearchPrincipals(params *iam_service.IamServiceSearchPrincipalsParams, authInfo runtime.ClientAuthInfoWriter, opts ...iam_service.ClientOption) (*iam_service.IamServiceSearchPrincipalsOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params, authInfo)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *iam_service.IamServiceSearchPrincipalsOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceSearchPrincipalsParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) (*iam_service.IamServiceSearchPrincipalsOK, error)); ok {
		return rf(params, authInfo, opts...)
	}
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceSearchPrincipalsParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) *iam_service.IamServiceSearchPrincipalsOK); ok {
		r0 = rf(params, authInfo, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam_service.IamServiceSearchPrincipalsOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*iam_service.IamServiceSearchPrincipalsParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) error); ok {
		r1 = rf(params, authInfo, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IamServiceUpdateWebConsolePreferences provides a mock function with given fields: params, authInfo, opts
func (_m *ClientService) IamServiceUpdateWebConsolePreferences(params *iam_service.IamServiceUpdateWebConsolePreferencesParams, authInfo runtime.ClientAuthInfoWriter, opts ...iam_service.ClientOption) (*iam_service.IamServiceUpdateWebConsolePreferencesOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params, authInfo)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *iam_service.IamServiceUpdateWebConsolePreferencesOK
	var r1 error
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceUpdateWebConsolePreferencesParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) (*iam_service.IamServiceUpdateWebConsolePreferencesOK, error)); ok {
		return rf(params, authInfo, opts...)
	}
	if rf, ok := ret.Get(0).(func(*iam_service.IamServiceUpdateWebConsolePreferencesParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) *iam_service.IamServiceUpdateWebConsolePreferencesOK); ok {
		r0 = rf(params, authInfo, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam_service.IamServiceUpdateWebConsolePreferencesOK)
		}
	}

	if rf, ok := ret.Get(1).(func(*iam_service.IamServiceUpdateWebConsolePreferencesParams, runtime.ClientAuthInfoWriter, ...iam_service.ClientOption) error); ok {
		r1 = rf(params, authInfo, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetTransport provides a mock function with given fields: transport
func (_m *ClientService) SetTransport(transport runtime.ClientTransport) {
	_m.Called(transport)
}

// NewClientService creates a new instance of ClientService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClientService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ClientService {
	mock := &ClientService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
