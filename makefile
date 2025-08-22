.PHONY: up down

up:
	podman-compose up --build

down:
	podman-compose down
