ALTER TABLE prestadores
    ADD COLUMN IF NOT EXISTS email_responsavel VARCHAR(255),
    ADD COLUMN IF NOT EXISTS removido_em TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_prestadores_removido_em ON prestadores (removido_em);

CREATE TABLE IF NOT EXISTS solicitacoes_cadastro_prestador (
    id UUID PRIMARY KEY,
    nome VARCHAR(120) NOT NULL,
    categoria VARCHAR(80) NOT NULL,
    descricao TEXT NOT NULL DEFAULT '',
    whatsapp VARCHAR(20) NOT NULL,
    bairro VARCHAR(120) NOT NULL DEFAULT '',
    email VARCHAR(255) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    codigo_hash CHAR(64) NOT NULL,
    tentativas INTEGER NOT NULL DEFAULT 0,
    expira_em TIMESTAMPTZ NOT NULL,
    confirmado_em TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_solicitacoes_cadastro_email ON solicitacoes_cadastro_prestador (email);
CREATE INDEX IF NOT EXISTS idx_solicitacoes_cadastro_expira_em ON solicitacoes_cadastro_prestador (expira_em);

CREATE TABLE IF NOT EXISTS solicitacoes_remocao_prestador (
    id UUID PRIMARY KEY,
    prestador_id UUID REFERENCES prestadores(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    codigo_hash CHAR(64) NOT NULL,
    tentativas INTEGER NOT NULL DEFAULT 0,
    expira_em TIMESTAMPTZ NOT NULL,
    confirmado_em TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_solicitacoes_remocao_email ON solicitacoes_remocao_prestador (email);
CREATE INDEX IF NOT EXISTS idx_solicitacoes_remocao_expira_em ON solicitacoes_remocao_prestador (expira_em);