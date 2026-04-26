---
name: deploy-frontend-vercel
description: 'Use quando precisar criar, configurar, publicar ou diagnosticar o frontend do Mapa do Corre na Vercel com Vite, React e VITE_API_URL.'
argument-hint: 'acao desejada: link, env, deploy, dominio ou diagnostico'
---

# Deploy Frontend Vercel

## Quando Usar
- Criar ou linkar o projeto `mapa-do-corre` na Vercel.
- Configurar `VITE_API_URL` para a API publica.
- Publicar o frontend React/Vite em producao.
- Diagnosticar build ou variaveis de ambiente na Vercel.

## Procedimento
1. Rode `cd frontend && npm run build` antes de publicar.
2. Confira login com `vercel whoami`.
3. Linke o projeto a partir de `frontend` com `vercel link --yes --project mapa-do-corre`.
4. Configure `VITE_API_URL` em producao e preview apontando para a API publica.
5. Publique com `cd frontend && vercel --prod --yes`.
6. Use a URL final da Vercel como `FRONTEND_ORIGIN` no backend.

## Regras
- Nao versionar `.vercel/`.
- Nao deixar `VITE_API_URL` apontando para `localhost` em producao.
- Se a API estiver sob prefixo de proxy, a URL deve ficar sem barra final, por exemplo `https://api.lotoscore.com.br/mapa-do-corre`.
