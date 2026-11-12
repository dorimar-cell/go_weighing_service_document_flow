-- 0003: связи и индексы основной схемы.
--
-- Связи сделаны DEFERRABLE INITIALLY DEFERRED: это позволяет ETL загружать
-- таблицы в произвольном порядке в одной транзакции.
-- Если данные легаси содержат «битые» ссылки, замените ADD CONSTRAINT на
-- ADD CONSTRAINT ... NOT VALID и проверьте отдельно.

-- Справочники и модели
ALTER TABLE dir_type_of_cargo
    ADD CONSTRAINT fk_type_of_cargo_cargo FOREIGN KEY (cargo_id) REFERENCES dir_cargo(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE dir_type_of_transport_models
    ADD CONSTRAINT fk_ttm_type FOREIGN KEY (id_dir_type_of_transport) REFERENCES dir_type_of_transport(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_ttm_suspension FOREIGN KEY (id_dir_type_of_car_suspension) REFERENCES dir_type_of_car_suspension(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE dir_type_of_trailer_models
    ADD CONSTRAINT fk_ttrm_type FOREIGN KEY (id_dir_type_of_trailer) REFERENCES dir_type_of_trailer(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_ttrm_suspension FOREIGN KEY (id_dir_type_of_car_suspension) REFERENCES dir_type_of_car_suspension(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE dir_transport_model_axle_characters
    ADD CONSTRAINT fk_tmac_model FOREIGN KEY (id_type_of_ts_models) REFERENCES dir_type_of_transport_models(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE dir_trailer_model_axle_characters
    ADD CONSTRAINT fk_trmac_model FOREIGN KEY (id_type_of_ts_models) REFERENCES dir_type_of_trailer_models(id) DEFERRABLE INITIALLY DEFERRED;

-- Пользователи и права
ALTER TABLE dir_users
    ADD CONSTRAINT fk_users_weight_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_users_group FOREIGN KEY (id_dir_user_groups) REFERENCES dir_user_group(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE com_dir_user_groups_access_permissions
    ADD CONSTRAINT fk_ugap_group FOREIGN KEY (id_dir_user_group) REFERENCES dir_user_group(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_ugap_permission FOREIGN KEY (id_system_dir_access_permission) REFERENCES system_dir_access_permissions(id) DEFERRABLE INITIALLY DEFERRED;

-- Транспорт / водитель / прицеп
ALTER TABLE com_dir_transport_driver_trailer
    ADD CONSTRAINT fk_cdtdt_driver FOREIGN KEY (id_dir_driver) REFERENCES dir_driver(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cdtdt_transport FOREIGN KEY (id_dir_transport) REFERENCES dir_transport(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cdtdt_trailer FOREIGN KEY (id_dir_trailer) REFERENCES dir_trailer(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cdtdt_type_transport FOREIGN KEY (id_dir_type_of_transport) REFERENCES dir_type_of_transport(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cdtdt_type_trailer FOREIGN KEY (id_dir_type_of_trailer) REFERENCES dir_type_of_trailer(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cdtdt_transport_model FOREIGN KEY (id_transport_model) REFERENCES dir_type_of_transport_models(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cdtdt_trailer_model FOREIGN KEY (id_trailer_model) REFERENCES dir_type_of_trailer_models(id) DEFERRABLE INITIALLY DEFERRED;

-- Накладная и её строки
ALTER TABLE dir_invoice
    ADD CONSTRAINT fk_invoice_carrier FOREIGN KEY (id_dir_carrier) REFERENCES dir_carrier(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_invoice_counterparty FOREIGN KEY (id_dir_counterparty) REFERENCES dir_counterparty(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_invoice_operation FOREIGN KEY (id_dir_type_of_operation) REFERENCES dir_type_of_operation(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE com_dir_invoice_dir_storage
    ADD CONSTRAINT fk_cids_cargo FOREIGN KEY (id_dir_cargo) REFERENCES dir_cargo(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cids_provider FOREIGN KEY (id_dir_provider) REFERENCES dir_provider(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cids_invoice FOREIGN KEY (id_dir_invoice) REFERENCES dir_invoice(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cids_storage FOREIGN KEY (id_dir_storage) REFERENCES dir_storage(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cids_storage_rec FOREIGN KEY (id_dir_storage_rec) REFERENCES dir_storage(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cids_placement FOREIGN KEY (id_dir_placement) REFERENCES dir_placement(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cids_placement_rec FOREIGN KEY (id_dir_placement_rec) REFERENCES dir_placement(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cids_pack FOREIGN KEY (id_dir_type_of_pack) REFERENCES dir_type_of_pack(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cids_load_point FOREIGN KEY (id_dir_load_point) REFERENCES dir_load_point(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cids_upload_point FOREIGN KEY (id_dir_upload_point) REFERENCES dir_upload_point(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_cids_type_cargo FOREIGN KEY (id_dir_type_of_cargo) REFERENCES dir_type_of_cargo(id) DEFERRABLE INITIALLY DEFERRED;

-- Пропуска
ALTER TABLE reg_weight
    ADD CONSTRAINT fk_rw_type_regweight FOREIGN KEY (id_system_dir_type_of_regweight) REFERENCES system_dir_type_of_regweight(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_identifier_state FOREIGN KEY (id_system_dir_identifier_state) REFERENCES system_dir_identifier_state(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_pass_type FOREIGN KEY (id_system_dir_type_of_pass) REFERENCES system_dir_type_of_pass(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_plantation FOREIGN KEY (id_dir_kfkh_plantation) REFERENCES dir_kfkh_plantation(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_characteristics FOREIGN KEY (id_dir_additional_characteristics) REFERENCES dir_additional_characteristics(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_fund FOREIGN KEY (id_dir_fund) REFERENCES dir_fund(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_cargo FOREIGN KEY (id_dir_cargo) REFERENCES dir_cargo(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_type_cargo FOREIGN KEY (id_dir_type_of_cargo) REFERENCES dir_type_of_cargo(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_storage FOREIGN KEY (id_dir_storage) REFERENCES dir_storage(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_storage_rec FOREIGN KEY (id_dir_storage_rec) REFERENCES dir_storage(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_placement FOREIGN KEY (id_dir_placement) REFERENCES dir_placement(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_placement_rec FOREIGN KEY (id_dir_placement_rec) REFERENCES dir_placement(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_driver FOREIGN KEY (id_dir_driver) REFERENCES dir_driver(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_transport FOREIGN KEY (id_dir_transport) REFERENCES dir_transport(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_trailer FOREIGN KEY (id_dir_trailer) REFERENCES dir_trailer(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_type_transport FOREIGN KEY (id_dir_type_of_transport) REFERENCES dir_type_of_transport(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_type_trailer FOREIGN KEY (id_dir_type_of_trailer) REFERENCES dir_type_of_trailer(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_transport_model FOREIGN KEY (id_dir_type_of_transport_models) REFERENCES dir_type_of_transport_models(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_trailer_model FOREIGN KEY (id_dir_type_of_trailer_model) REFERENCES dir_type_of_trailer_models(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_carrier FOREIGN KEY (id_dir_carrier) REFERENCES dir_carrier(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_counterparty FOREIGN KEY (id_dir_counterparty) REFERENCES dir_counterparty(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_rw_operation FOREIGN KEY (id_dir_type_of_operation) REFERENCES dir_type_of_operation(id) DEFERRABLE INITIALLY DEFERRED;

-- Рейс и взвешивания
ALTER TABLE weighing_info
    ADD CONSTRAINT fk_wi_invoice FOREIGN KEY (id_dir_invoice) REFERENCES dir_invoice(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_wi_laboratory FOREIGN KEY (id_dir_laboratory) REFERENCES dir_laboratory(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_wi_characteristics FOREIGN KEY (id_dir_additional_characteristics) REFERENCES dir_additional_characteristics(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_wi_transport_driver_trailer FOREIGN KEY (id_com_transport_driver_trailer) REFERENCES com_dir_transport_driver_trailer(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_wi_fund FOREIGN KEY (id_dir_fund) REFERENCES dir_fund(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE weighing
    ADD CONSTRAINT fk_weighing_info FOREIGN KEY (id_weighting_info) REFERENCES weighing_info(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_weighing_reg_weight FOREIGN KEY (id_reg_weight) REFERENCES reg_weight(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_weighing_first_user FOREIGN KEY (id_first_weight_user) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_weighing_second_user FOREIGN KEY (id_second_weight_user) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_weighing_created_user FOREIGN KEY (created_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE weighting_photo
    ADD CONSTRAINT fk_photo_weighing FOREIGN KEY (id_weighing) REFERENCES weighing(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_photo_camera FOREIGN KEY (id_settings_video_surveillance) REFERENCES settings_video_surveillance(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE dir_lab_register
    ADD CONSTRAINT fk_lab_register_weighing FOREIGN KEY (id_weighing) REFERENCES weighing(id) DEFERRABLE INITIALLY DEFERRED;

-- Аудит и нераспознанные
ALTER TABLE audit_weighing
    ADD CONSTRAINT fk_audit_user FOREIGN KEY (created_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_audit_weight_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE audit_weighting_photo
    ADD CONSTRAINT fk_audit_photo_audit FOREIGN KEY (id_audit) REFERENCES audit_weighing(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_audit_photo_camera FOREIGN KEY (id_settings_video_surveillance) REFERENCES settings_video_surveillance(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE unrec_weighing
    ADD CONSTRAINT fk_unrec_user FOREIGN KEY (created_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_unrec_weight_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE unrec_weighting_photo
    ADD CONSTRAINT fk_unrec_photo_weighing FOREIGN KEY (id_unrec_weighing) REFERENCES unrec_weighing(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_unrec_photo_camera FOREIGN KEY (id_settings_video_surveillance) REFERENCES settings_video_surveillance(id) DEFERRABLE INITIALLY DEFERRED;

-- Автонумерация
ALTER TABLE auto_invoice
    ADD CONSTRAINT fk_auto_invoice_weighing FOREIGN KEY (id_weighing) REFERENCES weighing(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE auto_invoice_everyday
    ADD CONSTRAINT fk_auto_invoice_everyday_weighing FOREIGN KEY (id_weighing) REFERENCES weighing(id) DEFERRABLE INITIALLY DEFERRED;

-- Документы
ALTER TABLE ttn_document
    ADD CONSTRAINT fk_ttn_weighing FOREIGN KEY (id_weighing) REFERENCES weighing(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_ttn_carriage FOREIGN KEY (id_system_dir_type_of_carriage) REFERENCES system_dir_type_of_carriage(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_ttn_carriage_document FOREIGN KEY (id_system_dir_type_of_carriage_document) REFERENCES system_dir_type_of_carriage_document(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_ttn_packing FOREIGN KEY (id_system_dir_kind_of_packing) REFERENCES system_dir_kind_of_packing(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_ttn_pack FOREIGN KEY (id_dir_kind_of_packing) REFERENCES dir_type_of_pack(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE invoice_document
    ADD CONSTRAINT fk_invoice_doc_weighing FOREIGN KEY (id_weighing) REFERENCES weighing(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE invoice_torg12_documents
    ADD CONSTRAINT fk_torg12_weighing FOREIGN KEY (id_weighing) REFERENCES weighing(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_torg12_consignee FOREIGN KEY (id_consignee_org) REFERENCES dir_counterparty(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_torg12_provider FOREIGN KEY (id_provider_org) REFERENCES dir_counterparty(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE acceptance_certificates
    ADD CONSTRAINT fk_acceptance_weighing FOREIGN KEY (id_weighing) REFERENCES weighing(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_acceptance_passport FOREIGN KEY (id_passport) REFERENCES dir_passport(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_acceptance_shipper FOREIGN KEY (id_shipper) REFERENCES dir_counterparty(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_acceptance_consignee FOREIGN KEY (id_consignee) REFERENCES dir_counterparty(id) DEFERRABLE INITIALLY DEFERRED;

-- КФХ
ALTER TABLE com_kfkh
    ADD CONSTRAINT fk_kfkh_combine FOREIGN KEY (id_dir_combine) REFERENCES dir_kfkh_combine(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_kfkh_combiner FOREIGN KEY (id_dir_combiner) REFERENCES dir_kfkh_combiner(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_kfkh_plantation FOREIGN KEY (id_dir_plantation) REFERENCES dir_kfkh_plantation(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_kfkh_weighing_info FOREIGN KEY (id_weighing_info) REFERENCES weighing_info(id) DEFERRABLE INITIALLY DEFERRED;

-- Склад и лаборатория
ALTER TABLE dir_storage_ostatki
    ADD CONSTRAINT fk_ostatki_cargo FOREIGN KEY (id_dir_cargo) REFERENCES dir_cargo(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_ostatki_storage FOREIGN KEY (id_dir_storage) REFERENCES dir_storage(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE dir_spisanie_ostatkov
    ADD CONSTRAINT fk_spisanie_cargo FOREIGN KEY (id_dir_cargo) REFERENCES dir_cargo(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_spisanie_storage FOREIGN KEY (id_dir_storage) REFERENCES dir_storage(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE dir_laboratory
    ADD CONSTRAINT fk_laboratory_weighing_info FOREIGN KEY (id_weighing_info) REFERENCES weighing_info(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_laboratory_user FOREIGN KEY (created_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE dir_deleted_laboratories_records
    ADD CONSTRAINT fk_deleted_lab_weighing_info FOREIGN KEY (id_weighing_info) REFERENCES weighing_info(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_deleted_lab_user FOREIGN KEY (deleted_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_deleted_lab_lab FOREIGN KEY (id_lab) REFERENCES dir_laboratory(id) DEFERRABLE INITIALLY DEFERRED;

-- Оси
ALTER TABLE dir_axle_real_weights
    ADD CONSTRAINT fk_axle_weighing FOREIGN KEY (id_weighings) REFERENCES weighing(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_axle_transport_model FOREIGN KEY (id_dir_type_of_transport_models) REFERENCES dir_type_of_transport_models(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_axle_trailer_model FOREIGN KEY (id_dir_type_of_trailer_models) REFERENCES dir_type_of_trailer_models(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE dir_terminal_axle
    ADD CONSTRAINT fk_terminal_axle_weighing FOREIGN KEY (id_weighings) REFERENCES weighing(id) DEFERRABLE INITIALLY DEFERRED;

-- Настройки: привязка к весовой и пользователям
ALTER TABLE settings_terminal
    ADD CONSTRAINT fk_settings_terminal_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_settings_terminal_user FOREIGN KEY (updated_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE settings_skud
    ADD CONSTRAINT fk_settings_skud_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_settings_skud_user FOREIGN KEY (updated_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE settings_scoreboard
    ADD CONSTRAINT fk_settings_scoreboard_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_settings_scoreboard_user FOREIGN KEY (updated_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE settings_speech
    ADD CONSTRAINT fk_settings_speech_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_settings_speech_user FOREIGN KEY (updated_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE settings_axle_load
    ADD CONSTRAINT fk_settings_axle_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_settings_axle_user FOREIGN KEY (updated_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE settings_audit_load
    ADD CONSTRAINT fk_settings_audit_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE settings_recognition_module
    ADD CONSTRAINT fk_settings_recognition_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_settings_recognition_user FOREIGN KEY (updated_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE settings_recognition_m_pixel
    ADD CONSTRAINT fk_settings_pixel_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_settings_pixel_camera FOREIGN KEY (id_settings_video_surveillance) REFERENCES settings_video_surveillance(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE settings_video_surveillance
    ADD CONSTRAINT fk_settings_video_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_settings_video_user FOREIGN KEY (updated_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE settings_bar_code
    ADD CONSTRAINT fk_settings_barcode_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_settings_barcode_user FOREIGN KEY (updated_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE settings_usb_key_reader
    ADD CONSTRAINT fk_settings_usb_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_settings_usb_user FOREIGN KEY (updated_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE settings_laboratory
    ADD CONSTRAINT fk_settings_lab_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_settings_lab_user FOREIGN KEY (modified_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE settings_laboratory2
    ADD CONSTRAINT fk_settings_lab2_room FOREIGN KEY (id_dir_weight_room) REFERENCES dir_weight_room(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_settings_lab2_user FOREIGN KEY (modified_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE settings_app
    ADD CONSTRAINT fk_settings_app_company FOREIGN KEY (id_company) REFERENCES dir_counterparty(id) DEFERRABLE INITIALLY DEFERRED,
    ADD CONSTRAINT fk_settings_app_user FOREIGN KEY (modified_by_user_id) REFERENCES dir_users(id) DEFERRABLE INITIALLY DEFERRED;

-- Индексы под журналы и выборки
CREATE INDEX IF NOT EXISTS idx_weighing_first_date ON weighing (date_time_first_weight);
CREATE INDEX IF NOT EXISTS idx_weighing_second_date ON weighing (date_time_second_weight);
CREATE INDEX IF NOT EXISTS idx_weighing_info ON weighing (id_weighting_info);
CREATE INDEX IF NOT EXISTS idx_weighing_reg_weight ON weighing (id_reg_weight);
CREATE INDEX IF NOT EXISTS idx_weighing_deleted ON weighing (deleted);

CREATE INDEX IF NOT EXISTS idx_reg_weight_name ON reg_weight (name);
CREATE INDEX IF NOT EXISTS idx_reg_weight_transport ON reg_weight (id_dir_transport);
CREATE INDEX IF NOT EXISTS idx_reg_weight_driver ON reg_weight (id_dir_driver);
CREATE INDEX IF NOT EXISTS idx_reg_weight_deleted ON reg_weight (deleted);

CREATE INDEX IF NOT EXISTS idx_weighting_info_transport ON weighing_info (id_com_transport_driver_trailer);
CREATE INDEX IF NOT EXISTS idx_weighting_info_invoice ON weighing_info (id_dir_invoice);
CREATE INDEX IF NOT EXISTS idx_weighting_photo_weighing ON weighting_photo (id_weighing);
CREATE INDEX IF NOT EXISTS idx_weighting_photo_audit ON audit_weighting_photo (id_audit);

CREATE INDEX IF NOT EXISTS idx_audit_weighing_date ON audit_weighing (created_date);
CREATE INDEX IF NOT EXISTS idx_audit_weighing_room ON audit_weighing (id_dir_weight_room);

CREATE INDEX IF NOT EXISTS idx_ttn_weighing ON ttn_document (id_weighing);
CREATE INDEX IF NOT EXISTS idx_invoice_doc_weighing ON invoice_document (id_weighing);
CREATE INDEX IF NOT EXISTS idx_torg12_weighing ON invoice_torg12_documents (id_weighing);
CREATE INDEX IF NOT EXISTS idx_acceptance_weighing ON acceptance_certificates (id_weighing);

CREATE INDEX IF NOT EXISTS idx_users_name ON dir_users (name);
CREATE INDEX IF NOT EXISTS idx_users_room ON dir_users (id_dir_weight_room);
CREATE INDEX IF NOT EXISTS idx_cids_invoice ON com_dir_invoice_dir_storage (id_dir_invoice);
CREATE INDEX IF NOT EXISTS idx_lab_register_weighing ON dir_lab_register (id_weighing);
CREATE INDEX IF NOT EXISTS idx_storage_ostatki_cargo ON dir_storage_ostatki (id_dir_cargo, created_date);

CREATE INDEX IF NOT EXISTS idx_settings_terminal_room ON settings_terminal (id_dir_weight_room, platform_number);
CREATE INDEX IF NOT EXISTS idx_settings_skud_room ON settings_skud (id_dir_weight_room, platform_number);
CREATE INDEX IF NOT EXISTS idx_settings_scoreboard_room ON settings_scoreboard (id_dir_weight_room, platform_number);
CREATE INDEX IF NOT EXISTS idx_settings_speech_room ON settings_speech (id_dir_weight_room, platform_number);
CREATE INDEX IF NOT EXISTS idx_settings_axle_room ON settings_axle_load (id_dir_weight_room, platform_number);
CREATE INDEX IF NOT EXISTS idx_settings_video_room ON settings_video_surveillance (id_dir_weight_room, platform_number, name);
