package test

import (
	"fmt"
	"raft/internal/comun/definitions"
	"raft/internal/comun/rpctimeout"
	"raft/pkg/cltraft"
	"strconv"
	"testing"
)

// TEST primer rango
func TestSliceLen(t *testing.T) { // (m *testing.M) {
	var a []int
	a = append(a,1)
	fmt.Println(len(a))
}

func TestSliceLen2(t *testing.T) { // (m *testing.M) {
	var a []int
	a = append(a,[]int{1,4,6,78,9,9}...)
	a = append(a,1)
	fmt.Println(len(a))
	fmt.Println(a)
	a = a[:2]
	fmt.Println(len(a))
	fmt.Println(a)
}

func TestElements(t *testing.T) { // (m *testing.M) {
	request1 := definitions.RequestAP{1,2,4,3,2,definitions.Entry{1,0,definitions.TipoOperacion{}}}
	//entry := definitions.Entry{}

	if (definitions.Entry{}) != request1.EntryElement {
		fmt.Println("hola")
	}
}


func TestDesesperadot (*testing.T) { // (m *testing.M) {
	nodosRaft := rpctimeout.StringArrayToHostPortArray([]string{"localhost:40001", "localhost:40002", "localhost:40003"})
	request := definitions.TipoOperacion{"leer","e1",""}
		request.Clave = request.Clave
		cliente := cltraft.NewClient("localhost", "40005", "leer", nodosRaft)
		cliente.SendOperationToLeader(request)
//	go send3operationsE1("escribir")
	//		go send3operationsE2("escribir")
	//		go send3operationsLe("leer")

	//	time.Sleep(time.Second * 30)
}
func send3operationsLe(op string)  {
	nodosRaft := rpctimeout.StringArrayToHostPortArray([]string{"localhost:40001", "localhost:40002", "localhost:40003"})
	request := definitions.TipoOperacion{op,"e1",""}
	for i := 0; i < 2; i++ {
		request.Clave = request.Clave + strconv.Itoa(i)
		cliente := cltraft.NewClient("localhost", "40005", op, nodosRaft)
		cliente.SendOperationToLeader(request)
	}
	request.Clave = "e2"
	for i := 0; i < 2; i++ {
		request.Clave = request.Clave + strconv.Itoa(i)
		cliente := cltraft.NewClient("localhost", "40005", op, nodosRaft)
		cliente.SendOperationToLeader(request)
	}
}
func send3operationsE1(op string)  {
	nodosRaft := rpctimeout.StringArrayToHostPortArray([]string{"localhost:40001", "localhost:40002", "localhost:40003"})
	request := definitions.TipoOperacion{op,"e1","Escritura1"}
	for i := 0; i < 2; i++ {
		request.Clave = request.Clave + strconv.Itoa(i)
		cliente := cltraft.NewClient("localhost", "40005", op, nodosRaft)
		cliente.SendOperationToLeader(request)
	}
}
func send3operationsE2(op string)  {
	nodosRaft := rpctimeout.StringArrayToHostPortArray([]string{"localhost:40001", "localhost:40002", "localhost:40003"})
	request := definitions.TipoOperacion{op,"e2","Escritura2"}
	for i := 0; i < 2; i++ {
		request.Clave = request.Clave + strconv.Itoa(i)
		cliente := cltraft.NewClient("localhost", "40005", op, nodosRaft)
		cliente.SendOperationToLeader(request)
	}
}


func TestCheckSLice(*testing.T){
	var a []int
	//a := make([]int, 1)
	fmt.Println(a)
	a = append(a,0)
	fmt.Println(a)
}

func TestMAuoria (*testing.T){
  a := 0

  for a < 5{
	  a++
	  fmt.Println("paso")
  }
}