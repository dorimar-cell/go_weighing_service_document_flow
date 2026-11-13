-- Этап 6: отчёты и печатные документы.
-- Новых таблиц нет: бланки используют legacy-таблицы ttn_document и
-- invoice_torg12_documents, реквизиты организации — dir_default_settings_for_print.

CREATE INDEX IF NOT EXISTS idx_ttn_document_weighing
    ON ttn_document (id_weighing);

CREATE INDEX IF NOT EXISTS idx_invoice_torg12_documents_weighing
    ON invoice_torg12_documents (id_weighing);

CREATE INDEX IF NOT EXISTS idx_dir_default_settings_for_print_field
    ON dir_default_settings_for_print (type_of_field);
