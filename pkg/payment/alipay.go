package payment

import (
	"github.com/smartwalle/alipay/v3"
)

func NewAlipayClient(appID, privateKey, aliPublicKey string, isProduction bool) (*alipay.Client, error) {
	client, err := alipay.New(appID, privateKey, isProduction)
	if err != nil {
		return nil, err
	}

	err = client.LoadAliPayPublicKey(aliPublicKey)
	if err != nil {
		return nil, err
	}

	return client, nil
}
