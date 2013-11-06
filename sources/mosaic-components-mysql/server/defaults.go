

package server


import "net"


const (
	DefaultWorkspacePath = "/tmp/mysql-component"
	DefaultDatabasesPath = DefaultWorkspacePath + "/databases"
	DefaultTemporaryPath = DefaultWorkspacePath + "/temporary"
	DefaultSocketPath = DefaultWorkspacePath + "/server.sock"
	DefaultPidPath = DefaultWorkspacePath + "/server.pid"
)

const (
	DefaultPackageBasePath string = "/usr"
	DefaultBinBasePath string = DefaultPackageBasePath + "/bin"
	DefaultDataBasePath string = DefaultPackageBasePath + "/share/mysql"
	DefaultExecutablePath = DefaultBinBasePath + "/mysqld"
	DefaultPluginsPath = DefaultDataBasePath + "/plugin"
	DefaultCharsetsPath = DefaultDataBasePath + "/charsets"
)

var DefaultSqlInitializationScriptPaths = []string {
		DefaultDataBasePath + "/mysql_system_tables.sql",
		DefaultDataBasePath + "/mysql_system_tables_data.sql",
		DefaultDataBasePath + "/mysql_performance_tables.sql",
}

var DefaultSqlEndpointIp = net.ParseIP ("0.0.0.0")
var DefaultSqlEndpointPort = uint16 (28203)
var DefaultSqlAdministratorLogin = "root"
var DefaultSqlAdministratorPassword = "31b21446c3cc6f36dabd19bfd6d8c6c1"
