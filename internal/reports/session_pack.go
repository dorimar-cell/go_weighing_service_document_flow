package reports

import (
	"fmt"
	"strings"

	"go_weighing_service/internal/dto"
)

// SessionPack — данные для комплекта документов по одному рейсу.
type SessionPack struct {
	Details *dto.SessionDetails
	TTN     dto.TTNDocument
	Torg12  dto.Torg12Document
	Invoice dto.InvoiceDocument
	Photos  []InlinePhoto
}

// renderSessionPack собирает комплект: карточка ТС, талон, ТТН, ТОРГ-12, счёт-фактура.
func renderSessionPack(pack SessionPack, settings dto.ReportSettings) []byte {
	details := pack.Details
	session := &details.Session

	var body strings.Builder
	body.WriteString("<h1>Комплект документов по рейсу №" + esc(int64Text(session.ID)) + "</h1>")
	body.WriteString("<p class=\"subtitle\">" + esc(session.TransportNumber) + " · " +
		esc(dateOnly(session.FirstMeasuredAt)) + "</p>")
	body.WriteString("<div class=\"meta\">")
	if settings.OrganizationName != "" {
		body.WriteString("<span><b>" + esc(settings.OrganizationName) + "</b></span>")
	}
	body.WriteString("<span>Документов: 7</span>")
	body.WriteString("<span>Сформирован: " + nowInLocation().Format("02.01.2006 15:04") + "</span>")
	body.WriteString("</div>")

	sections := []string{
		tscardBody(details, settings, pack.Photos),
		talonVariantBody(details, settings, "base"),
		cargoActBody(details, settings),
		uploadingBody(details, settings),
		ttnBody(details, pack.TTN, settings),
		torg12Body(details, pack.Torg12, settings),
		invoiceBody(details, pack.Invoice, settings),
	}

	for index, section := range sections {
		body.WriteString(fmt.Sprintf("<section class=\"talon-page\" data-document=\"%d\">", index+1))
		body.WriteString(section)
		body.WriteString("</section>")
	}

	body.WriteString("<p class=\"note\">Комплект сформирован веб-системой весовой. " +
		"Реквизиты ТТН, ТОРГ-12 и счёта берутся сохранённые на рейсе; для незаполненных " +
		"документов подставляются значения по умолчанию.</p>")

	return wrapHTML("Комплект документов №"+int64Text(session.ID), false, body.String())
}
