package main

import (
	"clienteRaft/cltraft"
	"clienteRaft/internal/comun/definitions"
	"clienteRaft/internal/comun/rpctimeout"
	"os"
	"strings"
)
// ./client operacion clave valor ip puerto nodos

func main() {
	var strSlice = []string{"r-0.raft.default.svc.cluster.local:6000", "r-1.raft.default.svc.cluster.local:6000",
		"r-2.raft.default.svc.cluster.local:6000", "r-3.raft.default.svc.cluster.local:6000"}
	// obtener entero de indice de este nodo
/*	var strSlice = []string{"localhost:6000", "localhost:6001",
		"localhost:6002", "localhost:6003"}*/
	conexion := strings.Split(os.Args[1], ":")
	ip := conexion[0]
	puerto := conexion[1]
	operacion := os.Args[2]
	clave := os.Args[3]
	valor := os.Args[4]


	var nodos []string
	// Resto de argumento son los end points como strings
	// De todas la replicas-> pasarlos a HostPort
	for _, endPoint := range strSlice {
		 nodos = append(nodos, endPoint)
	}
	nodosRaft := rpctimeout.StringArrayToHostPortArray(nodos)
	request := definitions.TipoOperacion{operacion,clave,valor}
	request.Clave = request.Clave
	cliente := cltraft.NewClient(ip, puerto, operacion, nodosRaft)
	cliente.SendOperationToLeader(request)
}
