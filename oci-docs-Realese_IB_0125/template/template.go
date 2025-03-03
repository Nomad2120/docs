package template

import _ "embed"

var (
	//go:embed test.html
	testTemplate []byte

	//go:embed offerta.html
	contractTemplate []byte

	//go:embed act.html
	actTemplate []byte

	//go:embed notification_to_debtor.docx
	notificationToDebitorTemplate []byte

	//go:embed notary.docx
	notaryTemplate []byte

	//go:embed qrpage.html
	qrPageTemplate []byte

	//go:embed act_old.html
	actOldTemplate []byte

	//go:embed accounts_report.html
	accountsReportTemplate []byte

	//go:embed accounts_report_kz.html
	accountsReportKZTemplate []byte
)

func GetTestTemplate() []byte {
	return testTemplate
}

func GetContractTemplate() []byte {
	return contractTemplate
}

func GetActTemplate() []byte {
	return actTemplate
}

func GetNotificationToDebitorTemplate() []byte {
	return notificationToDebitorTemplate
}

func GetNotaryTemplate() []byte {
	return notaryTemplate
}

func GetQrPageTemplate() []byte {
	return qrPageTemplate
}

func GetActOldTemplate() []byte {
	return actOldTemplate
}

func GetAccountsReportTemplate() []byte {
	return accountsReportTemplate
}

func GetAccountsReportKZTemplate() []byte {
	return accountsReportKZTemplate
}