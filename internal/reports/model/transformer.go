package model

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"github.com/zozole1818/nanny-contract/internal/reports/grpc"
	"google.golang.org/protobuf/proto"
	"unsafe"
)

var rSize = int(unsafe.Sizeof(ReportResponse{}))

type Transformer interface {
	ToByte(response ReportResponse) ([]byte, error)
	FromByte(bytes []byte) (ReportResponse, error)
}

type JSON struct {
}

type XML struct {
}

type Gob struct {
}

type ProtoBuff struct {
}

func convertToMessage(list []Contribution) []*grpc.ContributionMessage {
	var result []*grpc.ContributionMessage
	for _, c := range list {
		result = append(result, &grpc.ContributionMessage{
			Name:       c.Name,
			Value:      c.Value,
			Base:       c.Base,
			Percentage: c.Percentage,
		})
	}
	return result
}

func fromToMessage(list []*grpc.ContributionMessage) []Contribution {
	var result []Contribution
	for _, c := range list {
		result = append(result, Contribution{
			Name:       c.Name,
			Value:      c.Value,
			Base:       c.Base,
			Percentage: c.Percentage,
		})
	}
	return result
}

func (p ProtoBuff) ToByte(resp ReportResponse) ([]byte, error) {

	msg := grpc.ReportMessage{
		Name:                 resp.Name,
		NumberOfWorkingHours: resp.NumberOfWorkingHours,
		GrossPerHour:         resp.GrossPerHour,
		GrossIncome:          resp.GrossIncome,
		NetIncome:            resp.NetIncome,
		SicknessContribution: resp.SicknessContribution,
		BaseForCalculationOfContributionsPaidByCountry:  resp.BaseForCalculationOfContributionsPaidByCountry,
		BaseForCalculationOfContributionsPaidByEmployer: resp.BaseForCalculationOfContributionsPaidByEmployer,
		EmployerTotalCost:          resp.EmployerTotalCost,
		ContributionPaidByCountry:  convertToMessage(resp.ContributionPaidByCountry),
		ContributionPaidByEmployee: convertToMessage(resp.ContributionPaidByEmployee),
		ContributionPaidByEmployer: convertToMessage(resp.ContributionPaidByEmployer),
	}
	return proto.Marshal(&msg)
}

func (p ProtoBuff) FromByte(bytes []byte) (ReportResponse, error) {
	msg := grpc.ReportMessage{}
	err := proto.Unmarshal(bytes, &msg)
	if err != nil {
		return ReportResponse{}, err
	}
	resp := ReportResponse{
		Name:                 msg.Name,
		NumberOfWorkingHours: msg.NumberOfWorkingHours,
		GrossPerHour:         msg.GrossPerHour,
		GrossIncome:          msg.GrossIncome,
		NetIncome:            msg.NetIncome,
		SicknessContribution: msg.SicknessContribution,
		BaseForCalculationOfContributionsPaidByCountry:  msg.BaseForCalculationOfContributionsPaidByCountry,
		BaseForCalculationOfContributionsPaidByEmployer: msg.BaseForCalculationOfContributionsPaidByEmployer,
		ContributionPaidByEmployer:                      fromToMessage(msg.ContributionPaidByEmployer),
		ContributionPaidByEmployee:                      fromToMessage(msg.ContributionPaidByEmployee),
		ContributionPaidByCountry:                       fromToMessage(msg.ContributionPaidByCountry),
		EmployerTotalCost:                               msg.EmployerTotalCost,
	}
	return resp, nil
}

func (g Gob) ToByte(response ReportResponse) ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, rSize))
	defer buffer.Reset()
	enc := gob.NewEncoder(buffer)
	err := enc.Encode(response)
	return buffer.Bytes(), err
}

func (g Gob) FromByte(b []byte) (ReportResponse, error) {
	buffer := bytes.NewBuffer(b)
	defer buffer.Reset()
	dec := gob.NewDecoder(buffer)
	resp := ReportResponse{}
	err := dec.Decode(&resp)
	return resp, err
}

func (x XML) ToByte(response ReportResponse) ([]byte, error) {
	return xml.Marshal(response)
}

func (x XML) FromByte(bytes []byte) (ReportResponse, error) {
	resp := ReportResponse{}
	err := xml.Unmarshal(bytes, &resp)
	return resp, err
}

func (t JSON) ToByte(response ReportResponse) ([]byte, error) {
	return json.Marshal(response)
}

func (t JSON) FromByte(bytes []byte) (ReportResponse, error) {
	resp := ReportResponse{}
	err := json.Unmarshal(bytes, &resp)
	return resp, err
}
