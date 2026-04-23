run:
	docker compose up -d

build:
	docker compose up --build -d

stop:
	docker compose down

images:
	docker ps

all-logs:
	docker compose logs -f

logs-%:
	docker compose logs $* -f

restart:
	docker compose restart

restart-%:
	docker compose restart $*
