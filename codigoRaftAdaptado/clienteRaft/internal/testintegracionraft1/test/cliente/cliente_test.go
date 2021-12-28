package test

import (
	"raft/internal/comun/definitions"
	"raft/internal/comun/rpctimeout"
	"raft/pkg/cltraft"
	"testing"
)


//155.210.154.194:29250 155.210.154.197:29251 155.210.154.198:29252
func TestCLienteEscritura (*testing.T) { // (m *testing.M) {
	nodosRaft := rpctimeout.StringArrayToHostPortArray([]string{"155.210.154.194:29250", "155.210.154.197:29251", "155.210.154.198:29252"})
	request := definitions.TipoOperacion{"escribir","e1","hola"}
	request.Clave = request.Clave
	cliente := cltraft.NewClient("localhost", "29250", "escribir", nodosRaft)
	cliente.SendOperationToLeader(request)
}
func TestCLienteLectura (*testing.T) { // (m *testing.M) {
	nodosRaft := rpctimeout.StringArrayToHostPortArray([]string{"155.210.154.194:29250", "155.210.154.197:29251", "155.210.154.198:29252"})
	request := definitions.TipoOperacion{"leer","e1","hola"}
		request.Clave = request.Clave
		cliente := cltraft.NewClient("localhost", "29250", "leer", nodosRaft)
		cliente.SendOperationToLeader(request)
}

