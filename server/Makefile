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

integration:
	if [ -e "../integration/scripts/.env" ]; then rm ../integration/scripts/.env; fi
	touch ../integration/scripts/.env
	cat .env | grep INSTALLER_ID >> ../integration/scripts/.env
	cat .env | grep DECISION_ID >> ../integration/scripts/.env

	cat .env | grep INSTALLER_TOKEN >> ../integration/scripts/.env
	cat .env | grep DECISION_TOKEN >> ../integration/scripts/.env

	grep DECISION_SERVER_BASE_URL .env | cut -d '=' -f 2 | sed 's/^/URL=/' >> ../integration/scripts/.env
