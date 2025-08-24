.PHONY: up down

up:
	podman-compose up --build

up-d:
	podman-compose up --build -d

down:
	podman-compose down
