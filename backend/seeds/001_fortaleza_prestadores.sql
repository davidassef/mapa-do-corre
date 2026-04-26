DELETE FROM logs_cliques
WHERE prestador_id IN (
    SELECT id
    FROM prestadores
    WHERE nome = 'Teste Banco Fortaleza'
);

DELETE FROM prestadores
WHERE nome = 'Teste Banco Fortaleza';

INSERT INTO prestadores (
    id,
    nome,
    categoria,
    descricao,
    whatsapp,
    bairro,
    localizacao
)
VALUES
    (
        'a4a15fca-4e40-4a0c-95eb-06ec6fd34f01',
        'Marmita da Josi',
        'alimentacao',
        'Marmitas caseiras no almoco com entrega em Benfica, Damas e Parquelandia.',
        '5585988112233',
        'Benfica',
        ST_SetSRID(ST_MakePoint(-38.536936, -3.743269), 4326)::geography
    ),
    (
        'e13c7a55-c6fd-40d9-8da7-6a8081e7a102',
        'Seu Naldo Encanador',
        'encanador',
        'Atendimento residencial para vazamento, troca de torneira e revisao de caixa acoplada.',
        '5585988223344',
        'Montese',
        ST_SetSRID(ST_MakePoint(-38.545087, -3.767527), 4326)::geography
    ),
    (
        '32d349ab-10ab-4b6f-92f6-343b30370d03',
        'Dona Liduina Costuras',
        'costura',
        'Ajuste de roupa, barra, conserto de ziper e pequenas reformas sob medida.',
        '5585988334455',
        'Parquelandia',
        ST_SetSRID(ST_MakePoint(-38.559191, -3.744760), 4326)::geography
    ),
    (
        'c6d38b28-9658-44ee-8717-fd2bb7cc7e04',
        'Bikeboy do Centro',
        'entregas',
        'Entregas rapidas para pequenos volumes entre Centro, Jacarecanga e Benfica.',
        '5585988445566',
        'Centro',
        ST_SetSRID(ST_MakePoint(-38.526670, -3.727493), 4326)::geography
    ),
    (
        '31f2ef8d-e6d8-45bc-a78b-f6c40c978005',
        'Rita Faxina Express',
        'diarista',
        'Faxina residencial por diaria com agenda flexivel na regiao oeste de Fortaleza.',
        '5585988556677',
        'Farias Brito',
        ST_SetSRID(ST_MakePoint(-38.548907, -3.734215), 4326)::geography
    )
ON CONFLICT (id) DO UPDATE
SET
    nome = EXCLUDED.nome,
    categoria = EXCLUDED.categoria,
    descricao = EXCLUDED.descricao,
    whatsapp = EXCLUDED.whatsapp,
    bairro = EXCLUDED.bairro,
    localizacao = EXCLUDED.localizacao,
    updated_at = NOW();