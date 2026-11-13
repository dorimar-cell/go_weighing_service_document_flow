-- Этап 6: печатные документы, продолжение.
-- Счёт-фактура выбирается по рейсу.

CREATE INDEX IF NOT EXISTS idx_invoice_document_weighing
    ON invoice_document (id_weighing);
