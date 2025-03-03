package repository

import (
	"context"

	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
)

// Core -
type Core interface {
	SaveContract(context.Context, *model.SignContractRequest) (string, error)
	GetRegistration(context.Context, int) (*model.RegistrationInfo, error)
	GetServiceGroups(context.Context) ([]model.ServiceGroup, error)
	GetOSVCurrentMonth(context.Context, int) (*model.OSVResult, error)
	GetOSVByPeriod(context.Context, *model.OSVReportRequest) (*model.OSVResult, error)
	GetAbonentOSVAll(context.Context, int) ([]model.AbonentOSVResult, error)
	GetOSVForDebtors(context.Context, *model.OSVReportRequest) (*model.OSVResult, error)
	GetOSIInfo(context.Context, int) (*model.OSIResult, error)
	GetAllOSI(context.Context) ([]model.OSIResult, error)
	GetPayments(context.Context, *model.PaymentsReportRequest) ([]model.PaymentInfo, error)
	GetFixes(context.Context, *model.FixesRequest) ([]model.FixesInfo, error)
	GetAct(context.Context, int) (*model.GetActResult, error)
	GetSignedActs(context.Context, int) ([]model.SignedAct, error)
	GetActDocs(context.Context, int) ([]model.ActDoc, error)
	GetScan(context.Context, int) (string, error)
	SignAct(context.Context, *model.SignActRequest) error
	UnsignAct(context.Context, int) error
	GetPaymentOrders(context.Context, *model.PaymentOrdersRequest) ([]model.PaymentOrder, error)
	GetDebtorNotification(context.Context, *model.DebtorNotificationRequest) (*model.DebtorNotificationResult, error)
	GetNotaryApplication(context.Context, *model.NotaryApplicationRequest) (*model.NotaryApplicationResult, error)
	GetOSIAbonents(context.Context, int) ([]model.OSIAbonent, error)
	AddDoc(ctx context.Context, id int, req *model.AddDocRequest) (*model.AddDocResponse, error)
	GetAbonentInfo(context.Context, int) (*model.OSIAbonent, error)
	GetAccountReports(context.Context, int) (*model.AccountReportsResult, error)
}

// Ear
type Ear interface {
	GetBuildingInfo(ctx context.Context, atsID, id int) (*model.BuildingInfoResponse, error)
}
