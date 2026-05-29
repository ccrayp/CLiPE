run:
	@docker compose up -d
	@echo ""
	@printf "  Панель администратора: \033]8;;https://127.0.0.1\033\\\\https://clipe\033]8;;\033\\\\\n"

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
	cat .env | grep INSTALLER_TOKEN >> ../integration/scripts/.env

	grep DECISION_SERVER_BASE_URL .env | cut -d '=' -f 2 | sed 's/^/URL=/' >> ../integration/scripts/.env

rotate_secret:
	@INSTALLER_TOKEN=$$(openssl rand -base64 64 | tr -d '\n'); \
	DECISION_TOKEN=$$(openssl rand -base64 64 | tr -d '\n'); \
	sed -i.bak "s|^INSTALLER_TOKEN=.*|INSTALLER_TOKEN=$$INSTALLER_TOKEN|" .env; \
	sed -i.bak "s|^DECISION_TOKEN=.*|DECISION_TOKEN=$$DECISION_TOKEN|" .env; \
	rm -f .env.bak

	$(MAKE) integration

rotate-jwt:
	@JWT_SECRET_KEY=$$(openssl rand -base64 64 | tr -d '\n'); \
	sed -i.bak "s|^JWT_SECRET_KEY=.*|JWT_SECRET_KEY=$$JWT_SECRET_KEY|" .env; \
	rm -f .env.bak

rotate-all: rotate-secrets rotate-jwt