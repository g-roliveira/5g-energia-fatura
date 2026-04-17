# Neoenergia Private API

- gerado em: `2026-04-16T22:21:11.618548`
- cliente: `Paula Pereira Fernandes`
- documento: `*******7586`
- artefatos: `docs/neoenergia-private-api/live_20260416_222021`

## Fluxo observado

1. Login no frontend da Agência Virtual.
2. Captura do Bearer token no `localStorage` após autenticação.
3. Consumo dos endpoints privados com `Authorization: Bearer ...`.
4. Obtenção de protocolo e chamadas por UC para conta, faturas, histórico e PDF.

## Seleção de estado

Nenhum endpoint dedicado de seleção de estado foi observado no fluxo capturado. Após o login, as chamadas já seguem com `distribuidora=COELBA`, o que indica que a escolha de estado é resolvida no frontend ou por configuração local de sessão.

## Unidades consumidoras observadas

- `007098175908` | status `LIGADA` | endereço `RUA BAHIA, 130` | município `LAPAO`
- `007085489032` | status `DESLIGADA` | endereço `RUA BAHIA, 130` | município `LAPAO`

## Endpoints documentados

### `GET /areanaologada/2.0.0/lista-conteudo-servico`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "distribuidora-cms": "5",
  "distribuidora": "NEOENERGIA",
  "service_name": "home-anl"
}
```
response body:
```json
[
  {
    "contentId": 61,
    "contentName": "destaque-anl-link-03",
    "attributeName": "Texto",
    "contentVersion": {
      "contentVersionId": 34,
      "url": "/content/text/destaque-anl-link-03",
      "message": "https://www.neoenergia.com/economia-circular",
      "dtInicial": "20-05-2025 16:51:35",
      "dtFinal": null,
      "file": null
    }
  },
  {
    "contentId": 57,
    "contentName": "destaque-anl-link-02",
    "attributeName": "Texto",
    "contentVersion": {
      "contentVersionId": 33,
      "url": "/content/text/destaque-anl-link-02",
      "message": "https://www.neoenergia.com/seguranca-na-internet",
      "dtInicial": "20-05-2025 16:51:08",
      "dtFinal": null,
      "file": null
    }
  },
  {
    "contentId": 49,
    "contentName": "destaque-anl-texto-01",
    "attributeName": "Texto",
    "contentVersion": {
      "contentVersionId": 23,
      "url": "/content/text/conteudo-destaque1",
      "message": "<span style=\"font-size:12.0pt;line-height:115%;\r\nfont-family:&quot;Aptos&quot;,sans-serif;mso-ascii-theme-font:minor-latin;mso-fareast-font-family:\r\nAptos;mso-fareast-theme-font:minor-latin;mso-hansi-theme-font:minor-latin;\r\nmso-bidi-font-family:&quot;Times New Roman&quot;;mso-bidi-theme-font:minor-bidi;\r\nmso-ansi-language:PT-BR;mso-fareast-language:EN-US;mso-bidi-language:AR-SA\">Cadastre-se\r\ne tenha acesso a todos os serviços da Agência Virtual de maneira simples pelo\r\nnosso site. Veja como essa novidade pode facilitar os seus pagamentos e\r\ncontrole do seu consumo mensal.</span>",
      "dtInicial": "19-05-2025 18:17:17",
      "dtFinal": null,
      "file": null
    }
  },
  {
    "contentId": 47,
    "contentName": "destaque-anl-titulo-01",
    "attributeName": "Texto",
    "contentVersion": {
      "contentVersionId": 21,
      "url": "/content/text/titulo-destaque1",
      "message": "Novo login, como fazer?",
      "dtInicial": "15-05-2025 11:14:36",
      "dtFinal": null,
      "file": null
    }
  },
  {
    "contentId": 53,
    "contentName": "banner-anl-titulo-01",
    "attributeName": "Texto",
    "contentVersion": {
      "contentVersionId": 31,
      "url": "/content/text/banner-anl-titulo-01",
      "message": "Mais por você",
      "dtInicial": "20-05-2025 01:02:00",
      "dtFinal": null,
      "file": null
    }
  },
  {
    "contentId": 55,
    "contentName": "destaque-anl-titulo-02",
    "attributeName": "Texto",
    "contentVersion": {
      "contentVersionId": 25,
      "url": "/content/text/destaque-anl-titulo-02",
      "message": "Não caia em golpes na internet",
      "dtInicial": "20-05-2025 03:58:00",
      "dtFinal": null,
      "file": null
    }
  },
  {
    "contentId": 56,
    "contentName": "destaque-anl-texto-02",
    "attributeName": "Texto",
    "contentVersion": {
      "contentVersionId": 26,
      "url": "/content/text/destaque-anl-texto-02",
      "message": "<p class=\"MsoNormal\">Conheça os principais golpes na Internet e aprenda cinco\r\ndicas práticas para não cair em fraudes digitais.<o:p></o:p></p>",
      "dtInicial": "19-05-2025 18:58:52",
      "dtFinal": null,
      "file": null
    }
  },
  {
    "contentId": 58,
    "contentName": "destaque-anl-banner-03",
    "attributeName": "Imagem",
    "contentVersion": {
      "contentVersionId": 28,
      "url": "/content/image/destaque-anl-banner-03",
      "message": "",
      "dtInicial": "19-05-2025 19:00:13",
      "dtFinal": null,
      "file": "/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAUEBAUEAwUFBAUGBgUGCA4JCAcHCBEMDQoOFBEVFBMRExMWGB8bFhceFxMTGyUcHiAhIyMjFRomKSYiKR8iIyL/2wBDAQYGBggHCBAJCRAiFhMWIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiL/wgARCALYBLADAREAAhEBAxEB/8QAHQAAAQUBAQEBAAAAAAAAAAAAAwECBAUGAAcICf/EABoBAAMBAQEBAAAAAAAAAAAAAAABAgMEBQb/2gAMAwEAAhADEAAAAPRshQQai4fC5nI4FGjXJo1wIDRINAahrGgiGg0GiaxAYJoxNBCIzG3GOuNJnpu8NrPO4NTQ7YRqI4o7lkXKw6ZeHQbn3lc/RdY6Wyk4laLcG0lakbTGhNDYypGJrGORtLaNcFpMaiiq6K4UcChf51fhPCeOycy6D0FacxwlDgUFDgUOBQ4OBQ4OFw+BA4XD4ABRZXm8NWiNRNpSRPY8EBoICoaOPLHLcMcuMK/3yn3L2jCqMHXc2llasdEdBUyyy0iWrjpi22k7OAaEZyGg4CsQFaYmCTgkUPZ8+Zrg4ahwcHC4ODmcmgcJGICA0GggNBoNQ1jRNBEMY0GAMANRRZmjM6RcQ9zz73OW1fU5Pr5RAEQhTuPut+Hum8uws9QzpMJsCZTRUEachGNBoMoYkOhoNaa0y0+4kaRI0gbUcIAQGorRam1Jt3M5qwHZOZVEho1IjTmlBzFQocCgocHAocHBwcCBwcJR8KOOhyvNZaKiYydSI01DwcxARPhNGxMEvkxpxktBvlLqQyCgZNMzocqZSIER3Em4VFO2m0ytouNEeUmdVkPorQJaBxqOrKfak0nAjShw1D57zXAo+Dg4XD4FGguBA4EBGmggNBAaNgmg0GiaJoNGwGMGgDkIVTda0UWjx3usda65yfTysTE0IJnN12PH2dnddcw9ZBeTglw7vK7aLdNNhsTRUipJfTSJqn0vm1Eak+ksjW2IA2Np1SaolaZFuZlxZuJ9KVcS9ImaQdojHCcCsVHAocChwcHAocCBwICgggDocrzeOpWprJtSQBJuBzSj4TR8hqAqmICqiov9sTNUXNULMstXJbhQ0gqafmHS8ntES5jOWEdSLIwAUotICmTN7bm6fVObX0lz6H3c1jpBWOEg2B4BBwuGocHC4fC4fBwILgRjQQEBoNBgNQjGg0TQaDGMQIQxiaE3HJijETcY7XWW1bU5Lp5hAOpGm1DVcdpgLNBpCcGT1HN0aXHWQmNAVUcdUVR1o5ElTYyTpOlon01ypYaKuVNVKjmCGFONaVzayrUU7TK86ua+6cLHSDg9pwKCgocHAocHAocHBwIHBwcEdFDnecw0kUWFKU04EBwOBAaNqEQoBTFLjqoiL7fKBL8/5dKVufqpjYJK6zzfrww/RgoGQVBQemgNaVMbgLQRKFlnrqefo9Yl+2dEb28pDTUIHgUiggKHAglGoIjmILmINBIDQ
```

### `GET /imoveis/1.1.0/clientes/03021937586/ucs`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "documento": "*******7586",
  "canalSolicitante": "AGC",
  "distribuidora": "COELBA",
  "usuario": "***",
  "indMaisUcs": "X",
  "protocolo": "123",
  "opcaoSSOS": "S",
  "tipoPerfil": "1"
}
```
response body:
```json
{
  "ucs": [
    {
      "status": "LIGADA",
      "uc": "007098175908",
      "nomeCliente": "PAULA PEREIRA FERNANDES",
      "instalacao": "0009043917",
      "local": {
        "endereco": "RUA BAHIA, 130",
        "bairro": "CENTRO-LAPAO",
        "municipio": "LAPAO",
        "cep": "44905-000",
        "uf": "BA"
      },
      "enderecoEntrega": {
        "endereco": "RUA BAHIA 130",
        "bairro": "CENTRO-LAPAO",
        "municipio": "LAPAO",
        "cep": "44905-000",
        "uf": "BA"
      },
      "isGrupo": null,
      "nomeGrupo": null,
      "grupoTensao": "B",
      "bOptante": null,
      "contrato": "0023513361",
      "dt_inicio": "2026-02-20",
      "dt_fim": "9999-12-31",
      "medidor": null,
      "indMaisUcs": null,
      "ucColetiva": null,
      "indCCColetiva": null
    },
    {
      "status": "DESLIGADA",
      "uc": "007085489032",
      "nomeCliente": "PAULA PEREIRA FERNANDES",
      "instalacao": "0009043917",
      "local": {
        "endereco": "RUA BAHIA, 130",
        "bairro": "CENTRO-LAPAO",
        "municipio": "LAPAO",
        "cep": "44905-000",
        "uf": "BA"
      },
      "enderecoEntrega": {
        "endereco": "RUA BAHIA 130",
        "bairro": "CENTRO-LAPAO",
        "municipio": "LAPAO",
        "cep": "44905-000",
        "uf": "BA"
      },
      "isGrupo": null,
      "nomeGrupo": null,
      "grupoTensao": "B",
      "bOptante": null,
      "contrato": "0022215892",
      "dt_inicio": "2024-07-31",
      "dt_fim": "2024-12-19",
      "medidor": null,
      "indMaisUcs": null,
      "ucColetiva": null,
      "indCCColetiva": null
    }
  ],
  "e_resultado": "X",
  "retorno": {
    "id": "ZATCWS",
    "tipo": "S",
    "numero": "029",
    "mensagem": "Executado com sucesso."
  }
}
```

### `GET /multilogin/2.0.0/agv/cliente/03021937586/COELBA/grupo-de-cliente`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "tipoPerfil": "0"
}
```
response body:
```json
{
  "dados": [
    {
      "funcionalidades": [
        {
          "codigo": "F035"
        },
        {
          "codigo": "F000"
        },
        {
          "codigo": "F041"
        },
        {
          "codigo": "F069"
        },
        {
          "codigo": "F023"
        },
        {
          "codigo": "F013"
        },
        {
          "codigo": "F088"
        },
        {
          "codigo": "F003"
        },
        {
          "codigo": "F066"
        },
        {
          "codigo": "F009"
        },
        {
          "codigo": "F060"
        },
        {
          "codigo": "F089"
        },
        {
          "codigo": "F021"
        },
        {
          "codigo": "F062"
        },
        {
          "codigo": "F087"
        },
        {
          "codigo": "F027"
        },
        {
          "codigo": "F005"
        },
        {
          "codigo": "F059"
        },
        {
          "codigo": "F065"
        },
        {
          "codigo": "F075"
        },
        {
          "codigo": "F067"
        },
        {
          "codigo": "F015"
        },
        {
          "codigo": "F011"
        },
        {
          "codigo": "F007"
        },
        {
          "codigo": "F063"
        },
        {
          "codigo": "F076"
        },
        {
          "codigo": "F057"
        },
        {
          "codigo": "F081"
        },
        {
          "codigo": "F086"
        },
        {
          "codigo": "F072"
        },
        {
          "codigo": "F022"
        },
        {
          "codigo": "F002"
        },
        {
          "codigo": "F080"
        },
        {
          "codigo": "F058"
        },
        {
          "codigo": "F056"
        },
        {
          "codigo": "F061"
        },
        {
          "codigo": "F085"
        },
        {
          "codigo": "F068"
        },
        {
          "codigo": "F004"
        },
        {
          "codigo": "F064"
        },
        {
          "codigo": "F012"
        },
        {
          "codigo": "F006"
        },
        {
          "codigo": "F008"
        },
        {
          "codigo": "F001"
        },
        {
          "codigo": "F055"
        }
      ],
      "codigo": 1,
      "nome": "Comum",
      "ativo": true
    }
  ],
  "retorno": {
    "tipo": "Backoffice MultiLogin",
    "id": 200,
    "numero": "200_OK",
    "mensagem": "Sucesso"
  }
}
```

### `GET /multilogin/2.0.0/servicos/datacerta/ucs/007085489032/datacerta`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "codigo": "007085489032",
  "canalSolicitante": "AGC",
  "usuario": "***",
  "operacao": "CON",
  "tipoPerfil": "1",
  "documentoSolicitante": "***",
  "distribuidora": "COELBA"
}
```
response body:
```json
{
  "e_resultado": "X",
  "dataAtual": "Data Normal Regulada",
  "retorno": {
    "id": "ZATCWS",
    "tipo": "S",
    "numero": "029",
    "mensagem": "Executado com sucesso."
  }
}
```

### `GET /multilogin/2.0.0/servicos/datacerta/ucs/007098175908/datacerta`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "codigo": "007098175908",
  "canalSolicitante": "AGC",
  "usuario": "***",
  "operacao": "CON",
  "tipoPerfil": "1",
  "documentoSolicitante": "***",
  "distribuidora": "COELBA"
}
```
response body:
```json
{
  "e_resultado": "X",
  "possuiDataBoa": "X",
  "dataAtual": "06",
  "retorno": {
    "id": "ZATCWS",
    "tipo": "S",
    "numero": "029",
    "mensagem": "Executado com sucesso."
  }
}
```

### `GET /multilogin/2.0.0/servicos/debito-automatico/conta-cadastrada-debito`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "codigo": "007098175908",
  "codCliente": "1014945628",
  "canalSolicitante": "AGC",
  "usuario": "***",
  "valida": "",
  "distribuidora": "COELBA",
  "tipoPerfil": "1",
  "documentoSolicitante": "***"
}
```
response body:
```json
{
  "e_resultado": "X",
  "retorno": {
    "id": "ZATCWS",
    "tipo": "S",
    "numero": "145",
    "mensagem": "UC não possui débito automático cadastrado"
  }
}
```

### `GET /multilogin/2.0.0/servicos/fatura-digital/ucs/fatura-digital`
- chamadas observadas: `2`
- status HTTP: `[200, 400]`
```json
{
  "codigo": "007098175908",
  "canalSolicitante": "AGC",
  "usuario": "***",
  "distribuidora": "COELBA",
  "tipoPerfil": "1"
}
```
response body:
```json
{
  "e_resultado": "X",
  "PossuiFaturaDigital": null,
  "emailFatura": null,
  "emailCadastro": "paulinha_fernandes7@hotmail.com",
  "dominioWhatsapp": null,
  "dominioSMS": null,
  "faturaBraile": null,
  "faturaEntregaAlternativa": null,
  "retorno": {
    "tipo": "S",
    "id": "ZATCWS",
    "numero": "029",
    "mensagem": "Executado com sucesso."
  }
}
```

### `GET /multilogin/2.0.0/servicos/faturas/334075735546/dados-pagamento`
- chamadas observadas: `1`
- status HTTP: `[400]`
```json
{
  "codigo": "007085489032",
  "protocolo": "20260416280227754",
  "usuario": "***",
  "canalSolicitante": "AGC",
  "distribuidora": "COELBA",
  "regiao": "NE",
  "tipoPerfil": "1",
  "byPassActiv": "X",
  "documentoSolicitante": "*******7586",
  "documento": "*******7586"
}
```
response body:
```json
{
  "retorno": {
    "tipo": "E",
    "id": "ZATCWS",
    "numero": "448",
    "mensagem": "Fatura vencida a mais de 365 dias. Favor procurar demais canais."
  }
}
```

### `GET /multilogin/2.0.0/servicos/faturas/334075735546/pdf`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "codigo": "007085489032",
  "protocolo": "20260416280227754",
  "tipificacao": "1031602",
  "usuario": "***",
  "canalSolicitante": "AGC",
  "motivo": "02",
  "distribuidora": "COELBA",
  "regiao": "NE",
  "tipoPerfil": "1",
  "documento": "*******7586",
  "documentoSolicitante": "*******7586",
  "byPassActiv": ""
}
```
response body:
```json
{
  "fileName": "007085489032",
  "fileSize": "51724 ",
  "fileData": "<base64:68968 chars>",
  "fileExtension": ".pdf",
  "retorno": {
    "tipo": "S",
    "id": "ZATCWS",
    "numero": "066",
    "mensagem": "Documento gerado com sucesso.",
    "e_resultado": "X"
  }
}
```

### `GET /multilogin/2.0.0/servicos/faturas/339800707843/dados-pagamento`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "codigo": "007098175908",
  "protocolo": "20260416280227695",
  "usuario": "***",
  "canalSolicitante": "AGC",
  "distribuidora": "COELBA",
  "regiao": "NE",
  "tipoPerfil": "1",
  "byPassActiv": "X",
  "documentoSolicitante": "*******7586",
  "documento": "*******7586"
}
```
response body:
```json
{
  "codBarras": "838700000052215300300078098175908215056066285939",
  "retorno": {
    "e_resultado": "X"
  }
}
```

### `GET /multilogin/2.0.0/servicos/faturas/339800707843/pdf`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "codigo": "007098175908",
  "protocolo": "20260416280227695",
  "tipificacao": "1031602",
  "usuario": "***",
  "canalSolicitante": "AGC",
  "motivo": "02",
  "distribuidora": "COELBA",
  "regiao": "NE",
  "tipoPerfil": "1",
  "documento": "*******7586",
  "documentoSolicitante": "*******7586",
  "byPassActiv": ""
}
```
response body:
```json
{
  "fileName": "007098175908",
  "fileSize": "62657 ",
  "fileData": "<base64:83544 chars>",
  "fileExtension": ".pdf",
  "retorno": {
    "tipo": "S",
    "id": "ZATCWS",
    "numero": "066",
    "mensagem": "Documento gerado com sucesso.",
    "e_resultado": "X"
  }
}
```

### `GET /multilogin/2.0.0/servicos/faturas/lista-motivo-segundavia`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "usuario": "***",
  "canalSolicitante": "AGC",
  "distribuidora": "COELBA",
  "regiao": "NE",
  "tipoPerfil": "1",
  "documentoSolicitante": "*******7586",
  "codigo": "007098175908"
}
```
response body:
```json
{
  "motivos": [
    {
      "idMotivo": "02",
      "descricao": "NÃO RECEBI - CLIENTE"
    },
    {
      "idMotivo": "03",
      "descricao": "FATURA DANIFICADA - CLIENTE"
    },
    {
      "idMotivo": "04",
      "descricao": "COMPROVAR RESIDÊNCIA - CLIENTE"
    },
    {
      "idMotivo": "07",
      "descricao": "MUDANÇA NA MODALIDADE DE PAGAMENTO"
    },
    {
      "idMotivo": "10",
      "descricao": "NÃO ESTOU COM FATURA EM MÃOS"
    }
  ],
  "retorno": {
    "e_resultado": "X"
  }
}
```

### `GET /multilogin/2.0.0/servicos/faturas/ucs/faturas`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "codigo": "007098175908",
  "documento": "*******7586",
  "canalSolicitante": "AGC",
  "usuario": "***",
  "protocolo": "20260416280227695",
  "tipificacao": "",
  "byPassActiv": "X",
  "documentoSolicitante": "*******7586",
  "documentoCliente": "*******7586",
  "distribuidora": "COELBA",
  "tipoPerfil": "1"
}
```
response body:
```json
{
  "entregaFaturas": {
    "codigoTipoEntrega": "01",
    "descricaoTipoEntrega": "Entrega no Imóvel",
    "enderecoEntrega": "RUA BAHIA 130  LAPAO CENTRO-LAPAO 44905-000",
    "codigoTipoArrecadacao": "01",
    "descricaoTipoArrecadacao": "Indefinida",
    "dataCertaValida": null,
    "dataVencimento": "06",
    "dataCorte": "0000-00-00"
  },
  "faturas": [
    {
      "tipoFatura": {
        "codigo": "PR",
        "descricao": "Periodico"
      },
      "statusFatura": "A Vencer",
      "dataCompetencia": "2026-04-11",
      "dataEmissao": "2026-04-14",
      "dataPagamento": "0000-00-00",
      "dataVencimento": "2026-05-06",
      "grupoTensao": null,
      "mesReferencia": "2026/04",
      "numeroFatura": "339800707843",
      "origemFatura": "FAT",
      "situacaoComercial": "AB",
      "tipoArrecadacao": "Indefinida",
      "tipoEntrega": "Entrega no Imóvel",
      "tipoLeitura": "OSB",
      "uc": "007098175908",
      "valorEmissao": "521.53",
      "dataInicioPeriodo": "2026-03-13",
      "dataFimPeriodo": "2026-04-11",
      "emitidoFatAgrupadora": null,
      "nroFatAgrupadora": null,
      "vencFatAgrupada": "0000-00-00",
      "valorFatAgrupada": "0.00",
      "tipoDoc": "FA",
      "codigoCm": null,
      "numeroBoletoUnico": null,
      "agrupadorContaMinima": null,
      "valorTotalCMAgrupada": "0.00",
      "aceitaPix": "X"
    },
    {
      "tipoFatura": {
        "codigo": "PR",
        "descricao": "Periodico"
      },
      "statusFatura": "Vencida",
      "dataCompetencia": "2026-03-12",
      "dataEmissao": "2026-03-13",
      "dataPagamento": "0000-00-00",
      "dataVencimento": "2026-04-06",
      "grupoTensao": null,
      "mesReferencia": "2026/03",
      "numeroFatura": "337950698598",
      "origemFatura": "FAT",
      "situacaoComercial": "AB",
      "tipoArrecadacao": "Indefinida",
      "tipoEntrega": "Entrega no Imóvel",
      "tipoLeitura": "OSB",
      "uc": "007098175908",
      "valorEmissao": "460.37",
      "dataInicioPeriodo": "2026-02-20",
      "dataFimPeriodo": "2026-03-12",
      "emitidoFatAgrupadora": null,
      "nroFatAgrupadora": null,
      "vencFatAgrupada": "0000-00-00",
      "valorFatAgrupada": "0.00",
      "tipoDoc": "FA",
      "codigoCm": null,
      "numeroBoletoUnico": null,
      "agrupadorContaMinima": null,
      "valorTotalCMAgrupada": "0.00",
      "aceitaPix": "X"
    }
  ],
  "retorno": {
    "e_resultado": "X"
  }
}
```

### `GET /multilogin/2.0.0/servicos/historicos/ucs/007085489032/consumos`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "canalSolicitante": "AGC",
  "usuario": "***",
  "dataInicioPeriodoCalc": "2021-04-18T00:00:00",
  "protocoloSonda": "20260416280227754",
  "opcaoSSOS": "N",
  "protocolo": "20260416280227754",
  "documentoSolicitante": "*******7586",
  "byPassAtiv": "X",
  "distribuidora": "COELBA",
  "tipoPerfil": "1",
  "codigo": "007085489032"
}
```
response body:
```json
{
  "historicoConsumo": [
    {
      "dataPagamento": "06/02/2025",
      "dataVencimento": "08/01/2025",
      "dataLeitura": "19/12/2024",
      "consumoKw": "    96,00",
      "mesReferencia": "12/2024",
      "numeroLeitura": " 8.877,00",
      "tipoLeitura": "01",
      "dataInicioPeriodoCalc": "14/12/2024",
      "dataFimPeriodoCalc": "19/12/2024",
      "dataProxLeitura": "10/01/2025",
      "valorFatura": "113,12",
      "situacaoFatura": "PG",
      "origem": "FAT",
      "numeroFatura": "334075735546",
      "statusFatura": "PAGA",
      "indicativoCustoDisponibilidade": null,
      "indicativoMedia": null,
      "mensagemDisponibilidadeMedia": null,
      "mediadiaria": "16.00"
    },
    {
      "dataPagamento": "20/12/2024",
      "dataVencimento": "20/12/2024",
      "dataLeitura": "13/12/2024",
      "consumoKw": "   433,00",
      "mesReferencia": "12/2024",
      "numeroLeitura": " 8.877,00",
      "tipoLeitura": "01",
      "dataInicioPeriodoCalc": "13/11/2024",
      "dataFimPeriodoCalc": "13/12/2024",
      "dataProxLeitura": "10/01/2025",
      "valorFatura": "536,56",
      "situacaoFatura": "PG",
      "origem": "FAT",
      "numeroFatura": "333775733213",
      "statusFatura": "PAGA",
      "indicativoCustoDisponibilidade": null,
      "indicativoMedia": null,
      "mensagemDisponibilidadeMedia": null,
      "mediadiaria": "13.97"
    },
    {
      "dataPagamento": "20/12/2024",
      "dataVencimento": "21/11/2024",
      "dataLeitura": "12/11/2024",
      "consumoKw": "   424,00",
      "mesReferencia": "11/2024",
      "numeroLeitura": " 8.444,00",
      "tipoLeitura": "01",
      "dataInicioPeriodoCalc": "12/10/2024",
      "dataFimPeriodoCalc": "12/11/2024",
      "dataProxLeitura": "13/12/2024",
      "valorFatura": "582,88",
      "situacaoFatura": "PG",
      "origem": "FAT",
      "numeroFatura": "336850523976",
      "statusFatura": "PAGA",
      "indicativoCustoDisponibilidade": null,
      "indicativoMedia": null,
      "mensagemDisponibilidadeMedia": null,
      "mediadiaria": "13.25"
    },
    {
      "dataPagamento": "01/11/2024",
      "dataVencimento": "18/10/2024",
      "dataLeitura": "11/10/2024",
      "consumoKw": "   378,00",
      "mesReferencia": "10/2024",
      "numeroLeitura": " 8.020,00",
      "tipoLeitura": "01",
      "dataInicioPeriodoCalc": "14/09/2024",
      "dataFimPeriodoCalc": "11/10/2024",
      "dataProxLeitura": "12/11/2024",
      "valorFatura": "497,00",
      "situacaoFatura": "PG",
      "origem": "FAT",
      "numeroFatura": "338325515562",
      "statusFatura": "PAGA",
      "indicativoCustoDisponibilidade": null,
      "indicativoMedia": null,
      "mensagemDisponibilidadeMedia": null,
      "mediadiaria": "13.50"
    },
    {
      "dataPagamento": "01/11/2024",
      "dataVencimento": "20/09/2024",
      "dataLeitura": "13/09/2024",
      "consumoKw": "   677,00",
      "mesReferencia": "09/2024",
      "numeroLeitura": " 7.642,00",
      "tipoLeitura": "01",
      "dataInicioPeriodoCalc": "31/07/2024",
      "dataFimPeriodoCalc": "13/09/2024",
      "dataProxLeitura": "11/10/2024",
      "valorFatura": "808,03",
      "situacaoFatura": "PG",
      "origem": "FAT",
      "numeroFatura": "331450740612",
      "statusFatura": "PAGA",
      "indicativoCustoDisponibilidade": null,
      "indicativoMedia": null,
      "mensagemDisponibilidadeMedia": null,
      "mediadiaria": "15.04"
    }
  ],
  "mediamensal": "    401.6",
  "e_resultado": "X",
  "retorno": null
}
```

### `GET /multilogin/2.0.0/servicos/historicos/ucs/007098175908/consumos`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "canalSolicitante": "AGC",
  "usuario": "***",
  "dataInicioPeriodoCalc": "2021-04-18T00:00:00",
  "protocoloSonda": "20260416280227695",
  "opcaoSSOS": "N",
  "protocolo": "20260416280227695",
  "documentoSolicitante": "*******7586",
  "byPassAtiv": "X",
  "distribuidora": "COELBA",
  "tipoPerfil": "1",
  "codigo": "007098175908"
}
```
response body:
```json
{
  "historicoConsumo": [
    {
      "dataPagamento": "00/00/0000",
      "dataVencimento": "06/05/2026",
      "dataLeitura": "11/04/2026",
      "consumoKw": "   418,00",
      "mesReferencia": "04/2026",
      "numeroLeitura": "15.065,00",
      "tipoLeitura": "01",
      "dataInicioPeriodoCalc": "13/03/2026",
      "dataFimPeriodoCalc": "11/04/2026",
      "dataProxLeitura": "12/05/2026",
      "valorFatura": "521,53",
      "situacaoFatura": "AB",
      "origem": "FAT",
      "numeroFatura": "339800707843",
      "statusFatura": "AVENCER",
      "indicativoCustoDisponibilidade": null,
      "indicativoMedia": null,
      "mensagemDisponibilidadeMedia": null,
      "mediadiaria": "13.93"
    },
    {
      "dataPagamento": "00/00/0000",
      "dataVencimento": "06/04/2026",
      "dataLeitura": "12/03/2026",
      "consumoKw": "   370,00",
      "mesReferencia": "03/2026",
      "numeroLeitura": "14.647,00",
      "tipoLeitura": "01",
      "dataInicioPeriodoCalc": "20/02/2026",
      "dataFimPeriodoCalc": "12/03/2026",
      "dataProxLeitura": "11/04/2026",
      "valorFatura": "460,37",
      "situacaoFatura": "AB",
      "origem": "FAT",
      "numeroFatura": "337950698598",
      "statusFatura": "VENCIDA",
      "indicativoCustoDisponibilidade": null,
      "indicativoMedia": null,
      "mensagemDisponibilidadeMedia": null,
      "mediadiaria": "17.62"
    }
  ],
  "mediamensal": "      394",
  "e_resultado": "X",
  "retorno": null
}
```

### `GET /multilogin/2.0.0/servicos/imoveis/ucs/007085489032`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "usuario": "***",
  "canalSolicitante": "AGC",
  "distribuidora": "COELBA",
  "protocolo": "123",
  "tipoPerfil": "1",
  "opcaoSSOS": "N"
}
```
response body:
```json
{
  "codigo": "007085489032",
  "instalacao": "0009043917",
  "medidor": null,
  "fase": null,
  "local": {
    "endereco": "RUA BAHIA 130",
    "bairro": "CENTRO-LAPAO",
    "codLocalidade": "100000011936",
    "localidade": "LAPAO",
    "codMunicipio": "100000006161",
    "municipio": "LAPAO",
    "cep": "44905-000",
    "localizacao": {
      "sigla": "URBANO",
      "codigo": null,
      "descricao": null
    },
    "uf": "BA",
    "tipoLogradouro": "RU",
    "nomeLogradouro": "RUA BAHIA",
    "numero": "130",
    "complementoEndereco": null
  },
  "enderecoEntrega": {
    "endereco": "RUA BAHIA 130",
    "bairro": "CENTRO-LAPAO",
    "municipio": "LAPAO",
    "cep": "44905-000",
    "uf": "BA"
  },
  "situacao": {
    "codigo": "DS",
    "descricao": "DESLIGADA",
    "dataSituacaoUC": "20241219",
    "cargaInstalada": "470.0000000"
  },
  "dataLigacao": "20240731",
  "cliente": {
    "codigo": "1014945628",
    "nome": "PAULA PEREIRA FERNANDES",
    "dataAtualizacao": "20240731",
    "documento": {
      "tipo": {
        "codigo": "CPF",
        "descricao": "CADASTRO DE PESSOA FISICA"
      },
      "numero": "03021937586"
    },
    "segundoDocumento": {
      "uf": "BA",
      "orgaoExpedidor": {
        "codigo": "SSP",
        "descricao": "Carteira identidade"
      },
      "tipo": {
        "codigo": "Z_RG",
        "descricao": "Carteira identidade"
      },
      "numero": "1386263591"
    },
    "contato": {
      "email": "paulinha_fernandes7@hotmail.com",
      "celular": {
        "ddd": null,
        "numero": null
      },
      "telefone": {
        "ddd": null,
        "numero": "74-981350825"
      }
    },
    "dataNascimento": "19940626"
  },
  "servicos": {
    "baixaRenda": null,
    "faturaEmail": null,
    "dataCerta": null,
    "debitoAutomatico": null,
    "faturaBraile": null,
    "entregaAlternativa": null,
    "debitosVencidos": null,
    "contaMinima": "X",
    "parcelamentoAbertoUc": null
  },
  "caracteristicas": {
    "grandeCliente": null,
    "irrigacao": null,
    "fotovoltaico": null,
    "vip7": null,
    "espelho": null,
    "microMiniGeracaoD": null,
    "tarifaBranca": null
  },
  "mini_micro": null,
  "utd": null,
  "carga": null,
  "e_resultado": "X"
}
```

### `GET /multilogin/2.0.0/servicos/imoveis/ucs/007098175908`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "usuario": "***",
  "canalSolicitante": "AGC",
  "distribuidora": "COELBA",
  "protocolo": "123",
  "tipoPerfil": "1",
  "opcaoSSOS": "N"
}
```
response body:
```json
{
  "codigo": "007098175908",
  "instalacao": "0009043917",
  "medidor": "000000001204300816",
  "fase": "MONOFASE",
  "local": {
    "endereco": "RUA BAHIA 130",
    "bairro": "CENTRO-LAPAO",
    "codLocalidade": "100000011936",
    "localidade": "LAPAO",
    "codMunicipio": "100000006161",
    "municipio": "LAPAO",
    "cep": "44905-000",
    "localizacao": {
      "sigla": "URBANO",
      "codigo": null,
      "descricao": null
    },
    "uf": "BA",
    "tipoLogradouro": "RU",
    "nomeLogradouro": "RUA BAHIA",
    "numero": "130",
    "complementoEndereco": null
  },
  "enderecoEntrega": {
    "endereco": "RUA BAHIA 130",
    "bairro": "CENTRO-LAPAO",
    "municipio": "LAPAO",
    "cep": "44905-000",
    "uf": "BA"
  },
  "situacao": {
    "codigo": "LG",
    "descricao": "LIGADA",
    "dataSituacaoUC": "20260220",
    "cargaInstalada": "470.0000000"
  },
  "dataLigacao": "20260220",
  "cliente": {
    "codigo": "1014945628",
    "nome": "PAULA PEREIRA FERNANDES",
    "dataAtualizacao": "20240731",
    "documento": {
      "tipo": {
        "codigo": "CPF",
        "descricao": "CADASTRO DE PESSOA FISICA"
      },
      "numero": "03021937586"
    },
    "segundoDocumento": {
      "uf": "BA",
      "orgaoExpedidor": {
        "codigo": "SSP",
        "descricao": "Carteira identidade"
      },
      "tipo": {
        "codigo": "Z_RG",
        "descricao": "Carteira identidade"
      },
      "numero": "1386263591"
    },
    "contato": {
      "email": "paulinha_fernandes7@hotmail.com",
      "celular": {
        "ddd": null,
        "numero": null
      },
      "telefone": {
        "ddd": null,
        "numero": "74-981350825"
      }
    },
    "dataNascimento": "19940626"
  },
  "servicos": {
    "baixaRenda": null,
    "faturaEmail": null,
    "dataCerta": "X",
    "debitoAutomatico": null,
    "faturaBraile": null,
    "entregaAlternativa": null,
    "debitosVencidos": null,
    "contaMinima": "X",
    "parcelamentoAbertoUc": null
  },
  "caracteristicas": {
    "grandeCliente": null,
    "irrigacao": null,
    "fotovoltaico": null,
    "vip7": null,
    "espelho": null,
    "microMiniGeracaoD": null,
    "tarifaBranca": null
  },
  "mini_micro": null,
  "utd": null,
  "carga": null,
  "e_resultado": "X"
}
```

### `GET /multilogin/2.0.0/servicos/minha-conta`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "canalSolicitante": "AGC",
  "distribuidora": "COELBA",
  "usuario": "*******7586",
  "tipoPerfil": "1",
  "documentoSolicitante": "*******7586"
}
```
response body:
```json
{
  "nome": "Paula Pereira fernandes",
  "usuarioAcesso": "Paula",
  "email": "paulinha_fernandes7@hotmail.com",
  "celular": "(74) 98135-0825",
  "dtUltimaAtualizacao": "2026-04-15T23:42:19.601-03:00",
  "retorno": {
    "tipo": "WSO2",
    "id": "200",
    "numero": "200_OK",
    "mensagem": "Sucesso!"
  }
}
```

### `GET /multilogin/2.0.0/servicos/minha-conta/minha-conta-legado`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "canalSolicitante": "AGC",
  "usuario": "*******7586",
  "usuarioSap": "WSO2_CONEXAO",
  "usuarioSonda": "WSO2_CONEXAO",
  "distribuidora": "COELBA",
  "tipoPerfil": "1",
  "documentoSolicitante": "*******7586"
}
```
response body:
```json
{
  "nomeTitular": "PAULA PEREIRA FERNANDES",
  "clienteDocumentoSecundario": {
    "tipoDocumentoSecundario": "Z_RG",
    "documentoSecundario": "1386263591"
  },
  "dtNascimento": "1994-06-26",
  "emailCadastro": "paulinha_fernandes7@hotmail.com",
  "telefoneContato": "74-981350825",
  "e_resultado": "X",
  "retorno": {
    "tipo": "S",
    "id": "ZATCWS",
    "numero": "029",
    "mensagem": "Executado com sucesso."
  }
}
```

### `GET /protocolo/1.1.0/obterProtocolo`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "distribuidora": "COEL",
  "canalSolicitante": "AGC",
  "documento": "*******7586",
  "codCliente": "007098175908",
  "recaptchaAnl": "<redacted>",
  "regiao": "NE"
}
```
response body:
```json
{
  "protocoloSalesforce": 20260416280227695,
  "protocoloSalesforceStr": "20260416280227695",
  "protocoloLegado": 20260416280227695,
  "protocoloLegadoStr": "20260416280227695",
  "retorno": {
    "e_resultado": "X"
  }
}
```

### `POST /areanaologada/2.0.0/autentica`
- chamadas observadas: `1`
- status HTTP: `[200]`
request body:
```json
{
  "usuario": "*******7586",
  "senha": "<redacted>",
  "canalSolicitante": "AGU",
  "recaptcha": "<redacted>"
}
```
response body:
```json
{
  "token": "<redacted>",
  "retorno": {
    "tipo": "WSO2",
    "id": "200",
    "numero": "200_OK",
    "mensagem": "Sucesso"
  }
}
```
