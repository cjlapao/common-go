package mongodb

//TODO: Refactor implementation
import (
	"strings"

	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/guard"
)

var globalMongoDBService *MongoDBService
var globalFactory *MongoFactory
var tenantFactory *MongoFactory

type MongoDBServiceOptions struct {
	ConnectionString   string
	GlobalDatabaseName string
}

type MongoDBService struct {
	ConnectionString   string
	GlobalDatabaseName string
	TenantDatabaseName string
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
		ConnectionString:   options.ConnectionString,
		GlobalDatabaseName: options.GlobalDatabaseName,
	}
	if options.ConnectionString != "" && options.GlobalDatabaseName != "" {
		globalFactory = NewFactory(service.ConnectionString).WithDatabase(service.GlobalDatabaseName)
	}

	globalMongoDBService = &service
	return globalMongoDBService
}

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

func Get() *MongoDBService {
	if globalMongoDBService != nil {
		return globalMongoDBService
	}

	return New()
}

func (service *MongoDBService) WithDatabase(databaseName string) *MongoDBService {
	guard.FatalEmptyOrNil(databaseName, "Database name is empty")
	if !strings.EqualFold(service.GlobalDatabaseName, databaseName) {
		service.GlobalDatabaseName = databaseName
		Init()
	}

	return service
}

// GlobalDatabase Gets the global database factory and initiate it ready for consuption
// This will try to only keep a client per session to avoid starvation of the clients
func (service *MongoDBService) GlobalDatabase() *MongoFactory {
	if globalFactory == nil {
		logger.Info("Global factory not initiated, creating instance now.")
		Init()
	}

	return globalFactory
}

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
