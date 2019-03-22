.PHONY: docker-build docker-swarm-init docker-deploy-services

docker-swarm-init:
	docker swarm init

docker-swarm-leave:
	docker swarm leave --force

docker-build:
	docker build --build-arg moduleName=gateway -t farzadras/gateway -f ./docker/dockerfile-moduleName .
	docker build --build-arg moduleName=zombie_driver -t farzadras/zombie_driver -f ./docker/dockerfile-moduleName .
	docker build --build-arg moduleName=driver_location -t farzadras/driver_location -f ./docker/dockerfile-moduleName .
	docker rmi $$(docker images -q -f dangling=true)

docker-deploy-services:
	docker stack deploy --compose-file ./docker/docker-compose-s1.yaml services
	docker stack deploy --compose-file ./docker/docker-compose-s2.yaml services
	docker stack deploy --compose-file ./docker/docker-compose-s3.yaml services
	docker stack deploy --compose-file ./docker/docker-compose-s4.yaml services

docker-rm-services:
	docker stack rm services
	docker rm -f $$(docker ps -a -q)
	docker rmi $$(docker images -q -f dangling=true)

