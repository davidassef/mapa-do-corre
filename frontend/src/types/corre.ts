export type Coordenadas = {
  latitude: number;
  longitude: number;
};

export type Prestador = {
  id: string;
  nome: string;
  categoria: string;
  descricao: string;
  whatsApp: string;
  bairro: string;
  latitude: number;
  longitude: number;
  distanciaMetros: number;
};

export type BuscarCorresInput = Coordenadas & {
  raioMetros: number;
  categoria?: string;
};

export type CadastroPrestadorInput = {
  nome: string;
  categoria: string;
  descricao: string;
  whatsApp: string;
  bairro: string;
  email: string;
  latitude: number;
  longitude: number;
};

export type ConfirmacaoCodigoInput = {
  solicitacaoId: string;
  codigo: string;
};

export type SolicitacaoCodigo = {
  solicitacaoId: string;
  expiraEm: string;
  canalEntrega: 'email' | 'debug';
  debugCodigo?: string;
};

export type SolicitacaoCodigoPayload = {
  message: string;
  solicitacao: SolicitacaoCodigo;
};

export type SolicitacaoRemocaoPrestadorInput = {
  email: string;
};

export type ResumoCategoria = {
  categoria: string;
  totalPrestadores: number;
  totalCliques: number;
};

export type ResumoImpacto = {
  totalConexoesGeradas: number;
  totalPrestadores: number;
  totalPrestadoresRemovidos: number;
  categorias: ResumoCategoria[];
};