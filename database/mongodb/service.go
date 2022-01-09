package mongodb

import (
	"strings"

	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/guard"
)

var globalMongoDBService *MongoDBService

type MongoDBServiceOptions struct {
	ConnectionString   string
	GlobalDatabaseName string
}

type MongoDBService struct {
	ConnectionString string
	GlobalDatabase   string
	GlobalFactory    *MongoFactory
	TenantDatabase   string
	TenantFactory    *MongoFactory
}

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

func NewWithOptions(options MongoDBServiceOptions) *MongoDBService {
	service := MongoDBService{
		ConnectionString: options.ConnectionString,
		GlobalDatabase:   options.GlobalDatabaseName,
	}
	if options.ConnectionString != "" && options.GlobalDatabaseName != "" {
		service.GlobalFactory = NewFactory(service.ConnectionString).WithDatabase(service.GlobalDatabase)
	}

	globalMongoDBService = &service
	return globalMongoDBService
}

func Init() *MongoDBService {
	if globalMongoDBService != nil {
		if globalMongoDBService.ConnectionString != "" && globalMongoDBService.GlobalDatabase != "" {
			logger.Info("Initiating MongoDB Service for global database %v", globalMongoDBService.GlobalDatabase)
			globalMongoDBService.GlobalFactory = NewFactory(globalMongoDBService.ConnectionString).WithDatabase(globalMongoDBService.GlobalDatabase)
			logger.Info("MongoDB Service for global database %v initiated successfully", globalMongoDBService.GlobalDatabase)
		}
		return globalMongoDBService
	}

	return New()
}

func Get() *MongoDBService {
	if globalMongoDBService != nil {
		return globalMongoDBService
	}

	return New()
}

func (service *MongoDBService) WithDatabase(databaseName string) *MongoDBService {
	guard.FatalEmptyOrNil(databaseName, "Database name is empty")
	if !strings.EqualFold(service.GlobalDatabase, databaseName) {
		service.GlobalDatabase = databaseName
		Init()
	}

	return service
}

func (service *MongoDBService) GetGlobalFactory() *MongoFactory {
	return service.GlobalFactory
}

func (service *MongoDBService) GetTenantDatabase() (*MongoFactory, string) {
	ctx := execution_context.Get()
	tenantId := ctx.Authorization.TenantId
	if tenantId == "" || strings.ToLower(tenantId) == "global" {
		return service.GlobalFactory, service.GlobalDatabase
	}
	if !strings.EqualFold(tenantId, service.TenantDatabase) {
		service.TenantDatabase = tenantId
		logger.Info("Initiating MongoDB Service for tenant database %v", service.TenantDatabase)
		service.TenantFactory = NewFactory(service.ConnectionString).WithDatabase(service.TenantDatabase)
		logger.Info("MongoDB Service for tenant database %v initiated successfully", service.TenantDatabase)
	}

	return service.TenantFactory, service.TenantDatabase
}

func (service *MongoDBService) GetTenant() string {
	_, tenant := service.GetTenantDatabase()
	return tenant
}
