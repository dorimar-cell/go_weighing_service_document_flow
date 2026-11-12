-- 0002: основная схема весовой (без данных).
-- Сгенерировано из EF-модели WeightingApp.Data; ключи и связи проверяются вручную.

-- AcceptanceCertificates (AcceptanceCertificates)
CREATE TABLE IF NOT EXISTS acceptance_certificates (
    id bigserial PRIMARY KEY,
    id_weighing bigint NOT NULL,
    id_passport bigint,
    id_shipper bigint,
    id_consignee bigint NOT NULL,
    name varchar(55),
    date timestamptz NOT NULL,
    nds real,
    blockage real,
    law_basis varchar(255),
    description varchar(255),
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz
);

-- audit_weighing (audit_weighing)
CREATE TABLE IF NOT EXISTS audit_weighing (
    id bigserial PRIMARY KEY,
    audit_weight real NOT NULL,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    platform_number bigint,
    id_dir_weight_room bigint,
    audit_stable_weight real,
    modified_date timestamptz,
    transport_num varchar(45),
    is_registered boolean NOT NULL DEFAULT false
);

-- audit_weighting_photo (audit_weighting_photo)
CREATE TABLE IF NOT EXISTS audit_weighting_photo (
    id bigserial PRIMARY KEY,
    id_audit bigint NOT NULL,
    file_name varchar(255),
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    id_settings_video_surveillance bigint NOT NULL
);

-- auto_invoice (auto_invoice)
CREATE TABLE IF NOT EXISTS auto_invoice (
    id bigserial PRIMARY KEY,
    id_weighing bigint NOT NULL,
    auto_invoice_number bigint NOT NULL,
    auto_invoice_value text
);

-- auto_invoice_everyday (auto_invoice_everyday)
CREATE TABLE IF NOT EXISTS auto_invoice_everyday (
    id bigserial PRIMARY KEY,
    day timestamptz NOT NULL,
    time timestamptz NOT NULL,
    id_weighing bigint NOT NULL,
    cur_auto_inv_number bigint NOT NULL
);

-- com_dir_invoice_dir_storage (com_dir_invoice_dir_storage)
CREATE TABLE IF NOT EXISTS com_dir_invoice_dir_storage (
    id bigserial PRIMARY KEY,
    id_dir_cargo bigint,
    id_dir_provider bigint,
    id_dir_invoice bigint,
    id_dir_storage bigint,
    id_dir_placement bigint,
    id_dir_type_of_pack bigint,
    id_dir_load_point bigint,
    id_dir_upload_point bigint,
    pack_count bigint,
    pack_price real,
    deleted boolean NOT NULL DEFAULT false,
    id_dir_type_of_cargo bigint,
    doc_tara real,
    doc_brutto real,
    doc_netto real,
    doc_date timestamptz,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz,
    id_dir_storage_rec bigint,
    id_dir_placement_rec bigint
);

-- com_dir_transport_driver_trailer (com_dir_transport_driver_trailer)
CREATE TABLE IF NOT EXISTS com_dir_transport_driver_trailer (
    id bigserial PRIMARY KEY,
    id_dir_driver bigint,
    id_dir_transport bigint NOT NULL,
    id_dir_trailer bigint,
    deleted boolean NOT NULL DEFAULT false,
    id_dir_type_of_trailer bigint,
    id_dir_type_of_transport bigint,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz,
    recogn_transport_num_first_weight varchar(45),
    recogn_trailer_num_first_weight varchar(45),
    recogn_transport_num_second_weight varchar(45),
    recogn_trailer_num_second_weight varchar(45),
    id_transport_model bigint,
    id_trailer_model bigint
);

-- com_dir_UserGroups_AccessPermissions (com_dir_UserGroups_AccessPermissions)
CREATE TABLE IF NOT EXISTS com_dir_user_groups_access_permissions (
    id bigserial PRIMARY KEY,
    id_dir_user_group bigint NOT NULL,
    id_system_dir_access_permission bigint NOT NULL
);

-- com_kfkh (com_kfkh)
CREATE TABLE IF NOT EXISTS com_kfkh (
    id bigserial PRIMARY KEY,
    id_dir_combine bigint,
    id_dir_combiner bigint,
    id_dir_plantation bigint,
    use_scope real,
    deleted boolean NOT NULL DEFAULT false,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz,
    id_weighing_info bigint NOT NULL,
    talon text
);

-- dir_additional_characteristics (dir_additional_characteristics)
CREATE TABLE IF NOT EXISTS dir_additional_characteristics (
    id bigserial PRIMARY KEY,
    humidity real,
    sor real,
    marriage real,
    reason_of_marriage varchar(50),
    deleted boolean NOT NULL DEFAULT false,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz
);

-- dir_AxleRealWeights (dir_AxleRealWeights)
CREATE TABLE IF NOT EXISTS dir_axle_real_weights (
    id bigserial PRIMARY KEY,
    id_weighings bigint,
    weight varchar(150),
    id_dir_type_of_transport_models bigint,
    id_dir_type_of_trailer_models bigint,
    created_date timestamptz NOT NULL
);

-- dir_defaultSettingsForPrint (dir_defaultSettingsForPrint)
CREATE TABLE IF NOT EXISTS dir_default_settings_for_print (
    id bigserial PRIMARY KEY,
    type_of_field text,
    field_value text
);

-- dir_ExportFile (dir_ExportFile)
CREATE TABLE IF NOT EXISTS dir_export_file (
    id bigserial PRIMARY KEY,
    pattern_name varchar(150),
    file_name varchar(150),
    file_extension bigint NOT NULL,
    export boolean NOT NULL,
    command_text text,
    begin_time timestamptz,
    end_time timestamptz
);

-- dir_invoice (dir_invoice)
CREATE TABLE IF NOT EXISTS dir_invoice (
    id bigserial PRIMARY KEY,
    name varchar(45),
    id_dir_carrier bigint,
    id_dir_counterparty bigint,
    id_dir_type_of_operation bigint,
    date timestamptz NOT NULL,
    time time NOT NULL,
    price real,
    discount bigint,
    total_cost real,
    deleted boolean NOT NULL DEFAULT false,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz
);

-- dir_lab_register (dir_lab_register)
CREATE TABLE IF NOT EXISTS dir_lab_register (
    id bigserial PRIMARY KEY,
    id_weighing bigint NOT NULL,
    plombs varchar(200),
    deleted boolean NOT NULL DEFAULT false,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz
);

-- dir_loadPoint (dir_loadPoint)
CREATE TABLE IF NOT EXISTS dir_load_point (
    id bigserial PRIMARY KEY,
    name varchar(255),
    deleted boolean NOT NULL DEFAULT false,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz
);

-- dir_passport (dir_passport)
CREATE TABLE IF NOT EXISTS dir_passport (
    id bigserial PRIMARY KEY,
    fio varchar(255),
    series_of_passport varchar(4),
    id_of_passport varchar(6),
    issued_by varchar(255),
    issued_when timestamptz,
    birth_day timestamptz,
    birth_address varchar(255),
    residence_address varchar(255),
    deleted boolean NOT NULL DEFAULT false,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz,
    ref uuid NOT NULL,
    inn text
);

-- dir_storage_ostatki (dir_storage_ostatki)
CREATE TABLE IF NOT EXISTS dir_storage_ostatki (
    id bigserial PRIMARY KEY,
    ves real NOT NULL,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    id_dir_cargo bigint,
    id_dir_storage bigint
);

-- dir_table_presentation (dir_table_presentation)
CREATE TABLE IF NOT EXISTS dir_table_presentation (
    id bigserial PRIMARY KEY,
    width bigint NOT NULL,
    table_name text,
    column_name text
);

-- dir_TerminalAxle (dir_TerminalAxle)
CREATE TABLE IF NOT EXISTS dir_terminal_axle (
    id bigserial PRIMARY KEY,
    id_weighings bigint,
    weight_ts real,
    axles_count bigint,
    speed bigint,
    direction bigint,
    error_codes bigint,
    is_stabil bigint,
    axles varchar(150),
    created_date timestamptz NOT NULL
);

-- dir_TrailerModelAxleCharacters (dir_TrailerModelAxleCharacters)
CREATE TABLE IF NOT EXISTS dir_trailer_model_axle_characters (
    id bigserial PRIMARY KEY,
    id_type_of_ts_models bigint NOT NULL,
    wheel_pitch bigint,
    inter_axle_distance real,
    axle_group bigint,
    group_in_tones real,
    limit_in_tones real,
    deleted boolean NOT NULL DEFAULT false,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz
);

-- dir_TransportModelAxleCharacters (dir_TransportModelAxleCharacters)
CREATE TABLE IF NOT EXISTS dir_transport_model_axle_characters (
    id bigserial PRIMARY KEY,
    id_type_of_ts_models bigint NOT NULL,
    wheel_pitch bigint,
    inter_axle_distance real,
    axle_group bigint,
    group_in_tones real,
    limit_in_tones real,
    deleted boolean NOT NULL DEFAULT false,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz
);

-- dir_typeOfCarSuspension (dir_typeOfCarSuspension)
CREATE TABLE IF NOT EXISTS dir_type_of_car_suspension (
    id bigserial PRIMARY KEY,
    name varchar(50),
    deleted boolean NOT NULL DEFAULT false,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz
);

-- dir_typeOfPack (dir_typeOfPack)
CREATE TABLE IF NOT EXISTS dir_type_of_pack (
    id bigserial PRIMARY KEY,
    name varchar(45),
    deleted boolean NOT NULL DEFAULT false,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz
);

-- dir_typeOfTrailerModels (dir_typeOfTrailerModels)
CREATE TABLE IF NOT EXISTS dir_type_of_trailer_models (
    id bigserial PRIMARY KEY,
    name varchar(50),
    id_dir_type_of_trailer bigint NOT NULL,
    id_dir_type_of_car_suspension bigint NOT NULL,
    width real,
    height real,
    length real,
    deleted boolean NOT NULL DEFAULT false,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz,
    weight_allowed real
);

-- dir_typeOfTransportModels (dir_typeOfTransportModels)
CREATE TABLE IF NOT EXISTS dir_type_of_transport_models (
    id bigserial PRIMARY KEY,
    name varchar(50),
    id_dir_type_of_transport bigint NOT NULL,
    id_dir_type_of_car_suspension bigint NOT NULL,
    width real,
    height real,
    length real,
    deleted boolean NOT NULL DEFAULT false,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz,
    weight_allowed real
);

-- dir_uploadPoint (dir_uploadPoint)
CREATE TABLE IF NOT EXISTS dir_upload_point (
    id bigserial PRIMARY KEY,
    name varchar(255),
    price_delivery bigint,
    deleted boolean NOT NULL DEFAULT false,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz
);

-- dir_Users (dir_Users)
CREATE TABLE IF NOT EXISTS dir_users (
    id bigserial PRIMARY KEY,
    name varchar(100) NOT NULL,
    password_hash text,
    legacy_password_md5 varchar(45),
    deleted boolean NOT NULL DEFAULT false,
    id_dir_weight_room bigint NOT NULL,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz,
    id_dir_user_groups bigint,
    ref uuid NOT NULL
);

-- dp_messages (dp_messages)
CREATE TABLE IF NOT EXISTS dp_messages (
    id bigserial PRIMARY KEY,
    number bigint NOT NULL,
    value bigint NOT NULL
);

-- InvoiceDocument (InvoiceDocument)
CREATE TABLE IF NOT EXISTS invoice_document (
    id bigserial PRIMARY KEY,
    id_weighing bigint NOT NULL,
    name varchar(55),
    date timestamptz NOT NULL,
    nds real,
    document_name varchar(255),
    chief_name varchar(255),
    accountant_name varchar(255),
    businessman_name varchar(255),
    reg_ip_certificate varchar(255),
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz
);

-- InvoiceTorg12Documents (InvoiceTorg12Documents)
CREATE TABLE IF NOT EXISTS invoice_torg12_documents (
    id bigserial PRIMARY KEY,
    id_weighing bigint NOT NULL,
    name varchar(55),
    consignee_org_name varchar(455),
    id_consignee_org bigint,
    provider_org_name varchar(455),
    id_provider_org bigint,
    date timestamptz NOT NULL,
    nds real,
    struct_name varchar(255),
    foundation varchar(255),
    cargo_release_name varchar(255),
    cargo_release_work_as varchar(255),
    accountant_name varchar(255),
    shipper_name varchar(255),
    shipper_work_as varchar(255),
    cargo_pick_up_name varchar(255),
    cargo_pick_up_work_as varchar(255),
    consignee_name varchar(255),
    consignee_work_as varchar(255),
    proxy_name varchar(255),
    proxy_date timestamptz,
    proxy_given_by varchar(255),
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz
);

-- log_StartTerminalInfo (log_StartTerminalInfo)
CREATE TABLE IF NOT EXISTS log_start_terminal_info (
    id bigserial PRIMARY KEY,
    date varchar(10),
    time varchar(5),
    temperature bigint,
    supply_voltage real,
    number_of_sensors bigint,
    serial_number_terminal bigint,
    serial_number_libra bigint,
    date_verification_protocol_formation varchar(10),
    deleted boolean NOT NULL DEFAULT false,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz
);

-- ReferenceBooksTable (ReferenceBooksTable)
CREATE TABLE IF NOT EXISTS reference_books_table (
    id bigserial PRIMARY KEY,
    name varchar(150),
    name_table varchar(150),
    description text,
    hide_in_dictionaries_list boolean NOT NULL,
    access_tag varchar(5)
);

-- reg_weight (reg_weight)
CREATE TABLE IF NOT EXISTS reg_weight (
    id bigserial PRIMARY KEY,
    name varchar(45),
    document_id varchar(100),
    date_pass_expiry timestamptz,
    time_pass_expiry time,
    pass_enabled boolean NOT NULL,
    auto_save_tara boolean NOT NULL,
    tara real,
    doc_tara real,
    doc_brutto real,
    doc_netto real,
    doc_date timestamptz,
    invoice_name varchar(45),
    id_system_dir_type_of_regweight bigint NOT NULL,
    id_system_dir_identifier_state bigint NOT NULL,
    id_system_dir_type_of_pass bigint NOT NULL,
    id_dir_kfkh_plantation bigint,
    id_dir_additional_characteristics bigint,
    id_dir_fund bigint,
    id_dir_cargo bigint,
    id_dir_type_of_cargo bigint,
    id_dir_storage bigint,
    id_dir_storage_rec bigint,
    id_dir_placement bigint,
    id_dir_placement_rec bigint,
    id_dir_driver bigint,
    id_dir_transport bigint NOT NULL,
    id_dir_trailer bigint,
    id_dir_type_of_trailer bigint,
    id_dir_type_of_transport bigint,
    id_dir_carrier bigint,
    id_dir_counterparty bigint,
    id_dir_type_of_operation bigint,
    deleted boolean NOT NULL DEFAULT false,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz,
    id_dir_type_of_transport_models bigint,
    id_dir_type_of_trailer_model bigint
);

-- RusNameField (RusNameField)
CREATE TABLE IF NOT EXISTS rus_name_field (
    id bigserial PRIMARY KEY,
    name varchar(100),
    rus_name varchar(400)
);

-- Setting_Backup (Setting_Backup)
CREATE TABLE IF NOT EXISTS setting_backup (
    id bigserial PRIMARY KEY,
    flag_power_backup boolean NOT NULL,
    interval_backup bigint NOT NULL,
    path_backup text,
    time_backup timestamptz NOT NULL,
    restore_path varchar(200)
);

-- Settings_App (Settings_App)
CREATE TABLE IF NOT EXISTS settings_app (
    id bigserial PRIMARY KEY,
    ts_number_mask_use boolean NOT NULL,
    ts_number_mask varchar(30),
    need_operation_type boolean NOT NULL,
    need_fund boolean NOT NULL,
    need_required_fields boolean,
    placement_visible boolean,
    ves_doc_visible boolean,
    axle_load_visible boolean,
    need_cargo_type boolean NOT NULL,
    visible_cargo_type boolean,
    need_placement boolean NOT NULL,
    is_auto_invoice boolean DEFAULT false,
    auto_invoice_every_new_day boolean,
    visible_fund boolean NOT NULL,
    visible_placement boolean,
    modified_by_user_id bigint,
    modified_date timestamptz,
    id_company bigint,
    photo_directory varchar(255),
    conditional_zero real NOT NULL,
    app_font_size real,
    use_com_object boolean NOT NULL,
    use_region boolean NOT NULL,
    union_recognize_skud boolean NOT NULL,
    country bigint,
    invoice_number bigint
);

-- Settings_AuditLoad (Settings_AuditLoad)
CREATE TABLE IF NOT EXISTS settings_audit_load (
    id bigserial PRIMARY KEY,
    enable boolean NOT NULL,
    platform_number bigint NOT NULL,
    id_dir_weight_room bigint NOT NULL,
    limit_of_start real NOT NULL,
    photo_directory varchar(255),
    min_location bigint NOT NULL,
    max_location bigint NOT NULL,
    historical_taras_delta real NOT NULL,
    photo_period_sec bigint NOT NULL,
    driver_out_sec bigint NOT NULL
);

-- Settings_AxleLoad (Settings_AxleLoad)
CREATE TABLE IF NOT EXISTS settings_axle_load (
    id bigserial PRIMARY KEY,
    enable boolean NOT NULL,
    platform_number bigint NOT NULL,
    id_dir_weight_room bigint NOT NULL,
    updated_date timestamptz,
    updated_by_user_id bigint,
    interval bigint NOT NULL,
    axle_mode bigint NOT NULL,
    axle_notification bigint NOT NULL,
    limit_of_start double precision NOT NULL,
    inter_axle real NOT NULL,
    min_axle real NOT NULL,
    window_size bigint NOT NULL,
    period_delay bigint NOT NULL,
    store_data boolean NOT NULL,
    type bigint NOT NULL
);

-- Settings_BarCode (Settings_BarCode)
CREATE TABLE IF NOT EXISTS settings_bar_code (
    id bigserial PRIMARY KEY,
    id_dir_weight_room bigint NOT NULL,
    enable boolean NOT NULL,
    work_mode varchar(1),
    platform_number bigint NOT NULL,
    updated_date timestamptz NOT NULL,
    updated_by_user_id bigint NOT NULL
);

-- Settings_C1Integrator (Settings_C1Integrator)
CREATE TABLE IF NOT EXISTS settings_c1_integrator (
    id bigserial PRIMARY KEY,
    ip varchar(50),
    port bigint NOT NULL,
    enable boolean NOT NULL,
    user_name text,
    password text,
    server text,
    refr text,
    path_to_db text,
    path_to_image text
);

-- Settings_Dp (Settings_Dp)
CREATE TABLE IF NOT EXISTS settings_dp (
    id bigserial PRIMARY KEY,
    enable boolean NOT NULL,
    dp_mode bigint NOT NULL,
    full_packet_timer bigint NOT NULL,
    short_packet_timer bigint NOT NULL,
    full_fast_packet_timer bigint NOT NULL,
    max_weight_timer bigint NOT NULL,
    sensors_timer bigint NOT NULL,
    grafik_timer bigint NOT NULL,
    temp_timer bigint NOT NULL
);

-- Settings_Email (Settings_Email)
CREATE TABLE IF NOT EXISTS settings_email (
    id bigserial PRIMARY KEY,
    user_mail varchar(150),
    user_password varchar(150),
    adresant_name varchar(150),
    adresant_mail varchar(150),
    message_body varchar(200),
    message varchar(200),
    port_yandex bigint NOT NULL,
    port_mail bigint NOT NULL,
    port_g_mail bigint NOT NULL,
    enable boolean NOT NULL,
    time_period bigint NOT NULL,
    period_flag boolean NOT NULL,
    send_flag boolean NOT NULL,
    excel_flag boolean NOT NULL,
    event_type_id bigint NOT NULL,
    user_smtp text,
    user_use_ssl boolean NOT NULL,
    user_port bigint NOT NULL,
    custom_mail boolean NOT NULL
);

-- Settings_Export (Settings_Export)
CREATE TABLE IF NOT EXISTS settings_export (
    id bigserial PRIMARY KEY,
    enable boolean NOT NULL,
    for_defined_period boolean NOT NULL,
    with_period_beginning boolean NOT NULL,
    period_length bigint NOT NULL,
    period_beginning timestamptz NOT NULL,
    period_end timestamptz NOT NULL,
    file_path varchar(100),
    time_period bigint NOT NULL,
    is_query_confirmed boolean NOT NULL DEFAULT false,
    is_executed_after boolean NOT NULL DEFAULT false,
    is_protocol_switched boolean NOT NULL DEFAULT false,
    updated_date timestamptz,
    updated_by_user_id bigint
);

-- Settings_Integration (Settings_Integration)
CREATE TABLE IF NOT EXISTS settings_integration (
    id bigserial PRIMARY KEY,
    enable boolean NOT NULL,
    for_defined_period boolean,
    integration_mode bigint,
    delete_date text,
    period_length bigint,
    period_beginning timestamptz,
    period_end timestamptz,
    download_file_path varchar(200),
    upload_file_path varchar(200),
    time_period bigint NOT NULL,
    download_period_time bigint NOT NULL,
    is_upload_executed boolean NOT NULL DEFAULT false,
    is_download_executed boolean NOT NULL DEFAULT false,
    updated_date timestamptz,
    updated_by_user_id bigint,
    begin_time timestamptz,
    end_time timestamptz,
    with_counterparty boolean NOT NULL,
    with_additional_charachteristics boolean NOT NULL,
    with_cargo boolean NOT NULL,
    with_storage boolean NOT NULL,
    with_users boolean NOT NULL,
    with_weighing boolean NOT NULL
);

-- Settings_KFKH (Settings_KFKH)
CREATE TABLE IF NOT EXISTS settings_kfkh (
    id bigserial PRIMARY KEY,
    enable boolean NOT NULL
);

-- Settings_Laboratory (Settings_Laboratory)
CREATE TABLE IF NOT EXISTS settings_laboratory (
    id bigserial PRIMARY KEY,
    enable boolean NOT NULL,
    ves_doc boolean NOT NULL,
    is_netto_required boolean NOT NULL DEFAULT false,
    modified_date timestamptz,
    modified_by_user_id bigint,
    difference_kg bigint NOT NULL,
    start_date timestamptz NOT NULL,
    id_dir_weight_room bigint NOT NULL
);

-- Settings_Laboratory2 (Settings_Laboratory2)
CREATE TABLE IF NOT EXISTS settings_laboratory2 (
    id bigserial PRIMARY KEY,
    enable boolean NOT NULL,
    modified_date timestamptz,
    modified_by_user_id bigint,
    id_dir_weight_room bigint NOT NULL
);

-- Settings_PreAxleReport (Settings_PreAxleReport)
CREATE TABLE IF NOT EXISTS settings_pre_axle_report (
    id bigserial PRIMARY KEY,
    ppvk_number varchar(50),
    scale_number varchar(50),
    ppvk_owner varchar(150),
    certificate_number varchar(50),
    certificate_name varchar(100),
    certificate_org varchar(150),
    certificate_from timestamptz NOT NULL,
    certificate_to timestamptz NOT NULL,
    control_place varchar(150),
    weghing_mode varchar(50)
);

-- Settings_RecognitionModule (Settings_RecognitionModule)
CREATE TABLE IF NOT EXISTS settings_recognition_module (
    id bigserial PRIMARY KEY,
    id_dir_weight_room bigint NOT NULL,
    recognition_license_folder varchar(300),
    min_plate_width bigint NOT NULL,
    max_plate_width bigint NOT NULL,
    thread_count bigint NOT NULL,
    recog_time bigint NOT NULL,
    recognized_count bigint NOT NULL,
    un_recognized_count bigint NOT NULL,
    weight real NOT NULL,
    enable boolean NOT NULL,
    work_mode varchar(1),
    platform_number bigint NOT NULL,
    updated_date timestamptz NOT NULL,
    updated_by_user_id bigint NOT NULL
);

-- Settings_RecognitionMPixel (Settings_RecognitionMPixel)
CREATE TABLE IF NOT EXISTS settings_recognition_m_pixel (
    id bigserial PRIMARY KEY,
    id_settings_video_surveillance bigint,
    enable boolean NOT NULL,
    zone_visio_left bigint NOT NULL,
    zone_visio_top bigint NOT NULL,
    zone_visio_right bigint NOT NULL,
    zone_visio_bottom bigint NOT NULL,
    clear_number boolean NOT NULL,
    correction_symbols boolean NOT NULL,
    delete_zones boolean NOT NULL,
    check_height boolean NOT NULL,
    check_probability boolean NOT NULL,
    check_gipotize boolean NOT NULL,
    avoid_replay boolean NOT NULL,
    recognition_probability_threshold bigint NOT NULL,
    number_failed_recognition_cycles bigint NOT NULL,
    number_succesed_recognition_cycles bigint NOT NULL,
    car_threshold bigint NOT NULL,
    zone_threshold bigint NOT NULL,
    thhory bigint NOT NULL,
    thhorx bigint NOT NULL,
    thbin bigint NOT NULL,
    min_pipe_size bigint NOT NULL,
    max_pipe_size bigint NOT NULL,
    zone_mode_96x24 boolean NOT NULL,
    zone_mode_128x32 boolean NOT NULL,
    zone_mode_160x40 boolean NOT NULL,
    zone_mode_192x48 boolean NOT NULL,
    zone_mode_256x64 boolean NOT NULL,
    zone_mode_320x80 boolean NOT NULL,
    zone_mode_384x96 boolean NOT NULL,
    zone_mode_48x48 boolean NOT NULL,
    zone_mode_64x64 boolean NOT NULL,
    zone_mode_80x80 boolean NOT NULL,
    zone_mode_96x96 boolean NOT NULL,
    zone_mode_128x128 boolean NOT NULL,
    zone_mode_160x160 boolean NOT NULL,
    zone_mode_192x192 boolean NOT NULL,
    update_date timestamptz NOT NULL,
    update_by_user_id bigint NOT NULL,
    id_dir_weight_room bigint NOT NULL
);

-- Settings_Scoreboard (Settings_Scoreboard)
CREATE TABLE IF NOT EXISTS settings_scoreboard (
    id bigserial PRIMARY KEY,
    enable boolean NOT NULL,
    scoreboard_model_value varchar(3),
    com_name varchar(7),
    baudrate varchar(15),
    parity varchar(15),
    databits varchar(1),
    stopbits varchar(15),
    scoreboard_interval varchar(4),
    platform_number bigint NOT NULL,
    id_dir_weight_room bigint NOT NULL,
    updated_date timestamptz NOT NULL,
    updated_by_user_id bigint NOT NULL,
    interface_type varchar(15),
    rts boolean NOT NULL,
    dtr boolean NOT NULL,
    handshake varchar(30),
    use_socket boolean,
    socket_ip text,
    socket_port bigint NOT NULL
);

-- Settings_ServerReport (Settings_ServerReport)
CREATE TABLE IF NOT EXISTS settings_server_report (
    id bigserial PRIMARY KEY,
    enable boolean NOT NULL,
    photo_directory varchar(100),
    ip_address varchar(50),
    port bigint,
    updated_date timestamptz,
    updated_by_user_id bigint,
    weighing_name varchar(100),
    connection_string varchar(700),
    audit_photo_directory varchar(150)
);

-- Settings_SKUD (Settings_SKUD)
CREATE TABLE IF NOT EXISTS settings_skud (
    id bigserial PRIMARY KEY,
    enable_skud boolean NOT NULL,
    interface_type varchar(15),
    network_address varchar(2),
    com_name varchar(7),
    baudrate varchar(15),
    parity varchar(15),
    databits varchar(1),
    stopbits varchar(15),
    handshake varchar(30),
    enable_barrier boolean NOT NULL,
    use_traffic_light_delay boolean,
    traffic_light_delay_value bigint,
    barrier_in_open varchar(2),
    barrier_in_close varchar(2),
    barrier_out_open varchar(2),
    barrier_out_close varchar(2),
    barrier_work_mode varchar(1),
    enable_traffic_light boolean NOT NULL,
    traffic_light_in_near varchar(2),
    traffic_light_in_far varchar(2),
    traffic_light_out_near varchar(2),
    traffic_light_out_far varchar(2),
    traffic_light_work_mode varchar(1),
    enable_reader boolean NOT NULL,
    reader_in varchar(2),
    reader_center varchar(2),
    reader_out varchar(2),
    reader_view_mode varchar(1),
    reader_work_mode varchar(1),
    use_socket boolean,
    socket_port bigint,
    socket_ip text,
    enable_ir_sensors boolean NOT NULL,
    control_s1 boolean NOT NULL,
    control_s2 boolean NOT NULL,
    control_s3 boolean NOT NULL,
    control_s4 boolean NOT NULL,
    control_s5 boolean NOT NULL,
    control_s6 boolean NOT NULL,
    invert_s1 boolean NOT NULL,
    invert_s2 boolean NOT NULL,
    invert_s3 boolean NOT NULL,
    invert_s4 boolean NOT NULL,
    invert_s5 boolean NOT NULL,
    invert_s6 boolean NOT NULL,
    platform_number bigint NOT NULL,
    id_dir_weight_room bigint NOT NULL,
    updated_date timestamptz NOT NULL,
    updated_by_user_id bigint NOT NULL,
    barrier_state boolean NOT NULL,
    traffic_state boolean NOT NULL
);

-- Settings_Speech (Settings_Speech)
CREATE TABLE IF NOT EXISTS settings_speech (
    id bigserial PRIMARY KEY,
    speech_type bigint NOT NULL,
    enable boolean NOT NULL,
    volume bigint NOT NULL,
    rate bigint NOT NULL,
    platform_number bigint NOT NULL,
    id_dir_weight_room bigint NOT NULL,
    updated_date timestamptz,
    updated_by_user_id bigint,
    message_text varchar(150),
    message_text_two varchar(150),
    message_text_three varchar(150),
    message_text_four varchar(150),
    message_text_five varchar(150),
    message_text_six varchar(150),
    message_text_seven varchar(150),
    message_text_eight varchar(150)
);

-- Settings_Store (Settings_Store)
CREATE TABLE IF NOT EXISTS settings_store (
    id bigserial PRIMARY KEY,
    enable boolean NOT NULL,
    contract_date timestamptz
);

-- Settings_Terminal (Settings_Terminal)
CREATE TABLE IF NOT EXISTS settings_terminal (
    id bigserial PRIMARY KEY,
    terminal_model_value varchar(15),
    interface_type varchar(15),
    network_address varchar(2),
    multiplier_weight real NOT NULL,
    com_name varchar(7),
    baudrate varchar(15),
    parity varchar(15),
    databits varchar(1),
    stopbits varchar(15),
    handshake varchar(30),
    rts boolean NOT NULL,
    dtr boolean NOT NULL,
    platform_number bigint NOT NULL,
    id_dir_weight_room bigint NOT NULL,
    updated_date timestamptz NOT NULL,
    is_modbus boolean DEFAULT false,
    ip text,
    port text,
    use_socket boolean,
    socket_ip text,
    socket_port text,
    enable boolean NOT NULL,
    terminal_interval varchar(4),
    updated_by_user_id bigint NOT NULL,
    weight_conditional_zero real NOT NULL
);

-- Settings_USBKeyReader (Settings_USBKeyReader)
CREATE TABLE IF NOT EXISTS settings_usb_key_reader (
    id bigserial PRIMARY KEY,
    enable boolean NOT NULL,
    com_name varchar(7),
    baudrate varchar(15),
    parity varchar(15),
    databits varchar(1),
    stopbits varchar(15),
    usb_key_reader_interval varchar(4),
    id_dir_weight_room bigint NOT NULL,
    updated_date timestamptz NOT NULL,
    updated_by_user_id bigint NOT NULL
);

-- Settings_VideoSurveillance (Settings_VideoSurveillance)
CREATE TABLE IF NOT EXISTS settings_video_surveillance (
    id bigserial PRIMARY KEY,
    enable boolean NOT NULL,
    name varchar(35),
    ip varchar(200),
    vlc_directory varchar(200),
    port bigint NOT NULL,
    login varchar(30),
    password varchar(30),
    wait_time bigint NOT NULL,
    try_times bigint NOT NULL,
    stream_type_mini_window bigint NOT NULL,
    stream_type_full_window bigint NOT NULL,
    platform_number bigint NOT NULL,
    id_dir_weight_room bigint NOT NULL,
    updated_date timestamptz NOT NULL,
    updated_by_user_id bigint NOT NULL
);

-- system_dir_AccessPermissions (system_dir_AccessPermissions)
CREATE TABLE IF NOT EXISTS system_dir_access_permissions (
    id bigserial PRIMARY KEY,
    name varchar(200),
    control_tag varchar(10)
);

-- system_dir_DirectionTypeOperation (system_dir_DirectionTypeOperation)
CREATE TABLE IF NOT EXISTS system_dir_direction_type_operation (
    id bigserial PRIMARY KEY,
    name varchar(45)
);

-- system_dir_kindOfPacking (system_dir_kindOfPacking)
CREATE TABLE IF NOT EXISTS system_dir_kind_of_packing (
    id bigserial PRIMARY KEY,
    name varchar(45)
);

-- system_dir_typeOfCarriage (system_dir_typeOfCarriage)
CREATE TABLE IF NOT EXISTS system_dir_type_of_carriage (
    id bigserial PRIMARY KEY,
    name varchar(45)
);

-- system_dir_typeOfCarriageDocument (system_dir_typeOfCarriageDocument)
CREATE TABLE IF NOT EXISTS system_dir_type_of_carriage_document (
    id bigserial PRIMARY KEY,
    name varchar(45)
);

-- system_dir_typeOfRegweight (system_dir_typeOfRegweight)
CREATE TABLE IF NOT EXISTS system_dir_type_of_regweight (
    id bigserial PRIMARY KEY,
    name varchar(45)
);

-- TTNDocument (TTNDocument)
CREATE TABLE IF NOT EXISTS ttn_document (
    id bigserial PRIMARY KEY,
    id_weighing bigint NOT NULL,
    proxy_name varchar(255),
    proxy_date timestamptz,
    proxy_given_by varchar(255),
    dogovor varchar(50),
    field varchar(50),
    sort varchar(50),
    putevka varchar(50),
    prilozhenie varchar(255),
    distance varchar(50),
    ttn_number varchar(50),
    cargo_release_name varchar(255),
    cargo_release_work_as varchar(255),
    accountant_name varchar(255),
    cargo_pick_up_name varchar(255),
    cargo_pick_up_work_as varchar(255),
    cargo_pick_up_identity varchar(255),
    consignee_name varchar(255),
    consignee_work_as varchar(255),
    route_distance real,
    cargo_capacity real,
    trailer_number varchar(45),
    type_of_trailer varchar(45),
    carriage_document_number varchar(15),
    pack_seal_number varchar(45),
    load_point varchar(255),
    upload_point varchar(255),
    carrier_organization varchar(400),
    sub_carrier_name varchar(400),
    sub_carrier_id bigint,
    customer_organization varchar(400),
    shipper_organization varchar(400),
    shipper_organization_id bigint,
    consignee_organization varchar(400),
    payer_organization varchar(400),
    id_system_dir_type_of_carriage bigint,
    id_system_dir_type_of_carriage_document bigint,
    id_system_dir_kind_of_packing bigint,
    id_dir_kind_of_packing bigint,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz,
    cargo_pick_up_additional_info text,
    price_cargo real,
    price_transportation real,
    forwarder boolean,
    type_of_own text,
    own_dir text,
    own_position_dir text,
    carrier_dir text,
    carrier_position_dir text,
    customer_dir text,
    customer_position_dir text,
    shipper_dir text,
    shipper_position_dir text,
    consignee_dir text,
    consignee_position_dir text,
    payer_dir text,
    payer_position_dir text
);

-- unrec_weighing (unrec_weighing)
CREATE TABLE IF NOT EXISTS unrec_weighing (
    id bigserial PRIMARY KEY,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    platform_number bigint,
    id_dir_weight_room bigint,
    stable_weight real,
    modified_date timestamptz
);

-- unrec_weighting_photo (unrec_weighting_photo)
CREATE TABLE IF NOT EXISTS unrec_weighting_photo (
    id bigserial PRIMARY KEY,
    id_unrec_weighing bigint NOT NULL,
    file_name varchar(255),
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    id_settings_video_surveillance bigint NOT NULL
);

-- weighing (weighing)
CREATE TABLE IF NOT EXISTS weighing (
    id bigserial PRIMARY KEY,
    id_weighting_info bigint NOT NULL,
    first_weight real NOT NULL,
    handle_first_weight real,
    number_of_weight bigint,
    flag_weight boolean NOT NULL,
    deleted boolean NOT NULL DEFAULT false,
    second_weight real,
    handle_second_weight real,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz,
    id_first_weight_user bigint NOT NULL,
    id_second_weight_user bigint,
    id_reg_weight bigint,
    date_time_first_weight timestamptz NOT NULL,
    date_time_second_weight timestamptz,
    ref uuid NOT NULL,
    first_weight_1c boolean NOT NULL,
    is_flow_weight boolean NOT NULL DEFAULT false,
    has_edited boolean NOT NULL,
    second_weight_1c boolean NOT NULL
);

-- weighing_info (weighing_info)
CREATE TABLE IF NOT EXISTS weighing_info (
    id bigserial PRIMARY KEY,
    id_dir_invoice bigint,
    id_dir_laboratory bigint,
    id_dir_additional_characteristics bigint,
    deleted boolean NOT NULL DEFAULT false,
    id_com_transport_driver_trailer bigint NOT NULL,
    id_dir_fund bigint,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz
);

-- weighting_photo (weighting_photo)
CREATE TABLE IF NOT EXISTS weighting_photo (
    id bigserial PRIMARY KEY,
    id_weighing bigint NOT NULL,
    file_name varchar(255),
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    first_weight boolean NOT NULL,
    id_settings_video_surveillance bigint NOT NULL
);

-- Таблицы, которые создаются в легаси-коде, но отсутствуют в EF-модели.

CREATE TABLE IF NOT EXISTS dir_laboratory (
    id bigserial PRIMARY KEY,
    sediment real NOT NULL,
    color_number bigint NOT NULL,
    acidl_number real NOT NULL,
    moisture_contents real NOT NULL,
    peroxid_number real NOT NULL,
    mass_fraction_phosphorus real NOT NULL,
    non_fatty_impurities real NOT NULL,
    id_weighing_info bigint NOT NULL,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    modified_by_user_id bigint,
    modified_date timestamptz,
    deleted boolean DEFAULT false,
    is_not_success boolean DEFAULT true,
    is_completed boolean,
    is_left_somthing boolean NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS dir_deleted_laboratories_records (
    id bigserial PRIMARY KEY,
    sediment real NOT NULL,
    color_number bigint NOT NULL,
    acidl_number real NOT NULL,
    moisture_contents real NOT NULL,
    peroxid_number real NOT NULL,
    mass_fraction_phosphorus real NOT NULL,
    non_fatty_impurities real NOT NULL,
    id_weighing_info bigint NOT NULL,
    deleted_by_user_id bigint NOT NULL,
    deleted_date timestamptz NOT NULL,
    is_not_success boolean,
    is_completed boolean,
    is_left_somthing boolean NOT NULL DEFAULT false,
    id_lab bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS dir_spisanie_ostatkov (
    id bigserial PRIMARY KEY,
    description varchar(300),
    doc varchar(300) NOT NULL,
    ves real NOT NULL,
    created_by_user_id bigint NOT NULL,
    created_date timestamptz NOT NULL,
    on_date timestamptz NOT NULL,
    id_dir_cargo bigint NOT NULL,
    id_dir_storage bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS dir_refusing (
    id bigserial PRIMARY KEY,
    name varchar(200) NOT NULL
);

-- Печатные документы СП-9/СП-11/СП-14: снимки реквизитов на момент печати.
CREATE TABLE IF NOT EXISTS sp9_doc (
    id bigserial PRIMARY KEY,
    org_name varchar(150),
    dep_name varchar(150),
    tok_name varchar(150),
    warehouse_name varchar(150),
    sender_name varchar(150),
    rep_name varchar(150),
    sort_name varchar(150),
    talon varchar(50),
    fields varchar(250),
    driver_fio varchar(150),
    transport_model varchar(100),
    gos_num varchar(50),
    brutto real,
    netto real,
    tara real,
    combiner1 real,
    combiner2 real,
    combiner3 real,
    combiner4 real,
    combiner5 real,
    combiner6 real,
    combiner7 real,
    combiner8 real,
    combiner9 real,
    combiner10 real,
    date_from date,
    date_to date,
    w_id bigint,
    num bigint
);

CREATE TABLE IF NOT EXISTS sp11_doc (
    id bigserial PRIMARY KEY,
    org_name varchar(150),
    dep_name varchar(150),
    brigada_name varchar(150),
    placement_name varchar(150),
    conterparty_name varchar(150),
    sort_name varchar(150),
    invoice_name varchar(150),
    netto real,
    date_from date,
    counterparty bigint,
    counterparty_id bigint,
    date_to date,
    w_id bigint,
    num bigint
);

CREATE TABLE IF NOT EXISTS sp14_doc (
    id bigserial PRIMARY KEY,
    org_name varchar(150),
    cargo_name varchar(100),
    dep_name varchar(150),
    tok_name varchar(150),
    warehouse_name varchar(150),
    sender_name varchar(150),
    rep_name varchar(150),
    sort_name varchar(150),
    talon varchar(50),
    fields varchar(250),
    driver_fio varchar(150),
    transport_model varchar(100),
    gos_num varchar(50),
    brutto real,
    netto real,
    tara real,
    combiner_name1 varchar(20),
    combiner_name2 varchar(20),
    combiner_name3 varchar(20),
    combiner_name4 varchar(20),
    combiner_name5 varchar(20),
    combiner_name6 varchar(20),
    combiner1 real,
    combiner2 real,
    combiner3 real,
    combiner4 real,
    combiner5 real,
    combiner6 real,
    combiner7 real,
    combiner8 real,
    combiner9 real,
    combiner10 real,
    combiner11 real,
    combiner12 real,
    combiner13 real,
    combiner14 real,
    combiner_talon1 varchar(20),
    combiner_talon2 varchar(20),
    combiner_talon3 varchar(20),
    combiner_talon4 varchar(20),
    combiner_talon5 varchar(20),
    combiner_talon6 varchar(20),
    combiner_talon7 varchar(20),
    combiner_talon8 varchar(20),
    combiner_talon9 varchar(20),
    combiner_talon10 varchar(20),
    combiner_talon11 varchar(20),
    combiner_talon12 varchar(20),
    combiner_talon13 varchar(20),
    combiner_talon14 varchar(20),
    combiner_talon_n varchar(20),
    date_from date,
    date_to date,
    w_id bigint,
    num bigint
);
