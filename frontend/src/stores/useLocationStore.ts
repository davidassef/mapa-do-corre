import { create } from 'zustand';

import type { Coordenadas } from '../types/corre';

const centroFortaleza: Coordenadas = {
  latitude: -3.7327,
  longitude: -38.527,
};

type StatusPermissao = 'idle' | 'allowed' | 'blocked' | 'unsupported';

type LocationState = {
  centroFortaleza: Coordenadas;
  centroAtual: Coordenadas;
  raioMetros: number;
  isLocalizando: boolean;
  statusPermissao: StatusPermissao;
  definirCentro: (coordenadas: Coordenadas) => void;
  atualizarRaio: (raioMetros: number) => void;
  usarLocalizacaoAtual: () => Promise<void>;
};

export const useLocationStore = create<LocationState>((set) => ({
  centroFortaleza,
  centroAtual: centroFortaleza,
  raioMetros: 2500,
  isLocalizando: false,
  statusPermissao: 'idle',
  definirCentro: (coordenadas) => {
    set({ centroAtual: coordenadas });
  },
  atualizarRaio: (raioMetros) => {
    set({ raioMetros });
  },
  usarLocalizacaoAtual: async () => {
    if (!navigator.geolocation) {
      set({ statusPermissao: 'unsupported' });
      return;
    }

    set({ isLocalizando: true });

    await new Promise<void>((resolve) => {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          set({
            centroAtual: {
              latitude: position.coords.latitude,
              longitude: position.coords.longitude,
            },
            isLocalizando: false,
            statusPermissao: 'allowed',
          });
          resolve();
        },
        () => {
          set({
            isLocalizando: false,
            statusPermissao: 'blocked',
          });
          resolve();
        },
        {
          enableHighAccuracy: true,
          timeout: 10000,
        },
      );
    });
  },
}));