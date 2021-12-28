package cltraft

import (
	"clienteRaft/internal/comun/definitions"
	"clienteRaft/internal/comun/rpctimeout"
	"fmt"
	"log"
	"os"
	"time"
)
type Client struct {
	logger *log.Logger // Cada nodo Raft tiene su propio registro de trazas (logs)
	ip string
	port string
	operation string
	nodes []rpctimeout.HostPort
	lastIndexLeader int
}

func NewClient(ip string, port string, operation string, nodes []rpctimeout.HostPort) *Client {
	client := &Client{&log.Logger{}, ip, port, operation, nodes, 0}
	client.logger = log.New(os.Stdout, "cliente: " + ip,
		log.Lmicroseconds|log.Lshortfile)
	client.logger.Println("Client initialized")
	return client
}

func (cl *Client) SendOperationToLeader(op definitions.TipoOperacion){
	cl.connectToLeader(op)
}

func (cl *Client) connectToLeader(op definitions.TipoOperacion) {
/*	func (nrf *NodoRaft) SometerOperacionRaft(operacion definitions.TipoOperacion,
		reply *definitions.ResultadoRemoto) error {
*/
	replyCLient := &definitions.ResultadoRemoto{}
	//requestClient := &definitions.TipoOperacion{"add", "a", "b"}
	//requestClient := &definitions.TipoOperacion{cl.operation, "a", "b"}
	err := cl.nodes[cl.lastIndexLeader].CallTimeout("NodoRaft.SometerOperacionRaft", &op, replyCLient,
		50*time.Millisecond)/*
	err := cl.nodes[cl.lastIndexLeader].CallTimeout("NodoRaft.SometerOperacionRaft", requestClient, replyCLient,
		50*time.Millisecond)*/
	if err != nil {
		cl.updateIndexLeader()
		cl.connectToLeader(op)
	}else{
		if !replyCLient.IsLeader {
			cl.lastIndexLeader = replyCLient.IdLeader
			if cl.lastIndexLeader == -1 {
				cl.lastIndexLeader = 0
			}
			cl.connectToLeader(op)
		}else {
			fmt.Println("Entry commited")
		}
	}
}

func (cl *Client) updateIndexLeader() {
	cl.lastIndexLeader = cl.lastIndexLeader + 1
	fmt.Println(cl.lastIndexLeader)
	if cl.lastIndexLeader == len(cl.nodes) {
		cl.lastIndexLeader = 0
	}
}