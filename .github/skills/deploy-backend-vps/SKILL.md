---
name: deploy-backend-vps
description: 'Use quando precisar fazer deploy, atualizar, diagnosticar ou consultar logs do backend do Mapa do Corre na VPS em /opt/mapa-do-corre usando Docker Compose e Postgres compartilhado.'
argument-hint: 'acao desejada: bootstrap, deploy, logs, status ou diagnostico'
---

# Deploy Backend VPS

## Quando Usar
- Publicar a API do Mapa do Corre na VPS.
- Rebuildar o container `mapa-do-corre-backend`.
- Consultar logs, status ou health check da API.
- Conferir se o compose usa a rede compartilhada `shared-db-network`.

## Procedimento
1. Rode `make test` localmente.
2. Confira se `backend/.env` existe na VPS e nunca exponha seus valores.
3. Para o primeiro deploy, rode `make bootstrap-vps` depois do push inicial.
4. Rode `make deploy-backend` para atualizar por `git pull` e `docker compose`.
5. Verifique com `make status-vps` e `curl https://api.lotoscore.com.br/mapa-do-corre/health`.
6. Use `make logs-backend` se o container nao ficar saudavel.

## Regras
- Nao versionar `.env` nem imprimir segredos no terminal.
- Nao executar DDL ou migrations em producao sem permissao explicita em caixa alta.
- Reutilizar `shared-postgres` na rede `shared-db-network`; nao subir outro Postgres para producao.
- Deploy na VPS deve ser via push/pull, nunca copiando arquivos soltos.
