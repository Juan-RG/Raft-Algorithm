# Raft Consensus Algorithm

[![Version](https://img.shields.io/badge/version-1.0.0-brightgreen.svg)](https://github.com/user/linda-project/releases)

## Table of Contents

- [Introduction](#introduction)

---

## Use

<!--- //#kind delete cluster
#generar la imagenes locales con los comandos 
 --->
Generate local images:

docker build . -t localhost:5000/raft:latest
docker push localhost:5000/raft:latest
docker build . -t localhost:5000/cliente:latest
docker push localhost:5000/cliente:latest
<!--- 
#borramos el cluster si estaba creado anteriormente kind delete cluster
lanzamos el comando para registrar el cluser con las imagenes locales
--->

command to registry the images in the cluster:

./kind-with-register.sh
<!---  lanzamos el statefulset -->
Run the StatefullSet
kubectl apply -f statefulset_go.yaml

<!--- #para lanzar el cliente una vez lanzado raft --->
Run the client for raft:

kubectl apply -f pods_go.yaml
kubectl delete -f pods_go.yaml -> cada vez que eliminamos y ejecutamos 1 lectura y 1 escritura

<!--- Para lanzar los test lanzamos el comando -->
Run the test:

kubectl exec --stdin --tty c3 -- /bin/sh
<!--- para acceder por ssh y accedemos al pod CUIDADO AL LANZAR LOS test por tiempo de resolucion dns puede tardar demasiado mejor ejecutar de 1 en 1 -->
SSH connection to the Pod.
go test testintegracionraft1_test.go -timeout=30m -v

