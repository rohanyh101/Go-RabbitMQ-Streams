# Use the official RabbitMQ 3 image
FROM rabbitmq:3-management

# Set environment variables
ENV RABBITMQ_SERVER_ADDITIONAL_ERL_ARGS='-rabbitmq_stream advertised_host localhost -rabbitmq loopback_users "none"'

# Map ports
EXPOSE 5552
EXPOSE 5672
EXPOSE 15672

# Set the container name
RUN docker run --rm --name streaming_mq -p 5552:5552 -p 5672:5672 -p 15672:15672 rabbitmq:3-management --net rabbitmq-net