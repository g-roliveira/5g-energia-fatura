import { PrismaClient, TipoPessoa, ClientStatus, TipoCliente } from "@prisma/client"
import { randomUUID } from "crypto"

const db = new PrismaClient()

const ufs = ["BA", "SP", "RJ", "MG", "RS", "PR", "SC", "GO", "DF", "PE"]
const distribuidoras = ["COELBA", "CEMIG", "CPFL", "ENEL", "COPEL", "CELESC", "CELG", "CEB"]
const claseConsumo = ["Residencial", "Comercial", "Industrial", "Rural", "Poder Público"]
const tiposContrato = ["Mensal", "Anual", "Bienal", "Trimestral"]
const statusContrato = ["Ativo", "Suspenso", "Encerrado", "Em negociação"]

function cpf(n: number) {
  return `${String(n).padStart(3, "0")}.${String(n + 1).padStart(3, "0")}.${String(n + 2).padStart(3, "0")}-${String(n % 100).padStart(2, "0")}`
}

function cnpj(n: number) {
  return `${String(n).padStart(2, "0")}.${String(n + 1).padStart(3, "0")}.${String(n + 2).padStart(3, "0")}/0001-${String(n % 100).padStart(2, "0")}`
}

async function main() {
  console.log("Seeding database...")

  await db.integrationCredential.deleteMany()
  await db.consumerUnit.deleteMany()
  await db.commercialData.deleteMany()
  await db.clientAddress.deleteMany()
  await db.client.deleteMany()

  const statuses: ClientStatus[] = ["ativo", "inativo", "prospecto"]
  const tiposCliente: TipoCliente[] = ["residencial", "condominio", "empresa", "imobiliaria", "outro"]
  const tiposPessoa: TipoPessoa[] = ["PF", "PF", "PF", "PJ", "PJ"] // ~60% PF

  const clients = []
  for (let i = 0; i < 30; i++) {
    const tipoPessoa = tiposPessoa[i % tiposPessoa.length]
    const uf = ufs[i % ufs.length]
    clients.push(
      await db.client.create({
        data: {
          tipo_pessoa: tipoPessoa,
          nome_razao: tipoPessoa === "PF" ? `Cliente ${i + 1} Silva` : `Empresa ${i + 1} Ltda`,
          nome_fantasia: tipoPessoa === "PJ" ? `Fantasia ${i + 1}` : null,
          cpf_cnpj: tipoPessoa === "PF" ? cpf(i * 3 + 100) : cnpj(i * 2 + 10),
          email: `cliente${i + 1}@example.com`,
          telefone: `(71) 9${String(i + 1).padStart(4, "0")}-${String(i * 7 + 1000).padStart(4, "0")}`,
          status: statuses[i % statuses.length],
          tipo_cliente: tiposCliente[i % tiposCliente.length],
          observacoes: i % 5 === 0 ? `Observação do cliente ${i + 1}` : null,
          address: {
            create: {
              cep: `${String(40000 + i * 100).padStart(8, "0").replace(/(\d{5})(\d{3})/, "$1-$2")}`,
              logradouro: `Rua das Acácias, ${i + 1}`,
              numero: String(i * 3 + 1),
              complemento: i % 3 === 0 ? `Apto ${i + 1}` : null,
              bairro: `Bairro ${i % 10 + 1}`,
              cidade: uf === "BA" ? "Salvador" : uf === "SP" ? "São Paulo" : uf === "RJ" ? "Rio de Janeiro" : `Cidade ${i + 1}`,
              uf,
            },
          },
        },
      })
    )
  }

  // Commercial data for 3 clients
  for (let i = 0; i < 3; i++) {
    await db.commercialData.create({
      data: {
        client_id: clients[i].id,
        tipo_contrato: tiposContrato[i % tiposContrato.length],
        data_inicio: new Date(`2024-0${i + 1}-01`),
        data_fim: new Date(`2025-0${i + 1}-01`),
        status_contrato: statusContrato[i % statusContrato.length],
        observacoes_comerciais: `Contrato negociado com desconto ${i * 5}%`,
      },
    })
  }

  // 5 credentials spread across first 5 clients
  const credentials = []
  for (let i = 0; i < 5; i++) {
    credentials.push(
      await db.integrationCredential.create({
        data: {
          client_id: clients[i].id,
          label: `Credencial ${i + 1} - ${ufs[i]}`,
          documento_masked: i % 2 === 0 ? `***.${String(i + 100)}.***-**` : `**.***.***/0001-**`,
          uf: ufs[i],
          tipo_acesso: i % 3 === 0 ? "procurador" : "normal",
          go_credential_id: randomUUID(),
        },
      })
    )
  }

  // 15 UCs spread across 10 clients, 10 linked to credentials
  const ucClients = clients.slice(0, 10)
  for (let i = 0; i < 15; i++) {
    const client = ucClients[i % ucClients.length]
    const credential = i < 10 ? credentials[i % credentials.length] : null
    await db.consumerUnit.create({
      data: {
        client_id: client.id,
        uc_code: `UC${String(i + 1).padStart(6, "0")}`,
        distribuidora: distribuidoras[i % distribuidoras.length],
        apelido: `UC ${i + 1} - ${client.nome_razao.substring(0, 10)}`,
        classe_consumo: claseConsumo[i % claseConsumo.length],
        endereco_unidade: `Rua ${i + 1}, nº ${i * 2 + 1}`,
        cidade: ufs[i % ufs.length] === "BA" ? "Salvador" : `Cidade ${i + 1}`,
        uf: ufs[i % ufs.length],
        ativa: i % 7 !== 0,
        credential_id: credential?.id ?? null,
      },
    })
  }

  console.log(`Seed complete: 30 clients, 3 commercial records, 5 credentials, 15 UCs`)
}

main()
  .catch(console.error)
  .finally(() => db.$disconnect())
