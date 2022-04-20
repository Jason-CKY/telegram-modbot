MODBOT_VERSION ?= 1.11.2
BACKUP_DIR ?= ~/backup
MODBOT_BACKUP_FILE ?= modbot-backup.tar

format: 
	yapf -i -r -p telegram-reminderbot telegram-modbot

backup-modbot:
	docker stop modbot modbot_db
	docker run --rm --volumes-from modbot_db -v $(BACKUP_DIR):/backup ubuntu bash -c "cd /data/db && tar cvf /backup/$(MODBOT_BACKUP_FILE) ."
	docker start modbot modbot_db

restore-backup-modbot:
	ls $(BACKUP_DIR) | grep modbot-backup.tar
	docker volume create telegram-modbot_modbot-db
	docker run --rm -v telegram-modbot_modbot-db:/recover -v $(BACKUP_DIR):/backup ubuntu bash -c "cd /recover && tar xvf /backup/modbot-backup.tar"

start-prod:
	docker-compose pull modbot
	docker-compose up -d
	
start-dev:
	docker-compose -f docker-compose.dev.yml up -d

build-modbot:
	docker build --tag jasoncky96/telegram-modbot:latest -f ./compose/webserver/Dockerfile .
	docker build --tag jasoncky96/telegram-pollingserver:latest -f ./compose/pollingserver/Dockerfile .

deploy-pollingserver:
	docker buildx build --push --tag jasoncky96/telegram-pollingserver:latest --file ./compose/pollingserver/Dockerfile --platform linux/arm/v7,linux/arm64/v8,linux/amd64 .

deploy-modbot:
	docker buildx build --push --tag jasoncky96/telegram-modbot:$(MODBOT_VERSION) --file ./compose/webserver/Dockerfile --platform linux/arm/v7,linux/arm64/v8,linux/amd64 .
	docker buildx build --push --tag jasoncky96/telegram-modbot:latest --file ./compose/webserver/Dockerfile --platform linux/arm/v7,linux/arm64/v8,linux/amd64 .

stop-modbot:
	docker-compose down
	
destroy:
	docker-compose down -v
