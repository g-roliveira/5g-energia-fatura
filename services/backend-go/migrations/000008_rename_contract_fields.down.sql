-- Revert: remove valor_ip_com_desconto and rename fator_repasse_energia back.

ALTER TABLE public.contract DROP COLUMN valor_ip_com_desconto;

ALTER TABLE public.contract RENAME COLUMN fator_repasse_energia TO desconto_percentual;
