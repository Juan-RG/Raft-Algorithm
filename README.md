#kind delete cluster
#generar la imagenes locales con los comandos 
docker build . -t localhost:5000/raft:latest
docker push localhost:5000/raft:latest

docker build . -t localhost:5000/cliente:latest
docker push localhost:5000/cliente:latest
#borramos el cluster si estaba creado anteriormente kind delete cluster
lanzamos el comando para registrar el cluser con las imagenes locales
#./kind-with-register.sh
lanzamos el statefulset
#kubectl apply -f statefulset_go.yaml

#para lanzar el cliente una vez lanzado raft 
#kubectl apply -f pods_go.yaml
#kubectl delete -f pods_go.yaml -> cada vez que eliminamos y ejecutamos 1 lectura y 1 escritura

Para lanzar los test lanzamos el comando
kubectl exec --stdin --tty c3 -- /bin/sh
para acceder por ssh y accedemos al pod CUIDADO AL LANZAR LOS test por tiempo de resolucion dns puede tardar demasiado mejor ejecutar de 1 en 1
go test testintegracionraft1_test.go -timeout=30m -v

