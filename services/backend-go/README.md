# backend-go

Serviço principal do backend.

Responsabilidades:

- cliente HTTP da API privada Neoenergia
- scheduler e jobs
- sincronização de UCs e faturas
- persistência
- API pública para frontend e integração
- orquestração do `doc-extractor-py`

Estado atual:

- scaffold inicial
- sem regras de negócio migradas ainda

Próximo passo:

- implementar os endpoints internos de sync e consulta
- portar o cliente privado documentado em `docs/neoenergia-private-api`
