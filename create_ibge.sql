CREATE TABLE estados (
    codigo_ibge INT PRIMARY KEY,         -- Código numérico do IBGE para o estado. É a chave natural.
    nome VARCHAR(30) NOT NULL UNIQUE,           -- Nome completo do estado. Ex: "São Paulo".
    sigla CHAR(2) NOT NULL UNIQUE        -- Sigla do estado. Ex: "SP". A constraint UNIQUE garante a unicidade e acelera buscas pela sigla.
);

CREATE TABLE cidades (
    codigo_ibge INT PRIMARY KEY,         -- Código do IBGE para o município. Chave natural.
    nome VARCHAR(50) NOT NULL,          -- Nome completo da cidade.
    codigo_tom INT,              -- Código TOM (Tribunal de Contas dos Municípios), pode ser nulo se não aplicável.
    micro_regiao INT,           -- Nome da microrregião.
    regiao_imediata INT,        -- Nome da região imediata.
    estado_codigo_ibge INT NOT NULL,     -- Chave estrangeira referenciando o estado.

    -- Definindo a chave estrangeira para garantir a integridade relacional
    CONSTRAINT fk_estado
        FOREIGN KEY(estado_codigo_ibge) 
        REFERENCES estados(codigo_ibge)
);

-- Criar índice para otimizar a busca de cidades por estado.
CREATE INDEX idx_cidades_por_estado ON cidades(estado_codigo_ibge);