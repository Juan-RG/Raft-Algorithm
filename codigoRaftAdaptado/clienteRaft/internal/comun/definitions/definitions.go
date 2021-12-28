package definitions

import (
	"strconv"
)

/******************************* Structs Entrys ************************************/
//Defino un tipo de dato que sera lo que almacene el slice registro			//Todo: Llevar a definitions
//Duda de que esto igual es lo de aplicaOperacion
type Entry struct{
	Index int
	Term int
	Operation TipoOperacion
}
func (entry Entry) String() string {
	string := "Index " + strconv.Itoa(entry.Index) + " Term " +  strconv.Itoa(entry.Term) + " Op " + entry.Operation.String()
	return string
}
/******************************* Fin Structs Entrys ************************************/

/******************************* Structs Votacion ************************************/
type ReplyRequestVote struct {
	Term    int
	Granted bool
}


/************ Structs RPC call votation ******************/
type ArgsRequestVote struct {
	Term int// mandato + 1
	IdNode int// id candidtato
	LastIndexLog int// el ultimo indece de las peticiones comprometidas
	LastTermLog int // mandato de la ultima peticion comprometida
}

/************ fin Structs RPC call votation ******************/
/******************************* Fin Structs Votacion ************************************/


/******************************* Structs ApprendEntries ************************************/
type PairRequestReplay struct {
	Request RequestAP
	Reply ReplyAP
}

//Definimos un tipo de dato para pasar a appendEntries
type RequestAP struct {
	Term int
	LeaderId int
	PrevLogTerm int
	PrevLogIndex int
	LeaderCommit int
	EntryElement Entry
}

func (ra RequestAP) String() string {
	string := "Ter " + strconv.Itoa(ra.Term) + " LId " +  strconv.Itoa(ra.LeaderId) + " PLT " +  strconv.Itoa(ra.PrevLogTerm) +
		" PLIn " +  strconv.Itoa(ra.PrevLogIndex) +"  LComm "+strconv.Itoa(ra.LeaderCommit) + " Entry: " + ra.EntryElement.String()
	return string
}

//Definimos un tipo de dato para recibir el resultado de appendEntries
type ReplyAP struct {
	Term int 		//Para que el lider se actualice a si mismo????
	Success bool	//Para saber si se ha realizado con exito
	IDNode int
	MatchIndex bool
}

/******************************* fin Structs ApprendEntries ************************************/


/******************************* Structs client ************************************/
 /*
 En el codigo de unai lo de debajo
 	|
 	|
 	v
  */
/******************************* fin Structs client ************************************/
/****+*************** Inicio structs unai inicio ****************************************/

type EstadoParcial struct {
	Term	int
	IsLeader bool
	IdLeader	int
}

type EstadoRemoto struct {
	IdNodo	int
	EstadoParcial
}


type ResultadoRemoto struct {
	ValorADevolver string
	IndiceRegistro int
	EstadoParcial
}


type TipoOperacion struct {
	Operacion string  // La operaciones posibles son "leer" y "escribir"
	Clave string
	Valor string    // en el caso de la lectura Valor = ""
}
func (to TipoOperacion) String() string {
	return to.Operacion
}


// A medida que el nodo Raft conoce las operaciones de las  entradas de registro
// comprometidas, envía un AplicaOperacion, con cada una de ellas, al canal
// "canalAplicar" (funcion NuevoNodo) de la maquina de estados
type AplicaOperacion struct {
	Indice int  // en la entrada de registro
	Operacion TipoOperacion
}

type EstadoEntriesRemoto struct {
	Entries []Entry
	CommitIndex int
}
/****+*************** fin structs unai inicio ****************************************/