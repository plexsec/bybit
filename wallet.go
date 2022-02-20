package bybit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// WalletService :
type WalletService struct {
	Client *Client
}

// BalanceResponse :
type BalanceResponse struct {
	CommonResponse `json:",inline"`
	Result         BalanceResult `json:"result"`
}

// BalanceResult :
type BalanceResult struct {
	Balance map[Coin]Balance
}

// UnmarshalJSON :
func (r *BalanceResult) UnmarshalJSON(data []byte) error {
	parsedData := map[string]Balance{}
	if err := json.Unmarshal(data, &parsedData); err != nil {
		return err
	}
	r.Balance = map[Coin]Balance{}
	for coin, balanceData := range parsedData {
		r.Balance[Coin(coin)] = balanceData
	}
	return nil
}

// Balance :
type Balance struct {
	Equity           float64 `json:"equity"`
	AvailableBalance float64 `json:"available_balance"`
	UsedMargin       float64 `json:"used_margin"`
	OrderMargin      float64 `json:"order_margin"`
	PositionMargin   float64 `json:"position_margin"`
	OccClosingFee    float64 `json:"occ_closing_fee"`
	OccFundingFee    float64 `json:"occ_funding_fee"`
	WalletBalance    float64 `json:"wallet_balance"`
	RealisedPnl      float64 `json:"realised_pnl"`
	UnrealisedPnl    float64 `json:"unrealised_pnl"`
	CumRealisedPnl   float64 `json:"cum_realised_pnl"`
	GivenCash        float64 `json:"given_cash"`
	ServiceCash      float64 `json:"service_cash"`
}

// Balance :
func (s *WalletService) Balance(coin Coin) (*BalanceResponse, error) {
	var res BalanceResponse

	params := map[string]string{
		"coin": string(coin),
	}
	url, err := s.Client.BuildPrivateURL("/v2/private/wallet/balance", params)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

type TransferResponse struct {
	CommonResponse `json:",inline"`
	Result         TransferResult `json:"result"`
}

type TransferResult struct {
	TransferId string `json:"transfer_id"`
}

// Create Internal Transfer :
func (s *WalletService) InternalTransfer(coin Coin, amount float64, from, to AccountType) (*TransferResponse, error) {
	var res TransferResponse

	params := map[string]string{
		"coin":              string(coin),
		"amount":            fmt.Sprint(amount),
		"from_account_type": string(from),
		"to_account_type":   string(to),
	}

	url, err := s.Client.BuildPrivateURL("/asset/v1/private/transfer", params)
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("json marshal for InternalTransfer: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Create Subaccount Transfer :
func (s *WalletService) SubAccountTransfer(coin Coin, amount float64, subUserId string, typ TransferType) (*TransferResponse, error) {
	var res TransferResponse

	params := map[string]string{
		"coin":        string(coin),
		"amount":      fmt.Sprint(amount),
		"sub_user_id": subUserId,
		"type":        string(typ),
	}

	url, err := s.Client.BuildPrivateURL("/asset/v1/private/sub-member/transfe", params)
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("json marshal for SubAccountTransfer: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
