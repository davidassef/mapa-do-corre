import { useEffect, useState } from 'react';

import { CategoriaFilter } from '../components/CategoriaFilter';
import { CorreMap } from '../components/CorreMap';
import { confirmarRemocaoPrestador, solicitarRemocaoPrestador } from '../services/api';
import { useCorreStore } from '../stores/useCorreStore';
import { useLocationStore } from '../stores/useLocationStore';
import type { Prestador, SolicitacaoCodigo } from '../types/corre';

const formatadorDistancia = new Intl.NumberFormat('pt-BR');

export function MapaPage() {
  const centroAtual = useLocationStore((state) => state.centroAtual);
  const raioMetros = useLocationStore((state) => state.raioMetros);
  const isLocalizando = useLocationStore((state) => state.isLocalizando);
  const statusPermissao = useLocationStore((state) => state.statusPermissao);
  const usarLocalizacaoAtual = useLocationStore((state) => state.usarLocalizacaoAtual);
  const atualizarRaio = useLocationStore((state) => state.atualizarRaio);

  const isLoading = useCorreStore((state) => state.isLoading);
  const errorMessage = useCorreStore((state) => state.errorMessage);
  const origemDados = useCorreStore((state) => state.origemDados);
  const corres = useCorreStore((state) => state.corres);
  const categoriaAtiva = useCorreStore((state) => state.filtroCategoria);
  const carregarCorres = useCorreStore((state) => state.carregarCorres);
  const definirFiltroCategoria = useCorreStore((state) => state.definirFiltroCategoria);
  const registrarContato = useCorreStore((state) => state.registrarContato);

  const corresFiltrados = corres.filter((prestador) => {
    if (!categoriaAtiva) {
      return true;
    }

    return prestador.categoria === categoriaAtiva;
  });

  const categorias = Array.from(new Set(corres.map((prestador) => prestador.categoria))).sort((categoriaAtual, proximaCategoria) =>
    categoriaAtual.localeCompare(proximaCategoria),
  );

  useEffect(() => {
    void carregarCorres({
      latitude: centroAtual.latitude,
      longitude: centroAtual.longitude,
      raioMetros,
    });
  }, [carregarCorres, centroAtual.latitude, centroAtual.longitude, raioMetros]);

  async function recarregarMapa() {
    await carregarCorres({
      latitude: centroAtual.latitude,
      longitude: centroAtual.longitude,
      raioMetros,
    });
  }

  return (
    <section className="grid gap-6 xl:grid-cols-[1.4fr_0.9fr]">
      <div className="space-y-5">
        <div className="rounded-[2rem] border border-white/70 bg-white/90 p-6 shadow-mapa backdrop-blur">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
            <div className="space-y-2">
              <span className="inline-flex rounded-full bg-sol/20 px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-coral">
                Descoberta por proximidade
              </span>
              <h2 className="font-display text-3xl text-noite">Encontre quem resolve perto de voce</h2>
              <p className="max-w-2xl text-sm leading-6 text-noite/70">
                O mapa usa a sua localizacao atual ou um centro padrao em Fortaleza para listar corres dentro do raio definido.
              </p>
            </div>

            <button
              type="button"
              onClick={() => void usarLocalizacaoAtual()}
              className="rounded-full bg-coqueiro px-5 py-3 text-sm font-semibold text-white transition hover:bg-mar"
            >
              {isLocalizando ? 'Buscando localizacao...' : 'Usar minha localizacao'}
            </button>
          </div>

          <div className="mt-5 grid gap-4 md:grid-cols-[220px_1fr]">
            <label className="space-y-2 text-sm font-medium text-noite">
              Raio de busca
              <select
                value={raioMetros}
                onChange={(event) => atualizarRaio(Number(event.target.value))}
                className="w-full rounded-2xl border border-coqueiro/20 bg-areia px-4 py-3 text-noite outline-none focus:border-coqueiro"
              >
                <option value={1000}>1 km</option>
                <option value={2500}>2.5 km</option>
                <option value={5000}>5 km</option>
                <option value={8000}>8 km</option>
              </select>
            </label>

            <div className="space-y-2 text-sm font-medium text-noite">
              <span>Categoria</span>
              <CategoriaFilter
                categorias={categorias}
                categoriaAtiva={categoriaAtiva}
                onChange={definirFiltroCategoria}
              />
            </div>
          </div>

          <div className="mt-5 grid gap-4 md:grid-cols-3">
            <Indicador titulo="Centro da busca" valor="Fortaleza" detalhe={descreverPermissao(statusPermissao)} />
            <Indicador titulo="Prestadores visiveis" valor={String(corresFiltrados.length)} detalhe={origemDados === 'fallback' ? 'modo offline' : 'dados da API'} />
            <Indicador titulo="Raio atual" valor={`${formatadorDistancia.format(raioMetros)} m`} detalhe="sem novo fetch para categoria" />
          </div>
        </div>

        <CorreMap centro={centroAtual} raioMetros={raioMetros} prestadores={corresFiltrados} />
      </div>

      <div className="space-y-4">
        <div className="rounded-[2rem] border border-white/70 bg-white/90 p-6 shadow-mapa backdrop-blur">
          <h3 className="font-display text-2xl text-noite">Lista de corres</h3>
          <p className="mt-2 text-sm leading-6 text-noite/70">
            O filtro de categoria acontece no estado local para evitar nova ida ao backend quando so a categoria mudar.
          </p>
        </div>

        {isLoading ? <CardAviso mensagem="Carregando corres proximos..." /> : null}
        {errorMessage ? <CardAviso mensagem={errorMessage} tom="erro" /> : null}

        {corresFiltrados.map((prestador) => (
          <PrestadorCard key={prestador.id} prestador={prestador} onContato={registrarContato} onRemocaoConcluida={recarregarMapa} />
        ))}

        {!isLoading && corresFiltrados.length === 0 ? <CardAviso mensagem="Nenhum corre encontrado para este recorte." /> : null}
      </div>
    </section>
  );
}

type IndicadorProps = {
  titulo: string;
  valor: string;
  detalhe: string;
};

function Indicador({ titulo, valor, detalhe }: IndicadorProps) {
  return (
    <div className="rounded-[1.5rem] border border-coqueiro/10 bg-areia p-4">
      <span className="text-xs font-semibold uppercase tracking-[0.2em] text-coqueiro/70">{titulo}</span>
      <strong className="mt-2 block font-display text-2xl text-noite">{valor}</strong>
      <p className="mt-2 text-sm text-noite/65">{detalhe}</p>
    </div>
  );
}

type CardAvisoProps = {
  mensagem: string;
  tom?: 'neutro' | 'erro';
};

function CardAviso({ mensagem, tom = 'neutro' }: CardAvisoProps) {
  const classes = tom === 'erro'
    ? 'rounded-[1.75rem] border border-coral/20 bg-coral/10 p-5 text-sm text-coral'
    : 'rounded-[1.75rem] border border-white/70 bg-white/90 p-5 text-sm text-noite/70 shadow-mapa';

  return <div className={classes}>{mensagem}</div>;
}

type PrestadorCardProps = {
  prestador: Prestador;
  onContato: (prestadorId: string) => Promise<void>;
  onRemocaoConcluida: () => Promise<void>;
};

function PrestadorCard({ prestador, onContato, onRemocaoConcluida }: PrestadorCardProps) {
  const [isRemocaoAberta, setIsRemocaoAberta] = useState(false);
  const [emailRemocao, setEmailRemocao] = useState('');
  const [codigoRemocao, setCodigoRemocao] = useState('');
  const [solicitacaoRemocao, setSolicitacaoRemocao] = useState<SolicitacaoCodigo | null>(null);
  const [statusRemocao, setStatusRemocao] = useState<{ tom: 'idle' | 'success' | 'error'; mensagem: string } | null>(null);
  const [isSolicitandoRemocao, setIsSolicitandoRemocao] = useState(false);
  const [isConfirmandoRemocao, setIsConfirmandoRemocao] = useState(false);

  const linkWhatsApp = `https://wa.me/${prestador.whatsApp}?text=${encodeURIComponent(
    `Oi, ${prestador.nome}. Vi seu corre no Mapa do Corre e queria mais informacoes.`,
  )}`;

  async function handleContato() {
    await onContato(prestador.id);
    window.open(linkWhatsApp, '_blank', 'noopener,noreferrer');
  }

  async function handleSolicitarRemocao(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setIsSolicitandoRemocao(true);

    try {
      const payload = await solicitarRemocaoPrestador(prestador.id, {
        email: emailRemocao,
      });

      setSolicitacaoRemocao(payload.solicitacao);
      setCodigoRemocao('');
      setStatusRemocao({
        tom: 'idle',
        mensagem: payload.message,
      });
    } catch (error) {
      setStatusRemocao({
        tom: 'error',
        mensagem: extrairMensagemErro(error, 'Nao foi possivel solicitar a remocao agora.'),
      });
    } finally {
      setIsSolicitandoRemocao(false);
    }
  }

  async function handleConfirmarRemocao(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!solicitacaoRemocao) {
      return;
    }

    setIsConfirmandoRemocao(true);

    try {
      const mensagem = await confirmarRemocaoPrestador({
        solicitacaoId: solicitacaoRemocao.solicitacaoId,
        codigo: codigoRemocao,
      });

      setStatusRemocao({
        tom: 'success',
        mensagem,
      });
      await onRemocaoConcluida();
    } catch (error) {
      setStatusRemocao({
        tom: 'error',
        mensagem: extrairMensagemErro(error, 'Nao foi possivel confirmar a remocao agora.'),
      });
    } finally {
      setIsConfirmandoRemocao(false);
    }
  }

  function handleResetarRemocao() {
    setSolicitacaoRemocao(null);
    setCodigoRemocao('');
    setStatusRemocao(null);
  }

  return (
    <article className="rounded-[1.75rem] border border-white/70 bg-white/90 p-5 shadow-mapa backdrop-blur">
      <div className="flex items-start justify-between gap-4">
        <div>
          <span className="inline-flex rounded-full bg-mar/15 px-3 py-1 text-xs font-semibold uppercase tracking-[0.16em] text-coqueiro">
            {prestador.categoria}
          </span>
          <h4 className="mt-3 font-display text-xl text-noite">{prestador.nome}</h4>
          <p className="mt-2 text-sm leading-6 text-noite/70">{prestador.descricao}</p>
        </div>

        <div className="rounded-2xl bg-areia px-3 py-2 text-right text-xs font-semibold text-noite/70">
          <span className="block">Distancia</span>
          <strong className="font-display text-lg text-noite">{formatadorDistancia.format(prestador.distanciaMetros)} m</strong>
        </div>
      </div>

      <div className="mt-4 flex items-center justify-between gap-3 border-t border-coqueiro/10 pt-4">
        <div>
          <span className="block text-xs font-semibold uppercase tracking-[0.16em] text-coqueiro/70">Bairro</span>
          <span className="text-sm text-noite">{prestador.bairro}</span>
        </div>

        <div className="flex flex-wrap justify-end gap-2">
          <button
            type="button"
            onClick={() => setIsRemocaoAberta((estadoAtual) => !estadoAtual)}
            className="rounded-full border border-coqueiro/20 bg-white px-4 py-2 text-sm font-semibold text-coqueiro transition hover:border-coqueiro hover:bg-coqueiro hover:text-white"
          >
            {isRemocaoAberta ? 'Fechar remocao' : 'Solicitar remocao'}
          </button>

          <button
            type="button"
            onClick={() => void handleContato()}
            className="rounded-full bg-coral px-4 py-2 text-sm font-semibold text-white transition hover:bg-sol"
          >
            Chamar no WhatsApp
          </button>
        </div>
      </div>

      {isRemocaoAberta ? (
        <div className="mt-4 rounded-[1.5rem] border border-coqueiro/10 bg-areia p-4">
          <h5 className="font-display text-lg text-noite">Remover este corre</h5>
          <p className="mt-2 text-sm leading-6 text-noite/70">
            Use o mesmo e-mail informado no cadastro. Sem o codigo enviado para esse e-mail a remocao nao sera concluida.
          </p>

          {statusRemocao ? <div className={montarClasseRemocao(statusRemocao.tom)}>{statusRemocao.mensagem}</div> : null}

          {!solicitacaoRemocao ? (
            <form className="mt-4 space-y-3" onSubmit={handleSolicitarRemocao}>
              <label className="space-y-2 text-sm font-medium text-noite">
                E-mail do responsavel
                <input
                  type="email"
                  value={emailRemocao}
                  onChange={(event) => setEmailRemocao(event.target.value)}
                  placeholder="voce@email.com"
                  disabled={isSolicitandoRemocao}
                  className="w-full rounded-2xl border border-coqueiro/20 bg-white px-4 py-3 outline-none focus:border-coqueiro disabled:cursor-not-allowed disabled:opacity-70"
                />
              </label>

              <button
                type="submit"
                disabled={isSolicitandoRemocao}
                className="rounded-full bg-coqueiro px-5 py-3 text-sm font-semibold text-white transition hover:bg-mar disabled:cursor-not-allowed disabled:opacity-70"
              >
                {isSolicitandoRemocao ? 'Enviando codigo...' : 'Enviar codigo de remocao'}
              </button>
            </form>
          ) : (
            <form className="mt-4 space-y-3" onSubmit={handleConfirmarRemocao}>
              <p className="text-sm text-noite/70">Validade ate {formatarExpiracao(solicitacaoRemocao.expiraEm)}.</p>

              {solicitacaoRemocao.debugCodigo ? (
                <p className="rounded-2xl bg-noite px-4 py-3 text-sm font-semibold tracking-[0.18em] text-white">
                  Codigo local: {solicitacaoRemocao.debugCodigo}
                </p>
              ) : null}

              <label className="space-y-2 text-sm font-medium text-noite">
                Codigo de confirmacao
                <input
                  value={codigoRemocao}
                  onChange={(event) => setCodigoRemocao(event.target.value)}
                  placeholder="Digite os 6 digitos"
                  disabled={isConfirmandoRemocao}
                  className="w-full rounded-2xl border border-coqueiro/20 bg-white px-4 py-3 outline-none focus:border-coqueiro disabled:cursor-not-allowed disabled:opacity-70"
                />
              </label>

              <div className="flex flex-wrap gap-3">
                <button
                  type="submit"
                  disabled={isConfirmandoRemocao}
                  className="rounded-full bg-coral px-5 py-3 text-sm font-semibold text-white transition hover:bg-sol disabled:cursor-not-allowed disabled:opacity-70"
                >
                  {isConfirmandoRemocao ? 'Confirmando...' : 'Confirmar remocao'}
                </button>

                <button
                  type="button"
                  onClick={handleResetarRemocao}
                  className="rounded-full border border-coqueiro bg-white px-5 py-3 text-sm font-semibold text-coqueiro transition hover:bg-coqueiro hover:text-white"
                >
                  Trocar e-mail
                </button>
              </div>
            </form>
          )}
        </div>
      ) : null}
    </article>
  );
}

function montarClasseRemocao(tom: 'idle' | 'success' | 'error') {
  if (tom === 'success') {
    return 'mt-4 rounded-[1.25rem] border border-coqueiro/20 bg-coqueiro/10 p-4 text-sm text-coqueiro';
  }

  if (tom === 'error') {
    return 'mt-4 rounded-[1.25rem] border border-coral/20 bg-coral/10 p-4 text-sm text-coral';
  }

  return 'mt-4 rounded-[1.25rem] border border-coqueiro/10 bg-white p-4 text-sm text-noite/70';
}

function formatarExpiracao(expiraEm: string) {
  const data = new Date(expiraEm);
  if (Number.isNaN(data.getTime())) {
    return 'alguns minutos';
  }

  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  }).format(data);
}

function extrairMensagemErro(error: unknown, fallback: string) {
  if (error instanceof Error && error.message.trim()) {
    return error.message;
  }

  return fallback;
}

function descreverPermissao(statusPermissao: 'idle' | 'allowed' | 'blocked' | 'unsupported') {
  if (statusPermissao === 'allowed') {
    return 'localizacao liberada';
  }

  if (statusPermissao === 'blocked') {
    return 'usuario bloqueou a geolocalizacao';
  }

  if (statusPermissao === 'unsupported') {
    return 'navegador sem suporte';
  }

  return 'aguardando permissao';
}