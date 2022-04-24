MODBOT_VERSION ?= 1.11.3
BACKUP_DIR ?= ~/backup
MODBOT_BACKUP_FILE ?= modbot-backup.tar

.DEFAULT_GOAL := help

# declares .PHONY which will run the make command even if a file of the same name exists
.PHONY: help
help:			## Help command
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

lint:			## Lint check
	docker run --rm -v $(PWD):/src:Z \
	--workdir=/src odinuge/yapf:latest yapf \
	--style '{based_on_style: pep8, dedent_closing_brackets: true, coalesce_brackets: true}' \
	--no-local-style --verbose --recursive --diff --parallel app compose

format:			## Format code in place to conform to lint check
	docker run --rm -v $(PWD):/src:Z \
	--workdir=/src odinuge/yapf:latest yapf \
	--style '{based_on_style: pep8, dedent_closing_brackets: true, coalesce_brackets: true}' \
	--no-local-style --verbose --recursive --in-place --parallel app compose

pyflakes:		## Pyflakes check for any unused variables/classes
	docker run --rm -v $(PWD):/src:Z \
	--workdir=/src python:3.8 \
	/bin/bash -c "pip install --upgrade pyflakes && python -m pyflakes /src && echo 'pyflakes passed!'"

backup-modbot:		## Backup database volumes to BACKUP_DIR
	docker stop modbot modbot_db
	docker run --rm --volumes-from modbot_db -v $(BACKUP_DIR):/backup ubuntu bash -c "cd /data/db && tar cvf /backup/$(MODBOT_BACKUP_FILE) ."
	docker start modbot modbot_db

restore-backup-modbot:		## Restore volumes backup from BACKUP_DIR
	ls $(BACKUP_DIR) | grep modbot-backup.tar
	docker volume create telegram-modbot_modbot-db
	docker run --rm -v telegram-modbot_modbot-db:/recover -v $(BACKUP_DIR):/backup ubuntu bash -c "cd /recover && tar xvf /backup/modbot-backup.tar"

start-prod:		## Pull latest version of modbot image and run docker-compose up
	docker-compose pull modbot
	docker-compose up -d
	
start-dev:		## Run dev instance of modbot with live reload of api
	docker-compose -f docker-compose.dev.yml up -d

build-modbot:		## Build docker image for modbot and pollingserver
	docker build --tag jasoncky96/telegram-modbot:latest -f ./compose/webserver/Dockerfile .
	docker build --tag jasoncky96/telegram-pollingserver:latest -f ./compose/pollingserver/Dockerfile .

deploy-pollingserver:		## Deploy docker image for pollingserver
	docker buildx build --push --tag jasoncky96/telegram-pollingserver:latest --file ./compose/pollingserver/Dockerfile --platform linux/arm/v7,linux/arm64/v8,linux/amd64 .

deploy-modbot:			## Deploy docker image for modbot
	docker buildx build --push --tag jasoncky96/telegram-modbot:$(MODBOT_VERSION) --file ./compose/webserver/Dockerfile --platform linux/arm/v7,linux/arm64/v8,linux/amd64 .
	docker buildx build --push --tag jasoncky96/telegram-modbot:latest --file ./compose/webserver/Dockerfile --platform linux/arm/v7,linux/arm64/v8,linux/amd64 .

stop-modbot:		## Run docker-compose down
	docker-compose down
	
destroy:			## Run docker-compose down -v
	docker-compose down -v
