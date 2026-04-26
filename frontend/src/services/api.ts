import type {
  BuscarCorresInput,
  CadastroPrestadorInput,
  ConfirmacaoCodigoInput,
  Prestador,
  ResumoImpacto,
  SolicitacaoCodigoPayload,
  SolicitacaoRemocaoPrestadorInput,
} from '../types/corre';

const API_URL = import.meta.env.VITE_API_URL?.trim() || 'http://localhost:8080';

const corresFallback: Prestador[] = [
  {
    id: 'prestador-001',
    nome: 'Marmita da Josi',
    categoria: 'alimentacao',
    descricao: 'Marmitas caseiras para almoco em dias uteis.',
    whatsApp: '5585988112233',
    bairro: 'Benfica',
    latitude: -3.743269,
    longitude: -38.536936,
    distanciaMetros: 0,
  },
  {
    id: 'prestador-002',
    nome: 'Seu Naldo Encanador',
    categoria: 'encanador',
    descricao: 'Atendimento residencial para vazamentos e troca de torneiras.',
    whatsApp: '5585988223344',
    bairro: 'Montese',
    latitude: -3.767527,
    longitude: -38.545087,
    distanciaMetros: 0,
  },
  {
    id: 'prestador-003',
    nome: 'Dona Liduina Costuras',
    categoria: 'costura',
    descricao: 'Ajustes de roupa e pequenos consertos sob medida.',
    whatsApp: '5585988334455',
    bairro: 'Parquelandia',
    latitude: -3.74476,
    longitude: -38.559191,
    distanciaMetros: 0,
  },
  {
    id: 'prestador-004',
    nome: 'Bikeboy do Centro',
    categoria: 'entregas',
    descricao: 'Entregas rapidas para pequenos volumes no Centro e arredores.',
    whatsApp: '5585988445566',
    bairro: 'Centro',
    latitude: -3.727493,
    longitude: -38.52667,
    distanciaMetros: 0,
  },
  {
    id: 'prestador-005',
    nome: 'Rita Faxina Express',
    categoria: 'diarista',
    descricao: 'Faxina residencial por diaria com horario flexivel.',
    whatsApp: '5585988556677',
    bairro: 'Farias Brito',
    latitude: -3.734215,
    longitude: -38.548907,
    distanciaMetros: 0,
  },
];

export async function buscarCorres(input: BuscarCorresInput): Promise<{ prestadores: Prestador[]; origem: 'api' | 'fallback' }> {
  const params = new URLSearchParams({
    lat: String(input.latitude),
    lon: String(input.longitude),
    raioMetros: String(input.raioMetros),
  });

  if (input.categoria) {
    params.set('categoria', input.categoria);
  }

  try {
    const response = await fetch(`${API_URL}/corres?${params.toString()}`);
    if (!response.ok) {
      throw new Error('Falha ao buscar corres');
    }

    const payload = (await response.json()) as { prestadores: Prestador[] };

    return {
      prestadores: payload.prestadores,
      origem: 'api',
    };
  } catch {
    return {
      prestadores: filtrarFallback(input),
      origem: 'fallback',
    };
  }
}

export async function solicitarCadastroPrestador(input: CadastroPrestadorInput): Promise<SolicitacaoCodigoPayload> {
  const response = await fetch(`${API_URL}/prestadores`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  });

  if (!response.ok) {
    throw new Error(await lerMensagemErro(response, 'Falha ao cadastrar prestador'));
  }

  return (await response.json()) as SolicitacaoCodigoPayload;
}

export async function confirmarCadastroPrestador(input: ConfirmacaoCodigoInput): Promise<Prestador> {
  const response = await fetch(`${API_URL}/prestadores/confirmacoes`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  });

  if (!response.ok) {
    throw new Error(await lerMensagemErro(response, 'Falha ao confirmar cadastro'));
  }

  const payload = (await response.json()) as { prestador: Prestador };

  return payload.prestador;
}

export async function solicitarRemocaoPrestador(
  prestadorId: string,
  input: SolicitacaoRemocaoPrestadorInput,
): Promise<SolicitacaoCodigoPayload> {
  const response = await fetch(`${API_URL}/prestadores/${prestadorId}/remocao`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  });

  if (!response.ok) {
    throw new Error(await lerMensagemErro(response, 'Falha ao solicitar remocao'));
  }

  return (await response.json()) as SolicitacaoCodigoPayload;
}

export async function confirmarRemocaoPrestador(input: ConfirmacaoCodigoInput): Promise<string> {
  const response = await fetch(`${API_URL}/prestadores/remocoes/confirmacoes`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  });

  if (!response.ok) {
    throw new Error(await lerMensagemErro(response, 'Falha ao confirmar remocao'));
  }

  const payload = (await response.json()) as { message?: string };
  return payload.message?.trim() || 'Remocao confirmada com sucesso.';
}

export async function registrarClique(prestadorId: string): Promise<void> {
  try {
    await fetch(`${API_URL}/prestadores/${prestadorId}/cliques`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ origem: 'whatsapp' }),
    });
  } catch {
    return;
  }
}

export async function buscarResumoImpacto(): Promise<{ resumo: ResumoImpacto; origem: 'api' | 'fallback' }> {
  try {
    const response = await fetch(`${API_URL}/impacto/resumo`);
    if (!response.ok) {
      throw new Error('Falha ao buscar impacto');
    }

    const payload = (await response.json()) as { resumo: ResumoImpacto };

    return {
      resumo: payload.resumo,
      origem: 'api',
    };
  } catch {
    return {
      resumo: {
        totalConexoesGeradas: 17,
        totalPrestadores: corresFallback.length,
        totalPrestadoresRemovidos: 0,
        categorias: construirResumoFallback(),
      },
      origem: 'fallback',
    };
  }
}

function filtrarFallback(input: BuscarCorresInput): Prestador[] {
  const categoriaAtiva = input.categoria?.trim().toLowerCase() || '';

  return corresFallback
    .map((prestador) => ({
      ...prestador,
      distanciaMetros: Math.round(
        calcularDistanciaMetros(
          input.latitude,
          input.longitude,
          prestador.latitude,
          prestador.longitude,
        ),
      ),
    }))
    .filter((prestador) => prestador.distanciaMetros <= input.raioMetros)
    .filter((prestador) => {
      if (!categoriaAtiva) {
        return true;
      }

      return prestador.categoria.toLowerCase() === categoriaAtiva;
    })
    .sort((prestadorAtual, proximoPrestador) => prestadorAtual.distanciaMetros - proximoPrestador.distanciaMetros);
}

function construirResumoFallback() {
  return Array.from(new Set(corresFallback.map((prestador) => prestador.categoria)))
    .sort((categoriaAtual, proximaCategoria) => categoriaAtual.localeCompare(proximaCategoria))
    .map((categoria, indice) => ({
      categoria,
      totalPrestadores: corresFallback.filter((prestador) => prestador.categoria === categoria).length,
      totalCliques: 2 + indice * 3,
    }));
}

function calcularDistanciaMetros(
  latitudeOrigem: number,
  longitudeOrigem: number,
  latitudeDestino: number,
  longitudeDestino: number,
) {
  const raioTerraEmMetros = 6371000;
  const deltaLatitude = grausParaRad(latitudeDestino - latitudeOrigem);
  const deltaLongitude = grausParaRad(longitudeDestino - longitudeOrigem);
  const latitudeOrigemRad = grausParaRad(latitudeOrigem);
  const latitudeDestinoRad = grausParaRad(latitudeDestino);

  const a =
    Math.sin(deltaLatitude / 2) * Math.sin(deltaLatitude / 2) +
    Math.cos(latitudeOrigemRad) *
      Math.cos(latitudeDestinoRad) *
      Math.sin(deltaLongitude / 2) *
      Math.sin(deltaLongitude / 2);

  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));

  return raioTerraEmMetros * c;
}

function grausParaRad(graus: number) {
  return (graus * Math.PI) / 180;
}

async function lerMensagemErro(response: Response, fallback: string) {
  try {
    const payload = (await response.json()) as { message?: string };
    if (payload.message?.trim()) {
      return payload.message;
    }
  } catch {
    return fallback;
  }

  return fallback;
}