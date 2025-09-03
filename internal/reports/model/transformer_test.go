package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func getResp() ReportResponse {
	s := `{
    "name": "Umowa Aktywacyjna dla niani",
    "gross_income": 3600,
    "net_income": 3325.1,
    "sickness_contribution": true,
    "base_for_calculation_of_contributions_paid_by_country": 2333,
    "base_for_calculation_of_contributions_paid_by_employer": 1267,
    "contribution_paid_by_employer": [
        {
            "name": "składka emerytalna",
            "value": 123.66,
            "base": 1267,
            "percentage": 9.76,
            "PaidBy": ""
        },
        {
            "name": "składka rentowa",
            "value": 82.36,
            "base": 1267,
            "percentage": 6.5,
            "PaidBy": ""
        },
        {
            "name": "składka wypadkowa",
            "value": 21.16,
            "base": 1267,
            "percentage": 1.67,
            "PaidBy": ""
        }
    ],
    "contribution_paid_by_employee": [
        {
            "name": "składka emerytalna",
            "value": 123.66,
            "base": 1267,
            "percentage": 9.76,
            "PaidBy": ""
        },
        {
            "name": "składka rentowa",
            "value": 19.01,
            "base": 1267,
            "percentage": 1.5,
            "PaidBy": ""
        },
        {
            "name": "składka zdrowotna",
            "value": 101.19,
            "base": 1124.33,
            "percentage": 9,
            "PaidBy": ""
        },
        {
            "name": "składka chorobowa",
            "value": 31.04,
            "base": 1267,
            "percentage": 2.45,
            "PaidBy": ""
        }
    ],
    "contribution_paid_by_country": [
        {
            "name": "składka emerytalna",
            "value": 455.4,
            "base": 2333,
            "percentage": 19.52,
            "PaidBy": ""
        },
        {
            "name": "składka rentowa",
            "value": 186.64,
            "base": 2333,
            "percentage": 8,
            "PaidBy": ""
        },
        {
            "name": "składka wypadkowa",
            "value": 38.96,
            "base": 2333,
            "percentage": 1.67,
            "PaidBy": ""
        },
        {
            "name": "składka zdrowotna",
            "value": 209.97,
            "base": 2333,
            "percentage": 9,
            "PaidBy": ""
        }
    ],
    "employer_total_cost": 3827.18
}`
	resp := ReportResponse{}
	_ = json.Unmarshal([]byte(s), &resp)
	return resp
}

func BenchmarkJSON_ToByte(b *testing.B) {
	j := JSON{}
	resp := getResp()
	for b.Loop() {
		_, _ = j.ToByte(resp)
	}

}

func BenchmarkXML_ToByte(b *testing.B) {
	x := XML{}
	resp := getResp()
	for b.Loop() {
		_, _ = x.ToByte(resp)
	}

}

func BenchmarkGob_ToByte(b *testing.B) {
	x := Gob{}
	resp := getResp()
	for b.Loop() {
		_, _ = x.ToByte(resp)
	}

}

func BenchmarkProtoBuff_ToByte(b *testing.B) {
	x := ProtoBuff{}
	resp := getResp()
	for b.Loop() {
		_, _ = x.ToByte(resp)
	}

}

func Test_tmp(t *testing.T) {
	j := JSON{}
	x := XML{}
	g := Gob{}
	pb := ProtoBuff{}
	arr := []Transformer{j, x, g, pb}
	resp := getResp()
	for i, transformer := range arr {
		b, err := transformer.ToByte(resp)
		fmt.Println(i, ":", len(b), err)
	}

}
