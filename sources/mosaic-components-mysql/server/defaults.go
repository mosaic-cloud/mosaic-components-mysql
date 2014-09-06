

package server


import "net"
import "os"
import "strconv"


const (
	DefaultWorkspacePath = "/tmp/mosaic-components-mysql"
	DefaultPackageBasePath string = "/usr"
)

const (
	DefaultSqlEndpointIp = "0.0.0.0"
	DefaultSqlEndpointPort = "28203"
	DefaultSqlAdministratorLogin = "root"
	DefaultSqlAdministratorPassword = "31b21446c3cc6f36dabd19bfd6d8c6c1"
)


func ResolveDefaultWorkspacePath () (string) {
	return resolveDefaultValue (
			"mosaic_components_mysql__workspace",
			DefaultWorkspacePath)
}

func ResolveDefaultDatabasesPath () (string) {
	return resolveDefaultValue (
			"mosaic_components_mysql__databases",
			ResolveDefaultWorkspacePath () + "/databases")
}

func ResolveDefaultTemporaryPath () (string) {
	return resolveDefaultValue (
			"mosaic_components_mysql__temporary",
			ResolveDefaultWorkspacePath () + "/temporary")
}

func ResolveDefaultSocketPath () (string) {
	return resolveDefaultValue (
			"mosaic_components_mysql__server_sock",
			ResolveDefaultWorkspacePath () + "/server.sock")
}

func ResolveDefaultPidPath () (string) {
	return resolveDefaultValue (
			"mosaic_components_mysql__server_pid",
			ResolveDefaultWorkspacePath () + "/server.pid")
}


func ResolveDefaultPackageBasePath () (string) {
	return resolveDefaultValue (
			"mosaic_components_mysql__package",
			DefaultPackageBasePath)
}

func ResolveDefaultBinBasePath () (string) {
	return resolveDefaultValue (
			"mosaic_components_mysql__package_bin",
			ResolveDefaultPackageBasePath () + "/bin")
}

func ResolveDefaultLibBasePath () (string) {
	return resolveDefaultValue (
			"mosaic_components_mysql__package_lib",
			ResolveDefaultPackageBasePath () + "/lib/mysql")
}

func ResolveDefaultDataBasePath () (string) {
	return resolveDefaultValue (
			"mosaic_components_mysql__package_data",
			ResolveDefaultPackageBasePath () + "/share/mysql")
}


func ResolveDefaultExecutablePath () (string) {
	return resolveDefaultValue (
			"mosaic_components_mysql__server_executable",
			ResolveDefaultBinBasePath () + "/mysqld")
}

func ResolveDefaultPluginsPath () (string) {
	return resolveDefaultValue (
			"mosaic_components_mysql__server_plugins",
			ResolveDefaultLibBasePath () + "/plugin")
}

func ResolveDefaultCharsetsPath () (string) {
	return resolveDefaultValue (
			"mosaic_components_mysql__server_charsets",
			ResolveDefaultDataBasePath () + "/charsets")
}

func ResolveDefaultSqlScriptBasePath () (string) {
	return resolveDefaultValue (
			"mosaic_components_mysql__sql_scripts",
			ResolveDefaultDataBasePath ())
}

func ResolveDefaultSqlInitializationScriptPaths () ([]string) {
	return []string {
		ResolveDefaultSqlScriptBasePath () + "/mysql_system_tables.sql",
		ResolveDefaultSqlScriptBasePath () + "/mysql_system_tables_data.sql",
		ResolveDefaultSqlScriptBasePath () + "/mysql_performance_tables.sql",
	}
}


func ResolveDefaultSqlEndpointIp () (net.IP) {
	return net.ParseIP (
			resolveDefaultValue (
					"mosaic_components_mysql__server_endpoint_ip",
					DefaultSqlEndpointIp))
}

func ResolveDefaultSqlEndpointPort () (uint16) {
	_port, _error := strconv.ParseUint (
			resolveDefaultValue (
					"mosaic_components_mysql__server_endpoint_port",
					DefaultSqlEndpointPort),
			10, 16)
	if _error != nil {
		panic (_error)
	}
	return uint16 (_port)
}

func ResolveDefaultSqlAdministratorLogin () (string) {
	return resolveDefaultValue (
			"mosaic_components_mysql__server_administrator_login",
			DefaultSqlAdministratorLogin)
}

func ResolveDefaultSqlAdministratorPassword () (string) {
	return resolveDefaultValue (
			"mosaic_components_mysql__server_administrator_password",
			DefaultSqlAdministratorPassword)
}


func resolveDefaultValue (_variable string, _default string) (string) {
	_value := os.Getenv (_variable)
	if _value == "" {
		_value = _default
	}
	return _value
}
