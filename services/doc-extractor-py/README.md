# doc-extractor-py

Serviço documental responsável por enriquecer os dados da fatura a partir do PDF.

Responsabilidades:

- parse local com `PyMuPDF`
- fallback com `Mistral OCR`
- `source_map`
- `confidence_map`
- resposta compatível com `packages/contracts`

Estado atual:

- scaffold inicial
- endpoint de health
- endpoint de extração ainda em stub

Base de migração:

- `src/fatura/parser_pdf.py`
- `src/fatura/mistral_ocr.py`
