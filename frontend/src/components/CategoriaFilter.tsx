type CategoriaFilterProps = {
  categorias: string[];
  categoriaAtiva: string;
  onChange: (categoria: string) => void;
};

export function CategoriaFilter({ categorias, categoriaAtiva, onChange }: CategoriaFilterProps) {
  return (
    <div className="flex flex-wrap gap-2">
      <button
        type="button"
        onClick={() => onChange('')}
        className={montarClasseBotao(categoriaAtiva === '')}
      >
        Todas
      </button>

      {categorias.map((categoria) => (
        <button
          key={categoria}
          type="button"
          onClick={() => onChange(categoria)}
          className={montarClasseBotao(categoriaAtiva === categoria)}
        >
          {categoria}
        </button>
      ))}
    </div>
  );
}

function montarClasseBotao(isAtivo: boolean) {
  if (isAtivo) {
    return 'rounded-full border border-coqueiro bg-coqueiro px-4 py-2 text-sm font-semibold text-white transition';
  }

  return 'rounded-full border border-coqueiro/30 bg-white px-4 py-2 text-sm font-semibold text-noite transition hover:border-coqueiro hover:text-coqueiro';
}