package tests

import (
	"bytes"
	"encoding/json"
	"gymshark-interview/internal/server"
	"net/http"
	"testing"
)

func TestCalculatePackageForInexistentProduct(t *testing.T) {
	resp, err := http.Post(hostname+"/v1/products/some-prod-id/calculate/12345", "", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected status Not Found, got %d", resp.StatusCode)
	}
}

func TestCalculatePackageForProductWithPackageSizes(t *testing.T) {
	resp, err := http.Post(hostname+"/v1/products/0196b5d3-c52c-7e50-ac45-f83b35ee9e3d/calculate/12345", "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %d", resp.StatusCode)
	}
	var calculateResponse server.CalculatePackageSizeResponseBody
	err = json.NewDecoder(resp.Body).Decode(&calculateResponse)
	if err != nil {
		t.Fatalf("Failed decoding: %v", err)
	}

	want := []server.PackageResponseBody{
		{Amount: 1, Size: 500},
		{Amount: 1, Size: 2000},
		{Amount: 2, Size: 5000},
	}

	if len(want) != len(calculateResponse.Packages) {
		t.Fatalf("Unexpected response: want %v got %v", want, calculateResponse.Packages)
	}

	for i := range want {
		equals := false
		for j := range calculateResponse.Packages {
			if want[i].Amount == calculateResponse.Packages[j].Amount &&
				want[i].Size == calculateResponse.Packages[j].Size {
				equals = true
				break
			}
		}
		if !equals {
			t.Fatalf("Unexpected response: want %v got %v", want, calculateResponse.Packages)
		}
	}
}
