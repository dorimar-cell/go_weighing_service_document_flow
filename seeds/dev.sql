-- Демонстрационные данные для локальной разработки.
-- НЕ применять на продуктиве. Запуск:
--   Get-Content seeds/dev.sql -Raw | docker exec -i weighing_postgres psql -U weighing -d weighing_db

INSERT INTO dir_transport (name, black_list)
SELECT value, false FROM (VALUES ('А123ВС77'), ('В456ЕК50'), ('М789ОР99'), ('Т321НА77'), ('К654РУ50')) AS seed(value)
WHERE NOT EXISTS (SELECT 1 FROM dir_transport);

INSERT INTO dir_trailer (name, black_list)
SELECT value, false FROM (VALUES ('ЕА1234-77'), ('ЕН5678-50'), ('АК9012-99')) AS seed(value)
WHERE NOT EXISTS (SELECT 1 FROM dir_trailer);

INSERT INTO dir_driver (name, phone, inn, black_list)
SELECT name, phone, inn, false
FROM (VALUES
    ('Иванов Иван Иванович', '+7 (900) 100-20-10', '7701100137'),
    ('Петров Пётр Сергеевич', '+7 (900) 101-21-11', '7701100274'),
    ('Сидоров Алексей Николаевич', '+7 (900) 102-22-12', '7701100411')
) AS seed(name, phone, inn)
WHERE NOT EXISTS (SELECT 1 FROM dir_driver);

INSERT INTO dir_cargo (name, price_per_tonne, selling_price, is_danger)
SELECT name, price, selling, false
FROM (VALUES
    ('Пшеница 3 класса', 12000, 14500),
    ('Кукуруза', 11000, 13200),
    ('Ячмень', 10500, 12800),
    ('Свёкла сахарная', 4500, 5600),
    ('Подсолнечник', 27000, 31000)
) AS seed(name, price, selling)
WHERE NOT EXISTS (SELECT 1 FROM dir_cargo);

INSERT INTO dir_storage (name, company, phone, address)
SELECT name, company, phone, address
FROM (VALUES
    ('Склад №1 (зерно)', 'ООО «Предприятие 1»', '+7 (900) 201-01-01', 'г. Москва, ул. Промышленная, д. 1'),
    ('Склад №2 (фураж)', 'ООО «Предприятие 1»', '+7 (900) 202-02-02', 'г. Москва, ул. Промышленная, д. 2'),
    ('Склад №3 (свёкла)', 'ООО «Предприятие 2»', '+7 (900) 203-03-03', 'г. Москва, ул. Промышленная, д. 3')
) AS seed(name, company, phone, address)
WHERE NOT EXISTS (SELECT 1 FROM dir_storage);

INSERT INTO dir_counterparty (name, inn, kpp, phone, director, address)
SELECT name, inn, kpp, phone, director, address
FROM (VALUES
    ('ООО «Русское поле»', '7701234567', '770101001', '+7 (495) 100-10-10', 'Иванов Иван Иванович', 'г. Москва, ул. Тверская, д. 1'),
    ('АО «Агрокомплекс»', '7702345678', '770201001', '+7 (495) 200-20-20', 'Петров Пётр Сергеевич', 'г. Москва, ул. Арбат, д. 2'),
    ('ООО «Зернотрейд»', '7703456789', '770301001', '+7 (495) 300-30-30', 'Сидоров Алексей Николаевич', 'г. Москва, ул. Садовая, д. 3')
) AS seed(name, inn, kpp, phone, director, address)
WHERE NOT EXISTS (SELECT 1 FROM dir_counterparty);

INSERT INTO dir_carrier (name, inn, kpp, phone, address)
SELECT name, inn, kpp, phone, address
FROM (VALUES
    ('ООО «АгроТранс»', '7704567890', '770401001', '+7 (495) 400-40-40', 'г. Москва, ул. Колхозная, д. 4'),
    ('ООО «Зерновоз»', '7705678901', '770501001', '+7 (495) 500-50-50', 'г. Москва, ул. Полевая, д. 5')
) AS seed(name, inn, kpp, phone, address)
WHERE NOT EXISTS (SELECT 1 FROM dir_carrier);
