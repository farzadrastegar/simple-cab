.PHONY: docker-build docker-swarm-init docker-deploy-services

docker-swarm-init:
	docker swarm init

docker-swarm-leave:
	docker swarm leave --force

docker-build:
	docker build --build-arg moduleName=gateway -t farzadras/gateway -f ./support/docker/simple-cab/Dockerfile-moduleName .
	docker build --build-arg moduleName=zombie_driver -t farzadras/zombie_driver -f ./support/docker/simple-cab/Dockerfile-moduleName .
	docker build --build-arg moduleName=driver_location -t farzadras/driver_location -f ./support/docker/simple-cab/Dockerfile-moduleName .
	docker rmi $$(docker images -q -f dangling=true)

docker-deploy-services:
	docker stack deploy --compose-file ./support/docker/services/swarm/docker-compose-s1.yaml services
	docker stack deploy --compose-file ./support/docker/services/swarm/docker-compose-s2.yaml services
	docker stack deploy --compose-file ./support/docker/services/swarm/docker-compose-s3.yaml services
	docker stack deploy --compose-file ./support/docker/services/swarm/docker-compose-s4.yaml services

docker-rm-services:
	docker stack rm services
	docker rm -f $$(docker ps -a -q)
	docker rmi $$(docker images -q -f dangling=true)

docker-build-rabbitmq:
	docker build -t farzadras/rabbitmq -f ./support/docker/rabbitmq/Dockerfile .

docker-build-configserver:
	cd support/config-server/; ./run.bash; cd -
	docker build -t farzadras/configserver support/config-server/
	docker rmi $$(docker images -q -f dangling=true)

deploy-local-services:
	nohup docker-compose -f support/docker/services/localhost/docker-compose.yaml up >nohup-docker-compose-up 2>&1 &
	#cd cmd; nohup go run main.go -configServerUrl=http://localhost:8888 -profile=dev -configBranch=master >nohup-main 2>&1 &; cd -

rm-local-services:
	docker-compose -f support/docker/services/localhost/docker-compose.yaml down

