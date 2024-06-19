run-consumer:
	@go run ./consumer/main.go

run-pconsumer:
	@go run ./consumer_plugin/main.go

run-producer:
	@go run ./producer/main.go

run-pproducer:
	@go run ./producer_plugin/main.go

sniff:
# @docker run -it --net=container:streaming_mq tcpdump tcpdump port 5552 -A -n -w rabbitmq.pcap
	@docker run --name tcpdump --rm -it --net=container:streaming_mq tcpdump tcpdump -w rabbitmq.pcap port 5552 -A -n

copy:
	@docker cp tcpdump:rabbitmq.pcap /mnt/d/xProject/go-dev/projects/MICROSERVICES/RabbitMQ/rabbitmq-streams

trace:
	sniff && copy

docker-rmq:
	@docker run --rm --name rabbitmq --net rabbitmq-net -p 5672:5672 -p 15672:15672 rabbitmq:3-management

docker-rmqp:
	@docker run --rm --name streaming_mq --net rabbitmq-net -p 5552:5552  -p 5672:5672 -p 15672:15672 -e RABBITMQ_SERVER_ADDITIONAL_ERL_ARGS='-rabbitmq_stream advertised_host localhost -rabbitmq loopback_users "none"' rabbitmq:3-management