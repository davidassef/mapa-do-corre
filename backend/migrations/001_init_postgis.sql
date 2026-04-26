CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS prestadores (
    id UUID PRIMARY KEY,
    nome VARCHAR(120) NOT NULL,
    categoria VARCHAR(80) NOT NULL,
    descricao TEXT NOT NULL DEFAULT '',
    whatsapp VARCHAR(20) NOT NULL,
    bairro VARCHAR(120) NOT NULL DEFAULT '',
    localizacao GEOGRAPHY(POINT, 4326) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS logs_cliques (
    id UUID PRIMARY KEY,
    prestador_id UUID NOT NULL REFERENCES prestadores(id) ON DELETE CASCADE,
    origem VARCHAR(50) NOT NULL DEFAULT 'whatsapp',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_prestadores_categoria ON prestadores (categoria);
CREATE INDEX IF NOT EXISTS idx_prestadores_localizacao ON prestadores USING GIST (localizacao);
CREATE INDEX IF NOT EXISTS idx_logs_cliques_prestador_id ON logs_cliques (prestador_id);