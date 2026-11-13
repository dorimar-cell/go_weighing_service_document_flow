-- Dev-снимок данных весовой на 13.11.2026 (тестовые рейсы, пропуска, фото-метаданные).
-- НЕ применять на продуктиве. Восстановление после миграций:
--   docker cp seeds/dev-data.sql weighing_postgres:/tmp/dev-data.sql
--   docker exec weighing_postgres psql -U weighing -d weighing_db -1 -f /tmp/dev-data.sql
--
-- PostgreSQL database dump
--

-- Dumped from database version 16.4 (Debian 16.4-1.pgdg120+2)
-- Dumped by pg_dump version 16.4 (Debian 16.4-1.pgdg120+2)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: dir_driver; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.dir_driver VALUES (1, 'Иванов Иван Иванович', '+7 (900) 100-20-10', '7701100137', NULL, NULL, NULL, NULL, false, NULL, '2026-11-12 16:36:17.87995+00', '2026-11-12 16:36:17.87995+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_driver VALUES (2, 'Петров Пётр Сергеевич', '+7 (900) 101-21-11', '7701100274', NULL, NULL, NULL, NULL, false, NULL, '2026-11-12 16:36:17.87995+00', '2026-11-12 16:36:17.87995+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_driver VALUES (3, 'Сидоров Алексей Николаевич', '+7 (900) 102-22-12', '7701100411', NULL, NULL, NULL, NULL, false, NULL, '2026-11-12 16:36:17.87995+00', '2026-11-12 16:36:17.87995+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_driver VALUES (6, 'Ivanov Ivan', NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, '2026-11-12 18:51:06.309389+00', '2026-11-12 18:51:06.309389+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_driver VALUES (7, 'E2E Ivanov', NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, '2026-11-12 18:55:55.341354+00', '2026-11-12 18:55:55.341354+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_driver VALUES (8, 'Ref Driver', NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, '2026-11-12 19:04:08.371704+00', '2026-11-12 19:04:08.371704+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: dir_trailer; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.dir_trailer VALUES (1, 'ЕА1234-77', false, NULL, '2026-11-12 16:36:17.874658+00', '2026-11-12 16:36:17.874658+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_trailer VALUES (2, 'ЕН5678-50', false, NULL, '2026-11-12 16:36:17.874658+00', '2026-11-12 16:36:17.874658+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_trailer VALUES (3, 'АК9012-99', false, NULL, '2026-11-12 16:36:17.874658+00', '2026-11-12 16:36:17.874658+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_trailer VALUES (6, 'EA1234-77', false, NULL, '2026-11-12 18:51:06.309389+00', '2026-11-12 18:51:06.309389+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: dir_transport; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.dir_transport VALUES (1, 'А123ВС77', false, NULL, '2026-11-12 16:36:17.868475+00', '2026-11-12 16:36:17.868475+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_transport VALUES (2, 'В456ЕК50', false, NULL, '2026-11-12 16:36:17.868475+00', '2026-11-12 16:36:17.868475+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_transport VALUES (3, 'М789ОР99', false, NULL, '2026-11-12 16:36:17.868475+00', '2026-11-12 16:36:17.868475+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_transport VALUES (4, 'Т321НА77', false, NULL, '2026-11-12 16:36:17.868475+00', '2026-11-12 16:36:17.868475+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_transport VALUES (5, 'К654РУ50', false, NULL, '2026-11-12 16:36:17.868475+00', '2026-11-12 16:36:17.868475+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_transport VALUES (8, 'A123BC77', false, NULL, '2026-11-12 18:51:06.309389+00', '2026-11-12 18:51:06.309389+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_transport VALUES (9, 'B456EK50', false, NULL, '2026-11-12 18:51:28.55445+00', '2026-11-12 18:51:28.55445+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_transport VALUES (10, 'E2E-A123BC77', false, NULL, '2026-11-12 18:55:55.341354+00', '2026-11-12 18:55:55.341354+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_transport VALUES (11, 'REF-A123BC77', false, NULL, '2026-11-12 19:04:08.371704+00', '2026-11-12 19:04:08.371704+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_transport VALUES (12, 'PASS-A123BC77', false, NULL, '2026-11-12 19:33:52.860892+00', '2026-11-12 19:33:52.860892+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_transport VALUES (13, 'FLOW-A123BC77', false, NULL, '2026-11-12 19:42:28.879128+00', '2026-11-12 19:42:28.879128+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_transport VALUES (14, 'PERM-A123BC77', false, NULL, '2026-11-12 20:05:26.424661+00', '2026-11-12 20:05:26.424661+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_transport VALUES (15, 'PHOTO-A123BC77', false, NULL, '2026-11-12 20:14:04.527374+00', '2026-11-12 20:14:04.527374+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_transport VALUES (16, 'J-A123BC77', false, NULL, '2026-11-13 05:20:58.270124+00', '2026-11-13 05:20:58.270124+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: dir_type_of_car_suspension; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_type_of_trailer; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_type_of_trailer_models; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_type_of_transport; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_type_of_transport_models; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: com_dir_transport_driver_trailer; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.com_dir_transport_driver_trailer VALUES (3, 6, 8, 6, false, NULL, NULL, 1, '2026-11-12 18:51:06.309389+00', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_transport_driver_trailer VALUES (4, NULL, 9, NULL, false, NULL, NULL, 1, '2026-11-12 18:51:28.55445+00', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_transport_driver_trailer VALUES (5, 7, 10, 6, false, NULL, NULL, 1, '2026-11-12 18:55:55.341354+00', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_transport_driver_trailer VALUES (6, 8, 11, 6, false, NULL, NULL, 1, '2026-11-12 19:04:08.371704+00', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_transport_driver_trailer VALUES (7, NULL, 12, NULL, false, NULL, NULL, 1, '2026-11-12 19:34:01.126323+00', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_transport_driver_trailer VALUES (8, NULL, 13, NULL, false, NULL, NULL, 1, '2026-11-12 19:42:37.001972+00', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_transport_driver_trailer VALUES (9, NULL, 15, NULL, false, NULL, NULL, 1, '2026-11-12 20:14:04.527374+00', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL) ON CONFLICT DO NOTHING;


--
-- Data for Name: dir_additional_characteristics; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_cargo; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.dir_cargo VALUES (1, 'Пшеница 3 класса', 12000.000, 14500.000, false, NULL, '2026-11-12 16:36:17.884986+00', '2026-11-12 16:36:17.884986+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_cargo VALUES (2, 'Кукуруза', 11000.000, 13200.000, false, NULL, '2026-11-12 16:36:17.884986+00', '2026-11-12 16:36:17.884986+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_cargo VALUES (3, 'Ячмень', 10500.000, 12800.000, false, NULL, '2026-11-12 16:36:17.884986+00', '2026-11-12 16:36:17.884986+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_cargo VALUES (4, 'Свёкла сахарная', 4500.000, 5600.000, false, NULL, '2026-11-12 16:36:17.884986+00', '2026-11-12 16:36:17.884986+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_cargo VALUES (5, 'Подсолнечник', 27000.000, 31000.000, false, NULL, '2026-11-12 16:36:17.884986+00', '2026-11-12 16:36:17.884986+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: dir_carrier; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.dir_carrier VALUES (1, 'ООО «АгроТранс»', '7704567890', '770401001', '+7 (495) 400-40-40', 'г. Москва, ул. Колхозная, д. 4', NULL, '2026-11-12 16:36:17.901132+00', '2026-11-12 16:36:17.901132+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_carrier VALUES (2, 'ООО «Зерновоз»', '7705678901', '770501001', '+7 (495) 500-50-50', 'г. Москва, ул. Полевая, д. 5', NULL, '2026-11-12 16:36:17.901132+00', '2026-11-12 16:36:17.901132+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: dir_counterparty; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.dir_counterparty VALUES (1, 'ООО «Русское поле»', '7701234567', '770101001', '+7 (495) 100-10-10', 'Иванов Иван Иванович', 'г. Москва, ул. Тверская, д. 1', NULL, '2026-11-12 16:36:17.896009+00', '2026-11-12 16:36:17.896009+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_counterparty VALUES (2, 'АО «Агрокомплекс»', '7702345678', '770201001', '+7 (495) 200-20-20', 'Петров Пётр Сергеевич', 'г. Москва, ул. Арбат, д. 2', NULL, '2026-11-12 16:36:17.896009+00', '2026-11-12 16:36:17.896009+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_counterparty VALUES (3, 'ООО «Зернотрейд»', '7703456789', '770301001', '+7 (495) 300-30-30', 'Сидоров Алексей Николаевич', 'г. Москва, ул. Садовая, д. 3', NULL, '2026-11-12 16:36:17.896009+00', '2026-11-12 16:36:17.896009+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: dir_fund; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_type_of_operation; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.dir_type_of_operation VALUES (1, 'Приход', 1, true, false, NULL, '2026-11-12 16:24:21.097771+00', '2026-11-12 16:24:21.097771+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_type_of_operation VALUES (2, 'Расход', 2, true, false, NULL, '2026-11-12 16:24:21.097771+00', '2026-11-12 16:24:21.097771+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_type_of_operation VALUES (3, 'Перемещение', 3, false, false, NULL, '2026-11-12 16:24:21.097771+00', '2026-11-12 16:24:21.097771+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: dir_invoice; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.dir_invoice VALUES (1, '', 1, 1, 1, '2026-11-12 00:00:00+00', '19:04:08.371704', NULL, NULL, NULL, false, 1, '2026-11-12 19:04:08.371704+00', NULL, NULL) ON CONFLICT DO NOTHING;


--
-- Data for Name: dir_kfkh_plantation; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_passport; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_placement; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_storage; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.dir_storage VALUES (1, 'Склад №1 (зерно)', 'ООО «Предприятие 1»', '+7 (900) 201-01-01', 'г. Москва, ул. Промышленная, д. 1', NULL, '2026-11-12 16:36:17.890756+00', '2026-11-12 16:36:17.890756+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_storage VALUES (2, 'Склад №2 (фураж)', 'ООО «Предприятие 1»', '+7 (900) 202-02-02', 'г. Москва, ул. Промышленная, д. 2', NULL, '2026-11-12 16:36:17.890756+00', '2026-11-12 16:36:17.890756+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_storage VALUES (3, 'Склад №3 (свёкла)', 'ООО «Предприятие 2»', '+7 (900) 203-03-03', 'г. Москва, ул. Промышленная, д. 3', NULL, '2026-11-12 16:36:17.890756+00', '2026-11-12 16:36:17.890756+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: dir_type_of_cargo; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_user_group; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.dir_user_group VALUES (1, 'Администраторы', true, NULL, '2026-11-12 18:49:56.448971+00', '2026-11-12 18:49:56.448971+00') ON CONFLICT DO NOTHING;
INSERT INTO public.dir_user_group VALUES (2, 'Операторы', false, NULL, '2026-11-12 20:05:11.123486+00', '2026-11-12 20:05:11.123486+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: dir_weight_room; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.dir_weight_room VALUES (1, 'Весовая №1', 'Пост по умолчанию', NULL, '2026-11-12 18:49:56.448971+00', '2026-11-12 18:49:56.448971+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: dir_users; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.dir_users VALUES (1, 'system', NULL, NULL, false, 1, 1, '2026-11-12 18:49:56.448971+00', NULL, NULL, 1, '64e94b28-1ffa-45d8-964b-0c57c768668e', NULL, 'system') ON CONFLICT DO NOTHING;


--
-- Data for Name: system_dir_identifier_state; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.system_dir_identifier_state VALUES (1, 'Действителен', NULL, '2026-11-12 16:24:21.097771+00', '2026-11-12 16:24:21.097771+00') ON CONFLICT DO NOTHING;
INSERT INTO public.system_dir_identifier_state VALUES (2, 'Просрочен', NULL, '2026-11-12 16:24:21.097771+00', '2026-11-12 16:24:21.097771+00') ON CONFLICT DO NOTHING;
INSERT INTO public.system_dir_identifier_state VALUES (3, 'Выдан', NULL, '2026-11-12 16:24:21.097771+00', '2026-11-12 16:24:21.097771+00') ON CONFLICT DO NOTHING;
INSERT INTO public.system_dir_identifier_state VALUES (4, 'Не выдан', NULL, '2026-11-12 16:24:21.097771+00', '2026-11-12 16:24:21.097771+00') ON CONFLICT DO NOTHING;
INSERT INTO public.system_dir_identifier_state VALUES (5, 'Деактивирован', NULL, '2026-11-12 16:24:21.097771+00', '2026-11-12 16:24:21.097771+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: system_dir_type_of_pass; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.system_dir_type_of_pass VALUES (1, 'Постоянный', NULL, '2026-11-12 16:24:21.097771+00', '2026-11-12 16:24:21.097771+00') ON CONFLICT DO NOTHING;
INSERT INTO public.system_dir_type_of_pass VALUES (2, 'Временный', NULL, '2026-11-12 16:24:21.097771+00', '2026-11-12 16:24:21.097771+00') ON CONFLICT DO NOTHING;
INSERT INTO public.system_dir_type_of_pass VALUES (3, 'Разовый', NULL, '2026-11-12 16:24:21.097771+00', '2026-11-12 16:24:21.097771+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: system_dir_type_of_regweight; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.system_dir_type_of_regweight VALUES (1, 'Ключ') ON CONFLICT DO NOTHING;
INSERT INTO public.system_dir_type_of_regweight VALUES (2, 'Распознавание номера') ON CONFLICT DO NOTHING;
INSERT INTO public.system_dir_type_of_regweight VALUES (3, 'Штрихкод') ON CONFLICT DO NOTHING;
INSERT INTO public.system_dir_type_of_regweight VALUES (4, 'Въезд/выезд') ON CONFLICT DO NOTHING;


--
-- Data for Name: reg_weight; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.reg_weight VALUES (1, 'E2E-P-001', NULL, NULL, NULL, false, false, NULL, NULL, NULL, NULL, NULL, NULL, 1, 1, 3, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 12, NULL, NULL, NULL, NULL, NULL, NULL, false, 1, '2026-11-12 19:33:52.860892+00', NULL, NULL, NULL, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.reg_weight VALUES (2, 'E2E-FLOW-001', NULL, NULL, NULL, false, false, NULL, NULL, NULL, NULL, NULL, NULL, 1, 1, 3, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 13, NULL, NULL, NULL, NULL, NULL, NULL, false, 1, '2026-11-12 19:42:28.879128+00', NULL, NULL, NULL, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.reg_weight VALUES (3, 'PERM-TEST-001', NULL, NULL, NULL, true, false, NULL, NULL, NULL, NULL, NULL, NULL, 1, 1, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 14, NULL, NULL, NULL, NULL, NULL, NULL, false, 1, '2026-11-12 20:05:26.424661+00', NULL, NULL, NULL, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.reg_weight VALUES (4, 'PASS-J-001', NULL, NULL, NULL, true, false, NULL, NULL, NULL, NULL, NULL, NULL, 1, 1, 2, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 16, NULL, NULL, NULL, NULL, NULL, NULL, false, 1, '2026-11-13 05:20:58.270124+00', 1, '2026-11-13 05:20:58.47036+00', NULL, NULL) ON CONFLICT DO NOTHING;


--
-- Data for Name: weighing_info; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.weighing_info VALUES (2, NULL, NULL, NULL, false, 3, NULL, 1, '2026-11-12 18:51:06.309389+00', NULL, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_info VALUES (3, NULL, NULL, NULL, true, 4, NULL, 1, '2026-11-12 18:51:28.55445+00', NULL, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_info VALUES (4, NULL, NULL, NULL, false, 5, NULL, 1, '2026-11-12 18:55:55.341354+00', NULL, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_info VALUES (5, 1, NULL, NULL, false, 6, NULL, 1, '2026-11-12 19:04:08.371704+00', NULL, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_info VALUES (6, NULL, NULL, NULL, false, 7, NULL, 1, '2026-11-12 19:34:01.126323+00', NULL, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_info VALUES (7, NULL, NULL, NULL, false, 8, NULL, 1, '2026-11-12 19:42:37.001972+00', NULL, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_info VALUES (8, NULL, NULL, NULL, false, 9, NULL, 1, '2026-11-12 20:14:04.527374+00', NULL, NULL) ON CONFLICT DO NOTHING;


--
-- Data for Name: weighing; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.weighing VALUES (1, 2, 42339.6, NULL, NULL, true, false, 18336.4, NULL, 1, '2026-11-12 18:51:06.309389+00', 1, '2026-11-12 18:51:17.448777+00', 1, 1, NULL, '2026-11-12 18:51:06.211441+00', '2026-11-12 18:51:17.200404+00', 'd5e73238-c128-4d12-9669-dc304d56c15e', false, false, false, false) ON CONFLICT DO NOTHING;
INSERT INTO public.weighing VALUES (2, 3, 24999.8, NULL, NULL, false, true, NULL, NULL, 1, '2026-11-12 18:51:28.55445+00', 1, '2026-11-12 18:51:28.562519+00', 1, NULL, NULL, '2026-11-12 18:51:28.210128+00', NULL, '65b8f60a-2e44-40bc-9618-c997f94b934a', false, false, false, false) ON CONFLICT DO NOTHING;
INSERT INTO public.weighing VALUES (3, 4, 42338.7, NULL, NULL, true, false, 18340.8, NULL, 1, '2026-11-12 18:55:55.341354+00', 1, '2026-11-12 18:56:06.44554+00', 1, 1, NULL, '2026-11-12 18:55:54.865616+00', '2026-11-12 18:56:06.360894+00', 'e04ddf4f-b8df-4c4c-bd30-8fa3f42c2bb8', false, false, false, false) ON CONFLICT DO NOTHING;
INSERT INTO public.weighing VALUES (4, 5, 42338.6, NULL, NULL, true, false, 18337.8, NULL, 1, '2026-11-12 19:04:08.371704+00', 1, '2026-11-12 19:04:19.863389+00', 1, 1, NULL, '2026-11-12 19:04:08.176512+00', '2026-11-12 19:04:19.6791+00', 'ca83d05a-6266-4e03-9927-722ee3ef6523', false, false, false, false) ON CONFLICT DO NOTHING;
INSERT INTO public.weighing VALUES (5, 6, 42342.7, NULL, NULL, true, false, 18338.7, NULL, 1, '2026-11-12 19:34:01.126323+00', 1, '2026-11-12 19:34:12.189091+00', 1, 1, 1, '2026-11-12 19:34:00.611953+00', '2026-11-12 19:34:12.111373+00', '9ba176a8-9862-4cfd-99c4-f452c843142e', false, false, false, false) ON CONFLICT DO NOTHING;
INSERT INTO public.weighing VALUES (6, 7, 42339.6, NULL, NULL, true, false, 15000, 15000, 1, '2026-11-12 19:42:37.001972+00', 1, '2026-11-12 19:42:37.032713+00', 1, 1, 2, '2026-11-12 19:42:36.830163+00', '2026-11-12 19:42:37.032634+00', '27415cc4-79c6-4f01-b623-9602fa7bf8dc', false, false, false, false) ON CONFLICT DO NOTHING;
INSERT INTO public.weighing VALUES (7, 8, 42339.1, NULL, NULL, true, false, 18343.8, NULL, 1, '2026-11-12 20:14:04.527374+00', 1, '2026-11-12 20:14:15.779508+00', 1, 1, NULL, '2026-11-12 20:14:04.425137+00', '2026-11-12 20:14:15.423236+00', '647cd02b-134e-4bd1-93f7-886f7ae27eb5', false, false, false, false) ON CONFLICT DO NOTHING;


--
-- Data for Name: acceptance_certificates; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: audit_weighing; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_video_surveillance; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.settings_video_surveillance VALUES (1, true, 'Emulator', '127.0.0.1', NULL, 0, NULL, NULL, 3000, 3, 0, 0, 1, 1, '2026-11-12 20:13:20.959528+00', 1) ON CONFLICT DO NOTHING;


--
-- Data for Name: audit_weighting_photo; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: auto_invoice; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: auto_invoice_everyday; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_load_point; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_provider; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_type_of_pack; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_upload_point; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: com_dir_invoice_dir_storage; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.com_dir_invoice_dir_storage VALUES (1, 1, NULL, 1, 1, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, 1, '2026-11-12 19:04:08.371704+00', NULL, NULL, NULL, NULL) ON CONFLICT DO NOTHING;


--
-- Data for Name: system_dir_access_permissions; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.system_dir_access_permissions VALUES (1, 'Журнал взвешиваний', 'W005') ON CONFLICT DO NOTHING;
INSERT INTO public.system_dir_access_permissions VALUES (2, 'Управление пропусками', 'W007') ON CONFLICT DO NOTHING;
INSERT INTO public.system_dir_access_permissions VALUES (3, 'Первое взвешивание', 'W001') ON CONFLICT DO NOTHING;
INSERT INTO public.system_dir_access_permissions VALUES (4, 'Поточное взвешивание', 'W003') ON CONFLICT DO NOTHING;
INSERT INTO public.system_dir_access_permissions VALUES (5, 'Отмена взвешивания', 'W004') ON CONFLICT DO NOTHING;
INSERT INTO public.system_dir_access_permissions VALUES (6, 'Второе взвешивание', 'W002') ON CONFLICT DO NOTHING;
INSERT INTO public.system_dir_access_permissions VALUES (7, 'Экспорт журнала', 'W006') ON CONFLICT DO NOTHING;


--
-- Data for Name: com_dir_user_groups_access_permissions; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.com_dir_user_groups_access_permissions VALUES (1, 1, 1) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_user_groups_access_permissions VALUES (2, 1, 2) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_user_groups_access_permissions VALUES (3, 1, 3) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_user_groups_access_permissions VALUES (4, 1, 4) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_user_groups_access_permissions VALUES (5, 1, 5) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_user_groups_access_permissions VALUES (6, 1, 6) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_user_groups_access_permissions VALUES (7, 1, 7) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_user_groups_access_permissions VALUES (8, 2, 2) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_user_groups_access_permissions VALUES (9, 2, 6) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_user_groups_access_permissions VALUES (10, 2, 4) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_user_groups_access_permissions VALUES (11, 2, 1) ON CONFLICT DO NOTHING;
INSERT INTO public.com_dir_user_groups_access_permissions VALUES (12, 2, 3) ON CONFLICT DO NOTHING;


--
-- Data for Name: dir_kfkh_combine; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_kfkh_combiner; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: com_kfkh; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_axle_real_weights; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_default_settings_for_print; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_laboratory; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_deleted_laboratories_records; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_export_file; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_lab_register; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_refusing; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_spisanie_ostatkov; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_storage_ostatki; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_table_presentation; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_terminal_axle; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_trailer_model_axle_characters; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dir_transport_model_axle_characters; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: dp_messages; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: invoice_document; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: invoice_torg12_documents; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: log_start_terminal_info; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: reference_books_table; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: rus_name_field; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: setting_backup; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_app; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_audit_load; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_axle_load; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_bar_code; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_c1_integrator; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_dp; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_email; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_export; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_integration; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_kfkh; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_laboratory; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_laboratory2; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_pre_axle_report; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_recognition_m_pixel; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_recognition_module; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_scoreboard; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_server_report; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_skud; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_speech; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_store; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_terminal; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: settings_usb_key_reader; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: sp11_doc; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: sp14_doc; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: sp9_doc; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: system_dir_direction_type_operation; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: system_dir_kind_of_packing; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: system_dir_type_of_carriage; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: system_dir_type_of_carriage_document; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: ttn_document; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: unrec_weighing; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: unrec_weighting_photo; Type: TABLE DATA; Schema: public; Owner: weighing
--



--
-- Data for Name: weighing_measurement; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.weighing_measurement VALUES (1, 1, 1, '324674d9-489a-4aaa-b204-4048bcf4b055', 42339.600, true, 'Emulator', 1, '2026-11-12 18:51:06.211441+00', '20ac1ffba7d98bacd6109706f509810dec5173b569bb087bd980c22c9de95081', '2026-11-12 18:51:06.309389+00') ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_measurement VALUES (2, 1, 2, '13d165da-3a87-4b92-ad6c-205c9d36926d', 18336.400, true, 'Emulator', 1, '2026-11-12 18:51:17.200404+00', '37b0e96f1aed58e0b2ab26ff3d554ab88f041ee31475e4cada08338869c8f4dd', '2026-11-12 18:51:17.448777+00') ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_measurement VALUES (3, 2, 1, 'a564f50a-16ed-487b-8a56-546e71fb1f51', 24999.800, true, 'Emulator', 1, '2026-11-12 18:51:28.210128+00', 'f3895a66300bb15bd5e43e21e19faf7d7b8d69053d944fff817b7647f413196c', '2026-11-12 18:51:28.55445+00') ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_measurement VALUES (4, 3, 1, '5025c979-92f0-4d77-9df1-432a5fc7275b', 42338.700, true, 'Emulator', 1, '2026-11-12 18:55:54.865616+00', 'c0a1761e3fd835ecb0c003f0bee5ab3bb95bfe50a7bd96ea7259d1ed33d248b3', '2026-11-12 18:55:55.341354+00') ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_measurement VALUES (5, 3, 2, '59ed6fad-5686-457d-a56e-a05e6e738cee', 18340.800, true, 'Emulator', 1, '2026-11-12 18:56:06.360894+00', '02ed2eb6f1e1f4712016576d5d0a86114bb327d5c434c63f44fdfa14371902bb', '2026-11-12 18:56:06.44554+00') ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_measurement VALUES (6, 4, 1, '5a5711fb-3d98-47ce-94fb-6f35ec043bf6', 42338.600, true, 'Emulator', 1, '2026-11-12 19:04:08.176512+00', 'b2277ec58bfa9313dc938d715c08950665dc34f8a0274ff6a52865e4ca93cf72', '2026-11-12 19:04:08.371704+00') ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_measurement VALUES (7, 4, 2, '253323cf-0e7a-4eff-bf9b-91ee3085613b', 18337.800, true, 'Emulator', 1, '2026-11-12 19:04:19.6791+00', '0c859ae78a4c86d027975a7aad8b3460e4b3c705f564af3eb02b21fd1c22075a', '2026-11-12 19:04:19.863389+00') ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_measurement VALUES (8, 5, 1, '1330ec74-4e2b-4e09-b7dd-2fe77e2b7ce4', 42342.700, true, 'Emulator', 1, '2026-11-12 19:34:00.611953+00', 'a5decf6cdc887a5565827511106c2be4bb5a72ae8b6619bdd8103c0f23740add', '2026-11-12 19:34:01.126323+00') ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_measurement VALUES (9, 5, 2, '9f741bbc-0ad8-4eab-8ea3-0a5977e34c51', 18338.700, true, 'Emulator', 1, '2026-11-12 19:34:12.111373+00', '7b13caff27549bb0bee706aa3ea86dbf1f3a1a1968bd7cad412f95597300631e', '2026-11-12 19:34:12.189091+00') ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_measurement VALUES (10, 6, 1, '25a01180-b2de-4a68-8b99-ab675aa3fb5b', 42339.600, true, 'Emulator', 1, '2026-11-12 19:42:36.830163+00', 'b1cc21f8d6eea85ead6f3e2ec6ee6287f5e6993d9acb16077acd100f88879fd9', '2026-11-12 19:42:37.001972+00') ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_measurement VALUES (11, 7, 1, '135c80af-ba4e-4ce9-95b4-8465a8cf11b0', 42339.100, true, 'Emulator', 1, '2026-11-12 20:14:04.425137+00', '873fc7518b09cfa47330b69ca0cef2c26df42d5645cd799c4eb640dc3dcd788c', '2026-11-12 20:14:04.527374+00') ON CONFLICT DO NOTHING;
INSERT INTO public.weighing_measurement VALUES (12, 7, 2, '056dc7cd-804b-459c-943c-acd018e78aaf', 18343.800, true, 'Emulator', 1, '2026-11-12 20:14:15.423236+00', '88e7e6fbfcf50b5b9a15599f3965fb3821b57b0d316707c4f07a991f36a0a49b', '2026-11-12 20:14:15.779508+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: weighting_photo; Type: TABLE DATA; Schema: public; Owner: weighing
--

INSERT INTO public.weighting_photo VALUES (1, 7, '7/ebed3e0c12a5e048.jpg', 1, '2026-11-12 20:14:04.70161+00', true, 1) ON CONFLICT DO NOTHING;
INSERT INTO public.weighting_photo VALUES (2, 7, '7/f9936432d4a5cf27.jpg', 1, '2026-11-12 20:14:15.809421+00', false, 1) ON CONFLICT DO NOTHING;


--
-- Name: acceptance_certificates_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.acceptance_certificates_id_seq', 1, false);


--
-- Name: audit_weighing_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.audit_weighing_id_seq', 1, false);


--
-- Name: audit_weighting_photo_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.audit_weighting_photo_id_seq', 1, false);


--
-- Name: auto_invoice_everyday_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.auto_invoice_everyday_id_seq', 1, false);


--
-- Name: auto_invoice_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.auto_invoice_id_seq', 1, false);


--
-- Name: com_dir_invoice_dir_storage_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.com_dir_invoice_dir_storage_id_seq', 1, true);


--
-- Name: com_dir_transport_driver_trailer_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.com_dir_transport_driver_trailer_id_seq', 9, true);


--
-- Name: com_dir_user_groups_access_permissions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.com_dir_user_groups_access_permissions_id_seq', 12, true);


--
-- Name: com_kfkh_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.com_kfkh_id_seq', 1, false);


--
-- Name: dir_additional_characteristics_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_additional_characteristics_id_seq', 1, false);


--
-- Name: dir_axle_real_weights_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_axle_real_weights_id_seq', 1, false);


--
-- Name: dir_cargo_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_cargo_id_seq', 5, true);


--
-- Name: dir_carrier_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_carrier_id_seq', 2, true);


--
-- Name: dir_counterparty_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_counterparty_id_seq', 3, true);


--
-- Name: dir_default_settings_for_print_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_default_settings_for_print_id_seq', 1, false);


--
-- Name: dir_deleted_laboratories_records_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_deleted_laboratories_records_id_seq', 1, false);


--
-- Name: dir_driver_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_driver_id_seq', 8, true);


--
-- Name: dir_export_file_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_export_file_id_seq', 1, false);


--
-- Name: dir_fund_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_fund_id_seq', 1, false);


--
-- Name: dir_invoice_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_invoice_id_seq', 1, true);


--
-- Name: dir_kfkh_combine_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_kfkh_combine_id_seq', 1, false);


--
-- Name: dir_kfkh_combiner_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_kfkh_combiner_id_seq', 1, false);


--
-- Name: dir_kfkh_plantation_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_kfkh_plantation_id_seq', 1, false);


--
-- Name: dir_lab_register_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_lab_register_id_seq', 1, false);


--
-- Name: dir_laboratory_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_laboratory_id_seq', 1, false);


--
-- Name: dir_load_point_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_load_point_id_seq', 1, false);


--
-- Name: dir_passport_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_passport_id_seq', 1, false);


--
-- Name: dir_placement_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_placement_id_seq', 1, false);


--
-- Name: dir_provider_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_provider_id_seq', 1, false);


--
-- Name: dir_refusing_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_refusing_id_seq', 1, false);


--
-- Name: dir_spisanie_ostatkov_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_spisanie_ostatkov_id_seq', 1, false);


--
-- Name: dir_storage_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_storage_id_seq', 3, true);


--
-- Name: dir_storage_ostatki_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_storage_ostatki_id_seq', 1, false);


--
-- Name: dir_table_presentation_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_table_presentation_id_seq', 1, false);


--
-- Name: dir_terminal_axle_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_terminal_axle_id_seq', 1, false);


--
-- Name: dir_trailer_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_trailer_id_seq', 6, true);


--
-- Name: dir_trailer_model_axle_characters_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_trailer_model_axle_characters_id_seq', 1, false);


--
-- Name: dir_transport_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_transport_id_seq', 16, true);


--
-- Name: dir_transport_model_axle_characters_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_transport_model_axle_characters_id_seq', 1, false);


--
-- Name: dir_type_of_car_suspension_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_type_of_car_suspension_id_seq', 1, false);


--
-- Name: dir_type_of_cargo_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_type_of_cargo_id_seq', 1, false);


--
-- Name: dir_type_of_operation_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_type_of_operation_id_seq', 3, true);


--
-- Name: dir_type_of_pack_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_type_of_pack_id_seq', 1, false);


--
-- Name: dir_type_of_trailer_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_type_of_trailer_id_seq', 1, false);


--
-- Name: dir_type_of_trailer_models_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_type_of_trailer_models_id_seq', 1, false);


--
-- Name: dir_type_of_transport_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_type_of_transport_id_seq', 1, false);


--
-- Name: dir_type_of_transport_models_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_type_of_transport_models_id_seq', 1, false);


--
-- Name: dir_upload_point_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_upload_point_id_seq', 1, false);


--
-- Name: dir_user_group_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_user_group_id_seq', 2, true);


--
-- Name: dir_users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_users_id_seq', 1, true);


--
-- Name: dir_weight_room_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dir_weight_room_id_seq', 1, true);


--
-- Name: dp_messages_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.dp_messages_id_seq', 1, false);


--
-- Name: invoice_document_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.invoice_document_id_seq', 1, false);


--
-- Name: invoice_torg12_documents_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.invoice_torg12_documents_id_seq', 1, false);


--
-- Name: log_start_terminal_info_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.log_start_terminal_info_id_seq', 1, false);


--
-- Name: reference_books_table_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.reference_books_table_id_seq', 1, false);


--
-- Name: reg_weight_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.reg_weight_id_seq', 4, true);


--
-- Name: rus_name_field_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.rus_name_field_id_seq', 1, false);


--
-- Name: setting_backup_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.setting_backup_id_seq', 1, false);


--
-- Name: settings_app_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_app_id_seq', 1, false);


--
-- Name: settings_audit_load_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_audit_load_id_seq', 1, false);


--
-- Name: settings_axle_load_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_axle_load_id_seq', 1, false);


--
-- Name: settings_bar_code_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_bar_code_id_seq', 1, false);


--
-- Name: settings_c1_integrator_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_c1_integrator_id_seq', 1, false);


--
-- Name: settings_dp_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_dp_id_seq', 1, false);


--
-- Name: settings_email_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_email_id_seq', 1, false);


--
-- Name: settings_export_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_export_id_seq', 1, false);


--
-- Name: settings_integration_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_integration_id_seq', 1, false);


--
-- Name: settings_kfkh_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_kfkh_id_seq', 1, false);


--
-- Name: settings_laboratory2_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_laboratory2_id_seq', 1, false);


--
-- Name: settings_laboratory_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_laboratory_id_seq', 1, false);


--
-- Name: settings_pre_axle_report_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_pre_axle_report_id_seq', 1, false);


--
-- Name: settings_recognition_m_pixel_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_recognition_m_pixel_id_seq', 1, false);


--
-- Name: settings_recognition_module_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_recognition_module_id_seq', 1, false);


--
-- Name: settings_scoreboard_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_scoreboard_id_seq', 1, false);


--
-- Name: settings_server_report_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_server_report_id_seq', 1, false);


--
-- Name: settings_skud_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_skud_id_seq', 1, false);


--
-- Name: settings_speech_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_speech_id_seq', 1, false);


--
-- Name: settings_store_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_store_id_seq', 1, false);


--
-- Name: settings_terminal_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_terminal_id_seq', 1, false);


--
-- Name: settings_usb_key_reader_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_usb_key_reader_id_seq', 1, false);


--
-- Name: settings_video_surveillance_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.settings_video_surveillance_id_seq', 1, true);


--
-- Name: sp11_doc_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.sp11_doc_id_seq', 1, false);


--
-- Name: sp14_doc_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.sp14_doc_id_seq', 1, false);


--
-- Name: sp9_doc_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.sp9_doc_id_seq', 1, false);


--
-- Name: system_dir_access_permissions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.system_dir_access_permissions_id_seq', 7, true);


--
-- Name: system_dir_direction_type_operation_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.system_dir_direction_type_operation_id_seq', 1, false);


--
-- Name: system_dir_identifier_state_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.system_dir_identifier_state_id_seq', 5, true);


--
-- Name: system_dir_kind_of_packing_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.system_dir_kind_of_packing_id_seq', 1, false);


--
-- Name: system_dir_type_of_carriage_document_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.system_dir_type_of_carriage_document_id_seq', 1, false);


--
-- Name: system_dir_type_of_carriage_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.system_dir_type_of_carriage_id_seq', 1, false);


--
-- Name: system_dir_type_of_pass_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.system_dir_type_of_pass_id_seq', 3, true);


--
-- Name: system_dir_type_of_regweight_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.system_dir_type_of_regweight_id_seq', 4, true);


--
-- Name: ttn_document_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.ttn_document_id_seq', 1, false);


--
-- Name: unrec_weighing_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.unrec_weighing_id_seq', 1, false);


--
-- Name: unrec_weighting_photo_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.unrec_weighting_photo_id_seq', 1, false);


--
-- Name: weighing_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.weighing_id_seq', 7, true);


--
-- Name: weighing_info_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.weighing_info_id_seq', 8, true);


--
-- Name: weighing_measurement_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.weighing_measurement_id_seq', 12, true);


--
-- Name: weighting_photo_id_seq; Type: SEQUENCE SET; Schema: public; Owner: weighing
--

SELECT pg_catalog.setval('public.weighting_photo_id_seq', 2, true);


--
-- PostgreSQL database dump complete
--

