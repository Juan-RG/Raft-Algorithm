package raft

import (
	"fmt"
	"raft/internal/comun/definitions"
	"time"
)

type stateMachine struct {
	database map[string]string
}

func newStateMachine() *stateMachine {
	machine := &stateMachine{}
	machine.database = make(map[string]string)
	return  machine
}

func (sm *stateMachine) applyOperation(operation definitions.TipoOperacion) {
	if operation.Operacion == "leer" {
		fmt.Println("Leo")
		fmt.Println(sm.database[operation.Clave])
	}else if operation.Operacion == "escribir"{
		fmt.Println("escribo")
		sm.database[operation.Clave] = operation.Valor
	}else {
		fmt.Errorf("Operacion no definida en la maquina de estados")
	}
}

func (nrf *NodoRaft) checkForApplyOperation() {
	for  {
		time.Sleep(500 * time.Millisecond)
		if nrf.commitIndex > nrf.lastApplied {
			nrf.lastApplied++
			nrf.stateMachine.applyOperation(nrf.logEntries[nrf.lastApplied].Operation)
		}
	}
}
