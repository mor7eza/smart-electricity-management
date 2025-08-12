emqx-init:
	# port 1883 uses for tcp and 18083 for dashboard
	docker run -d --name emqx -p 1883:1883 -p 18083:18083 emqx/emqx

emqx-start:
	docker run -d emqx

transmitter:
	@ cd ./transmitter && go run . -count=100 -interval=60

api-gateway:
	@ cd ./services/api-gateway && go run ./cmd/server
