

package main


import "fmt"
import "net"
import "os"
import "strings"

import "mosaic-components/libraries/backend"
import "mosaic-components/libraries/channels"
import "mosaic-components-mysql/server"
import "vgl/transcript"

import . "mosaic-components/libraries/messages"


func main () () {
	
	var _channelEndpoint string
	if false {
		_channelEndpoint = "stdio"
	} else {
		_channelEndpoint = "tcp:127.0.0.1:24704"
	}
	
	if _error := execute (_channelEndpoint); _error != nil {
		panic (_error)
	}
}


func execute (_channelEndpoint string) (error) {
	
	_callbacks := & callbacks {
			server : nil,
			backend : nil,
			transcript : nil,
	}
	_callbacks.transcript = transcript.NewTranscript (_callbacks, packageTranscript)
	_transcript := packageTranscript
	var _error error
	
	_transcript.TraceInformation ("initializing...")
	
	_transcript.TraceInformation ("creating the component backend...")
	var _backend backend.Backend
	var _backendChannelCallbacks channels.Callbacks
	if _backend, _backendChannelCallbacks, _error = backend.Create (_callbacks); _error != nil {
		panic (_error)
	}
	
	_transcript.TraceInformation ("creating the component channel...")
	var _channel channels.Channel
	if _channelEndpoint == "stdio" {
		_transcript.TraceInformation ("  * using the stdio endpoint;")
		_inboundStream := os.Stdin
		_outboundStream := os.Stdout
		if _channel, _error = channels.Create (_backendChannelCallbacks, _inboundStream, _outboundStream, nil); _error != nil {
			panic (_error)
		}
	} else if strings.HasPrefix (_channelEndpoint, "tcp:") {
		_channelTcpEndpoint := _channelEndpoint[4:]
		_transcript.TraceInformation ("  * usig the TCP endpoint `%s`;", _channelTcpEndpoint)
		if _channel, _error = channels.CreateAndDial (_backendChannelCallbacks, "tcp", _channelTcpEndpoint); _error != nil {
			panic (_error)
		}
	} else {
		_transcript.TraceError ("invalid component channel endpoint; aborting!")
		panic ("failed")
	}
	
	_transcript.TraceInformation ("executing...")
	
	_transcript.TraceInformation ("waiting for the termination of the component backend...")
	if _error := _backend.WaitTerminated (); _error != nil {
		panic (_error)
	}
	
	_transcript.TraceInformation ("terminating the component channel...")
	if _error := _channel.Terminate (); _error != nil {
		panic (_error)
	}
	
	_transcript.TraceInformation ("terminated.")
	return nil
}


type callbacks struct {
	server server.Server
	serverConfiguration *server.ServerConfiguration
	backend backend.Controller
	terminated chan error
	transcript transcript.Transcript
}


func (_callbacks *callbacks) Initialized (_backend backend.Controller) (error) {
	
	_callbacks.transcript.TraceInformation ("initializing the component...")
	_callbacks.backend = _backend
	
	_callbacks.transcript.TraceInformation ("acquiring the SQL endpoint...")
	var _ip net.IP
	var _port uint16
	if _ip_1, _port_1, _fqdn, _error := backend.TcpSocketAcquireSync (_callbacks.backend, ResourceIdentifier ("tests")); _error != nil {
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
	if _error := _callbacks.backend.ComponentRegisterSync (componentGroup); _error != nil {
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
	
	if _operation == "mosaic-mysql:get-sql-endpoint" {
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
		return nil
	} else {
		_callbacks.transcript.TraceError ("invoked invalid component call operation `%s`; ignoring!", _operation)
		return nil
	}
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


var packageTranscript = transcript.NewPackageTranscript ()
var componentGroup = ComponentGroup ("be149e7b52c7cbe0695e208081ffaefbbc5778a7")
