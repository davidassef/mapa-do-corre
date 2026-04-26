import { useEffect, useState } from 'react';

import { buscarResumoImpacto } from '../services/api';
import type { ResumoImpacto } from '../types/corre';

const formatadorNumero = new Intl.NumberFormat('pt-BR');

export function ImpactoPage() {
  const [resumo, setResumo] = useState<ResumoImpacto | null>(null);
  const [origemDados, setOrigemDados] = useState<'api' | 'fallback'>('fallback');

  useEffect(() => {
    let isMounted = true;

    void (async () => {
      const payload = await buscarResumoImpacto();
      if (!isMounted) {
        return;
      }

      setResumo(payload.resumo);
      setOrigemDados(payload.origem);
    })();

    return () => {
      isMounted = false;
    };
  }, []);

  if (!resumo) {
    return (
      <section className="rounded-[2rem] border border-white/70 bg-white/90 p-6 shadow-mapa backdrop-blur">
        Carregando indicadores de impacto...
      </section>
    );
  }

  return (
    <section className="space-y-6">
      <div className="rounded-[2rem] border border-white/70 bg-white/90 p-6 shadow-mapa backdrop-blur">
        <div className="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
          <div>
            <span className="inline-flex rounded-full bg-coral/15 px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em] text-coral">
              Dashboard de impacto
            </span>
            <h2 className="mt-3 font-display text-3xl text-noite">Conexoes geradas pela rede local</h2>
            <p className="mt-2 text-sm leading-6 text-noite/70">
              Este painel resume o volume de contatos disparados a partir do app e serve de base para o relatorio de extensao.
            </p>
          </div>

          <div className="rounded-[1.5rem] bg-areia px-4 py-3 text-sm text-noite/70">
            Fonte atual: {origemDados === 'api' ? 'API' : 'fallback local'}
          </div>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <ImpactoCard titulo="Total de conexoes" valor={formatadorNumero.format(resumo.totalConexoesGeradas)} detalhe="Cliques rastreados no botao de contato" />
        <ImpactoCard titulo="Prestadores ativos" valor={formatadorNumero.format(resumo.totalPrestadores)} detalhe="Cadastros atualmente visiveis" />
        <ImpactoCard titulo="Corres removidos" valor={formatadorNumero.format(resumo.totalPrestadoresRemovidos)} detalhe="Remocoes seguras sem apagar o historico" />
        <ImpactoCard titulo="Categorias ativas" valor={formatadorNumero.format(resumo.categorias.length)} detalhe="Diversidade da rede territorial" />
      </div>

      <div className="grid gap-4 lg:grid-cols-2">
        {resumo.categorias.map((categoria) => (
          <article key={categoria.categoria} className="rounded-[1.75rem] border border-white/70 bg-white/90 p-5 shadow-mapa backdrop-blur">
            <div className="flex items-center justify-between gap-4">
              <div>
                <span className="text-xs font-semibold uppercase tracking-[0.18em] text-coqueiro/70">Categoria</span>
                <h3 className="mt-2 font-display text-2xl capitalize text-noite">{categoria.categoria}</h3>
              </div>

              <div className="rounded-2xl bg-mar/10 px-4 py-3 text-right">
                <span className="block text-xs font-semibold uppercase tracking-[0.16em] text-coqueiro/70">Cliques</span>
                <strong className="font-display text-2xl text-coqueiro">{formatadorNumero.format(categoria.totalCliques)}</strong>
              </div>
            </div>

            <p className="mt-4 text-sm leading-6 text-noite/70">
              {formatadorNumero.format(categoria.totalPrestadores)} prestadores cadastrados nesta categoria.
            </p>
          </article>
        ))}
      </div>
    </section>
  );
}

type ImpactoCardProps = {
  titulo: string;
  valor: string;
  detalhe: string;
};

function ImpactoCard({ titulo, valor, detalhe }: ImpactoCardProps) {
  return (
    <article className="rounded-[1.75rem] border border-white/70 bg-white/90 p-5 shadow-mapa backdrop-blur">
      <span className="text-xs font-semibold uppercase tracking-[0.2em] text-coqueiro/70">{titulo}</span>
      <strong className="mt-3 block font-display text-4xl text-noite">{valor}</strong>
      <p className="mt-2 text-sm text-noite/70">{detalhe}</p>
    </article>
  );
}