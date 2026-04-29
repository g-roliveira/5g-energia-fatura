-- Adiciona coluna pdf_generated à cycle_consumer_unit
ALTER TABLE public.cycle_consumer_unit
ADD COLUMN pdf_generated BOOLEAN NOT NULL DEFAULT FALSE;
