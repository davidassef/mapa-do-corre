import { useState } from 'react';

import { CadastroPage } from './pages/CadastroPage';
import { ImpactoPage } from './pages/ImpactoPage';
import { MapaPage } from './pages/MapaPage';

type Aba = 'mapa' | 'cadastro' | 'impacto';

const abas: Array<{ id: Aba; label: string; descricao: string }> = [
  {
    id: 'mapa',
    label: 'Mapa',
    descricao: 'Busca por proximidade e contato direto',
  },
  {
    id: 'cadastro',
    label: 'Cadastro',
    descricao: 'Publicacao aberta para novos corres',
  },
  {
    id: 'impacto',
    label: 'Impacto',
    descricao: 'Painel simples de conexoes geradas',
  },
];

export default function App() {
  const [abaAtual, setAbaAtual] = useState<Aba>('mapa');

  return (
    <div className="min-h-screen bg-[radial-gradient(circle_at_top_left,_rgba(244,162,97,0.34),_transparent_26%),linear-gradient(135deg,_#f7f1e8_0%,_#f0e4d3_48%,_#d8efe6_100%)] px-4 py-6 text-noite md:px-8 lg:px-10">
      <div className="mx-auto max-w-7xl space-y-6">
        <header className="overflow-hidden rounded-[2.25rem] border border-white/60 bg-noite px-6 py-7 text-white shadow-mapa md:px-8">
          <div className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr] lg:items-end">
            <div className="space-y-4">
              <span className="inline-flex rounded-full border border-white/20 px-4 py-2 text-xs font-semibold uppercase tracking-[0.22em] text-white/70">
                Plataforma territorial para Fortaleza
              </span>
              <div>
                <h1 className="font-display text-4xl leading-tight md:text-5xl">Mapa do Corre</h1>
                <p className="mt-3 max-w-2xl text-sm leading-7 text-white/72 md:text-base">
                  Uma vitrine de proximidade para microempreendedores informais e autonomos serem encontrados sem intermediarios, com atrito minimo e contato direto pelo WhatsApp.
                </p>
              </div>
            </div>

            <div className="grid gap-4 rounded-[1.75rem] border border-white/10 bg-white/5 p-4 backdrop-blur">
              <span className="text-xs font-semibold uppercase tracking-[0.18em] text-white/55">Nucleo do MVP</span>
              <div className="grid gap-3 md:grid-cols-3 lg:grid-cols-1 xl:grid-cols-3">
                <PainelResumo titulo="Busca georreferenciada" valor="Raio local" />
                <PainelResumo titulo="Contato direto" valor="wa.me" />
                <PainelResumo titulo="Metrica de extensao" valor="tracking" />
              </div>
            </div>
          </div>
        </header>

        <nav className="flex flex-wrap gap-3">
          {abas.map((aba) => (
            <button
              key={aba.id}
              type="button"
              onClick={() => setAbaAtual(aba.id)}
              className={
                abaAtual === aba.id
                  ? 'rounded-full bg-coqueiro px-5 py-3 text-left text-sm font-semibold text-white shadow-mapa'
                  : 'rounded-full border border-white/70 bg-white/80 px-5 py-3 text-left text-sm font-semibold text-noite shadow-mapa backdrop-blur transition hover:border-coqueiro/30'
              }
            >
              <span className="block">{aba.label}</span>
              <span className="mt-1 block text-xs font-medium text-inherit/70">{aba.descricao}</span>
            </button>
          ))}
        </nav>

        {abaAtual === 'mapa' ? <MapaPage /> : null}
        {abaAtual === 'cadastro' ? <CadastroPage onAbrirMapa={() => setAbaAtual('mapa')} /> : null}
        {abaAtual === 'impacto' ? <ImpactoPage /> : null}
      </div>
    </div>
  );
}

type PainelResumoProps = {
  titulo: string;
  valor: string;
};

function PainelResumo({ titulo, valor }: PainelResumoProps) {
  return (
    <div className="rounded-[1.25rem] border border-white/10 bg-white/5 p-4">
      <span className="block text-xs font-semibold uppercase tracking-[0.16em] text-white/50">{titulo}</span>
      <strong className="mt-2 block font-display text-2xl text-white">{valor}</strong>
    </div>
  );
}