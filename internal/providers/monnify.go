package providers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)


type Monnify struct {
	Url string
	ApiKey string
	SecretKey string
}


type CreateInvoice struct {
	Amount float64 `json:"amount"`
	CurrencyCode  string `json:"currencyCode"`
	Reference	string `json:"invoiceReference"`
	CustomerName string	`json:"customerName"`
	CustomerEmail string `json:"customerEmail"`
	ContractCode string  `json:"contractCode"`
	Description string	`json:"description"`
	ExpiryDate string	`json:"expiryDate"`
	RedirectUrl string  `json:"redirectUrl"`
	PaymentMethod []string `json:"paymentMethod"`
}


type CreateTransaction struct {
	Amount float64 `json:"amount"`
	CustomerName string `json:"customerName"`
	CustomerEmail string `json:"customerEmail"`
	PaymentReference string `json:"paymentReference"`
	PaymentDescription string `json:"paymentDescription"`
	CurrencyCode string `json:"currencyCode"`
	ContractCode string `json:"contractCode"`
	RedirectUrl string `json:"redirectUrl"`
	PaymentMethod string `json:"paymentMethod"`
}


func (m *Monnify) GetAccessToken() (map[string]any, error,) {

	credentials := m.ApiKey + ":" + m.SecretKey
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))

	url :=  m.Url + "/api/v1/auth/login"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
		if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Set Authorization header
	req.Header.Set("Authorization", "Basic "+encoded)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]any;
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}



func (m *Monnify) CreateInvoice(payload *CreateInvoice)  (map[string]any, error){

	accessTokenResult, err := m.GetAccessToken()

	if err != nil {
		return nil, err;
	}

   responseBody, ok := accessTokenResult["responseBody"].(map[string]any)
    if !ok {
        return nil, fmt.Errorf("invalid response format: expected responseBody to be a map")
    }

    accessToken, ok := responseBody["accessToken"].(string)
    if !ok {
        return nil, fmt.Errorf("invalid response format: expected accessToken to be a string")
    }

	url :=  m.Url + "/api/v1/invoice/create"

	// Marshal the payload to JSON
    data, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal payload: %w", err)
    }

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Set required headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

	// Parse the response
    var result map[string]any
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return result, nil
}



func (m  *Monnify) ConfirmInvoicePayment(invoiceRef string) error {
	accessTokenResult, err := m.GetAccessToken()

	if err != nil {
		return err;
	}

   responseBody, ok := accessTokenResult["responseBody"].(map[string]any)
    if !ok {
        return fmt.Errorf("invalid response format: expected responseBody to be a map")
    }

    accessToken, ok := responseBody["accessToken"].(string)
    if !ok {
        return fmt.Errorf("invalid response format: expected accessToken to be a string")
    }

	url := m.Url + fmt.Sprintf("/api/v1/invoice/%s/details", invoiceRef)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return err;
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err;
	}
	defer resp.Body.Close()

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	resultBody, ok := result["responseBody"].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid response format: expected responseBody to be a map")
	}

	invoiceStatus, ok := resultBody["invoiceStatus"].(string)
	if !ok {
		return fmt.Errorf("invalid response format: expected invoiceStatus to be a string")
	}
	
	if invoiceStatus != "PAID" {
		return fmt.Errorf("invoice is not paid")
	}

	return nil
}

func (m *Monnify) InitiateTransaction(payload *CreateTransaction)  (map[string]any, error) {
	accessTokenResult, err := m.GetAccessToken()

	if err != nil {
		return nil, err;
	}

   responseBody, ok := accessTokenResult["responseBody"].(map[string]any)
    if !ok {
        return nil, fmt.Errorf("invalid response format: expected responseBody to be a map")
    }

    accessToken, ok := responseBody["accessToken"].(string)
    if !ok {
        return nil, fmt.Errorf("invalid response format: expected accessToken to be a string")
    }

	url := m.Url + "/api/v1/merchant/transactions/init-transaction"

	// Marshal the payload to JSON
    data, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal payload: %w", err)
    }

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	
	// Set required headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+accessToken)


	// Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

	// Parse the response
    var result map[string]any
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return result, nil
}

