package database

type DatabaseContext struct {
	CurrentDatabaseName string
	ConnectionString    string
	CurrentCollection   string
}
