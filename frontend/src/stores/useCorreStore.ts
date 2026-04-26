import { create } from 'zustand';

import { buscarCorres, registrarClique } from '../services/api';
import type { BuscarCorresInput, Prestador } from '../types/corre';

export type OrigemDados = 'api' | 'fallback' | 'idle';

type CorreState = {
  corres: Prestador[];
  isLoading: boolean;
  errorMessage: string | null;
  filtroCategoria: string;
  origemDados: OrigemDados;
  carregarCorres: (input: BuscarCorresInput) => Promise<void>;
  definirFiltroCategoria: (categoria: string) => void;
  registrarContato: (prestadorId: string) => Promise<void>;
};

export const useCorreStore = create<CorreState>((set) => ({
  corres: [],
  isLoading: false,
  errorMessage: null,
  filtroCategoria: '',
  origemDados: 'idle',
  carregarCorres: async (input) => {
    set({ isLoading: true, errorMessage: null });

    try {
      const payload = await buscarCorres(input);
      set({
        corres: payload.prestadores,
        origemDados: payload.origem,
        isLoading: false,
      });
    } catch {
      set({
        errorMessage: 'Nao foi possivel carregar os corres agora.',
        isLoading: false,
      });
    }
  },
  definirFiltroCategoria: (categoria) => {
    set({ filtroCategoria: categoria });
  },
  registrarContato: async (prestadorId) => {
    await registrarClique(prestadorId);
  },
}));