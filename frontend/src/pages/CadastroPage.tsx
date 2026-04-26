import { useState } from 'react';

import { CorreMap } from '../components/CorreMap';
import { confirmarCadastroPrestador, solicitarCadastroPrestador } from '../services/api';
import { useLocationStore } from '../stores/useLocationStore';
import type { CadastroPrestadorInput, Coordenadas, SolicitacaoCodigo } from '../types/corre';

type StatusCadastro = {
  tom: 'idle' | 'success' | 'error';
  mensagem: string;
};

const categoriaOptions = ['alimentacao', 'costura', 'diarista', 'encanador', 'entregas'];

type CadastroPageProps = {
  onAbrirMapa: () => void;
};

export function CadastroPage({ onAbrirMapa }: CadastroPageProps) {
  const centroFortaleza = useLocationStore((state) => state.centroFortaleza);
  const [posicaoSelecionada, setPosicaoSelecionada] = useState<Coordenadas | null>(centroFortaleza);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isConfirmando, setIsConfirmando] = useState(false);
  const [codigoConfirmacao, setCodigoConfirmacao] = useState('');
  const [solicitacaoAtiva, setSolicitacaoAtiva] = useState<SolicitacaoCodigo | null>(null);
  const [ultimoPrestadorCadastrado, setUltimoPrestadorCadastrado] = useState<string | null>(null);
  const [statusCadastro, setStatusCadastro] = useState<StatusCadastro>({
    tom: 'idle',
    mensagem: 'Clique no mapa, preencha os dados e valide o e-mail para publicar.',
  });
  const [formulario, setFormulario] = useState<CadastroPrestadorInput>({
    nome: '',
    categoria: 'alimentacao',
    descricao: '',
    whatsApp: '',
    bairro: '',
    email: '',
    latitude: centroFortaleza.latitude,
    longitude: centroFortaleza.longitude,
  });

  const isFluxoCodigoAtivo = solicitacaoAtiva !== null;

  function atualizarCampo<K extends keyof CadastroPrestadorInput>(campo: K, valor: CadastroPrestadorInput[K]) {
    setFormulario((estadoAtual) => ({
      ...estadoAtual,
      [campo]: valor,
    }));
  }

  function handleSelecionarPosicao(coordenadas: Coordenadas) {
    setPosicaoSelecionada(coordenadas);
    setFormulario((estadoAtual) => ({
      ...estadoAtual,
      latitude: coordenadas.latitude,
      longitude: coordenadas.longitude,
    }));
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setIsSubmitting(true);

    try {
      const payload = await solicitarCadastroPrestador(formulario);
      setSolicitacaoAtiva(payload.solicitacao);
      setCodigoConfirmacao('');
      setUltimoPrestadorCadastrado(null);
      setStatusCadastro({
        tom: 'idle',
        mensagem: payload.message,
      });

      return;
    } catch (error) {
      setStatusCadastro({
        tom: 'error',
        mensagem: extrairMensagemErro(error),
      });
      setUltimoPrestadorCadastrado(null);
    } finally {
      setIsSubmitting(false);
    }
  }

  async function handleConfirmarCadastro(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!solicitacaoAtiva) {
      return;
    }

    setIsConfirmando(true);

    try {
      const prestador = await confirmarCadastroPrestador({
        solicitacaoId: solicitacaoAtiva.solicitacaoId,
        codigo: codigoConfirmacao,
      });

      setUltimoPrestadorCadastrado(prestador.nome);
      setSolicitacaoAtiva(null);
      setCodigoConfirmacao('');
      setStatusCadastro({
        tom: 'success',
        mensagem: `Cadastro publicado com sucesso: ${prestador.nome}.`,
      });
      setFormulario({
        nome: '',
        categoria: 'alimentacao',
        descricao: '',
        whatsApp: '',
        bairro: '',
        email: '',
        latitude: posicaoSelecionada?.latitude ?? centroFortaleza.latitude,
        longitude: posicaoSelecionada?.longitude ?? centroFortaleza.longitude,
      });
    } catch (error) {
      setStatusCadastro({
        tom: 'error',
        mensagem: extrairMensagemErro(error),
      });
    } finally {
      setIsConfirmando(false);
    }
  }

  function handleEditarDados() {
    setSolicitacaoAtiva(null);
    setCodigoConfirmacao('');
    setStatusCadastro({
      tom: 'idle',
      mensagem: 'Revise os dados, solicite um novo codigo e tente novamente.',
    });
  }

  return (
    <section className="grid gap-5 xl:grid-cols-[minmax(0,0.95fr)_minmax(0,1.05fr)] xl:items-start">
      <div className="rounded-lg border border-noite/10 bg-white p-5 shadow-sm">
        <div className="space-y-2">
          <span className="inline-flex rounded-full bg-mar/15 px-3 py-1 text-xs font-semibold uppercase text-coqueiro">
            Cadastro aberto
          </span>
          <h2 className="font-display text-2xl text-noite md:text-3xl">Publique seu corre</h2>
          <p className="text-sm leading-6 text-noite/70">
            Informe os dados principais, confirme o e-mail e marque o ponto de atendimento no mapa.
          </p>
        </div>

        <form className="mt-6 space-y-4" onSubmit={handleSubmit}>
          <CampoTexto
            label="Nome do corre"
            value={formulario.nome}
            onChange={(value) => atualizarCampo('nome', value)}
            placeholder="Ex.: Rita Faxina Express"
          />

          <div className="grid gap-4 md:grid-cols-2">
            <label className="space-y-2 text-sm font-medium text-noite">
              Categoria
              <select
                value={formulario.categoria}
                onChange={(event) => atualizarCampo('categoria', event.target.value)}
                className="w-full rounded-md border border-coqueiro/20 bg-[#f8f5ef] px-3 py-3 outline-none focus:border-coqueiro"
              >
                {categoriaOptions.map((categoria) => (
                  <option key={categoria} value={categoria}>
                    {categoria}
                  </option>
                ))}
              </select>
            </label>

            <CampoTexto
              label="WhatsApp"
              value={formulario.whatsApp}
              onChange={(value) => atualizarCampo('whatsApp', value)}
              placeholder="85999998888"
              disabled={isFluxoCodigoAtivo || isSubmitting || isConfirmando}
            />
          </div>

          <CampoTexto
            label="Bairro"
            value={formulario.bairro}
            onChange={(value) => atualizarCampo('bairro', value)}
            placeholder="Benfica"
            disabled={isFluxoCodigoAtivo || isSubmitting || isConfirmando}
          />

          <CampoTexto
            label="E-mail para validacao e remocao"
            value={formulario.email}
            onChange={(value) => atualizarCampo('email', value)}
            placeholder="voce@email.com"
            type="email"
            disabled={isFluxoCodigoAtivo || isSubmitting || isConfirmando}
          />

          <label className="space-y-2 text-sm font-medium text-noite">
            Descricao curta
            <textarea
              value={formulario.descricao}
              onChange={(event) => atualizarCampo('descricao', event.target.value)}
              disabled={isFluxoCodigoAtivo || isSubmitting || isConfirmando}
              className="min-h-28 w-full rounded-md border border-coqueiro/20 bg-[#f8f5ef] px-3 py-3 outline-none focus:border-coqueiro disabled:cursor-not-allowed disabled:opacity-70"
              placeholder="Diga o que voce resolve e em que horario atende."
            />
          </label>

          <div className="grid gap-4 md:grid-cols-2">
            <CampoTexto
              label="Latitude"
              value={String(formulario.latitude)}
              onChange={(value) => atualizarCampo('latitude', Number(value))}
              disabled={isFluxoCodigoAtivo || isSubmitting || isConfirmando}
            />
            <CampoTexto
              label="Longitude"
              value={String(formulario.longitude)}
              onChange={(value) => atualizarCampo('longitude', Number(value))}
              disabled={isFluxoCodigoAtivo || isSubmitting || isConfirmando}
            />
          </div>

          <div className={montarClasseStatus(statusCadastro.tom)}>
            {statusCadastro.mensagem}
          </div>

          {solicitacaoAtiva ? (
            <div className="rounded-lg border border-coqueiro/15 bg-[#f8f5ef] p-4 text-sm text-noite/75">
              <p>
                Codigo solicitado para <strong>{formulario.email}</strong>.
              </p>
              <p className="mt-2">Validade ate {formatarExpiracao(solicitacaoAtiva.expiraEm)}.</p>
              {solicitacaoAtiva.debugCodigo ? (
                <p className="mt-3 rounded-md bg-noite px-4 py-3 font-semibold text-white">
                  Codigo local: {solicitacaoAtiva.debugCodigo}
                </p>
              ) : null}
            </div>
          ) : null}

          {ultimoPrestadorCadastrado ? (
            <button
              type="button"
              onClick={onAbrirMapa}
              className="rounded-md border border-coqueiro bg-white px-4 py-3 text-sm font-semibold text-coqueiro transition hover:bg-coqueiro hover:text-white"
            >
              Ver {ultimoPrestadorCadastrado} no mapa
            </button>
          ) : null}

          {!solicitacaoAtiva ? (
            <button
              type="submit"
              disabled={isSubmitting || isConfirmando}
              className="rounded-md bg-coral px-4 py-3 text-sm font-semibold text-white transition hover:bg-sol disabled:cursor-not-allowed disabled:opacity-70"
            >
              {isSubmitting ? 'Enviando codigo...' : 'Enviar codigo por e-mail'}
            </button>
          ) : null}
        </form>

        {solicitacaoAtiva ? (
          <form className="mt-5 space-y-4 border-t border-coqueiro/10 pt-5" onSubmit={handleConfirmarCadastro}>
            <CampoTexto
              label="Codigo de confirmacao"
              value={codigoConfirmacao}
              onChange={setCodigoConfirmacao}
              placeholder="Digite os 6 digitos"
              disabled={isConfirmando || isSubmitting}
            />

            <div className="flex flex-wrap gap-3">
              <button
                type="submit"
                disabled={isConfirmando}
                className="rounded-md bg-coqueiro px-4 py-3 text-sm font-semibold text-white transition hover:bg-mar disabled:cursor-not-allowed disabled:opacity-70"
              >
                {isConfirmando ? 'Confirmando...' : 'Confirmar publicacao'}
              </button>

              <button
                type="button"
                onClick={handleEditarDados}
                className="rounded-md border border-coqueiro bg-white px-4 py-3 text-sm font-semibold text-coqueiro transition hover:bg-coqueiro hover:text-white"
              >
                Editar dados
              </button>
            </div>
          </form>
        ) : null}
      </div>

      <div className="space-y-4">
        <CorreMap
          centro={posicaoSelecionada ?? centroFortaleza}
          prestadores={[]}
          posicaoSelecionada={posicaoSelecionada}
          onSelecionarPosicao={handleSelecionarPosicao}
        />
        <div className="rounded-lg border border-noite/10 bg-white p-5 shadow-sm">
          <span className="text-xs font-semibold uppercase text-coqueiro/70">Ponto selecionado</span>
          <h3 className="mt-2 font-display text-xl text-noite">Referencia de atendimento</h3>
          <div className="mt-4 grid gap-3 sm:grid-cols-2">
            <div className="rounded-md bg-[#f8f5ef] px-3 py-3">
              <span className="block text-xs font-semibold uppercase text-coqueiro/70">Latitude</span>
              <strong className="mt-1 block text-sm text-noite">{formulario.latitude.toFixed(6)}</strong>
            </div>
            <div className="rounded-md bg-[#f8f5ef] px-3 py-3">
              <span className="block text-xs font-semibold uppercase text-coqueiro/70">Longitude</span>
              <strong className="mt-1 block text-sm text-noite">{formulario.longitude.toFixed(6)}</strong>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

type CampoTextoProps = {
  label: string;
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  type?: 'text' | 'email';
  disabled?: boolean;
};

function CampoTexto({ label, value, onChange, placeholder, type = 'text', disabled = false }: CampoTextoProps) {
  return (
    <label className="space-y-2 text-sm font-medium text-noite">
      {label}
      <input
        type={type}
        value={value}
        onChange={(event) => onChange(event.target.value)}
        placeholder={placeholder}
        disabled={disabled}
        className="w-full rounded-md border border-coqueiro/20 bg-[#f8f5ef] px-3 py-3 outline-none focus:border-coqueiro disabled:cursor-not-allowed disabled:opacity-70"
      />
    </label>
  );
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

function montarClasseStatus(tom: StatusCadastro['tom']) {
  if (tom === 'success') {
    return 'rounded-md border border-coqueiro/20 bg-coqueiro/10 p-4 text-sm text-coqueiro';
  }

  if (tom === 'error') {
    return 'rounded-md border border-coral/20 bg-coral/10 p-4 text-sm text-coral';
  }

  return 'rounded-md border border-coqueiro/10 bg-[#f8f5ef] p-4 text-sm text-noite/70';
}

function extrairMensagemErro(error: unknown) {
  if (error instanceof Error && error.message.trim()) {
    return error.message;
  }

  return 'Nao foi possivel publicar o cadastro agora.';
}