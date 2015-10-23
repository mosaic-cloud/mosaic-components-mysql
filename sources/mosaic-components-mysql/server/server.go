

package server


import "bufio"
import "fmt"
import "io"
import "os"
import "syscall"
import "time"

import "vgl/transcript"


type Server interface {
	Initialize (_bootstrap bool) (error)
	Terminate () (error)
}


type server struct {
	configuration *ServerConfiguration
	state serverState
	isolates chan func () ()
	process *os.Process
	transcript transcript.Transcript
}

type serverState int
const (
	invalidServerStateMin serverState = iota
	serverCreated
	serverRunning
	serverTerminated
	serverFailed
	invalidServerStateMax
)

const isolatesBufferSize = 16


func Create (_configuration *ServerConfiguration) (Server, error) {
	
	_server := & server {
		configuration : _configuration,
		state : serverCreated,
		isolates : make (chan func () (), isolatesBufferSize),
		process : nil,
		transcript : nil,
	}
	
	_server.transcript = transcript.NewTranscript (_server, packageTranscript)
	_server.transcript.TraceDebugging ("created mysql server controller.")
	
	go _server.executeLoop ()
	
	return _server, nil
}


func (_server *server) Initialize (_bootstrap bool) (error) {
	_completion := make (chan error, 1)
	defer close (_completion)
	_server.isolates <- func () () {
		if _bootstrap {
			if _error := _server.handleBootstrap (); _error != nil {
				_completion <- _error
				return
			}
		}
		if _error := _server.handleStart (); _error != nil {
			_completion <- _error
			return
		}
		_completion <- nil
	}
	return <- _completion
}


func (_server *server) Terminate () (error) {
	_completion := make (chan error, 1)
	defer close (_completion)
	_server.isolates <- func () () {
		if _server.state == serverTerminated {
			_completion <- fmt.Errorf ("illegal-state")
			return
		} else if _server.state != serverRunning {
			_completion <- fmt.Errorf ("illegal-state")
			return
		}
		if _error := _server.handleStop (); _error != nil {
			_completion <- _error
			return
		}
		_completion <- nil
	}
	return <- _completion
}


func (_server *server) handleBootstrap () (error) {
	
	if _server.state != serverCreated {
		return fmt.Errorf ("illegal-state")
	}
	
	_server.transcript.TraceDebugging ("bootstrapping...")
	
	_markerPath := _server.configuration.GenericConfiguration.DatabasesPath + "/.bootstrapp.marker"
	
	if _error := os.MkdirAll (_server.configuration.GenericConfiguration.WorkspacePath, 0700); _error != nil {
		_server.transcript.TraceDebugging ("bootstrap failed (while creating the workspace folder): `%s`!", _error.Error ())
		return _error
	}
	if _error := os.MkdirAll (_server.configuration.GenericConfiguration.DatabasesPath, 0700); _error != nil {
		_server.transcript.TraceDebugging ("bootstrap failed (while creating the databases folder): `%s`!", _error.Error ())
		return _error
	}
	if _error := os.MkdirAll (_server.configuration.GenericConfiguration.TemporaryPath, 0700); _error != nil {
		_server.transcript.TraceDebugging ("bootstrap failed (while creating the temporary folder): `%s`!", _error.Error ())
		return _error
	}
	
	var _markerFile *os.File
	defer func () () {
		if _markerFile == nil {
			return
		}
		if _, _error := _markerFile.Write ([]byte ("failed!\n")); _error != nil {
			panic (_error)
		}
		if _error := _markerFile.Close (); _error != nil {
			panic (_error)
		}
	} ()
	if _markerFile_1, _error := os.OpenFile (_markerPath, os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0444); _error != nil {
		if os.IsExist (_error) {
			_server.transcript.TraceDebugging ("already bootstrapped!")
			return nil
		} else {
			_server.transcript.TraceDebugging ("bootstrap failed (while creating the marker file): `%s`!", _error.Error ())
			return _error
		}
	} else {
		_markerFile = _markerFile_1
	}
	if _, _error := _markerFile.Write ([]byte ("pending...\n")); _error != nil {
		panic (_error)
	}
	
	_executable, _arguments, _environment, _directory := prepareBootstrapExecution (_server.configuration)
	_console := _server.prepareConsole ()
	
	var _scriptFile *os.File
	defer func () () {
		if _scriptFile == nil {
			return
		}
		if _error := _scriptFile.Close (); _error != nil {
			panic (_error)
		}
	} ()
	if _scriptFile_1, _error := prepareBootstrapScript (_server.configuration); _error != nil {
		return _error
	} else {
		_scriptFile = _scriptFile_1
	}
	
	_attributes := & os.ProcAttr {
			Env : _environment,
			Dir : _directory,
			Files : []*os.File {
					_scriptFile,
					nil,
					_console,
			},
	}
	
	if usePdeathSignal {
		_attributes.Sys = & syscall.SysProcAttr {
				Pdeathsig : syscall.SIGTERM,
		}
	}
	
	_server.transcript.TraceDebugging ("  * process executable: `%v`", _executable)
	_server.transcript.TraceDebugging ("  * process arguments: `%v`", _arguments)
	
	var _process *os.Process
	if _process_1, _error := os.StartProcess (_executable, _arguments, _attributes); _error != nil {
		_server.transcript.TraceDebugging ("bootstrap failed (while starting the process): `%s`!", _error.Error ())
		return _error
	} else {
		_process = _process_1
	}
	
	if _error := _scriptFile.Close (); _error != nil {
		panic (_error)
	}
	_scriptFile = nil
	
	if _state, _error := _process.Wait (); _error != nil {
		_server.transcript.TraceDebugging ("bootstrap failed (while waiting for the process): `%s`!", _error.Error ())
		return _error
	} else if !_state.Success () {
		_exit := _state.Sys () .(syscall.WaitStatus)
		_server.transcript.TraceDebugging ("bootstrap failed (process failed): exit code `%d`, exit signal `%d`!", _exit.ExitStatus (), _exit.Signal ())
		time.Sleep (consoleFlushTimeout)
		if !ignoreBootstrappExitCode {
			return fmt.Errorf ("bootstrapping failed")
		}
	}
	
	if _error := _markerFile.Truncate (0); _error != nil {
		panic (_error)
	}
	if _error := _markerFile.Close (); _error != nil {
		panic (_error)
	}
	_markerFile = nil
	
	_server.transcript.TraceDebugging ("bootstrapped.")
	return nil
}


func (_server *server) handleStart () (error) {
	
	if _server.state != serverCreated {
		return fmt.Errorf ("illegal-state")
	}
	
	_server.transcript.TraceDebugging ("starting...")
	
	_executable, _arguments, _environment, _directory := prepareServerExecution (_server.configuration)
	_console := _server.prepareConsole ()
	
	_attributes := & os.ProcAttr {
			Env : _environment,
			Dir : _directory,
			Files : []*os.File {
					nil,
					nil,
					_console,
			},
	}
	
	if usePdeathSignal {
		_attributes.Sys = & syscall.SysProcAttr {
				Pdeathsig : syscall.SIGTERM,
		}
	}
	
	_server.transcript.TraceDebugging ("  * process executable: `%v`", _executable)
	_server.transcript.TraceDebugging ("  * process arguments: `%v`", _arguments)
	
	var _process *os.Process
	if _process_1, _error := os.StartProcess (_executable, _arguments, _attributes); _error != nil {
		_server.transcript.TraceDebugging ("staring failed (while starting the process): `%s`!", _error.Error ())
		return _error
	} else {
		_process = _process_1
	}
	
	_server.process = _process
	_server.state = serverRunning
	
	_server.transcript.TraceDebugging ("started.")
	return nil
}


func (_server *server) handleStop () (error) {
	
	_server.transcript.TraceDebugging ("stopping...")
	
	if _server.state != serverRunning {
		return fmt.Errorf ("illegal-state")
	}
	
	if _error := _server.process.Signal (syscall.SIGTERM); _error != nil {
		_server.transcript.TraceDebugging ("stopping failed: `%s`!", _error.Error ())
		return _error
	}
	
	if _state, _error := _server.process.Wait (); _error != nil {
		_server.transcript.TraceDebugging ("stopping failed: `%s`!", _error.Error ())
		return _error
	} else if !_state.Success () {
		_exit := _state.Sys () .(syscall.WaitStatus)
		_server.transcript.TraceDebugging ("stopping failed: exit code `%d`, exit signal `%d`!", _exit.ExitStatus (), _exit.Signal ())
		time.Sleep (consoleFlushTimeout)
		if !ignoreServerExitCode {
			return fmt.Errorf ("stopping failed")
		}
	}
	
	_server.state = serverTerminated
	
	_server.transcript.TraceDebugging ("stopped...")
	return nil
}


func (_server *server) executeLoop () () {
	for {
		_isolate, _ok := <- _server.isolates
		if !_ok {
			_server.isolates = nil
			break
		}
		_isolate ()
	}
	if _server.isolates != nil {
		close (_server.isolates)
		_server.isolates = nil
	}
}


func prepareBootstrapScript (_configuration *ServerConfiguration) (*os.File, error) {
	
	_scriptContents := make ([][]byte, 0, 16)
	
	_scriptContents = append (_scriptContents, []byte ("CREATE DATABASE mysql;"))
	_scriptContents = append (_scriptContents, []byte ("USE mysql;"))
	
	for _, _scriptPath := range _configuration.SqlInitializationScriptPaths {
		// FIXME: Prevent file descriptor leak on error!
		if _scriptFile, _error := os.Open (_scriptPath); _error != nil {
			return nil, _error
		} else if _scriptStat, _error := _scriptFile.Stat (); _error != nil {
			return nil, _error
		} else {
			// FIXME: Enforce a "sane" file size!
			_scriptSize := int (_scriptStat.Size ())
			_scriptReader := bufio.NewReaderSize (_scriptFile, _scriptSize)
			for {
				if _scriptLine, _error := _scriptReader.ReadBytes ('\n'); _error == io.EOF {
					break
				} else if _error != nil {
					return nil, _error
				} else {
					if _scriptLine[len (_scriptLine) - 1] == '\n' {
						_scriptLine = _scriptLine[: len (_scriptLine) - 1]
					}
					_scriptContents = append (_scriptContents, _scriptLine)
				}
			}
			if _error := _scriptFile.Close (); _error != nil {
				panic (_error)
			}
		}
	}
	
	_scriptContents = append (_scriptContents,
			[]byte (fmt.Sprintf (
						`INSERT INTO mysql.user VALUES ('%%','root','','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','','','','',0,0,0,0,'','');`)),
			[]byte (fmt.Sprintf (
						`UPDATE mysql.user SET password = PASSWORD ('%s') WHERE user = 'root';`, _configuration.SqlAdministratorPassword)),
	)
	
	var _reader, _writer *os.File
	if _reader_1, _writer_1, _error := os.Pipe (); _error != nil {
		return nil, _error
	} else {
		_reader = _reader_1
		_writer = _writer_1
	}
	
	go func () () {
		for _, _scriptContent := range _scriptContents {
			// packageTranscript.TraceDebugging ("pushing script chunk: `%s`...", _scriptContent)
			if _, _error := _writer.Write (_scriptContent); _error != nil {
				panic (_error)
			}
			if _, _error := _writer.Write ([]byte ("\n")); _error != nil {
				panic (_error)
			}
		}
		if _error := _writer.Close (); _error != nil {
			panic (_error)
		}
	} ()
	
	return _reader, nil
}


func (_server *server) prepareConsole () (*os.File) {
	
	var _reader, _writer *os.File
	if _reader_1, _writer_1, _error := os.Pipe (); _error != nil {
		panic (_error)
	} else {
		_reader = _reader_1
		_writer = _writer_1
	}
	
	go func () () {
		// FIXME: Handle errors!
		_scanner := bufio.NewScanner (_reader)
		for _scanner.Scan () {
			_server.transcript.TraceInformation (">>  %s", _scanner.Text ())
		}
		_reader.Close ()
	} ()
	
	return _writer
}


func prepareServerExecution (_configuration *ServerConfiguration) (string, []string, []string, string) {
	
	_executable, _arguments, _environment, _directory := prepareGenericExecution (_configuration)
	
	pushStringf (&_arguments, "--bind-address=%s", _configuration.SqlEndpointIp.String ())
	pushStringf (&_arguments, "--port=%d", _configuration.SqlEndpointPort)
	
	pushStrings (&_arguments, "--extra-port=0")
	pushStrings (&_arguments, "--skip-ssl", "--skip-name-resolve", "--skip-host-cache")
	
	return _executable, _arguments, _environment, _directory
}

func prepareBootstrapExecution (_configuration *ServerConfiguration) (string, []string, []string, string) {
	
	_executable, _arguments, _environment, _directory := prepareGenericExecution (_configuration)
	
	pushStrings (&_arguments, "--bootstrap")
	pushStrings (&_arguments, "--skip-grant")
	pushStrings (&_arguments, "--skip-networking")
	pushStrings (&_arguments, "--thread-handling=no-threads")
	
	return _executable, _arguments, _environment, _directory
}

func prepareGenericExecution (_configuration *ServerConfiguration) (string, []string, []string, string) {
	
	_executable := _configuration.GenericConfiguration.ExecutablePath
	_directory := _configuration.GenericConfiguration.TemporaryPath
	_arguments := make ([]string, 0, 128)
	_environment := make ([]string, 0, 128)
	
	if useStrace {
		pushStrings (&_arguments, "/usr/bin/strace")
		pushStrings (&_arguments, "-f", "-v", "-x", "-s", "1024", "-o", "/tmp/strace.txt")
		pushStrings (&_arguments, "--", _executable)
		_executable = "/usr/bin/strace"
	} else {
		pushStrings (&_arguments, _executable)
	}
	
	pushStrings (&_arguments, "--no-defaults")
	
	pushStringf (&_arguments, "--character-sets-dir=%s", _configuration.GenericConfiguration.CharsetsPath)
	pushStringf (&_arguments, "--plugin-dir=%s", _configuration.GenericConfiguration.PluginsPath)
	pushStringf (&_arguments, "--datadir=%s", _configuration.GenericConfiguration.DatabasesPath)
	pushStringf (&_arguments, "--tmpdir=%s", _configuration.GenericConfiguration.TemporaryPath)
	pushStringf (&_arguments, "--socket=%s", _configuration.GenericConfiguration.SocketPath)
	pushStringf (&_arguments, "--pid-file=%s", _configuration.GenericConfiguration.PidPath)
	pushStringf (&_arguments, "--basedir=%s", _configuration.GenericConfiguration.PackagePath)
	
	pushStrings (&_arguments, "--console")
	pushStrings (&_arguments, "--log-warnings")
	
	if useMemlock {
		pushStrings (&_arguments, "--memlock")
	}
	
	if os.Getuid () == 0 {
		pushStrings (&_arguments, "--user=root")
	}
	
	return _executable, _arguments, _environment, _directory
}

func pushStrings (_collection *[]string, _values ... string) () {
	*_collection = append (*_collection, _values ...)
}

func pushStringf (_collection *[]string, _format string, _parts ... interface{}) () {
	pushStrings (_collection, fmt.Sprintf (_format, _parts ...))
}


const usePdeathSignal = false
const useMemlock = false
const useStrace = false
const ignoreBootstrappExitCode = false
const ignoreServerExitCode = true
const consoleFlushTimeout = 4 * time.Second
