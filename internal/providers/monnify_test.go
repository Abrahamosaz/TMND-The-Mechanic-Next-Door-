package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	// "github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// setupTestEnv sets up the test environment and returns a cleanup function
func setupTestEnv() func() {
	// Set environment variables
	os.Setenv("MONNIFY_BASE_URL", "https://api.monnify.com")
	os.Setenv("MONNIFY_API_KEY", "MK_PROD_A1TZFL9979")
	os.Setenv("MONNIFY_SECRET_KEY", "0AHHQC6D6UAQBF8W2X3SCAHY81019GC9")
	os.Setenv("MONNIFY_CONTRACT_CODE", "068243131842")

	// Return clesanup function
	cleanup := func() {
		os.Unsetenv("MONNIFY_BASE_URL")
		os.Unsetenv("MONNIFY_API_KEY")
		os.Unsetenv("MONNIFY_SECRET_KEY")
		os.Unsetenv("MONNIFY_CONTRACT_CODE")
	}

	return cleanup
}


func TestMonnify_GetAccessToken(t *testing.T) {
	// Setup test environment
    cleanup := setupTestEnv()
    defer cleanup()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/auth/login", r.URL.Path)

		// Verify Authorization header
		authHeader := r.Header.Get("Authorization")
		assert.Contains(t, authHeader, "Basic ")

		// Mock response
		mockResponse := map[string]interface{}{
			"requestSuccessful": true,
			"responseMessage":   "success",
			"responseCode":      "0",
			"responseBody": map[string]interface{}{
				"accessToken": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsibW9ubmlmeS12YWx1ZS1hZGRlZC1zZXJ2aWNlIiwibW9ubmlmeS1wYXltZW50LWVuZ2luZSIsIm1vbm5pZnktZGlzYnVyc2VtZW50LXNlcnZpY2UiLCJtb25uaWZ5LW9mZmxpbmUtcGF5bWVudC1zZXJ2aWNlIl0sInNjb3BlIjpbInByb2ZpbGUiXSwiZXhwIjoxNjU4MjIwNTk1LCJhdXRob3JpdGllcyI6WyJNUEVfTUFOQUdFX0xJTUlUX1BST0ZJTEUiLCJNUEVfVVBEQVRFX1JFU0VSVkVEX0FDQ09VTlQiLCJNUEVfSU5JVElBTElaRV9QQVlNRU5UIiwiTVBFX1JFU0VSVkVfQUNDT1VOVCIsIk1QRV9DQU5fUkVUUklFVkVfVFJBTlNBQ1RJT04iLCJNUEVfUkVUUklFVkVfUkVTRVJWRURfQUNDT1VOVCIsIk1QRV9ERUxFVEVfUkVTRVJWRURfQUNDT1VOVCIsIk1QRV9SRVRSSUVWRV9SRVNFUlZFRF9BQ0NPVU5UX1RSQU5TQUNUSU9OUyJdLCJqdGkiOiJhZmYwMTMyMi1kMmNmLTQzOGYtYWVkMC0wOGY0NTBhNmVhZWYiLCJjbGllbnRfaWQiOiJNS19URVNUX0pSUUFaUkZEMlcifQ.KOkMyFJIMKc6aZViDB3ekoiCPU2647eW_1RpySy9t_OXzfDSER2nn2QzXZsGPn8GlPJ4ZdhtuhevHQCLnEKQYIScgB9EQTPHTldiBG7uf9ta3NQK6Sxzq2pkM8heuO1v87tqtPLbBF7LlmI21liUPAPrxJVnw5PhCAOKbInIG-sj0BWbwMXJ2E8Cgz6yUeQ46-C0bmkyc0U6-FqENPSNl7oNXNY-nWaQfXd7tw2ybE4YoEzoRXcLTy0KLjlf3sxA5PizyT3nh9YqjVIRli-3sggnpvdjmIzw8784g_knNYmA8rj1Of7ROa_enDYmH6dWHbB1luchU9Y67FMhlBDWLg",
				"expiresIn": 3567,
			},
		}

		// Send mock response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	
	defer server.Close()
	
	// Create Monnify instance
	m := Monnify{
		Url: os.Getenv("MONNIFY_BASE_URL"),
		ApiKey: os.Getenv("MONNIFY_API_KEY"),
		SecretKey: os.Getenv("MONNIFY_SECRET_KEY"),
	}
	
	// Call the method
	result, err := m.GetAccessToken()


	fmt.Println("result value", result)
	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, true, result["requestSuccessful"])
	assert.Equal(t, "success", result["responseMessage"])
	assert.Equal(t, "0", result["responseCode"])

	responseBody, ok := result["responseBody"].(map[string]any)
	assert.True(t, ok)
	assert.NotEmpty(t, responseBody["accessToken"])
	assert.NotEmpty(t, responseBody["expiresIn"])
}


func TestMonnify_CreateInvoice(t *testing.T) {
    // Setup test environment
    cleanup := setupTestEnv()
    defer cleanup()

    // Create Monnify instance
    m := Monnify{
        Url: os.Getenv("MONNIFY_BASE_URL"),
        ApiKey: os.Getenv("MONNIFY_API_KEY"),
        SecretKey: os.Getenv("MONNIFY_SECRET_KEY"),
    }

	contractCode := os.Getenv("MONNIFY_CONTRACT_CODE")
    // Create test payload
    payload := &CreateInvoice{
        Amount: 1000,
        CurrencyCode: "NGN",
		Reference: fmt.Sprintf("%d", time.Now().Unix()),
        CustomerName: "Abraham Osazee",
        CustomerEmail: "abrahamosazee3@gmail.com",
        ContractCode: contractCode,
        Description: "test invoice",
		ExpiryDate: time.Now().Add(10 * time.Minute).Format("2006-01-02 15:04:05"),
        RedirectUrl: "http://localhost:3000",
		PaymentMethod: []string{"ACCOUNT_TRANSFER"},
    }

    // Call the method
    result, err := m.CreateInvoice(payload)

	fmt.Println("result value", result)
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, true, result["requestSuccessful"])
    assert.Equal(t, "success", result["responseMessage"])
    assert.Equal(t, "0", result["responseCode"])

    // Verify response body
    responseBody, ok := result["responseBody"].(map[string]any)
	assert.True(t, ok)
    assert.NotEmpty(t, responseBody["accountNumber"])
    assert.NotEmpty(t, responseBody["accountName"])
}




func TestMonnify_InitiateTransaction(t *testing.T) {
    // Setup test environment
    cleanup := setupTestEnv()
    defer cleanup()

    // Create Monnify instance
    m := Monnify{
        Url: os.Getenv("MONNIFY_BASE_URL"),
        ApiKey: os.Getenv("MONNIFY_API_KEY"),
        SecretKey: os.Getenv("MONNIFY_SECRET_KEY"),
    }

    // Create test payload
    payload := &CreateTransaction{
        Amount: 1000,
        CurrencyCode: "NGN",
		PaymentReference: fmt.Sprintf("%d", time.Now().Unix()),
        CustomerName: "Abraham Osazee",
        CustomerEmail: "abrahamosazee3@gmail.com",
        ContractCode: os.Getenv("MONNIFY_CONTRACT_CODE"),
        PaymentDescription: "test invoice",
        RedirectUrl: "http://localhost:3000",
    }

    // Call the method
    result, err := m.InitiateTransaction(payload)

	fmt.Println("result value", result)
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, true, result["requestSuccessful"])
    assert.Equal(t, "success", result["responseMessage"])
    assert.Equal(t, "0", result["responseCode"])

    // Verify response body
    responseBody, ok := result["responseBody"].(map[string]any)
	assert.True(t, ok)
    assert.NotEmpty(t, responseBody["checkoutUrl"])
}


// TestMain can be used for global setup/teardown
func TestMain(m *testing.M) {
	// Global setup
	// ...

	// Run tests
	code := m.Run()

	// Global teardown
	// ...

	os.Exit(code)
}