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
    <div className="min-h-screen bg-[#eef3ef] px-4 py-5 text-noite md:px-8 lg:px-10">
      <div className="mx-auto max-w-7xl space-y-5">
        <header className="rounded-lg border border-noite/10 bg-noite px-5 py-5 text-white shadow-mapa md:px-6">
          <div className="grid gap-5 lg:grid-cols-[minmax(0,1fr)_minmax(420px,0.72fr)] lg:items-end">
            <div className="max-w-3xl">
              <span className="inline-flex rounded-full border border-white/20 px-3 py-1 text-xs font-semibold uppercase text-white/70">
                Fortaleza
              </span>
              <div>
                <h1 className="mt-3 font-display text-3xl leading-tight md:text-4xl">Mapa do Corre</h1>
                <p className="mt-2 max-w-2xl text-sm leading-6 text-white/75 md:text-base">
                  Encontre microempreendedores e autonomos perto de voce, com contato direto pelo WhatsApp.
                </p>
              </div>
            </div>

            <dl className="grid overflow-hidden rounded-lg border border-white/10 bg-white/10 sm:grid-cols-3">
              <PainelResumo titulo="Busca" valor="Raio local" />
              <PainelResumo titulo="Contato" valor="wa.me" />
              <PainelResumo titulo="Impacto" valor="Cliques" />
            </dl>
          </div>
        </header>

        <nav className="grid gap-1 rounded-lg border border-noite/10 bg-white/85 p-1 shadow-sm backdrop-blur md:grid-cols-3">
          {abas.map((aba) => (
            <button
              key={aba.id}
              type="button"
              onClick={() => setAbaAtual(aba.id)}
              className={
                abaAtual === aba.id
                  ? 'rounded-md bg-coqueiro px-4 py-3 text-left text-sm font-semibold text-white shadow-sm'
                  : 'rounded-md px-4 py-3 text-left text-sm font-semibold text-noite/75 transition hover:bg-coqueiro/10 hover:text-coqueiro'
              }
            >
              <span className="block">{aba.label}</span>
              <span className="mt-1 block text-xs font-medium opacity-70">{aba.descricao}</span>
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
    <div className="border-b border-white/10 px-4 py-3 last:border-b-0 sm:border-b-0 sm:border-r sm:last:border-r-0">
      <dt className="text-xs font-semibold uppercase text-white/55">{titulo}</dt>
      <dd className="mt-1 font-display text-xl text-white">{valor}</dd>
    </div>
  );
}