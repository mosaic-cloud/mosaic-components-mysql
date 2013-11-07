

package server


import "net"
import "os"


type GenericConfiguration struct {
	
	WorkspacePath string
	DatabasesPath string
	TemporaryPath string
	SocketPath string
	PidPath string
	
	PackageBasePath string
	ExecutablePath string
	PluginsPath string
	CharsetsPath string
}

type ServerConfiguration struct {
	
	SqlEndpointIp net.IP
	SqlEndpointPort uint16
	SqlAdministratorLogin string
	SqlAdministratorPassword string
	SqlInitializationScriptPaths []string
	
	GenericConfiguration *GenericConfiguration
}


func ResolveDefaultServerConfiguration () (*ServerConfiguration, error) {
	
	var _genericConfiguration *GenericConfiguration
	if _genericConfiguration_1, _error := ResolveDefaultGenericConfiguration (); _error != nil {
		return nil, _error
	} else {
		_genericConfiguration = _genericConfiguration_1
	}
	
	_configuration := & ServerConfiguration {
			SqlEndpointIp : DefaultSqlEndpointIp,
			SqlEndpointPort : DefaultSqlEndpointPort,
			SqlAdministratorLogin : DefaultSqlAdministratorLogin,
			SqlAdministratorPassword : DefaultSqlAdministratorPassword,
			SqlInitializationScriptPaths : DefaultSqlInitializationScriptPaths,
			GenericConfiguration : _genericConfiguration,
	}
	
	return _configuration, nil
}

func ResolveDefaultGenericConfiguration () (*GenericConfiguration, error) {
	
	_configuration := & GenericConfiguration {
			WorkspacePath : DefaultWorkspacePath,
			DatabasesPath : DefaultDatabasesPath,
			TemporaryPath : DefaultTemporaryPath,
			SocketPath : DefaultSocketPath,
			PidPath : DefaultPidPath,
			PackageBasePath : DefaultPackageBasePath,
			ExecutablePath : DefaultExecutablePath,
			PluginsPath : DefaultPluginsPath,
			CharsetsPath : DefaultCharsetsPath,
	}
	
	_workspace := os.Getenv ("mosaic_component_temporary")
	if _workspace != "" {
		_configuration.WorkspacePath = _workspace
		_configuration.DatabasesPath = _workspace + "/databases"
		_configuration.TemporaryPath = _workspace + "/temporary"
		_configuration.SocketPath = _workspace + "/server.sock"
		_configuration.PidPath = _workspace + "/server.pid"
	}
	
	return _configuration, nil
}
