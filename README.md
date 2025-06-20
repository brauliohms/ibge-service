# API de Cidades do IBGE

## Este é um micro serviço para consulta de estados e cidades do Brasil, baseado nos dados do IBGE

### Endpoints

- `/api/v1/docs` - Documentação Swagger da API.

- `/api/v1/estados` - Retorna uma lista de estados brasileiros.

- `/api/v1/estados/{sigla}` - Retorna os dados de um estados brasileiro pelo sigla do estado.

- `/api/v1/estados/{sigla}/cidades` - Retorna uma lista de cidades de um estado específico pelo sigla do estado.

- `/api/v1/cidades/{codigo_ibge}` - Retorna os dados de uma cidade brasileira pelo código IBGE.

- `/api/v1/cidades/{codigo_tom}/tom` - Retorna os dados de uma cidade brasileira pelo código TOM.
