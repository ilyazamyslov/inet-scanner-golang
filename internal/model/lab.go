package model

type LabResp struct {
	Q1 []Host `json:"q1"`
	Q2 []Host `json:"q2"`
	Q3 []Host `json:"q3"`
	Q4 int    `json:"q4"`
}
