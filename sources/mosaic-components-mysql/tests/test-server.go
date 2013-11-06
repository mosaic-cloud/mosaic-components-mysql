
package main


import "time"

import . "mosaic-components-mysql/server"


func main () () {
	
	var _configuration *ServerConfiguration
	var _server Server
	var _error error
	
	if _configuration, _error = ResolveDefaultServerConfiguration (); _error != nil {
		panic (_error)
	}
	
	if _server, _error = Create (_configuration); _error != nil {
		panic (_error)
	}
	
	if _error = _server.Initialize (true); _error != nil {
		panic (_error)
	}
	
	time.Sleep (2 * time.Second)
	
	if _error = _server.Terminate (); _error != nil {
		panic (_error)
	}
}
