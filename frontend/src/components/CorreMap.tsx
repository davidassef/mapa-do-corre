import { Circle, CircleMarker, MapContainer, Popup, TileLayer, useMapEvents } from 'react-leaflet';

import type { Coordenadas, Prestador } from '../types/corre';

type CorreMapProps = {
  centro: Coordenadas;
  raioMetros?: number;
  prestadores: Prestador[];
  posicaoSelecionada?: Coordenadas | null;
  onSelecionarPosicao?: (coordenadas: Coordenadas) => void;
};

export function CorreMap({
  centro,
  raioMetros,
  prestadores,
  posicaoSelecionada = null,
  onSelecionarPosicao,
}: CorreMapProps) {
  return (
    <div className="overflow-hidden rounded-lg border border-noite/10 bg-white shadow-sm">
      <MapContainer
        center={[centro.latitude, centro.longitude]}
        zoom={14}
        scrollWheelZoom
        className="h-[380px] w-full md:h-[440px] xl:h-[500px]"
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />

        {raioMetros ? (
          <Circle
            center={[centro.latitude, centro.longitude]}
            radius={raioMetros}
            pathOptions={{ color: '#0f8b8d', fillColor: '#0f8b8d', fillOpacity: 0.12 }}
          />
        ) : null}

        {prestadores.map((prestador) => (
          <CircleMarker
            key={prestador.id}
            center={[prestador.latitude, prestador.longitude]}
            radius={10}
            pathOptions={{ color: '#e76f51', fillColor: '#f4a261', fillOpacity: 0.9 }}
          >
            <Popup>
              <div className="space-y-2 text-sm text-noite">
                <strong className="block text-base">{prestador.nome}</strong>
                <span className="block text-coqueiro">{prestador.categoria}</span>
                <p>{prestador.descricao}</p>
                <span className="block">{prestador.bairro}</span>
              </div>
            </Popup>
          </CircleMarker>
        ))}

        {posicaoSelecionada ? (
          <CircleMarker
            center={[posicaoSelecionada.latitude, posicaoSelecionada.longitude]}
            radius={12}
            pathOptions={{ color: '#155e63', fillColor: '#155e63', fillOpacity: 1 }}
          >
            <Popup>Local selecionado para o cadastro</Popup>
          </CircleMarker>
        ) : null}

        {onSelecionarPosicao ? <MapClickHandler onSelecionarPosicao={onSelecionarPosicao} /> : null}
      </MapContainer>
    </div>
  );
}

type MapClickHandlerProps = {
  onSelecionarPosicao: (coordenadas: Coordenadas) => void;
};

function MapClickHandler({ onSelecionarPosicao }: MapClickHandlerProps) {
  useMapEvents({
    click(event) {
      onSelecionarPosicao({
        latitude: event.latlng.lat,
        longitude: event.latlng.lng,
      });
    },
  });

  return null;
}