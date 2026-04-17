# Clientes + UCs + Sincronização/Faturas — Specification

## Problem Statement

Operadores do sistema 5g-energia precisam cadastrar clientes, vincular suas unidades consumidoras (UCs) à concessionária e monitorar faturas sincronizadas automaticamente. Hoje não existe nenhuma interface para isso — o processo é manual e desconectado. O módulo é o núcleo operacional do produto.

## Goals

- [ ] CRUD completo de clientes com endereços e dados comerciais no PostgreSQL local
- [ ] Cadastro e gestão de UCs vinculadas a clientes e credenciais de integração
- [ ] Fluxo de sincronização de UC via backend Go (sem chamada direta do browser)
- [ ] Leitura de invoices e auditoria de sync a partir de dados já persistidos no Go

## Out of Scope

| Feature | Reason |
|---------|--------|
| Importação CSV de clientes | Complexidade extra — badge "Em breve" |
| Portal do cliente final | Outro produto |
| Agendamento automático de sync | Fase 2 |
| Notificações push/email | Fase 2 |
| Multi-tenancy/isolamento por empresa | Fase 2 |
| Pagamento de faturas online | Fora do produto |

---

## User Stories

### P1: Listar e buscar clientes ⭐ MVP — CLNT-01

**User Story**: Como operador, quero ver todos os clientes em uma tabela com busca e filtros, para localizar rapidamente quem preciso atender.

**Acceptance Criteria**:
1. WHEN acesso `/clientes` THEN sistema SHALL exibir tabela com colunas: nome/razão, CPF/CNPJ, tipo, status, qtd UCs, cidade/UF, cadastrado em, ações
2. WHEN digito no campo busca THEN sistema SHALL filtrar por nome/razão social e CPF/CNPJ em tempo real (debounce 300ms)
3. WHEN seleciono filtro de status (ativo/inativo/prospecto) THEN sistema SHALL aplicar filtro via URL params
4. WHEN seleciono filtro de tipo (residencial/condomínio/empresa/imobiliária/outro) THEN sistema SHALL aplicar filtro via URL params
5. WHEN tabela carrega THEN sistema SHALL mostrar skeleton de 5 linhas durante fetch
6. WHEN não há clientes THEN sistema SHALL exibir empty state com CTA "Criar primeiro cliente"
7. WHEN clico em "Importar CSV" THEN sistema SHALL mostrar botão desabilitado com badge "Em breve"

**Independent Test**: Acessar `/clientes`, ver tabela com dados do seed, filtrar por status "ativo", ver contagem reduzida.

---

### P1: Criar e editar cliente ⭐ MVP — CLNT-02

**User Story**: Como operador, quero cadastrar um novo cliente com todos seus dados, para iniciar o relacionamento comercial.

**Acceptance Criteria**:
1. WHEN clico "Novo cliente" THEN sistema SHALL navegar para `/clientes/novo`
2. WHEN submeto formulário com campos inválidos THEN sistema SHALL exibir erros inline por campo (Zod)
3. WHEN submeto formulário válido (POST /api/clients) THEN sistema SHALL criar cliente e redirecionar para detalhe
4. WHEN acesso `/clientes/:id/editar` THEN sistema SHALL pré-preencher formulário com dados existentes
5. WHEN submeto edição válida (PATCH /api/clients/:id) THEN sistema SHALL atualizar e mostrar toast "Cliente atualizado"
6. WHEN CPF/CNPJ já existe THEN sistema SHALL exibir erro "CPF/CNPJ já cadastrado" no campo
7. WHEN `tipo_pessoa` é PF THEN sistema SHALL validar CPF (11 dígitos); PJ valida CNPJ (14 dígitos)

**Independent Test**: Criar cliente PF com CPF válido, editar nome, verificar alteração na lista.

---

### P1: Ver detalhe do cliente ⭐ MVP — CLNT-03

**User Story**: Como operador, quero ver todos os dados de um cliente em um painel organizado por abas, para ter visão completa sem navegar por várias telas.

**Acceptance Criteria**:
1. WHEN acesso `/clientes/:id` THEN sistema SHALL exibir header com nome, status badge, tipo, CPF/CNPJ
2. WHEN navego para aba "Dados cadastrais" THEN sistema SHALL exibir todos os campos preenchidos
3. WHEN navego para aba "Endereço" THEN sistema SHALL exibir dados de endereço
4. WHEN navego para aba "UCs" THEN sistema SHALL exibir lista de UCs com status
5. WHEN navego para aba "Comercial" THEN sistema SHALL exibir dados contratuais
6. WHEN navego para aba "Integração" THEN sistema SHALL exibir credenciais de integração mascaradas
7. WHEN clico "Arquivar cliente" THEN sistema SHALL exibir diálogo de confirmação antes de executar

**Independent Test**: Abrir detalhe de cliente seed, navegar por todas as abas sem erro.

---

### P1: Gerenciar UCs do cliente ⭐ MVP — UC-01

**User Story**: Como operador, quero cadastrar e ver as UCs de um cliente, para saber quais unidades estão sendo gerenciadas.

**Acceptance Criteria**:
1. WHEN clico "Adicionar UC" no detalhe do cliente THEN sistema SHALL abrir formulário de nova UC
2. WHEN submeto UC com `uc_code` duplicado THEN sistema SHALL rejeitar com erro "UC já cadastrada"
3. WHEN UC é cadastrada THEN sistema SHALL aparecer na lista de UCs do cliente
4. WHEN visualizo lista de UCs THEN sistema SHALL exibir: código, apelido, distribuidora, status local, status último sync

**Independent Test**: Adicionar UC a cliente seed, ver UC na lista com status "Pendente sync".

---

### P1: Criar credencial de integração ⭐ MVP — CRED-01

**User Story**: Como operador, quero vincular credenciais da concessionária a uma UC, para habilitar sincronização.

**Acceptance Criteria**:
1. WHEN preencho label, documento, senha, UF e tipo_acesso THEN sistema SHALL enviar para backend Go via POST /api/integration/credentials
2. WHEN credencial criada no Go THEN sistema SHALL salvar `go_credential_id` + documento mascarado no PostgreSQL local
3. WHEN exibo credencial THEN sistema SHALL nunca mostrar senha — apenas documento mascarado e label
4. WHEN erro na criação no Go THEN sistema SHALL exibir mensagem de erro e não salvar localmente

**Independent Test**: Criar credencial com CPF e senha mock, ver que local DB tem apenas go_credential_id e CPF mascarado.

---

### P1: Sincronizar UC ⭐ MVP — SYNC-01

**User Story**: Como operador, quero disparar sincronização de uma UC manualmente, para obter a última fatura sem precisar acessar o portal da concessionária.

**Acceptance Criteria**:
1. WHEN clico "Sincronizar agora" THEN sistema SHALL chamar POST /api/integration/ucs/:uc/sync
2. WHEN sync iniciado THEN sistema SHALL mostrar spinner no botão com texto "Sincronizando..."
3. WHEN sync retorna `sync_run_id` THEN sistema SHALL atualizar UI com status do sync run
4. WHEN sync conclui com sucesso THEN sistema SHALL exibir toast "Sincronização concluída" e atualizar dados da UC
5. WHEN sync falha THEN sistema SHALL exibir toast com mensagem de erro do `error_message`
6. WHEN sync em andamento THEN sistema SHALL desabilitar botão para evitar duplo disparo

**Independent Test**: Disparar sync na UC seed, observar estado "Sincronizando...", ver resultado.

---

### P1: Listar invoices da UC ⭐ MVP — INV-01

**User Story**: Como operador, quero ver o histórico de faturas sincronizadas de uma UC, para acompanhar pagamentos e vencimentos.

**Acceptance Criteria**:
1. WHEN acesso faturas de uma UC THEN sistema SHALL buscar dados de GET /api/integration/ucs/:uc/invoices (dados do Go backend, não Neoenergia direto)
2. WHEN tabela carrega THEN sistema SHALL exibir: número fatura, referência, valor, vencimento, status, completude, atualizado em
3. WHEN filtro por status aplicado THEN sistema SHALL filtrar corretamente
4. WHEN clico em fatura THEN sistema SHALL navegar para detalhe da invoice
5. WHEN sem invoices THEN sistema SHALL exibir empty state com CTA "Sincronizar UC"

**Independent Test**: Ver invoices de UC sincronizada no seed, ver pelo menos 1 linha com dados completos.

---

### P1: Ver detalhe da invoice ⭐ MVP — INV-02

**User Story**: Como operador, quero ver todos os dados de uma fatura, incluindo itens e confiabilidade da extração.

**Acceptance Criteria**:
1. WHEN acesso `/clientes/:id/ucs/:ucId/faturas/:faturaId` THEN sistema SHALL buscar GET /api/integration/invoices/:id
2. WHEN exibo invoice THEN sistema SHALL mostrar `billing_record` e `document_record` em seções separadas
3. WHEN `items` disponíveis THEN sistema SHALL listar itens da fatura em tabela
4. WHEN `completeness` é `partial` ou `failed` THEN sistema SHALL exibir badge de alerta
5. WHEN campos têm `source_map` THEN sistema SHALL indicar origem (API / PDF / OCR)

**Independent Test**: Abrir detalhe de fatura do seed, ver billing_record e items renderizados.

---

### P1: Ver auditoria de sync ⭐ MVP — SYNC-02

**User Story**: Como operador, quero ver o log de uma sincronização, para diagnosticar falhas e reprocessar se necessário.

**Acceptance Criteria**:
1. WHEN acesso detalhe de sync run THEN sistema SHALL buscar GET /api/integration/sync-runs/:id
2. WHEN exibo sync run THEN sistema SHALL mostrar: status, created_at, error_message (se houver)
3. WHEN `raw_response` disponível THEN sistema SHALL exibir em bloco colapsável (não expandido por padrão)
4. WHEN status é `failed` THEN sistema SHALL exibir CTA "Reprocessar"

**Independent Test**: Abrir sync run do seed com status "failed", ver error_message e botão reprocessar.

---

### P2: Arquivar cliente — CLNT-04

**User Story**: Como operador, quero arquivar um cliente inativo, para manter a lista limpa sem perder histórico.

**Acceptance Criteria**:
1. WHEN confirmo arquivamento THEN sistema SHALL chamar POST /api/clients/:id/archive
2. WHEN arquivado THEN `archived_at` SHALL ser preenchido e status SHALL ser `inativo`
3. WHEN listando clientes THEN clientes arquivados SHALL ser ocultados por padrão (filtro "Mostrar arquivados")

---

### P2: Editar UC — UC-02

**User Story**: Como operador, quero editar apelido e dados locais de uma UC, para manter informações atualizadas.

**Acceptance Criteria**:
1. WHEN edito UC (PATCH /api/ucs/:id) THEN sistema SHALL atualizar campos locais
2. WHEN edito `uc_code` THEN sistema SHALL barrar com erro (código é imutável após criação)

---

### P3: Filtros avançados na lista de clientes — CLNT-05

**Acceptance Criteria**:
1. WHEN combino múltiplos filtros THEN sistema SHALL aplicar como AND
2. WHEN limpo todos os filtros THEN sistema SHALL restaurar lista completa

---

## Edge Cases

- WHEN backend Go está indisponível THEN BFF SHALL retornar 503 com `{ error: "Serviço de integração indisponível" }`
- WHEN sync timeout (>30s) THEN BFF SHALL retornar 504 com mensagem adequada
- WHEN senha enviada via query string THEN BFF SHALL rejeitar request (validação Zod no body)
- WHEN `uc_code` contém caracteres especiais THEN sistema SHALL sanitizar antes de enviar ao Go
- WHEN cliente arquivado tenta adicionar UC THEN sistema SHALL bloquear com erro

---

## Requirement Traceability

| Req ID | Story | Phase | Status |
|--------|-------|-------|--------|
| CLNT-01 | Lista de clientes | Design | Pending |
| CLNT-02 | Criar/editar cliente | Design | Pending |
| CLNT-03 | Detalhe do cliente | Design | Pending |
| CLNT-04 | Arquivar cliente | Design | Pending |
| CLNT-05 | Filtros avançados | Design | Pending |
| UC-01 | Gerenciar UCs | Design | Pending |
| UC-02 | Editar UC | Design | Pending |
| CRED-01 | Criar credencial | Design | Pending |
| SYNC-01 | Sincronizar UC | Design | Pending |
| SYNC-02 | Auditoria de sync | Design | Pending |
| INV-01 | Listar invoices | Design | Pending |
| INV-02 | Detalhe invoice | Design | Pending |

---

## Success Criteria

- [ ] Operador cria cliente completo (PF ou PJ) em < 2 minutos
- [ ] Fluxo completo credencial → sync → invoice funciona end-to-end via BFF
- [ ] Zero chamadas diretas do browser ao backend Go (verificável via DevTools)
- [ ] Tabela de clientes com 100+ registros seed carrega em < 2s
- [ ] Todos os estados de erro exibem feedback legível ao usuário
