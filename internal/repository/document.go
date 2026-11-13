package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go_weighing_service/internal/dto"
)

// DocumentRepository — данные для отчётов и печатных документов.
type DocumentRepository struct {
	pool *pgxpool.Pool
}

func NewDocumentRepository(pool *pgxpool.Pool) *DocumentRepository {
	return &DocumentRepository{pool: pool}
}

const sessionDetailsQuery = `
SELECT w.id,
       t.name,
       COALESCE(tr.name, ''),
       COALESCE(d.name, ''),
       w.first_weight::float8,
       w.second_weight::float8,
       w.date_time_first_weight,
       w.date_time_second_weight,
       w.flag_weight,
       w.deleted,
       w.created_date,
       COALESCE(cp.id, 0),
       COALESCE(cp.name, ''),
       COALESCE(cr.id, 0),
       COALESCE(cr.name, ''),
       COALESCE(cg.id, 0),
       COALESCE(cg.name, ''),
       COALESCE(st.id, 0),
       COALESCE(st.name, ''),
       COALESCE(ot.id, 0),
       COALESCE(ot.name, ''),
       COALESCE(rw.id, 0),
       COALESCE(rw.name, ''),
       COALESCE(concat_ws(' ', d.license_series, d.license_number), ''),
       COALESCE(cp.inn, ''),
       COALESCE(cp.kpp, ''),
       COALESCE(cp.address, ''),
       COALESCE(cp.director, ''),
       COALESCE(cr.inn, ''),
       COALESCE(cr.kpp, ''),
       COALESCE(cr.address, ''),
       cg.price_per_tonne,
       cg.selling_price,
       COALESCE(u.name, '')
FROM weighing w
JOIN weighing_info wi ON wi.id = w.id_weighting_info
JOIN com_dir_transport_driver_trailer c ON c.id = wi.id_com_transport_driver_trailer
JOIN dir_transport t ON t.id = c.id_dir_transport
LEFT JOIN dir_trailer tr ON tr.id = c.id_dir_trailer
LEFT JOIN dir_driver d ON d.id = c.id_dir_driver
LEFT JOIN reg_weight rw ON rw.id = w.id_reg_weight
LEFT JOIN dir_invoice inv ON inv.id = wi.id_dir_invoice
LEFT JOIN dir_carrier cr ON cr.id = inv.id_dir_carrier
LEFT JOIN dir_counterparty cp ON cp.id = inv.id_dir_counterparty
LEFT JOIN dir_type_of_operation ot ON ot.id = inv.id_dir_type_of_operation
LEFT JOIN dir_users u ON u.id = COALESCE(w.id_second_weight_user, w.id_first_weight_user)
LEFT JOIN LATERAL (
    SELECT cg.id, cg.name, cg.price_per_tonne, cg.selling_price, st.id AS storage_id
    FROM com_dir_invoice_dir_storage cids
    LEFT JOIN dir_cargo cg ON cg.id = cids.id_dir_cargo
    LEFT JOIN dir_storage st ON st.id = cids.id_dir_storage
    WHERE cids.id_dir_invoice = inv.id AND cids.deleted = false
    ORDER BY cids.id
    LIMIT 1
) cg ON true
LEFT JOIN dir_storage st ON st.id = cg.storage_id
WHERE w.id = $1`

// SessionDetails возвращает расширенные реквизиты рейса для бланков.
func (r *DocumentRepository) SessionDetails(ctx context.Context, id int64) (*dto.SessionDetails, error) {
	row := r.pool.QueryRow(ctx, sessionDetailsQuery, id)

	var details dto.SessionDetails
	var session dto.SessionDTO
	var secondWeight sql.NullFloat64
	var secondAt sql.NullTime
	var flagWeight, deleted bool
	var counterpartyID, carrierID, cargoID, storageID, operationTypeID, passID int64
	var pricePerTonne, sellingPrice sql.NullFloat64

	err := row.Scan(
		&session.ID,
		&session.TransportNumber,
		&session.TrailerNumber,
		&session.DriverName,
		&session.FirstWeightKg,
		&secondWeight,
		&session.FirstMeasuredAt,
		&secondAt,
		&flagWeight,
		&deleted,
		&session.CreatedAt,
		&counterpartyID,
		&session.CounterpartyName,
		&carrierID,
		&session.CarrierName,
		&cargoID,
		&session.CargoName,
		&storageID,
		&session.StorageName,
		&operationTypeID,
		&session.OperationTypeName,
		&passID,
		&session.PassNumber,
		&details.DriverLicense,
		&details.CounterpartyINN,
		&details.CounterpartyKPP,
		&details.CounterpartyAddress,
		&details.CounterpartyDirector,
		&details.CarrierINN,
		&details.CarrierKPP,
		&details.CarrierAddress,
		&pricePerTonne,
		&sellingPrice,
		&details.OperatorName,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("session details: %w", err)
	}

	if counterpartyID != 0 {
		session.CounterpartyID = &counterpartyID
	}
	if carrierID != 0 {
		session.CarrierID = &carrierID
	}
	if cargoID != 0 {
		session.CargoID = &cargoID
	}
	if storageID != 0 {
		session.StorageID = &storageID
	}
	if operationTypeID != 0 {
		session.OperationTypeID = &operationTypeID
	}
	if passID != 0 {
		session.PassID = &passID
	}
	if secondWeight.Valid {
		value := secondWeight.Float64
		session.SecondWeightKg = &value
		session.NetWeightKg = ptrFloat(math.Abs(session.FirstWeightKg - value))
	}
	if secondAt.Valid {
		value := secondAt.Time
		session.SecondMeasuredAt = &value
	}
	if pricePerTonne.Valid {
		details.CargoPricePerTonne = ptrFloat(pricePerTonne.Float64)
	}
	if sellingPrice.Valid {
		details.CargoSellingPrice = ptrFloat(sellingPrice.Float64)
	}

	switch {
	case deleted:
		session.Status = "cancelled"
	case flagWeight:
		session.Status = "closed"
	default:
		session.Status = "open"
	}

	details.Session = session
	return &details, nil
}

// PrintSettings — реквизиты организации из dir_default_settings_for_print.
func (r *DocumentRepository) PrintSettings(ctx context.Context) (map[string]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT type_of_field, field_value FROM dir_default_settings_for_print`)
	if err != nil {
		return nil, fmt.Errorf("print settings: %w", err)
	}
	defer rows.Close()

	values := make(map[string]string)
	for rows.Next() {
		var key string
		var value sql.NullString
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("scan print settings: %w", err)
		}
		if value.Valid {
			values[key] = value.String
		}
	}

	return values, rows.Err()
}

// SavePrintSettings обновляет реквизиты организации.
func (r *DocumentRepository) SavePrintSettings(ctx context.Context, values map[string]string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin print settings: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for key, value := range values {
		tag, err := tx.Exec(ctx,
			`UPDATE dir_default_settings_for_print SET field_value = $2 WHERE type_of_field = $1`,
			key, value,
		)
		if err != nil {
			return fmt.Errorf("update print setting %s: %w", key, err)
		}
		if tag.RowsAffected() > 0 {
			continue
		}

		if _, err := tx.Exec(ctx,
			`INSERT INTO dir_default_settings_for_print (type_of_field, field_value) VALUES ($1, $2)`,
			key, value,
		); err != nil {
			return fmt.Errorf("insert print setting %s: %w", key, err)
		}
	}

	return tx.Commit(ctx)
}

const ttnSelect = `
SELECT id,
       COALESCE(ttn_number, ''),
       COALESCE(dogovor, ''),
       COALESCE(putevka, ''),
       COALESCE(prilozhenie, ''),
       COALESCE(customer_organization, ''),
       COALESCE(shipper_organization, ''),
       COALESCE(consignee_organization, ''),
       COALESCE(payer_organization, ''),
       COALESCE(carrier_organization, ''),
       COALESCE(trailer_number, ''),
       COALESCE(type_of_trailer, ''),
       COALESCE(load_point, ''),
       COALESCE(upload_point, ''),
       COALESCE(CASE WHEN trim(distance) ~ '^[0-9]+([.,][0-9]+)?$'
                     THEN replace(trim(distance), ',', '.')::float8 END, 0)::float8,
       COALESCE(route_distance, 0)::float8,
       COALESCE(cargo_capacity, 0)::float8,
       COALESCE(cargo_release_name, ''),
       COALESCE(cargo_release_work_as, ''),
       COALESCE(accountant_name, ''),
       COALESCE(cargo_pick_up_name, ''),
       COALESCE(cargo_pick_up_work_as, ''),
       COALESCE(consignee_name, ''),
       COALESCE(consignee_work_as, ''),
       COALESCE(cargo_pick_up_additional_info, ''),
       COALESCE(price_cargo, 0)::float8,
       COALESCE(price_transportation, 0)::float8,
       COALESCE(forwarder, false),
       COALESCE(modified_date, created_date)
FROM ttn_document
WHERE id_weighing = $1
ORDER BY id DESC
LIMIT 1`

// TTN возвращает сохранённые реквизиты ТТН или nil, если документ ещё не заполняли.
func (r *DocumentRepository) TTN(ctx context.Context, sessionID int64) (*dto.TTNDocument, error) {
	row := r.pool.QueryRow(ctx, ttnSelect, sessionID)

	var doc dto.TTNDocument
	err := row.Scan(
		&doc.ID,
		&doc.TTNNumber,
		&doc.Dogovor,
		&doc.Putevka,
		&doc.Prilozhenie,
		&doc.CustomerOrganization,
		&doc.ShipperOrganization,
		&doc.ConsigneeOrganization,
		&doc.PayerOrganization,
		&doc.CarrierOrganization,
		&doc.TrailerNumber,
		&doc.TypeOfTrailer,
		&doc.LoadPoint,
		&doc.UploadPoint,
		&doc.Distance,
		&doc.RouteDistance,
		&doc.CargoCapacity,
		&doc.CargoReleaseName,
		&doc.CargoReleaseWorkAs,
		&doc.AccountantName,
		&doc.CargoPickUpName,
		&doc.CargoPickUpWorkAs,
		&doc.ConsigneeName,
		&doc.ConsigneeWorkAs,
		&doc.AdditionalInfo,
		&doc.PriceCargo,
		&doc.PriceTransportation,
		&doc.Forwarder,
		&doc.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ttn document: %w", err)
	}

	doc.SessionID = sessionID
	doc.Saved = true
	return &doc, nil
}

// SaveTTN создаёт или обновляет реквизиты ТТН.
func (r *DocumentRepository) SaveTTN(
	ctx context.Context,
	doc *dto.TTNDocument,
	userID int64,
) (*dto.TTNDocument, error) {
	existing, err := r.TTN(ctx, doc.SessionID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		var id int64
		if err := r.pool.QueryRow(ctx, `
			INSERT INTO ttn_document (
				id_weighing, ttn_number, dogovor, putevka, prilozhenie,
				customer_organization, shipper_organization, consignee_organization,
				payer_organization, carrier_organization, trailer_number, type_of_trailer,
				load_point, upload_point, distance, route_distance, cargo_capacity,
				cargo_release_name, cargo_release_work_as, accountant_name,
				cargo_pick_up_name, cargo_pick_up_work_as, consignee_name, consignee_work_as,
				cargo_pick_up_additional_info, price_cargo, price_transportation, forwarder,
				created_by_user_id, created_date
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
				$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, now()
			) RETURNING id`,
			doc.SessionID, doc.TTNNumber, doc.Dogovor, doc.Putevka, doc.Prilozhenie,
			doc.CustomerOrganization, doc.ShipperOrganization, doc.ConsigneeOrganization,
			doc.PayerOrganization, doc.CarrierOrganization, doc.TrailerNumber, doc.TypeOfTrailer,
			doc.LoadPoint, doc.UploadPoint, floatText(doc.Distance), doc.RouteDistance, doc.CargoCapacity,
			doc.CargoReleaseName, doc.CargoReleaseWorkAs, doc.AccountantName,
			doc.CargoPickUpName, doc.CargoPickUpWorkAs, doc.ConsigneeName, doc.ConsigneeWorkAs,
			doc.AdditionalInfo, doc.PriceCargo, doc.PriceTransportation, doc.Forwarder,
			userID,
		).Scan(&id); err != nil {
			return nil, fmt.Errorf("insert ttn document: %w", err)
		}
		return r.TTN(ctx, doc.SessionID)
	}

	if _, err := r.pool.Exec(ctx, `
		UPDATE ttn_document
		SET ttn_number = $2, dogovor = $3, putevka = $4, prilozhenie = $5,
		    customer_organization = $6, shipper_organization = $7,
		    consignee_organization = $8, payer_organization = $9,
		    carrier_organization = $10, trailer_number = $11, type_of_trailer = $12,
		    load_point = $13, upload_point = $14, distance = $15, route_distance = $16,
		    cargo_capacity = $17, cargo_release_name = $18, cargo_release_work_as = $19,
		    accountant_name = $20, cargo_pick_up_name = $21, cargo_pick_up_work_as = $22,
		    consignee_name = $23, consignee_work_as = $24, cargo_pick_up_additional_info = $25,
		    price_cargo = $26, price_transportation = $27, forwarder = $28,
		    modified_by_user_id = $29, modified_date = now()
		WHERE id = $1`,
		existing.ID,
		doc.TTNNumber, doc.Dogovor, doc.Putevka, doc.Prilozhenie,
		doc.CustomerOrganization, doc.ShipperOrganization, doc.ConsigneeOrganization,
		doc.PayerOrganization, doc.CarrierOrganization, doc.TrailerNumber, doc.TypeOfTrailer,
		doc.LoadPoint, doc.UploadPoint, floatText(doc.Distance), doc.RouteDistance, doc.CargoCapacity,
		doc.CargoReleaseName, doc.CargoReleaseWorkAs, doc.AccountantName,
		doc.CargoPickUpName, doc.CargoPickUpWorkAs, doc.ConsigneeName, doc.ConsigneeWorkAs,
		doc.AdditionalInfo, doc.PriceCargo, doc.PriceTransportation, doc.Forwarder,
		userID,
	); err != nil {
		return nil, fmt.Errorf("update ttn document: %w", err)
	}

	return r.TTN(ctx, doc.SessionID)
}

const torg12Select = `
SELECT id,
       COALESCE(name, ''),
       to_char(date, 'YYYY-MM-DD'),
       COALESCE(consignee_org_name, ''),
       COALESCE(provider_org_name, ''),
       COALESCE(struct_name, ''),
       COALESCE(foundation, ''),
       COALESCE(nds, 0)::float8,
       COALESCE(cargo_release_name, ''),
       COALESCE(cargo_release_work_as, ''),
       COALESCE(accountant_name, ''),
       COALESCE(shipper_name, ''),
       COALESCE(shipper_work_as, ''),
       COALESCE(cargo_pick_up_name, ''),
       COALESCE(cargo_pick_up_work_as, ''),
       COALESCE(consignee_name, ''),
       COALESCE(consignee_work_as, ''),
       COALESCE(proxy_name, ''),
       COALESCE(to_char(proxy_date, 'YYYY-MM-DD'), ''),
       COALESCE(proxy_given_by, ''),
       COALESCE(modified_date, created_date)
FROM invoice_torg12_documents
WHERE id_weighing = $1
ORDER BY id DESC
LIMIT 1`

// Torg12 возвращает сохранённые реквизиты ТОРГ-12 или nil.
func (r *DocumentRepository) Torg12(ctx context.Context, sessionID int64) (*dto.Torg12Document, error) {
	row := r.pool.QueryRow(ctx, torg12Select, sessionID)

	var doc dto.Torg12Document
	var date string
	err := row.Scan(
		&doc.ID,
		&doc.Name,
		&date,
		&doc.ConsigneeOrg,
		&doc.ProviderOrg,
		&doc.StructName,
		&doc.Foundation,
		&doc.NDS,
		&doc.CargoReleaseName,
		&doc.CargoReleaseWorkAs,
		&doc.AccountantName,
		&doc.ShipperName,
		&doc.ShipperWorkAs,
		&doc.CargoPickUpName,
		&doc.CargoPickUpWorkAs,
		&doc.ConsigneeName,
		&doc.ConsigneeWorkAs,
		&doc.ProxyName,
		&doc.ProxyDate,
		&doc.ProxyGivenBy,
		&doc.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("torg12 document: %w", err)
	}

	doc.Date = date
	doc.SessionID = sessionID
	doc.Saved = true
	return &doc, nil
}

// SaveTorg12 создаёт или обновляет реквизиты ТОРГ-12.
func (r *DocumentRepository) SaveTorg12(
	ctx context.Context,
	doc *dto.Torg12Document,
	date time.Time,
	proxyDate *time.Time,
	userID int64,
) (*dto.Torg12Document, error) {
	existing, err := r.Torg12(ctx, doc.SessionID)
	if err != nil {
		return nil, err
	}

	proxyDateArg := any(nil)
	if proxyDate != nil {
		proxyDateArg = *proxyDate
	}

	if existing == nil {
		var id int64
		if err := r.pool.QueryRow(ctx, `
			INSERT INTO invoice_torg12_documents (
				id_weighing, name, consignee_org_name, provider_org_name, date, nds,
				struct_name, foundation, cargo_release_name, cargo_release_work_as,
				accountant_name, shipper_name, shipper_work_as, cargo_pick_up_name,
				cargo_pick_up_work_as, consignee_name, consignee_work_as,
				proxy_name, proxy_date, proxy_given_by, created_by_user_id, created_date
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
				$15, $16, $17, $18, $19, $20, $21, now()
			) RETURNING id`,
			doc.SessionID, doc.Name, doc.ConsigneeOrg, doc.ProviderOrg, date, doc.NDS,
			doc.StructName, doc.Foundation, doc.CargoReleaseName, doc.CargoReleaseWorkAs,
			doc.AccountantName, doc.ShipperName, doc.ShipperWorkAs, doc.CargoPickUpName,
			doc.CargoPickUpWorkAs, doc.ConsigneeName, doc.ConsigneeWorkAs,
			doc.ProxyName, proxyDateArg, doc.ProxyGivenBy, userID,
		).Scan(&id); err != nil {
			return nil, fmt.Errorf("insert torg12 document: %w", err)
		}
		return r.Torg12(ctx, doc.SessionID)
	}

	if _, err := r.pool.Exec(ctx, `
		UPDATE invoice_torg12_documents
		SET name = $2, consignee_org_name = $3, provider_org_name = $4, date = $5,
		    nds = $6, struct_name = $7, foundation = $8, cargo_release_name = $9,
		    cargo_release_work_as = $10, accountant_name = $11, shipper_name = $12,
		    shipper_work_as = $13, cargo_pick_up_name = $14, cargo_pick_up_work_as = $15,
		    consignee_name = $16, consignee_work_as = $17, proxy_name = $18,
		    proxy_date = $19, proxy_given_by = $20,
		    modified_by_user_id = $21, modified_date = now()
		WHERE id = $1`,
		existing.ID,
		doc.Name, doc.ConsigneeOrg, doc.ProviderOrg, date, doc.NDS,
		doc.StructName, doc.Foundation, doc.CargoReleaseName, doc.CargoReleaseWorkAs,
		doc.AccountantName, doc.ShipperName, doc.ShipperWorkAs, doc.CargoPickUpName,
		doc.CargoPickUpWorkAs, doc.ConsigneeName, doc.ConsigneeWorkAs,
		doc.ProxyName, proxyDateArg, doc.ProxyGivenBy, userID,
	); err != nil {
		return nil, fmt.Errorf("update torg12 document: %w", err)
	}

	return r.Torg12(ctx, doc.SessionID)
}

func ptrFloat(value float64) *float64 {
	return &value
}

func floatText(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
