package repository

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/common"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/config"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
	"go.uber.org/zap"
)

type earHTTPRepo struct {
	cfg        *config.Config
	log        *zap.SugaredLogger
	httpClient *http.Client
}

// NewEarHTTPRepo -
func NewEarHTTPRepo(cfg *config.Config, log *zap.SugaredLogger, httpClient *http.Client) Ear {
	return &earHTTPRepo{cfg: cfg, log: log, httpClient: httpClient}
}

func (r *earHTTPRepo) GetBuildingInfo(ctx context.Context, atsID, id int) (*model.BuildingInfoResponse, error) {
	url := fmt.Sprintf("%s/api/Search/Building?id=%d&atsId=%d", r.cfg.EarURL, id, atsID)
	var resp model.BuildingInfoResponse
	status, body, err := common.RequestJSON(ctx, r.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetBuildingInfo. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetBuildingInfo. Invalid Response Status Code: %d", status)
	}

	return &resp, nil
}
