-- CreateEnum
CREATE TYPE "TipoPessoa" AS ENUM ('PF', 'PJ');

-- CreateEnum
CREATE TYPE "ClientStatus" AS ENUM ('ativo', 'inativo', 'prospecto');

-- CreateEnum
CREATE TYPE "TipoCliente" AS ENUM ('residencial', 'condominio', 'empresa', 'imobiliaria', 'outro');

-- CreateTable
CREATE TABLE "Client" (
    "id" TEXT NOT NULL,
    "tipo_pessoa" "TipoPessoa" NOT NULL,
    "nome_razao" TEXT NOT NULL,
    "nome_fantasia" TEXT,
    "cpf_cnpj" TEXT NOT NULL,
    "email" TEXT,
    "telefone" TEXT,
    "status" "ClientStatus" NOT NULL DEFAULT 'prospecto',
    "tipo_cliente" "TipoCliente" NOT NULL,
    "observacoes" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,
    "archived_at" TIMESTAMP(3),

    CONSTRAINT "Client_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "ClientAddress" (
    "id" TEXT NOT NULL,
    "client_id" TEXT NOT NULL,
    "cep" TEXT,
    "logradouro" TEXT,
    "numero" TEXT,
    "complemento" TEXT,
    "bairro" TEXT,
    "cidade" TEXT,
    "uf" CHAR(2),
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "ClientAddress_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "ConsumerUnit" (
    "id" TEXT NOT NULL,
    "client_id" TEXT NOT NULL,
    "uc_code" TEXT NOT NULL,
    "distribuidora" TEXT,
    "apelido" TEXT,
    "classe_consumo" TEXT,
    "endereco_unidade" TEXT,
    "cidade" TEXT,
    "uf" CHAR(2),
    "ativa" BOOLEAN NOT NULL DEFAULT true,
    "credential_id" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "ConsumerUnit_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "CommercialData" (
    "id" TEXT NOT NULL,
    "client_id" TEXT NOT NULL,
    "tipo_contrato" TEXT,
    "data_inicio" TIMESTAMP(3),
    "data_fim" TIMESTAMP(3),
    "status_contrato" TEXT,
    "observacoes_comerciais" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "CommercialData_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "IntegrationCredential" (
    "id" TEXT NOT NULL,
    "client_id" TEXT NOT NULL,
    "label" TEXT NOT NULL,
    "documento_masked" TEXT NOT NULL,
    "uf" CHAR(2) NOT NULL,
    "tipo_acesso" TEXT NOT NULL DEFAULT 'normal',
    "go_credential_id" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "IntegrationCredential_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "Client_cpf_cnpj_key" ON "Client"("cpf_cnpj");

-- CreateIndex
CREATE UNIQUE INDEX "ClientAddress_client_id_key" ON "ClientAddress"("client_id");

-- CreateIndex
CREATE UNIQUE INDEX "ConsumerUnit_uc_code_key" ON "ConsumerUnit"("uc_code");

-- CreateIndex
CREATE UNIQUE INDEX "CommercialData_client_id_key" ON "CommercialData"("client_id");

-- CreateIndex
CREATE UNIQUE INDEX "IntegrationCredential_go_credential_id_key" ON "IntegrationCredential"("go_credential_id");

-- AddForeignKey
ALTER TABLE "ClientAddress" ADD CONSTRAINT "ClientAddress_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "Client"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ConsumerUnit" ADD CONSTRAINT "ConsumerUnit_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "Client"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ConsumerUnit" ADD CONSTRAINT "ConsumerUnit_credential_id_fkey" FOREIGN KEY ("credential_id") REFERENCES "IntegrationCredential"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "CommercialData" ADD CONSTRAINT "CommercialData_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "Client"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "IntegrationCredential" ADD CONSTRAINT "IntegrationCredential_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "Client"("id") ON DELETE CASCADE ON UPDATE CASCADE;
