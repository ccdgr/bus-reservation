package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PayPalClient struct {
	ClientID string
	Secret   string
	BaseURL  string
	Client   *http.Client
}

func NewPayPalClient(clientID, secret string, isSandbox bool) *PayPalClient {
	baseURL := "https://api-m.paypal.com"
	if isSandbox {
		baseURL = "https://api-m.sandbox.paypal.com"
	}
	return &PayPalClient{
		ClientID: clientID,
		Secret:   secret,
		BaseURL:  baseURL,
		Client:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *PayPalClient) getAccessToken(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", p.BaseURL+"/v1/oauth2/token", bytes.NewBufferString("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(p.ClientID, p.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := p.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("failed to get token: %d", res.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}

type CreateOrderParams struct {
	ReturnURL string
	CancelURL string
	Amount    string
}

func (p *PayPalClient) CreateOrder(ctx context.Context, params CreateOrderParams) (string, error) {
	token, err := p.getAccessToken(ctx)
	if err != nil {
		return "", err
	}

	body := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"amount": map[string]interface{}{
					"currency_code": "CNY",
					"value":         params.Amount,
				},
			},
		},
		"application_context": map[string]interface{}{
			"return_url":  params.ReturnURL,
			"cancel_url":  params.CancelURL,
			"user_action": "PAY_NOW",
		},
	}

	b, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", p.BaseURL+"/v2/checkout/orders", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := p.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("create order failed with status: %d", res.StatusCode)
	}

	var result struct {
		ID    string `json:"id"`
		Links []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}

	for _, link := range result.Links {
		if link.Rel == "approve" {
			return link.Href, nil
		}
	}
	return "", fmt.Errorf("approve link not found")
}

func (p *PayPalClient) CaptureOrder(ctx context.Context, orderID string) error {
	token, err := p.getAccessToken(ctx)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v2/checkout/orders/%s/capture", p.BaseURL, orderID), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := p.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return fmt.Errorf("capture failed with status: %d", res.StatusCode)
	}
	return nil
}
