package testintegracionraft1

import (
	"fmt"
	"os"
	"raft/internal/comun/check"
	"raft/internal/comun/definitions"
	"raft/pkg/cltraft"

	"path/filepath"
	"strconv"
	"testing"
	"time"

	"raft/internal/comun/rpctimeout"
	"raft/internal/despliegue"
	"raft/internal/raft"
)

const (
	//hosts
	/*
	MAQUINA1      = "155.210.154.194"
	MAQUINA2      = "155.210.154.197"
	MAQUINA3      = "155.210.154.198"
	*/
	MAQUINA1      = "localhost"
	MAQUINA2      = "localhost"
	MAQUINA3      = "localhost"
	//puertos
	PUERTOREPLICA1 = "29250"
	PUERTOREPLICA2 = "29251"
	PUERTOREPLICA3 = "29252"

	//nodos replicas
	REPLICA1 = MAQUINA1 + ":" + PUERTOREPLICA1
	REPLICA2 = MAQUINA2 + ":" + PUERTOREPLICA2
	REPLICA3 = MAQUINA3 + ":" + PUERTOREPLICA3

	// paquete main de ejecutables relativos a PATH previo
	EXECREPLICA = "cmd/srvraft/main.go"

	// comandos completo a ejecutar en máquinas remota con ssh. Ejemplo :
	// 				cd $HOME/raft; go run cmd/srvraft/main.go 127.0.0.1:29001

	// Ubicar, en esta constante, nombre de fichero de vuestra clave privada local
	// emparejada con la clave pública en authorized_keys de máquinas remotas

	PRIVKEYFILE = "id_rsa"
)

// PATH de los ejecutables de modulo golang de servicio Raft
var PATH string = filepath.Join(os.Getenv("HOME"), "tmp", "p5", "raft")
//var PATH string = filepath.Join("/home/a805001", "tmp", "p5", "raft")
	// go run cmd/srvraft/main.go 0 127.0.0.1:29001 127.0.0.1:29002 127.0.0.1:29003
var EXECREPLICACMD string = "cd " + PATH + "; go run " + EXECREPLICA



// TEST primer rango
func TestPrimerasPruebas(t *testing.T) { // (m *testing.M) {
	// <setup code>
	// Crear canal de resultados de ejecuciones ssh en maquinas remotas
	cfg := makeCfgDespliegue(t,
							3,
							[]string{REPLICA1, REPLICA2, REPLICA3},
							[]bool{true, true, true})

	// tear down code
	// eliminar procesos en máquinas remotas
	defer cfg.stop()

	// Run test sequence

	// Test1 : No debería haber ningun primario, si SV no ha recibido aún latidos
	t.Run("T1:soloArranqueYparada",
		func(t *testing.T) { cfg.soloArranqueYparadaTest1(t) })

	// Test2 : No debería haber ningun primario, si SV no ha recibido aún latidos
	t.Run("T2:ElegirPrimerLider",
		func(t *testing.T) { cfg.elegirPrimerLiderTest2(t) })

	// Test3: tenemos el primer primario correcto
	t.Run("T3:FalloAnteriorElegirNuevoLider",
		func(t *testing.T) { cfg.falloAnteriorElegirNuevoLiderTest3(t) })

	// Test4: Tres operaciones comprometidas en configuración estable
	t.Run("T4:tresOperacionesComprometidasEstable",
		func(t *testing.T) { cfg.tresOperacionesComprometidasEstable(t) })

	// Test5: Se consigue acuerdo a pesar de desconexiones de seguidor
	t.Run("T5:AcuerdoAPesarDeDesconexionesDeSeguidor ",
		func(t *testing.T) { cfg.AcuerdoApesarDeSeguidor(t) })

	t.Run("T5:SinAcuerdoPorFallos ",
		func(t *testing.T) { cfg.SinAcuerdoPorFallos(t) })

	t.Run("T5:SometerConcurrentementeOperaciones ",
		func(t *testing.T) { cfg.SometerConcurrentementeOperaciones(t) })

}

/*
// TEST primer rango
func TestAcuerdosConFallos(t *testing.T) { // (m *testing.M) {
	// <setup code>
	// Crear canal de resultados de ejecuciones ssh en maquinas remotas
	cfg := makeCfgDespliegue(t,
							3,
							[]string{REPLICA1, REPLICA2, REPLICA3},
							[]bool{true, true, true})

	// tear down code
	// eliminar procesos en máquinas remotas
	defer cfg.stop()


}

*/
// ---------------------------------------------------------------------
// 
// Canal de resultados de ejecución de comandos ssh remotos
type canalResultados chan string

func (cr canalResultados) stop() {
	close(cr)

	// Leer las salidas obtenidos de los comandos ssh ejecutados
	for s := range cr {
		fmt.Println(s)
	}
}

// ---------------------------------------------------------------------
// Operativa en configuracion de despliegue y pruebas asociadas
type configDespliegue struct {
	t *testing.T
	conectados []bool
	numReplicas int
	nodosRaft []rpctimeout.HostPort
	cr canalResultados
}

// Crear una configuracion de despliegue
func makeCfgDespliegue(t *testing.T, n int, nodosraft []string,
										conectados []bool) *configDespliegue {
	cfg := &configDespliegue{}
	cfg.t = t
	cfg.conectados = conectados
	cfg.numReplicas = n
	cfg.nodosRaft = rpctimeout.StringArrayToHostPortArray(nodosraft)
	cfg.cr = make(canalResultados, 2000)
	
	return cfg
}

func (cfg *configDespliegue) stop() {
	cfg.stopDistributedProcesses()

	time.Sleep(50 * time.Millisecond)

	cfg.cr.stop()
}

// --------------------------------------------------------------------------
// FUNCIONES DE SUBTESTS

// Se pone en marcha una replica ?? - 3 NODOS RAFT
func (cfg *configDespliegue) soloArranqueYparadaTest1(t *testing.T) {
	//t.Skip("SKIPPED soloArranqueYparadaTest1")

	fmt.Println(t.Name(), ".....................")

	cfg.t = t  // Actualizar la estructura de datos de tests para errores


	// Poner en marcha replicas en remoto con un tiempo de espera incluido
	cfg.startDistributedProcesses()
	time.Sleep(time.Second * 8)
	// Comprobar estado replica 0
	cfg.comprobarEstadoRemoto (0, 0, false, -1)

	// Comprobar estado replica 1
	cfg.comprobarEstadoRemoto (1, 0, false, -1)

	// Comprobar estado replica 2
	cfg.comprobarEstadoRemoto (2, 0, false, -1)

	// Parar réplicas almacenamiento en remoto
	cfg.stopDistributedProcesses()

	fmt.Println(".............", t.Name(), "Superado")
}


// Primer lider en marcha - 3 NODOS RAFT
func (cfg *configDespliegue) elegirPrimerLiderTest2(t *testing.T) {
	//t.Skip("SKIPPED ElegirPrimerLiderTest2")

	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()
	time.Sleep(time.Second * 8)
	// Se ha elegido lider ?
	fmt.Printf("Probando lider en curso\n")
	cfg.pruebaUnLider(3)


	// Parar réplicas alamcenamiento en remoto
	cfg.stopDistributedProcesses()   // Parametros

	fmt.Println(".............", t.Name(), "Superado")
}

// Fallo de un primer lider y reeleccion de uno nuevo - 3 NODOS RAFT
func (cfg *configDespliegue) falloAnteriorElegirNuevoLiderTest3(t *testing.T) {
	//t.Skip("SKIPPED FalloAnteriorElegirNuevoLiderTest3")

	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()
	time.Sleep(time.Second * 8)
	fmt.Printf("Lider inicial\n")
	time.Sleep(time.Second *5)
	cfg.pruebaUnLider(3)


	// Desconectar lider
	cfg.stopDistributedLeaderProcesses()
	time.Sleep(time.Second *5)
	fmt.Printf("Comprobar nuevo lider\n")
	cfg.pruebaUnLider(3)
	

	// Parar réplicas almacenamiento en remoto
	cfg.stopDistributedProcesses()  //parametros

	fmt.Println(".............", t.Name(), "Superado")
}

// 3 operaciones comprometidas con situacion estable y sin fallos - 3 NODOS RAFT
func (cfg *configDespliegue) tresOperacionesComprometidasEstable(t *testing.T) {
	//t.Skip("SKIPPED tresOperacionesComprometidasEstable")

	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	fmt.Printf("Lider inicial\n")
	time.Sleep(time.Second *8)
	cfg.pruebaUnLider(3)

	cfg.enviarOperacionesCLiente()
	time.Sleep(time.Second * 5)
	fmt.Println(cfg.obtenerEntriesNodo(0))
	fmt.Println(cfg.obtenerEntriesNodo(1))
	fmt.Println(cfg.obtenerEntriesNodo(2))
	cfg.stopDistributedProcesses()  //parametros

	fmt.Println(".............", t.Name(), "Superado")
}

// Se consigue acuerdo a pesar de desconexiones de seguidor -- 3 NODOS RAFT
func(cfg *configDespliegue) AcuerdoApesarDeSeguidor(t *testing.T) {

	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	fmt.Printf("Lider inicial\n")
	time.Sleep(time.Second * 8)
	cfg.pruebaUnLider(3)

	// Comprometer una entrada
	cfg.enviarOperacionCLiente()
	time.Sleep(time.Second * 4)
	cfg.stopNumberOfDistributedProcesses(1)
	time.Sleep(time.Second * 6)
	//  Obtener un lider y, a continuación desconectar una de los nodos Raft

	_, e1Term,_,_ := cfg.obtenerEstadoRemotoWithNodesOff(0)
	_, e2Term,_,_ := cfg.obtenerEstadoRemotoWithNodesOff(1)
	_, e3Term,_,_ := cfg.obtenerEstadoRemotoWithNodesOff(2)
	fmt.Println(cfg.obtenerEstadoRemotoWithNodesOff(0))
	fmt.Println(cfg.obtenerEstadoRemotoWithNodesOff(1))
	fmt.Println(cfg.obtenerEstadoRemotoWithNodesOff(2))

	count := 0
	if e1Term != -1 {
		count++
	}
	if e2Term != -1 {
		count++
	}
	if e3Term != -1 {
		count++
	}
	if count > 2 {
		t.Errorf("Hay mas de 2 nodos arrancados")
	}

	// Comprobar varios acuerdos con una réplica desconectada
	cfg.enviarOperacionCLiente()
	fmt.Println(cfg.obtenerEntriesNodoWithNodeOut(0))
	fmt.Println(cfg.obtenerEntriesNodoWithNodeOut(1))
	fmt.Println(cfg.obtenerEntriesNodoWithNodeOut(2))
	cfg.enviarOperacionCLiente()
	fmt.Println(cfg.obtenerEntriesNodoWithNodeOut(0))
	fmt.Println(cfg.obtenerEntriesNodoWithNodeOut(1))
	fmt.Println(cfg.obtenerEntriesNodoWithNodeOut(2))
	cfg.enviarOperacionCLiente()
	fmt.Println(cfg.obtenerEntriesNodoWithNodeOut(0))
	fmt.Println(cfg.obtenerEntriesNodoWithNodeOut(1))
	fmt.Println(cfg.obtenerEntriesNodoWithNodeOut(2))

	// reconectar nodo Raft previamente desconectado y comprobar varios acuerdos
	cfg.startDistributedProcesses()//Lanzo todos los nodos de nuevo
	time.Sleep(time.Second * 10)

	fmt.Println(cfg.obtenerEntriesNodo(0))
	fmt.Println(cfg.obtenerEntriesNodo(1))
	fmt.Println(cfg.obtenerEntriesNodo(2))

	entrysNodo1, commitI1 := cfg.obtenerEntriesNodo(0)
	entrysNodo2, commitI2 := cfg.obtenerEntriesNodo(1)
	entrysNodo3, commitI3 := cfg.obtenerEntriesNodo(2)

	if commitI1 != commitI2 || commitI1 != commitI3 || commitI2 != commitI3 {
		t.Errorf("Indice de comiteadas diferente")
	}
	if len(entrysNodo1) != len(entrysNodo2) || len(entrysNodo1) != len(entrysNodo3) || len(entrysNodo2) != len(entrysNodo3){
		t.Errorf("Replycacion de nodos diferente")
	}else {
		for i,element := range entrysNodo1 {
			if element != entrysNodo2[i] || element != entrysNodo3[i] || entrysNodo2[i] != entrysNodo3[i]  {
				t.Errorf("Error Logs diferentes ")
			}
		}
	}

	cfg.stopDistributedProcesses()
}

// NO se consigue acuerdo al desconectarse mayoría de seguidores -- 3 NODOS RAFT
func(cfg *configDespliegue) SinAcuerdoPorFallos(t *testing.T) {
	//t.Skip("SKIPPED SinAcuerdoPorFallos")
	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	fmt.Printf("Lider inicial\n")
	time.Sleep(time.Second *8)
	cfg.pruebaUnLider(3)
	e1,_,e1Leader,_ := cfg.obtenerEstadoRemoto(0)
	e2,_,e2Leader,_ := cfg.obtenerEstadoRemoto(1)
	e3,_,e3Leader,_ := cfg.obtenerEstadoRemoto(2)

	if !e1Leader && !e2Leader && !e3Leader {
		t.Errorf("NO hay lider")
	}else {
		fmt.Println("Estado actual")
		fmt.Println("nodo ", e1, " isLeader: ", e1Leader)
		fmt.Println("nodo ", e2, " isLeader: ", e2Leader)
		fmt.Println("nodo ", e3, " isLeader: ", e3Leader)
	}
	// Comprometer una entrada
	cfg.enviarOperacionCLiente()
	time.Sleep(time.Second * 2)
	//  Obtener un lider y, a continuación desconectar 2 de los nodos Raft
	fmt.Println(cfg.obtenerEntriesNodo(0))
	fmt.Println(cfg.obtenerEntriesNodo(1))
	fmt.Println(cfg.obtenerEntriesNodo(2))

	cfg.stopNumberOfDistributedProcessesNoleader()
	time.Sleep(time.Second * 2)
	var e1Term, e2Term, e3Term  int
	e1, e1Term,e1Leader,_ = cfg.obtenerEstadoRemotoWithNodesOff(0)
	e2, e2Term,e2Leader,_ = cfg.obtenerEstadoRemotoWithNodesOff(1)
	e3, e3Term,e3Leader,_ = cfg.obtenerEstadoRemotoWithNodesOff(2)
	fmt.Println(cfg.obtenerEstadoRemotoWithNodesOff(0))
	fmt.Println(cfg.obtenerEstadoRemotoWithNodesOff(1))
	fmt.Println(cfg.obtenerEstadoRemotoWithNodesOff(2))

	count := 0
	leaderID := -1
	if e1Term != -1 {
		count++
		leaderID = e1
	}
	if e2Term != -1 {
		count++
		leaderID = e2
	}
	if e3Term != -1 {
		count++
		leaderID = e3
	}
	fmt.Println(count)
	if count > 1 {
		 t.Errorf("Hay mas de un nodo arrancado")
	}
	// Comprobar varios acuerdos con 2 réplicas desconectada
	// Comprometer una entrada
	cfg.enviarOperacionCLiente()
	fmt.Println(cfg.obtenerEntriesNodo(leaderID))
	// Comprometer una entrada
	cfg.enviarOperacionCLiente()
	fmt.Println(cfg.obtenerEntriesNodo(leaderID))
	// Comprometer una entrada
	cfg.enviarOperacionCLiente()
	fmt.Println(cfg.obtenerEntriesNodo(leaderID))
	time.Sleep(time.Second * 2)
	fmt.Println(cfg.obtenerEntriesNodo(leaderID))
	// reconectar lo2 nodos Raft  desconectados y probar varios acuerdos
	cfg.startDistributedProcesses()
	time.Sleep(time.Second * 10)

	fmt.Println(cfg.obtenerEntriesNodo(0))
	fmt.Println(cfg.obtenerEntriesNodo(1))
	fmt.Println(cfg.obtenerEntriesNodo(2))

	entrysNodo1, commitI1 := cfg.obtenerEntriesNodo(0)
	entrysNodo2, commitI2 := cfg.obtenerEntriesNodo(1)
	entrysNodo3, commitI3 := cfg.obtenerEntriesNodo(2)

	if commitI1 != commitI2 || commitI1 != commitI3 || commitI2 != commitI3 {
		t.Errorf("Indice de comiteadas diferente")
	}
	if len(entrysNodo1) != len(entrysNodo2) || len(entrysNodo1) != len(entrysNodo3) || len(entrysNodo2) != len(entrysNodo3){
		t.Errorf("Replycacion de nodos diferente")
	}else {
		for i,element := range entrysNodo1 {
			if element != entrysNodo2[i] || element != entrysNodo3[i] || entrysNodo2[i] != entrysNodo3[i]  {
				t.Errorf("Error Logs diferentes ")
			}
		}
	}
	cfg.stopDistributedProcesses()
}

// Se somete 5 operaciones de forma concurrente -- 3 NODOS RAFT
func(cfg *configDespliegue) SometerConcurrentementeOperaciones(t *testing.T) {
//	t.Skip("SKIPPED SometerConcurrentementeOperaciones")

	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	fmt.Printf("Esperamos Lider inicial\n")
	time.Sleep(time.Second * 8)
	cfg.pruebaUnLider(3)
	e1,_,e1Leader,_ := cfg.obtenerEstadoRemoto(0)
	e2,_,e2Leader,_ := cfg.obtenerEstadoRemoto(1)
	e3,_,e3Leader,_ := cfg.obtenerEstadoRemoto(2)

	if !e1Leader && !e2Leader && !e3Leader {
		t.Errorf("NO hay lider")
	}else {
		fmt.Println("Estado actual")
		fmt.Println("nodo ", e1, " isLeader: ", e1Leader)
		fmt.Println("nodo ", e2, " isLeader: ", e2Leader)
		fmt.Println("nodo ", e3, " isLeader: ", e3Leader)
	}
	// Obtener un lider y, a continuación someter una operacion
	fmt.Println("EStado previo del log")
	fmt.Println(cfg.obtenerEntriesNodo(0))
	fmt.Println(cfg.obtenerEntriesNodo(1))
	fmt.Println(cfg.obtenerEntriesNodo(2))
	cfg.SometerConcurrentementeOperacionesCliente()
	time.Sleep(time.Second * 4)
	// Comprobar estados de nodos Raft, sobre todo
	// el avance del mandato en curso e indice de registro de cada uno
	// que debe ser identico entre ellos
	fmt.Println(cfg.obtenerEntriesNodo(0))
	fmt.Println(cfg.obtenerEntriesNodo(1))
	fmt.Println(cfg.obtenerEntriesNodo(2))
	entrysNodo1, commitI1 := cfg.obtenerEntriesNodo(0)
	entrysNodo2, commitI2 := cfg.obtenerEntriesNodo(1)
	entrysNodo3, commitI3 := cfg.obtenerEntriesNodo(2)

	if commitI1 != commitI2 || commitI1 != commitI3 || commitI2 != commitI3 {
		t.Errorf("Indice de comiteadas diferente")
	}
	if len(entrysNodo1) != len(entrysNodo2) || len(entrysNodo1) != len(entrysNodo3) || len(entrysNodo2) != len(entrysNodo3){
		t.Errorf("Replycacion de nodos diferente")
	}else {
		for i,element := range entrysNodo1 {
			if element != entrysNodo2[i] || element != entrysNodo3[i] || entrysNodo2[i] != entrysNodo3[i]  {
				t.Errorf("Error Logs diferentes ")
			}
		}
	}
}




// --------------------------------------------------------------------------
// FUNCIONES DE APOYO
// Comprobar que hay un solo lider
// probar varias veces si se necesitan reelecciones
func (cfg *configDespliegue) pruebaUnLider(numreplicas int) int {
	for iters := 0; iters < 10; iters++ {
		time.Sleep(500 * time.Millisecond)
		mapaLideres := make(map[int][]int)
		for i := 0; i < numreplicas; i++ {
			if cfg.conectados[i] {
				if _, mandato, eslider, _ := cfg.obtenerEstadoRemoto(i);
																	  eslider {
																		  fmt.Println("nodo: ", i, " mandato: ", mandato, "es lider")
					mapaLideres[mandato] = append(mapaLideres[mandato], i)
				}
			}
		}

		ultimoMandatoConLider := -1
		for mandato, lideres := range mapaLideres {
			if len(lideres) > 1 {
				cfg.t.Fatalf("mandato %d tiene %d (>1) lideres",
														mandato, len(lideres))
			}
			if mandato > ultimoMandatoConLider {
				ultimoMandatoConLider = mandato
			}
		}

		if len(mapaLideres) != 0 {
			
			return mapaLideres[ultimoMandatoConLider][0]  // Termina
			
		}
	}
	cfg.t.Fatalf("un lider esperado, ninguno obtenido")
	
	return -1   // Termina
}

func (cfg *configDespliegue) obtenerEstadoRemoto(
										indiceNodo int) (int, int, bool, int) {
	var reply definitions.EstadoRemoto

	err := cfg.nodosRaft[indiceNodo].CallTimeout("NodoRaft.ObtenerEstadoNodo",
								raft.Vacio{}, &reply, 600 * time.Millisecond)

	check.CheckError(err, "Error en llamada RPC ObtenerEstadoRemoto")

	return reply.IdNodo, reply.Term, reply.IsLeader, reply.IdLeader
}
func (cfg *configDespliegue) obtenerEstadoRemotoWithNodesOff(
	indiceNodo int) (int, int, bool, int) {
	var reply definitions.EstadoRemoto

	err := cfg.nodosRaft[indiceNodo].CallTimeout("NodoRaft.ObtenerEstadoNodo",
		raft.Vacio{}, &reply, 600 * time.Millisecond)

	if err != nil {
		reply.Term = -1
		reply.IdLeader = -1
		reply.IsLeader = false
	}
	//check.CheckError(err, "Error en llamada RPC ObtenerEstadoRemoto")

	return reply.IdNodo, reply.Term, reply.IsLeader, reply.IdLeader
}
func (cfg *configDespliegue) obtenerEntriesNodo(
	indiceNodo int) ([]definitions.Entry, int) {
	var reply definitions.EstadoEntriesRemoto

	err := cfg.nodosRaft[indiceNodo].CallTimeout("NodoRaft.ObtenerEntriesNodo",
		raft.Vacio{}, &reply, 600 * time.Millisecond)

	check.CheckError(err, "Error en llamada RPC ObtenerEntriesNodo")

	return reply.Entries, reply.CommitIndex
}

func (cfg *configDespliegue) obtenerEntriesNodoWithNodeOut(
	indiceNodo int) ([]definitions.Entry, int) {
	var reply definitions.EstadoEntriesRemoto

	err := cfg.nodosRaft[indiceNodo].CallTimeout("NodoRaft.ObtenerEntriesNodo",
		raft.Vacio{}, &reply, 600 * time.Millisecond)

	if err != nil {
		reply.Entries = nil
		reply.CommitIndex = -1
	}
	//check.CheckError(err, "Error en llamada RPC ObtenerEntriesNodo")

	return reply.Entries, reply.CommitIndex
}

// start  gestor de vistas; mapa de replicas y maquinas donde ubicarlos;
// y lista clientes (host:puerto)
func (cfg *configDespliegue) startDistributedProcesses() {
 //cfg.t.Log("Before starting following distributed processes: ", cfg.nodosRaft)

	for i, endPoint := range cfg.nodosRaft {
		despliegue.ExecMutipleHosts( EXECREPLICACMD +
								" " + strconv.Itoa(i) + " " +
								rpctimeout.HostPortArrayToString(cfg.nodosRaft),
								[]string{endPoint.Host()}, cfg.cr, PRIVKEYFILE)

		// dar tiempo para se establezcan las replicas
		//time.Sleep(1000 * time.Millisecond)
	}

	// aproximadamente 500 ms para cada arranque por ssh en portatil
	time.Sleep(2500 * time.Millisecond)
	//
	//
}

//
func (cfg *configDespliegue) stopDistributedProcesses() {
	var reply raft.Vacio

	for _, endPoint := range cfg.nodosRaft {
		err := endPoint.CallTimeout("NodoRaft.ParaNodo",
								raft.Vacio{}, &reply, 600 * time.Millisecond)
		check.CheckError(err, "Error en llamada RPC Para nodo")
	}
}
func (cfg *configDespliegue) stopNumberOfDistributedProcesses(num int) {
	var reply raft.Vacio

	for _, endPoint := range cfg.nodosRaft {
		if num == 0 {
			return
		}
		err := endPoint.CallTimeout("NodoRaft.ParaNodo",
			raft.Vacio{}, &reply, 600 * time.Millisecond)
		check.CheckError(err, "Error en llamada RPC Para nodo")
		num--
	}
}

//
func (cfg *configDespliegue) stopDistributedLeaderProcesses() {
	var reply raft.Vacio
	for i := 0; i < len(cfg.nodosRaft); i++ {
		if cfg.conectados[i] {
			if _, _, eslider, _ := cfg.obtenerEstadoRemoto(i);
				eslider {
					//Paramos el lider
					err := cfg.nodosRaft[i].CallTimeout("NodoRaft.ParaNodo",
					raft.Vacio{}, &reply, 600 * time.Millisecond)
				check.CheckError(err, "Error en llamada RPC Para nodo")
					//reactivamos el lider
					time.Sleep(time.Millisecond * 250)
				despliegue.ExecMutipleHosts( EXECREPLICACMD +
					" " + strconv.Itoa(i) + " " +
					rpctimeout.HostPortArrayToString(cfg.nodosRaft),
					[]string{cfg.nodosRaft[i].Host()}, cfg.cr, PRIVKEYFILE)
					time.Sleep(time.Millisecond * 250)
			}
		}
	}

}

// Comprobar estado remoto de un nodo con respecto a un estado prefijado
func (cfg *configDespliegue) comprobarEstadoRemoto(idNodoDeseado int,
				 mandatoDeseado int, esLiderDeseado bool, IdLiderDeseado int) {
	idNodo, mandato, esLider, idLider := cfg.obtenerEstadoRemoto(idNodoDeseado)

	cfg.t.Log("Estado replica 0: ", idNodo, mandato, esLider, idLider, "\n")
	/*
		if idNodo != idNodoDeseado || mandato != mandatoDeseado || esLider != esLiderDeseado || idLider != IdLiderDeseado {
	  cfg.t.Fatalf("Estado incorrecto en replica %d en subtest %s",
													idNodoDeseado, cfg.t.Name())
	}

	 */

}

func (cfg *configDespliegue) enviarOperacionesCLiente() {
	for i := 0 ; i < 3 ; i++ {
		op := definitions.TipoOperacion{}
		op.Operacion = "Prueba"
		cliente := cltraft.NewClient("localhost","29250","op", cfg.nodosRaft)
		cliente.SendOperationToLeader(op)
	}
}

func (cfg *configDespliegue) enviarOperacionCLiente() {
		op := definitions.TipoOperacion{}
		op.Operacion = "Prueba"
		cliente := cltraft.NewClient("localhost","29250","op", cfg.nodosRaft)
		cliente.SendOperationToLeader(op)
}

func (cfg *configDespliegue) stopNumberOfDistributedProcessesNoleader() {
	var reply raft.Vacio
	for i := 0; i < len(cfg.nodosRaft); i++ {
		if cfg.conectados[i] {
			if _, _, eslider, _ := cfg.obtenerEstadoRemoto(i);
				!eslider {
					fmt.Println("es lider: ",eslider)
				//Paramos el lider
				err := cfg.nodosRaft[i].CallTimeout("NodoRaft.ParaNodo",
					raft.Vacio{}, &reply, 600 * time.Millisecond)
				check.CheckError(err, "Error en llamada RPC Para nodo")
			}
		}
	}
}

func (cfg *configDespliegue) SometerConcurrentementeOperacionesCliente() {
	for i:=0 ; i < 5; i++ {
		go cfg.enviarOperacionCLiente()
	}
}
