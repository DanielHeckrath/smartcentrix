package main

import (
	"testing"
	"time"

	"github.com/101loops/clock"
)

func TestGenerateToken(t *testing.T) {
	// mock clock because generate token uses clock.Now()
	work := clock.NewMock()
	clock.Work = work

	tests := []struct {
		testName string

		userID  string
		tokenID string
		time    int64

		token string
	}{
		{"should generate valid token", "4c2eed51-2c8d-4d3b-b8cf-c6ad694f8868", "e92f66b2-db84-44d5-9e8c-776492df652c", 1486253039, "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI0YzJlZWQ1MS0yYzhkLTRkM2ItYjhjZi1jNmFkNjk0Zjg4NjgiLCJleHAiOjE0ODg4NDUwMzksImp0aSI6ImU5MmY2NmIyLWRiODQtNDRkNS05ZThjLTc3NjQ5MmRmNjUyYyIsImlhdCI6MTQ4NjI1MzAzOSwibmJmIjoxNDg2MjUzMDM5fQ.AZwJy0XaEJI4Ka8GF3ttTmFf6ARUWtX78QskYDAf7sLydKwVZKAVrdUlXM7Q5888KarmWqvZJ6Gb1NQlPxwQ3E7TEnmdKuDZkGBntLGnxC62WcamnHkqHTCjgp-dRFUEA_w0jlwtjIAvzjJ9mS3fmouLvYiClm7_-CswWv6JkLnMYjkKrVYzF2g3mZ9Tv9c2-nTqFmEnDpHn-Ubd_KYx0LJuELt_mImtBcy2TpplhmKcY-ZcVGmMyAd949OviWNgObn5msiHwB1zcpdWdebKXHqkboI4Mb0aOV_PagtKIMqRFwMkqnYjYXqSIlWDToucJwLC8ZgeTgle-UfoYd4oYA"},
		{"token should be different on time change", "4c2eed51-2c8d-4d3b-b8cf-c6ad694f8868", "e92f66b2-db84-44d5-9e8c-776492df652c", 1486255069, "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI0YzJlZWQ1MS0yYzhkLTRkM2ItYjhjZi1jNmFkNjk0Zjg4NjgiLCJleHAiOjE0ODg4NDcwNjksImp0aSI6ImU5MmY2NmIyLWRiODQtNDRkNS05ZThjLTc3NjQ5MmRmNjUyYyIsImlhdCI6MTQ4NjI1NTA2OSwibmJmIjoxNDg2MjU1MDY5fQ.e3z0xyKV_N0Ma-Z17VumEw3A9TWgh5G6-PY5iKT15CSpZ6460yRNTJCOU29MxIXHABLCUvUkG63_Z46a4eAuadaPp_4eN0yG1hiWVT0f4CwbQXhYyUWV13dJOm9jCndojw_0RtY_ra3Mrf62srPiMx7wiawJYY0I3hbt9A7X-m-HUe6JSqGXxPvcVMyvU6IlhStHoXLaAvMGYBEER490KFbiGSxhhyQoGnLfrpHnw7crlS_vo5V5a8J5n4WYwvuIrMGA7yj-fcv05UB-M4MO8hkCrcIfxos3MpseaeFSLNvTuyeO3wr4s_Yg_fAG80ts6Kkk5RsDAfs3pi_nNhvlzw"},
	}

	for _, data := range tests {
		t.Run(data.testName, func(tt *testing.T) {
			work.Set(time.Unix(data.time, 0))
			token, err := generateTokenWithID(data.tokenID, data.userID)

			if err != nil {
				tt.Errorf("generateTokenWithID(%[4]s, %[1]s) => (%[2]s, nil) expected, but got (%[3]s, %[4]v)", data.userID, data.token, token, err, data.tokenID)
			}

			if token != data.token {
				tt.Errorf("generateTokenWithID(%[3]s, %[1]s) => (%[2]s, nil) expected, but got (%[3]s, nil)", data.userID, data.token, token, data.tokenID)
			}
		})

	}
}

func TestValidateToken(t *testing.T) {
	tests := []struct {
		testName string

		userID  string
		tokenID string
		time    int64

		token string
	}{
		{"should validate valid token", "4c2eed51-2c8d-4d3b-b8cf-c6ad694f8868", "e92f66b2-db84-44d5-9e8c-776492df652c", 1486253039, "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI0YzJlZWQ1MS0yYzhkLTRkM2ItYjhjZi1jNmFkNjk0Zjg4NjgiLCJleHAiOjE0ODg4NDUwMzksImp0aSI6ImU5MmY2NmIyLWRiODQtNDRkNS05ZThjLTc3NjQ5MmRmNjUyYyIsImlhdCI6MTQ4NjI1MzAzOSwibmJmIjoxNDg2MjUzMDM5fQ.AZwJy0XaEJI4Ka8GF3ttTmFf6ARUWtX78QskYDAf7sLydKwVZKAVrdUlXM7Q5888KarmWqvZJ6Gb1NQlPxwQ3E7TEnmdKuDZkGBntLGnxC62WcamnHkqHTCjgp-dRFUEA_w0jlwtjIAvzjJ9mS3fmouLvYiClm7_-CswWv6JkLnMYjkKrVYzF2g3mZ9Tv9c2-nTqFmEnDpHn-Ubd_KYx0LJuELt_mImtBcy2TpplhmKcY-ZcVGmMyAd949OviWNgObn5msiHwB1zcpdWdebKXHqkboI4Mb0aOV_PagtKIMqRFwMkqnYjYXqSIlWDToucJwLC8ZgeTgle-UfoYd4oYA"},
		{"should validate token with different claims", "4c2eed51-2c8d-4d3b-b8cf-c6ad694f8868", "e92f66b2-db84-44d5-9e8c-776492df652c", 1486255069, "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI0YzJlZWQ1MS0yYzhkLTRkM2ItYjhjZi1jNmFkNjk0Zjg4NjgiLCJleHAiOjE0ODg4NDcwNjksImp0aSI6ImU5MmY2NmIyLWRiODQtNDRkNS05ZThjLTc3NjQ5MmRmNjUyYyIsImlhdCI6MTQ4NjI1NTA2OSwibmJmIjoxNDg2MjU1MDY5fQ.e3z0xyKV_N0Ma-Z17VumEw3A9TWgh5G6-PY5iKT15CSpZ6460yRNTJCOU29MxIXHABLCUvUkG63_Z46a4eAuadaPp_4eN0yG1hiWVT0f4CwbQXhYyUWV13dJOm9jCndojw_0RtY_ra3Mrf62srPiMx7wiawJYY0I3hbt9A7X-m-HUe6JSqGXxPvcVMyvU6IlhStHoXLaAvMGYBEER490KFbiGSxhhyQoGnLfrpHnw7crlS_vo5V5a8J5n4WYwvuIrMGA7yj-fcv05UB-M4MO8hkCrcIfxos3MpseaeFSLNvTuyeO3wr4s_Yg_fAG80ts6Kkk5RsDAfs3pi_nNhvlzw"},
	}

	for _, data := range tests {
		t.Run(data.testName, func(tt *testing.T) {
			claims, err := validateToken(data.token)

			if err != nil {
				tt.Errorf("validateToken(%[1]s) => should not return a validation error, but got %[2]v", data.token, err)
			}

			if claims.UserID != data.userID {
				tt.Errorf("expected claims.UserID to be %[1]s, but got %[2]s", data.userID, claims.UserID)
			}

			if claims.Id != data.tokenID {
				tt.Errorf("expected claims.Id to be %[1]s, but got %[2]s", data.tokenID, claims.Id)
			}

			if claims.IssuedAt != data.time {
				tt.Errorf("expected claims.IssuedAt to be %[1]d, but got %[2]d", data.time, claims.IssuedAt)
			}
		})

	}
}

func TestInvalidToken(t *testing.T) {
	tests := []struct {
		testName string
		token    string
	}{
		{"should detect token with invalid signature", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI0YzJlZWQ1MS0yYzhkLTRkM2ItYjhjZi1jNmFkNjk0Zjg4NjgiLCJleHAiOjE0ODg4NDUwMzksImp0aSI6ImU5MmY2NmIyLWRiODQtNDRkNS05ZThjLTc3NjQ5MmRmNjUyYyIsImlhdCI6MTQ4NjI1MzAzOSwibmJmIjoxNDg2MjUzMDM5fQ.AZwJy0XaEJI4Ka8GF3ttTmFf6ARUWtX78QskYDAf7sLydKwVZKAVrdUlXM7Q5888KarmWqvZJ6Gb1NQlPxwQ3E7TEnmdKuDZkGBntLGnxC62WcamnHkqHTCjgp-dRFUEA_w0jlwtjIAvzjJ9mS3fmouLvYiClm7_-CswWv6JkLnMYjkKrVYzF2gasdfnTqFmEnDpHn-Ubd_KYx0LJuELt_mImtBcy2TpplhmKcY-ZcVGmMyAd949OviWNgObn5msiHwB1zcpdWdebKXHqkboI4Mb0aOV_PagtKIMqRFwMkqnYjYXqSIlWDToucJwLC8ZgeTgle-UfoYd4oYA"},
		{"should detect token with invalid claims", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI0YzJlZWQ1MS0yYzhkLTRkM2ItYjhjZi1jNmFkNjk0Zjg4dNjgiLCJleHAiOjE0ODg4NDcwNjksImp0aSI6ImU5MmY2NmIyLWRiODQtNDRkNS05ZThjLTc3NjQ5MmRmNjUyYyIsImlhdCIasdNTA2OSwibmJmIjoxNDg2MjU1MDY5fQ===.e3z0xyKV_N0Ma-Z17VumEw3A9TWgh5G6-PY5iKT15CSpZ6460yRNTJCOU29MxIXHABLCUvUkG63_Z46a4eAuadaPp_4eN0yG1hiWVT0f4CwbQXhYyUWV13dJOm9jCndojw_0RtY_ra3Mrf62srPiMx7wiawJYY0I3hbt9A7X-m-HUe6JSqGXxPvcVMyvU6IlhStHoXLaAvMGYBEER490KFbiGSxhhyQoGnLfrpHnw7crlS_vo5V5a8J5n4WYwvuIrMGA7yj-fcv05UB-M4MO8hkCrcIfxos3MpseaeFSLNvTuyeO3wr4s_Yg_fAG80ts6Kkk5RsDAfs3pi_nNhvlzw"},
	}

	for _, data := range tests {
		t.Run(data.testName, func(tt *testing.T) {
			_, err := validateToken(data.token)

			if err == nil {
				tt.Errorf("validateToken(%[1]s) => should return a validation error, but got nil", data.token)
			}
		})

	}
}
