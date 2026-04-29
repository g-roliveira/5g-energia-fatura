-- Rename desconto_percentual to fator_repasse_energia and add valor_ip_com_desconto
-- for the new billing model.

ALTER TABLE public.contract RENAME COLUMN desconto_percentual TO fator_repasse_energia;

-- Add valor_ip_com_desconto: contractual value for Iluminacao Publica
-- in the COM (repasse) scenario.
ALTER TABLE public.contract ADD COLUMN valor_ip_com_desconto NUMERIC(12,2) NOT NULL DEFAULT 0.00;
