package raft

import (
	"fmt"
	"raft/internal/comun/definitions"
	"raft/internal/comun/rpctimeout"
	"time"
)



//Funcion para llamar con RPC a appendEntries de cada uno de los nodos
func (nrf *NodoRaft) appendEntries() { //mriar si puedo meter los param de la funcion de una


	chReply := make(chan definitions.PairRequestReplay, len(nrf.Nodos)) // canal de request, reply
	request := definitions.RequestAP{}

	nrf.Mux.Lock()

	//Inicializamos los nodos
	request.Term = nrf.currentTerm
	request.LeaderId = nrf.Yo
	request.LeaderCommit = nrf.commitIndex
	request.PrevLogTerm = 0
	request.PrevLogIndex = 0

	for i, v := range nrf.Nodos {
		if i != nrf.Yo {
			nIndexOfNode := nrf.nextIndex[i] //listado de indices de insercion de los nodos
			if len(nrf.logEntries) != 0 && nIndexOfNode > 0 { // Si tenemos entries y el indice no es un caso corner sacamos termino y indice
				request.PrevLogTerm = nrf.logEntries[nIndexOfNode - 1].Term
				request.PrevLogIndex = nIndexOfNode - 1
			}else {
				//Si no indicamos que case corner con 0 y -1
				request.PrevLogTerm = 0
				request.PrevLogIndex = nIndexOfNode - 1
			}
			//Si hay Entry disponible para ese nodo le envio la entri
			if nIndexOfNode < (len(nrf.logEntries)) {//Si el indice es diferente a los entries actuales enviamos el entry
				request.EntryElement = nrf.logEntries[nIndexOfNode]
			}else{
				//si no heartbeat
				request.EntryElement = definitions.Entry{}
			}

			go nrf.connectRPCAppendEntries(chReply, v, request)
		}
	}
	nrf.Mux.Unlock()


	go nrf.waitAnswers(chReply)
	nrf.Mux.Lock()
	lastIndex := len(nrf.logEntries) - 1

	if lastIndex == -1  || nrf.hopeCommit[lastIndex] != 0 {
		nrf.Mux.Unlock()
		return
	}
	nrf.hopeCommit[lastIndex] = 1
	nrf.Mux.Unlock()
	go nrf.commitOperation(lastIndex)
}

func (nrf *NodoRaft) connectRPCAppendEntries(chReply chan definitions.PairRequestReplay, node rpctimeout.HostPort, request definitions.RequestAP) {
	timeAppendEntries := 300
	reply := definitions.ReplyAP{0, false,0, false}
		err := node.CallTimeout("NodoRaft.SometerOperacion", &request, &reply,
			time.Duration(timeAppendEntries)*time.Millisecond)
		if err != nil {
			nrf.Logger.Println("Error SometerOperacion", err)
			reply.Term = -1 // Coloco -1 para controlar el error
			reply.Success = false
			reply.IDNode = -1
		}
		//genero la pareja de [request, replay]
	pairRequestReply := definitions.PairRequestReplay{request, reply}
	chReply <- pairRequestReply
}

func (nrf *NodoRaft) waitAnswers(chReply chan definitions.PairRequestReplay) {
		for i := 0; i < len(nrf.Nodos) - 1; i++ { // espero la respuesta de los nodos
			select {
				case pairRequestReplay := <-chReply:
					if pairRequestReplay.Reply.IDNode == -1 && pairRequestReplay.Reply.Term == -1 {
						//error conexion no hacemos nada ya se reenviara en el siguiente intento
					}else{
						nrf.checkCurrentTerm(pairRequestReplay.Reply.Term) //comprobamos el termino
						if nrf.state == FOLLOWER {
							return
						}
						nrf.Mux.Lock() //actualizamos los valores con mutex
						//tenemos que estar en matchIndex y succes
						if pairRequestReplay.Reply.Success && pairRequestReplay.Reply.MatchIndex && (pairRequestReplay.Request.EntryElement !=
							(definitions.Entry{})) {  //si hemos tenido, aumentamos contador de mayoria
							nrf.nextIndex[pairRequestReplay.Reply.IDNode]++
							nrf.entriesVoted[pairRequestReplay.Request.PrevLogIndex + 1]++ //aumentamos el voto en su indice previo + 1
							//si no tenemos success y no hay maxindex disminuimos el indice hasta encontrar match
						}else if !pairRequestReplay.Reply.Success && !pairRequestReplay.Reply.MatchIndex {
							nrf.nextIndex[pairRequestReplay.Reply.IDNode]--
						}
						nrf.Mux.Unlock()
					}
			}
		}
}

func (nrf *NodoRaft) commitOperation(lastIndex int) {
	if lastIndex != -1 { //Comrpuebo si hay que commitear
		if lastIndex > nrf.commitIndex {

		votos := 0
		for votos <= len(nrf.Nodos)/2 { //espera activa a tener mayoria //Todo: MIrar en un futuro de hacer esto con gorutines
			nrf.Mux.Lock()
			votos = nrf.entriesVoted[lastIndex]
			nrf.Mux.Unlock()
		}
		<-nrf.previousCommit[lastIndex] //espero a la anterior
		nrf.previousCommit[lastIndex]<-true
		//espero a la anterior

		nrf.Mux.Lock()
		fmt.Println("Entro index ", lastIndex)
		nrf.commitIndex = lastIndex //actualizo el indice y aviso a la siguiente a commitear
		nrf.previousCommit[lastIndex + 1] <- true
		nrf.Mux.Unlock()
		}
	}
}

func (nrf *NodoRaft) SometerOperacion(op *definitions.RequestAP, reply *definitions.ReplyAP) error {
	reply.Term = nrf.currentTerm
	reply.Success = false
	reply.IDNode = nrf.Yo
	nrf.Mux.Lock()
	defer nrf.Mux.Unlock()
	//si el termino es menor no lo reconocemos como lider
	if op.Term < nrf.currentTerm { //Yo tengo mandato superior
		return nil
	}
	nrf.IdLider = op.LeaderId //si es igual o seperior a nosotros lo colocamos como lider y actualizamos
	nrf.hearbeatChannel <- 1 //Avisamos de que hemos recibido una peticion de un nodo con termino igual o superoir a nosotros
	nrf.state = FOLLOWER
	nrf.updateCommitIndex(op.LeaderCommit)// actualizo las entries commiteadas
	if op.Term > nrf.currentTerm {
		nrf.currentTerm = op.Term
		nrf.votedFor = -1
	}
	if (op.PrevLogIndex == -1 && op.EntryElement != (definitions.Entry{})) && len(nrf.logEntries) == 0 { // caso inicial
		reply.MatchIndex = true
		nrf.logEntries = append(nrf.logEntries, op.EntryElement)
		reply.Success = true
		return nil
	}else {
		if op.PrevLogIndex == -1 && op.EntryElement == (definitions.Entry{}) && len(nrf.logEntries) == 0{ // Primeros hearbeat
			reply.MatchIndex = true
			return nil
		}
		//si no tengo datos el prevlogIndex tiene que ser -1 caso corner
		if op.PrevLogIndex == 0 && len(nrf.logEntries) == 0 { //si no tengo elemento digo que estoy vacio
			reply.MatchIndex = false
			return nil
		}
		//si no tengo un dato pero el prevlogindex es 0
		if len(nrf.logEntries) - 1 < op.PrevLogIndex {
			//fmt.Println(" 5")
			//si ambos son 0 esta bien colocado
			if len(nrf.logEntries) == 0 && op.PrevLogIndex == 0 {	//Caso corner
				//fmt.Println(" 5.1")
				reply.MatchIndex = true
			//si no no estan bien colocados
			}else{
				//fmt.Println(" 5.2")
				reply.MatchIndex = false  //Indice incorrecto
			}
			return nil
		}else{ //si el prevlogIndex se puede indexar comparamos terminos
			if len(nrf.logEntries) - 1 >= op.PrevLogIndex + 1  {
				reply.MatchIndex = true
				return nil
			}
			termOfPrevLogIndex := nrf.logEntries[op.PrevLogIndex].Term
			//si el termino es diferente eliminamos hasta esa entrada
			if termOfPrevLogIndex != op.PrevLogTerm {
				nrf.logEntries = nrf.logEntries[:op.PrevLogIndex] //borro las que no coinciden
				reply.MatchIndex = false
				return nil
			}else{ //Si son iguales añadimos el nuevo registro
				if op.EntryElement != (definitions.Entry{}) {
					nrf.logEntries = append(nrf.logEntries, op.EntryElement)
					//Todo: añadir el commit del leader
					reply.Success = true
					reply.MatchIndex = true
					return nil
				}else { //si es hearbeat indicamos que el indice esta bien
					reply.MatchIndex = true
					return nil
				}
			}

		}
	}
	return nil
}

func (nrf *NodoRaft) updateCommitIndex(leaderCommit int) {
	//actualiza las commiteadas con el menor entre las commiteadas del lider o el tamaño del vector
	if leaderCommit > nrf.commitIndex {
		nrf.commitIndex = min(leaderCommit, len(nrf.logEntries) - 1)
	}
}

func min(a int, b int) int{
	if a < b {
		return a
	} else {
		return b
	}
}









