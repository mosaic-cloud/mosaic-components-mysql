

package main


import "errors"
import "fmt"
import "net"

import "mosaic-components/libraries/backend"
import "mosaic-components-mysql/server"
import "vgl/transcript"

import . "mosaic-components/libraries/messages"


var selfGroup = ComponentGroup ("be149e7b52c7cbe0695e208081ffaefbbc5778a7")


type callbacks struct {
	server server.Server
	serverConfiguration *server.ServerConfiguration
	backend backend.Controller
	transcript transcript.Transcript
	configuration map[string]interface{}
}


func (_callbacks *callbacks) Initialized (_backend backend.Controller) (error) {
	
	_callbacks.transcript.TraceInformation ("initializing the component...")
	_callbacks.backend = _backend
	
	_callbacks.transcript.TraceInformation ("acquiring the SQL endpoint...")
	var _ip net.IP
	var _port uint16
	if _ip_1, _port_1, _fqdn, _error := backend.TcpSocketAcquireSync (_callbacks.backend, ResourceIdentifier ("sql")); _error != nil {
		panic (_error)
	} else {
		_ip = _ip_1
		_port = _port_1
		_callbacks.transcript.TraceInformation ("  * aquired `%s:%d` / `%s`.", _ip, _port, _fqdn)
	}
	
	_callbacks.transcript.TraceInformation ("creating the MySQL server...")
	if _serverConfiguration, _error := server.ResolveDefaultServerConfiguration (); _error != nil {
		panic (_error)
	} else {
		_serverConfiguration.SqlEndpointIp = _ip
		_serverConfiguration.SqlEndpointPort = _port
		_callbacks.serverConfiguration = _serverConfiguration
	}
	if _server, _error := server.Create (_callbacks.serverConfiguration); _error != nil {
		panic (_error)
	} else {
		_callbacks.server = _server
	}
	
	_callbacks.transcript.TraceInformation ("initializing (and bootstrapping) the MySQL server...")
	if _error := _callbacks.server.Initialize (true); _error != nil {
		panic (_error)
	}
	
	_callbacks.transcript.TraceInformation ("registering the component...")
	if _error := _callbacks.backend.ComponentRegisterSync (selfGroup); _error != nil {
		panic (_error)
	}
	
	_callbacks.transcript.TraceInformation ("initialized the component.")
	
	return nil
}


func (_callbacks *callbacks) Terminated (_error error) (error) {
	
	_callbacks.transcript.TraceInformation ("terminating the component...")
	
	_callbacks.transcript.TraceInformation ("terminating the MySQL server...")
	if _error := _callbacks.server.Terminate (); _error != nil {
		panic (_error)
	}
	
	_callbacks.transcript.TraceInformation ("terminated the component.")
	return nil
}


func (_callbacks *callbacks) ComponentCallInvoked (_operation ComponentOperation, _inputs interface{}, _correlation Correlation, _attachment Attachment) (error) {
	
	switch _operation {
		
		case "mosaic-mysql:get-sql-endpoint" :
			
			_outputs := map[string]interface{} {
					"ip" : _callbacks.serverConfiguration.SqlEndpointIp.String (),
					"port" : _callbacks.serverConfiguration.SqlEndpointPort,
					"administrator-login" : _callbacks.serverConfiguration.SqlAdministratorLogin,
					"administrator-password" : _callbacks.serverConfiguration.SqlAdministratorPassword,
					"url" : fmt.Sprintf (
							"mysql://%s:%d/mysql?user=%s&password=%s",
							_callbacks.serverConfiguration.SqlEndpointIp.String (),
							_callbacks.serverConfiguration.SqlEndpointPort,
							_callbacks.serverConfiguration.SqlAdministratorLogin,
							_callbacks.serverConfiguration.SqlAdministratorPassword),
			}
			
			if _error := _callbacks.backend.ComponentCallSucceeded (_correlation, _outputs, nil); _error != nil {
				panic (_error)
			}
			
		default :
			
			_callbacks.transcript.TraceError ("invoked invalid component call operation `%s`; ignoring!", _operation)
			if _error := _callbacks.backend.ComponentCallFailed (_correlation, errors.New ("invalid-operation"), nil); _error != nil {
				panic (_error)
			}
	}
	
	return nil
}


func (_callbacks *callbacks) ComponentCastInvoked (_operation ComponentOperation, _inputs interface{}, _attachment Attachment) (error) {
	_callbacks.transcript.TraceError ("invoked invalid component cast operation `%s`; ignoring!", _operation)
	return nil
}

func (_callbacks *callbacks) ComponentCallSucceeded (_correlation Correlation, _outputs interface{}, _attachment Attachment) (error) {
	_callbacks.transcript.TraceError ("returned unexpected component call `%s`; ignoring!", _correlation)
	return nil
}

func (_callbacks *callbacks) ComponentCallFailed (_correlation Correlation, _error interface{}, _attachment Attachment) (error) {
	_callbacks.transcript.TraceError ("returned unexpected component call `%s`; ignoring!", _correlation)
	return nil
}

func (_callbacks *callbacks) ComponentRegisterSucceeded (_correlation Correlation) (error) {
	_callbacks.transcript.TraceError ("returned unexpected component register `%s`; ignoring!", _correlation)
	return nil
}

func (_callbacks *callbacks) ComponentRegisterFailed (_correlation Correlation, _error interface{}) (error) {
	_callbacks.transcript.TraceError ("returned unexpected component register `%s`; ignoring!", _correlation)
	return nil
}

func (_callbacks *callbacks) ResourceAcquireSucceeded (_correlation Correlation, _descriptor ResourceDescriptor) (error) {
	_callbacks.transcript.TraceError ("returned unexpected resource acquire `%s`; ignoring!", _correlation)
	return nil
}

func (_callbacks *callbacks) ResourceAcquireFailed (_correlation Correlation, _error interface{}) (error) {
	_callbacks.transcript.TraceError ("returned unexpected resource acquire `%s`; ignoring!", _correlation)
	return nil
}


func Main (_componentIdentifier string, _channelEndpoint string, _configuration map[string]interface{}) (error) {
	
	_callbacks := & callbacks {}
	_callbacks.transcript = transcript.NewTranscript (_callbacks, packageTranscript)
	_callbacks.configuration = _configuration
	
	return backend.Execute (_callbacks, _componentIdentifier, _channelEndpoint)
}

func main () () {
	backend.PreMain (Main)
}


var packageTranscript = transcript.NewPackageTranscript ()
