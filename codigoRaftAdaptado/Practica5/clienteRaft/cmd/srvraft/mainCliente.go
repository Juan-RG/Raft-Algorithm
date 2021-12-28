package main

import (
	"clienteRaft/cltraft"
	"clienteRaft/internal/comun/definitions"
	"clienteRaft/internal/comun/rpctimeout"
	"os"
)
// ./client operacion clave valor ip puerto nodos

func main() {

	// obtener entero de indice de este nodo
	operacion := os.Args[1]
	clave := os.Args[2]
	valor := os.Args[3]
	ip := os.Args[4]
	puerto := os.Args[5]

	var nodos []string
	// Resto de argumento son los end points como strings
	// De todas la replicas-> pasarlos a HostPort
	for _, endPoint := range os.Args[6:] {
		 nodos = append(nodos, endPoint)
	}
	nodosRaft := rpctimeout.StringArrayToHostPortArray(nodos)
	request := definitions.TipoOperacion{operacion,clave,valor}
	request.Clave = request.Clave
	cliente := cltraft.NewClient(ip, puerto, operacion, nodosRaft)
	cliente.SendOperationToLeader(request)
}
