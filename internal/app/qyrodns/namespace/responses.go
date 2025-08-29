package namespace

import (
	"github.com/qyrocloud/qyrodns/internal/pkg/apikey"
)

type ApiKeyAccessResponse struct {
	Access *ApiKeyAccess  `json:"access"`
	ApiKey *apikey.ApiKey `json:"api_key"`
}
