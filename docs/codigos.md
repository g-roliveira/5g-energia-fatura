200
{
  "code": 200,
  "code_message": "A requisição foi processada com sucesso.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.21-20240912103431",
    "product": "Consultas",
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cnpj_md5": "e379c59a3573be16480533ae3f53ea1b",
      "login_senha_md5": "e379c59a3573be16480533ae3f53ea1b",
      "uc": "111111111"
    },
    "client_name": "Minha Empresa",
    "token_name": "Token de Produção",
    "billable": true,
    "price": "0.3",
    "requested_at": "2024-09-12T11:57:12.000-03:00",
    "elapsed_time_in_milliseconds": 522,
    "remote_ip": "111.111.111.111",
    "signature": "U2FsdGVkX19jM3YCmonq4K/+Np1u06sCmAGfaICeP8PMx1O11Hlytvvcv/lXU6n0ALMGDyDsv6qmVpqCLQF9Fw=="
  },
  "data_count": 1,
  "data": [
    {
      "data_vencimento": "11/11/1111",
      "nome": "Exemplo de Nome",
      "normalizado_valor": 95.37,
      "ocr":
        {
          "mes": 9,
          "ano": 2024,
          "valor": "R$ 95,37",
          "normalizado_valor": 95.37,
          "vencimento": "06/10/2024",
          "leitura_anterior_data": "11/11/1111",
          "leitura_data": "11/11/1111",
          "leitura_proxima_data": "11/11/1111",
          "emissao_data": "11/11/1111",
          "controle_n": "11-11111111111111-11",
          "numero_dias": 34,
          "codigo_barras": "123456789012345678901234567890123456789012345678",
          "aviso": "Exemplo de Texto",
          "nota_fiscal": {
            "numero_serie": "Nº 111111111 Série X",
            "apresentacao_data": "11/11/1111"
          },
          "cliente": {
            "codigo": "111",
            "cpf": "123.456.789-01",
            "cnpj": "12.345.678/9012-34",
            "nome": "Exemplo de Nome",
            "classificacao": "RURAL-TRIFASICO",
            "tensao_nominal": "220/127",
            "limites_tensao": "116 a 133 / 201 a 231",
            "endereco": "Avenida Paulista, 1636. São Paulo. SP. Brasil."
          },
          "consumo": {
            "medidor": "1111111",
            "constante": "40",
            "leitura_anterior": null,
            "leitura_atual": "8"
          },
          "energia": {
            "historico_consumo": [
              {
                "periodo": "SET/24",
                "kwh": "100.0"
              },
              {
                "periodo": "AGO/24",
                "kwh": "0.0"
              },
              {
                "periodo": "JUL/24",
                "kwh": "0.0"
              }
            ]
          },
          "composicao_fornecimento": {
            "energia": "R$ 22,27",
            "normalizado_energia": 22.27,
            "encargos": "",
            "normalizado_encargos": 21.25,
            "distribuicao": "R$ 17,38",
            "normalizado_distribuicao": 17.38,
            "tributos": "R$ 20,21",
            "normalizado_tributos": 20.21,
            "transmissao": "R$ 9,49",
            "normalizado_transmissao": 9.49,
            "perdas": "R$ 4,77",
            "normalizado_perdas": 4.77
          },
          "informacoes_gerais": "",
          "itens_fatura": [
            {
              "codigo": "111",
              "descricao": "Exemplo de descrição",
              "quantidade": "100,00",
              "quantidade_residual": "",
              "quantidade_faturada": "",
              "tarifa": "0,307900",
              "valor": "30,79",
              "base_icms": "39,07",
              "aliq_icms": "18,00%",
              "icms": "7,03",
              "valor_total": "37,82"
            },
            {
              "codigo": "111",
              "descricao": "Exemplo de descrição",
              "quantidade": "100,00",
              "quantidade_residual": "",
              "quantidade_faturada": "",
              "tarifa": "0,434500",
              "valor": "43,45",
              "base_icms": "55,14",
              "aliq_icms": "18,00%",
              "icms": "9,92",
              "valor_total": "53,37"
            },
            {
              "codigo": "111",
              "descricao": "Exemplo de descrição",
              "quantidade": "100,00",
              "quantidade_residual": "",
              "quantidade_faturada": "",
              "tarifa": "0,009200",
              "valor": "0,92",
              "base_icms": "1,16",
              "aliq_icms": "18,00%",
              "icms": "0,20",
              "valor_total": "1,12"
            },
            {
              "codigo": "111",
              "descricao": "Exemplo de descrição",
              "quantidade": "",
              "quantidade_residual": "",
              "quantidade_faturada": "",
              "tarifa": "",
              "valor": "",
              "base_icms": "78,20",
              "aliq_icms": "3,21%",
              "icms": "",
              "valor_total": "2,51"
            },
            {
              "codigo": "111",
              "descricao": "Exemplo de descrição",
              "quantidade": "",
              "quantidade_residual": "",
              "quantidade_faturada": "",
              "tarifa": "",
              "valor": "",
              "base_icms": "78,20",
              "aliq_icms": "0,70%",
              "icms": "",
              "valor_total": "0,55"
            },
            {
              "codigo": "111",
              "descricao": "Exemplo de descrição",
              "quantidade": "",
              "quantidade_residual": "",
              "quantidade_faturada": "",
              "tarifa": "",
              "valor": "75,16",
              "base_icms": "",
              "aliq_icms": "",
              "icms": "17,15",
              "valor_total": "95,37"
            }
          ]
        }
      ,
      "uc": "111111111",
      "valor": "R$ 95,37",
      "site_receipt": "https://www.exemplo.com/exemplo-de-url"
    }
  ],
  "site_receipts": [
    "https://www.exemplo.com/exemplo-de-url"
  ]
}

{
  "code": 600,
  "code_message": "Um erro inesperado ocorreu e será analisado.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": false,
    "price": "0.0",
    "requested_at": "2026-04-15T23:36:08.356-03:00",
    "elapsed_time_in_milliseconds": 1719,
    "remote_ip": null,
    "signature": "YVgx6GLPrw8U19kSWjIE4kaElPQyE5Z5cUQuenX6kCikjfstc13NrnUJHF8eIlAz0PRsGhqzyFJk4erNFBnDJy33NDk8wl0NnhrO8Q=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 601,
  "code_message": "Não foi possível se autenticar com o token informado.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": false,
    "price": "0.0",
    "requested_at": "2026-04-15T23:36:08.392-03:00",
    "elapsed_time_in_milliseconds": 1223,
    "remote_ip": null,
    "signature": "NYzYbZ8mXbKn3QTxmE9y0TAd2fEZM9y8kmexfvf0lPtl87+b1lCmjINGVTZtSREL0sc0DS6B5pfZrRZnlngNQ/uR8jNaqK/3DShxuQ=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 602,
  "code_message": "O serviço informado na URL não é válido.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": false,
    "price": "0.0",
    "requested_at": "2026-04-15T23:36:08.427-03:00",
    "elapsed_time_in_milliseconds": 1048,
    "remote_ip": null,
    "signature": "WaXL+3jOIatk/8zwinbA5D7js+V5ZllRcJBIv0Loe1oNqf8FEnsfkkb/CNt6Y3wFWVLF3csFIgAhb/GeTEM8cT1FroP1bE/7x6asbQ=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 603,
  "code_message": "O token informado não tem autorização de acesso ao serviço. Verifique se ele continua ativo e se ele não possui algum limite de uso especificado.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": false,
    "price": "0.0",
    "requested_at": "2026-04-15T23:36:08.460-03:00",
    "elapsed_time_in_milliseconds": 1824,
    "remote_ip": null,
    "signature": "1nJeAI3DbRyRkQEHcdtJeDQKAD1fXFqfDg7iNcgh3K1kZVvdSqAWcNfPFGU1pXRCc1x8EgCx2leCO6GG9cnEw/EwzSkbKCYaMLIY7w=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 604,
  "code_message": "A consulta não foi validada antes de pesquisar a fonte de origem.",
  "errors": [
    "timeout não é um número"
  ],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": false,
    "price": "0.0",
    "requested_at": "2026-04-15T23:36:08.493-03:00",
    "elapsed_time_in_milliseconds": 1378,
    "remote_ip": null,
    "signature": "bScv8HkVky7IzlMPqd3o2E4zSblO90m4WtTf1jdGTZWcTgtJPNkrgscY376Atn7SyVbfTZtJfd1mxOeqcnCZqwCse0Y2Umshfgcu3A=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 605,
  "code_message": "A consulta não foi realizada dentro do tempo de limite de timeout especificado.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": false,
    "price": "0.0",
    "requested_at": "2026-04-15T23:36:08.534-03:00",
    "elapsed_time_in_milliseconds": 35000,
    "remote_ip": null,
    "signature": "6Ql3QsZW0kEg0yY2UGnqM/f5t49b+hwvC2tTCiZfWL18Jt6ZVEwfSTeVdq55jy8ZeE+5y1mzJJDo135trZDjHGQag+1gqUO9/xX4DQ=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 606,
  "code_message": "Parâmetros obrigatórios não foram enviados. Por favor, verifique a documentação de uso do serviço.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": true,
    "price": null,
    "requested_at": "2026-04-15T23:36:08.567-03:00",
    "elapsed_time_in_milliseconds": 1741,
    "remote_ip": null,
    "signature": "eCVLNf5qs9Iq3qBoQC5cD621c+WrfGzxje0+qg3CQ62/X13bqth5Wa+p4KbqufIRuqe0o3tAzIL4l42Tbb51+wNrBl9JNiRuWWEqnQ=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 607,
  "code_message": "Parâmetro(s) inválido(s).",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": true,
    "price": null,
    "requested_at": "2026-04-15T23:36:08.599-03:00",
    "elapsed_time_in_milliseconds": 1149,
    "remote_ip": null,
    "signature": "CAvnLPkaU0d6URmvJhdDbsf59f3GTHf59ADHYQ0v3CWBGa7g84/egsfG9Kv7gNly4ElCg/49aomPY/n08LLb9o5mmJSDRciJlvbEpQ=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 608,
  "code_message": "Os parâmetros foram recusados pelo site ou aplicativo de origem que processou esta consulta.",
  "errors": [
    "Login Cpf não é válido"
  ],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": true,
    "price": null,
    "requested_at": "2026-04-15T23:36:08.629-03:00",
    "elapsed_time_in_milliseconds": 1031,
    "remote_ip": null,
    "signature": "OWH5n5yX1poWER6CHnLqvulSJfBWI7ydUmtiIbkaMhCUMiqaWQ62/WiorB3jRQI6S6ptigVktsS1tA5FcWKJOXLgtkukvFXoS7wp5w=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 609,
  "code_message": "Tentativas de consultar o site ou aplicativo de origem excedidas.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": false,
    "price": "0.0",
    "requested_at": "2026-04-15T23:36:08.658-03:00",
    "elapsed_time_in_milliseconds": 1917,
    "remote_ip": null,
    "signature": "YJeF3tEZjzuOVG+CWx8KO8Snmkpe/5I7MflCRreU+Cu7dbPoUEksQAtpIUlaNxm3drY0xVxQlJxRzHIuLcVZXvQxnMpzuGVuwIoaMw=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 610,
  "code_message": "Falha em resolver algum tipo de CAPTCHA.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": false,
    "price": "0.0",
    "requested_at": "2026-04-15T23:36:08.687-03:00",
    "elapsed_time_in_milliseconds": 1107,
    "remote_ip": null,
    "signature": "VeoGc8I28/d2cjaHReJ0XKSOhqq3Gt3E2a/Vdb0cAhjE4vMX7E2XgIuYxMUpKXnMb5xGB8wvAICUM2J/Yz17spRVJmyuvaxsVWJmGg=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 611,
  "code_message": "Os dados estão incompletos no site ou aplicativo de origem e não puderam ser retornados.",
  "errors": [
    "Os dados não podem ser exibidos. Entre em contato."
  ],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": true,
    "price": null,
    "requested_at": "2026-04-15T23:36:08.717-03:00",
    "elapsed_time_in_milliseconds": 1138,
    "remote_ip": null,
    "signature": "bT59UozFGDhqnwYN75czJKD5idzsRBkYKm2UIG5/3tLErHayj7x3LHbsWg+1LeGQdHAiIS/ingwJgTIlEFsFRK5Ft9miitRzwh6H4w=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 612,
  "code_message": "A consulta não retornou dados no site ou aplicativo de origem no qual a automação foi executada.",
  "errors": [
    "Nenhum registro foi localizado com os parâmetros informados"
  ],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": true,
    "price": null,
    "requested_at": "2026-04-15T23:36:08.745-03:00",
    "elapsed_time_in_milliseconds": 1669,
    "remote_ip": null,
    "signature": "d9H3GnltPxyMhixcUfWRpfWYtnoAa+uVmeiwWLyORJBMUJHKkzTGfooaZMo+CGsc2v4d9LEpQXYJlUh3nJk8T3Iq6gx3gCJIhcIHgg=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 613,
  "code_message": "A consulta foi bloqueada pelo servidor do site ou aplicativo de origem. Por favor, tente novamente.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": false,
    "price": "0.0",
    "requested_at": "2026-04-15T23:36:08.774-03:00",
    "elapsed_time_in_milliseconds": 1643,
    "remote_ip": null,
    "signature": "MI3p4up1JR9EJ8ED2Rma+fEo+A34OZeJ23Cvr7hVEJwmksNh8u0FMcianLJkdUYkZmmaZKqXx7aNJLp54DINSmFLxGoyeTp73XbUpg=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 614,
  "code_message": "Um erro inesperado com o site ou aplicativo de origem ocorreu. Por favor, tente novamente.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": false,
    "price": "0.0",
    "requested_at": "2026-04-15T23:36:08.803-03:00",
    "elapsed_time_in_milliseconds": 1546,
    "remote_ip": null,
    "signature": "AqhDXXV9KgxSI/a0HkSJ29HBzAhTk5IBkDkz3eekpbkKCod8k1Dh3hVIMa9ciCDNldv7sKh/VxpSY/ADsSFEX1WfCIZisSk05d2Amw=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 615,
  "code_message": "O site ou aplicativo de origem parece estar indisponível.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": false,
    "price": "0.0",
    "requested_at": "2026-04-15T23:36:08.831-03:00",
    "elapsed_time_in_milliseconds": 1796,
    "remote_ip": null,
    "signature": "+x5THXIqJn5SkdQ1DMelzZ/xrHy86erzpg8iJ8CiCHRFm6rM15j/DiyIz/bJwJ8t97iAQV7xX+o2WSsMNX1Nsc1VdjJdu1CdE89vfw=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 617,
  "code_message": "Contate o prestador de serviço. Há uma sobrecarga de uso do serviço.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": false,
    "price": "0.0",
    "requested_at": "2026-04-15T23:36:08.861-03:00",
    "elapsed_time_in_milliseconds": 1259,
    "remote_ip": null,
    "signature": "RrsCG/8tQ7y+uTDKnoAwe5BVZ3LX+m55RGzrHHaYw/rOXv0mdLFTw7vcTZd1zFCvI6REhxlypaxcDE+jss9Yxd8PX7UQ42iMmw5cpA=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 618,
  "code_message": "O site ou aplicativo de origem está sobrecarregado. Tente novamente em alguns instantes.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": false,
    "price": "0.0",
    "requested_at": "2026-04-15T23:36:08.891-03:00",
    "elapsed_time_in_milliseconds": 1162,
    "remote_ip": null,
    "signature": "GJPox8K+Lxli/ZzLg452DnDasqjqx3WMRpF+AjH/NNoWc5NjusFFiSXQ1shGMwhncmDaEAbfDl4sHXpucJCVLWomn2D6+nlHlSi9yA=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 619,
  "code_message": "O parâmetro enviado sofreu alterações no site ou aplicativo de origem. Verifique a alteração diretamente no site ou aplicativo de origem.",
  "errors": [
    "O registro atualizou seus valores de identificação e não pode mais ser consultado com os parâmetros informados"
  ],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": true,
    "price": null,
    "requested_at": "2026-04-15T23:36:08.921-03:00",
    "elapsed_time_in_milliseconds": 1041,
    "remote_ip": null,
    "signature": "IMdEtRQgf19AgECxwxgaehXk74tb3GdPDe0FgpWwssV2XgUU7MdpP7WLU6PTcxCClMjCUdrzHDVE89EhPQlKqV5wpu90eH1C4nn/Gg=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 620,
  "code_message": "O site ou aplicativo de origem emitiu um erro que provavelmente não mudará em breve para esta consulta. Leia-o para saber mais.",
  "errors": [
    "A consulta não pode ser realizada pelo site usando os parâmetros informados. Entre em contato com um posto de atendimento."
  ],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": true,
    "price": null,
    "requested_at": "2026-04-15T23:36:08.950-03:00",
    "elapsed_time_in_milliseconds": 1042,
    "remote_ip": null,
    "signature": "yDYfrVmwk1YAUM1Z95NwMt3TpdHCNQhbPlKWzAchDOpVp7vsIJOiisL0xIdDFOX3SV8VeGDEPbw+eJoehN9o98qnKDnEy5trpwYM+w=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 621,
  "code_message": "Houve um erro ao tentar gerar o arquivo de visualização desta requisição.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": false,
    "price": "0.0",
    "requested_at": "2026-04-15T23:36:08.978-03:00",
    "elapsed_time_in_milliseconds": 1085,
    "remote_ip": null,
    "signature": "dfWs2q+wgiC3i7SchPxGdtfEQqzokt/jgqTwIyFCF043TahyxUSwdZxuwxr9DBH3P+vsrUkZ57lZjPolW8YJMyjO8RIT8n6OIDXQmg=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}
{
  "code": 622,
  "code_message": "Parece que você está tentando realizar a mesma consulta diversas vezes seguidas. Por favor, verifique se há algum problema em sua integração. Se acredita que está tudo certo, entre em contato com o suporte.",
  "errors": [],
  "header": {
    "api_version": "v2",
    "api_version_full": "2.2.36-20260415210821",
    "product": null,
    "service": "contas/neoenergia/download-ocr",
    "parameters": {
      "login_cpf": "",
      "login_cnpj": "",
      "login_senha": "",
      "uf": "",
      "uc": "",
      "mes_ano": "",
      "acesso_imobiliaria": ""
    },
    "client_name": "Minha Empresa LTDA",
    "token_name": "desenvolvimento",
    "billable": false,
    "price": "0.0",
    "requested_at": "2026-04-15T23:36:09.006-03:00",
    "elapsed_time_in_milliseconds": 1118,
    "remote_ip": null,
    "signature": "bgs+CacYbI40PMDbpgz3eOhTPFVqOP12bjO+71UeTHnDB5D3AXPDTcjAvrsYe1hczwM3et/23ZFY71FEm/FN6iajX6BV17aB7FsSDA=="
  },
  "data_count": 0,
  "data": [],
  "site_receipts": [
    "https://api.infosimples.com/exemplo-de-url"
  ]
}