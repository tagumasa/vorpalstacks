// Package appsync provides the storage layer for AWS AppSync service resources.
package appsync

import "errors"

// Sentinel errors for AppSync store operations.
var (
	ErrApiNotFound                  = errors.New("api not found")
	ErrChannelNamespaceNotFound     = errors.New("channel namespace not found")
	ErrApiAlreadyExists             = errors.New("api already exists")
	ErrChannelNamespaceExists       = errors.New("channel namespace already exists")
	ErrGraphqlApiNotFound           = errors.New("graphql api not found")
	ErrGraphqlApiAlreadyExists      = errors.New("graphql api already exists")
	ErrDataSourceNotFound           = errors.New("data source not found")
	ErrDataSourceAlreadyExists      = errors.New("data source already exists")
	ErrResolverNotFound             = errors.New("resolver not found")
	ErrResolverAlreadyExists        = errors.New("resolver already exists")
	ErrFunctionNotFound             = errors.New("function not found")
	ErrFunctionAlreadyExists        = errors.New("function already exists")
	ErrTypeNotFound                 = errors.New("type not found")
	ErrTypeAlreadyExists            = errors.New("type already exists")
	ErrSchemaCreationNotFound       = errors.New("schema creation status not found")
	ErrApiKeyNotFound               = errors.New("api key not found")
	ErrApiCacheNotFound             = errors.New("api cache not found")
	ErrDomainNameNotFound           = errors.New("domain name not found")
	ErrDomainNameAlreadyExists      = errors.New("domain name already exists")
	ErrApiAssociationNotFound       = errors.New("api association not found")
	ErrMergedApiAssociationNotFound = errors.New("merged api association not found")
)
