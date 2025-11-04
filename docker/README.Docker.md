### Building and running your application

When you're ready, start your application by running:
` docker compose -f docker/compose.yaml up --build -d`.
To access the inside of the container:
docker exec -it <container_id> /bin/sh

Your application will be available at http://localhost:5080.
