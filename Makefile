SSH_KEY ?= ~/.ssh/LightsailDefaultKey-ca-central-1.pem
VPS ?= ubuntu@16.52.255.233
VPS_DIR ?= /opt/mapa-do-corre
COMPOSE_FILE ?= docker-compose.vps.yml

SSH := ssh -i $(SSH_KEY) $(VPS)

.PHONY: help test build deploy-frontend deploy-backend bootstrap-vps logs-backend status-vps

help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | awk 'BEGIN {FS = ":.*?## "} {printf "%-18s %s\n", $$1, $$2}'

test: ## Executa testes do backend
	cd backend && go test ./...

build: ## Gera build do frontend
	cd frontend && npm run build

deploy-frontend: ## Publica o frontend na Vercel em producao
	cd frontend && vercel --prod --yes

bootstrap-vps: ## Prepara /opt/mapa-do-corre na VPS depois do push inicial
	$(SSH) 'sudo mkdir -p $(VPS_DIR) \
		&& sudo chown ubuntu:ubuntu $(VPS_DIR) \
		&& if [ ! -d $(VPS_DIR)/.git ]; then git clone git@github.com:davidassef/mapa-do-corre.git $(VPS_DIR); fi'

deploy-backend: ## Atualiza a API na VPS por git pull e docker compose
	git push origin main
	$(SSH) 'cd $(VPS_DIR) \
		&& git checkout main \
		&& git pull origin main \
		&& docker compose -f $(COMPOSE_FILE) up -d --build'

logs-backend: ## Mostra logs recentes da API na VPS
	$(SSH) 'docker logs mapa-do-corre-backend --tail 120 -f'

status-vps: ## Mostra estado do container da API na VPS
	$(SSH) 'docker ps --filter name=mapa-do-corre-backend --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"'
