package mongodb

import (
	"strings"

	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/guard"
)

// Global mongoDB service to keep single service for consumers
var globalMongoDBService *MongoDBService

// Global global database factory to keep a single mongodb client
var globalFactory *MongoFactory

// Global tenant database factory to keep a single mongodb client
var tenantFactory *MongoFactory

// MongoDBServiceOptions structure
type MongoDBServiceOptions struct {
	ConnectionString   string
	GlobalDatabaseName string
}

// MongoDBService structure
type MongoDBService struct {
	ConnectionString   string
	GlobalDatabaseName string
	TenantDatabaseName string
}

// New Creates a MongoDB service using the default configuration
// This uses the environment variables to define the connection
// and the database name, the variables are:
// MONGODB_CONNECTION_STRING: for the connection string
// MONGODB_DATABASENAME: for the database name
// returns a MongoDBService pointer
func New() *MongoDBService {
	ctx := execution_context.Get()
	connStr := ctx.Configuration.GetString("MONGODB_CONNECTION_STRING")
	globalDatabaseName := ctx.Configuration.GetString("MONGODB_DATABASENAME")

	options := MongoDBServiceOptions{
		ConnectionString:   connStr,
		GlobalDatabaseName: globalDatabaseName,
	}

	return NewWithOptions(options)
}

// NewWithOptions Creates a MongoDB service passing the options object
// returns a MongoDBService pointer
func NewWithOptions(options MongoDBServiceOptions) *MongoDBService {
	service := MongoDBService{
		ConnectionString:   options.ConnectionString,
		GlobalDatabaseName: options.GlobalDatabaseName,
	}
	if options.ConnectionString != "" && options.GlobalDatabaseName != "" {
		globalFactory = NewFactory(service.ConnectionString).WithDatabase(service.GlobalDatabaseName)
	}

	globalMongoDBService = &service
	return globalMongoDBService
}

// Init initiates the MongoDB service and global database factory
// returns a MongoDBService pointer
func Init() *MongoDBService {
	if globalMongoDBService != nil {
		if globalMongoDBService.ConnectionString != "" && globalMongoDBService.GlobalDatabaseName != "" {
			logger.Info("Initiating MongoDB Service for global database %v", globalMongoDBService.GlobalDatabaseName)
			globalFactory = NewFactory(globalMongoDBService.ConnectionString).WithDatabase(globalMongoDBService.GlobalDatabaseName)
			logger.Info("MongoDB Service for global database %v initiated successfully", globalMongoDBService.GlobalDatabaseName)
		}
		return globalMongoDBService
	}

	return New()
}

// Get Gets the current global service
// returns a MongoDBService pointer
func Get() *MongoDBService {
	if globalMongoDBService != nil {
		return globalMongoDBService
	}

	return New()
}

// WithDatabase Sets the global database name
// returns a MongoDBService pointer
func (service *MongoDBService) WithDatabase(databaseName string) *MongoDBService {
	guard.FatalEmptyOrNil(databaseName, "Database name is empty")
	if !strings.EqualFold(service.GlobalDatabaseName, databaseName) {
		service.GlobalDatabaseName = databaseName
		Init()
	}

	return service
}

// GlobalDatabase Gets the global database factory and initiate it ready for consuption.
// This will try to only keep a client per session to avoid starvation of the clients
// returns a MongoFactory pointer
func (service *MongoDBService) GlobalDatabase() *MongoFactory {
	if globalFactory == nil {
		logger.Info("Global factory not initiated, creating instance now.")
		Init()
	}

	return globalFactory
}

// TenantDatabase Gets the tenant database factory and initiate it ready for consumption
// if there is no tenant set this will bring the global database and we will treat it as
// a single tenant system.
// This will try to only keep a client per session to avoid starvation of the clients
// returns a MongoFactory pointer
func (service *MongoDBService) TenantDatabase() *MongoFactory {
	ctx := execution_context.Get()
	tenantId := ctx.Authorization.TenantId
	if tenantId == "" || strings.ToLower(tenantId) == "global" {
		return service.GlobalDatabase()
	}
	if !strings.EqualFold(tenantId, service.TenantDatabaseName) {
		service.TenantDatabaseName = tenantId
		logger.Info("Initiating MongoDB Service for tenant database %v", service.TenantDatabaseName)
		tenantFactory = NewFactory(service.ConnectionString).WithDatabase(service.TenantDatabaseName)
		logger.Info("MongoDB Service for tenant database %v initiated successfully", service.TenantDatabaseName)
	}

	return tenantFactory
}

func (service *MongoDBService) GetTenant() string {
	service.TenantDatabase()
	return service.TenantDatabaseName
}
