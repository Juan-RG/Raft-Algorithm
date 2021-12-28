package raft

import (
	"os"
	"raft/internal/comun/definitions"
	"time"
)

/**
* Metodo que devuelve el ultimo termino del array de entries
* Si slice vacio return -1
* si no return term de la ultima entrada
 */
func (nrf *NodoRaft) getTermLastEnty() int {
	var termLastCommitEntry int
	if len(nrf.logEntries) == 0 {
		termLastCommitEntry = -1
	} else {
		termLastCommitEntry = nrf.logEntries[len(nrf.logEntries) - 1].Term
	}
	return termLastCommitEntry
}

/*
 Funcion para añadir operaciones al entri
 */
func (nrf *NodoRaft) addEntryToLogEntries(operation definitions.TipoOperacion) {
	newEntry := definitions.Entry{len(nrf.logEntries),nrf.currentTerm, operation}
	nrf.Mux.Lock()
	nrf.logEntries = append(nrf.logEntries, newEntry)
	nrf.hopeCommit = append(nrf.hopeCommit,0)
//	nrf.prevLogIndex = len(nrf.nextIndex) - 1
	//Añadimos otro canal para la espera del commit (uno por entrada)
	nrf.entriesVoted = append(nrf.entriesVoted, 1)
	chCommit := make(chan bool, 1)
	nrf.previousCommit = append(nrf.previousCommit, chCommit)
	nrf.Mux.Unlock()


}


/*
	Metodo para actualizar el mandato
*/
func (nrf *NodoRaft) checkCurrentTerm(term int) {
	if nrf.currentTerm < term {
		nrf.currentTerm = term
		nrf.state = FOLLOWER
	}
}


/*
	Funcion para saber si es lider
 */
func (nrf *NodoRaft) isLeader() bool {
	if nrf.state == LEADER {
		return true
	}
	return false
}

// Metodo Para() utilizado cuando no se necesita mas al nodo
//
// Quizas interesante desactivar la salida de depuracion
// de este nodo
//
func (nrf *NodoRaft) para() {
	go func() {time.Sleep(5 * time.Millisecond); os.Exit(0) } ()
}

// Devuelve "yo", mandato en curso y si este nodo cree ser lider
//
// Primer valor devuelto es el indice de este  nodo Raft el el conjunto de nodos
// la operacion si consigue comprometerse.
// El segundo valor es el mandato en curso
// El tercer valor es true si el nodo cree ser el lider
// Cuarto valor es el lider, es el indice del líder si no es él
func (nrf *NodoRaft) obtenerEstado() (int, int, bool, int) {
	var yo int = nrf.Yo
	var mandato int
	var esLider bool
	var idLider int = nrf.IdLider

	if nrf.state == LEADER {
		esLider = true
		mandato = nrf.currentTerm
	}
	return yo, mandato, esLider, idLider
}
