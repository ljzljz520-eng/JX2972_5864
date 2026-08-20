package service

import (
	"emergency-claim-code/internal/model"
	"errors"
	"fmt"
)

func (s *Service) IssueCode(recordID, at string) (model.ClaimCode, error) {
	record, err := s.GetRecord(recordID)
	if err != nil {
		return model.ClaimCode{}, err
	}
	if record.Status != "approved" {
		return model.ClaimCode{}, fmt.Errorf("code requires approval")
	}
	code := model.ClaimCode{Code: fmt.Sprintf("%s-C%04d", record.BatchID, s.sequence+1), RecordID: recordID, BatchID: record.BatchID, State: "issued", IssuedAt: at}
	if errs := model.ValidateClaimCode(code); len(errs) > 0 {
		return model.ClaimCode{}, errors.New(errs[0])
	}
	if err := s.repo.SaveClaimCode(code); err != nil {
		return model.ClaimCode{}, err
	}
	return code, nil
}

func (s *Service) GetCode(code string) (model.ClaimCode, error) { return s.repo.GetClaimCode(code) }

func (s *Service) RedeemCode(code, actor, at string) (model.ClaimCode, error) {
	item, err := s.GetCode(code)
	if err != nil {
		return model.ClaimCode{}, err
	}
	if item.State != "issued" {
		return model.ClaimCode{}, fmt.Errorf("code is not issued")
	}
	item.State = "redeemed"
	if err := s.repo.SaveClaimCode(item); err != nil {
		return model.ClaimCode{}, err
	}
	record, err := s.GetRecord(item.RecordID)
	if err == nil {
		_ = s.audit(record, "redeem", actor, code, at)
	}
	return item, err
}

func (s *Service) RevokeCode(code, actor, at string) (model.ClaimCode, error) {
	item, err := s.GetCode(code)
	if err != nil {
		return model.ClaimCode{}, err
	}
	if item.State != "issued" {
		return model.ClaimCode{}, fmt.Errorf("code cannot be revoked")
	}
	item.State = "revoked"
	if err := s.repo.SaveClaimCode(item); err != nil {
		return model.ClaimCode{}, err
	}
	record, err := s.GetRecord(item.RecordID)
	if err == nil {
		_ = s.audit(record, "revoke-code", actor, code, at)
	}
	return item, err
}
