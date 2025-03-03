package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/common"
	config "gitlab.enterprise.qazafn.kz/oci/oci-docs/config"
	model "gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
	repository "gitlab.enterprise.qazafn.kz/oci/oci-docs/repository"
)

var coreRepo repository.Core

func init() {
	cfg := config.New("../")
	logger := common.InitLogger(cfg)
	httpClient := common.InitHTTPClient(cfg, logger)
	// converter, err := conv.NewWKDocConverter()
	// if err != nil {
	// 	logger.Fatalf("Error create doc converter: %v", err)
	// }
	coreRepo = repository.NewCoreHTTPRepo(cfg, logger, httpClient)
}

func TestSaveContract(t *testing.T) {
	doc := common.LoadFixture(t, "../testdata/oferta.base64")
	req := &model.SignContractRequest{
		ID:        30,
		DocBase64: string(doc),
		Extension: "pdf",
	}
	resp, err := coreRepo.SaveContract(context.Background(), req)
	require.NoError(t, err)
	fmt.Printf("%+v\n", resp)
}

func TestHtml2Pdf(t *testing.T) {
	doc := common.LoadFixture(t, "../testdata/html.base64")
	resp, err := coreRepo.Html2Pdf(context.Background(), string(doc))
	require.NoError(t, err)
	fmt.Printf("%+v\n", resp)
}
