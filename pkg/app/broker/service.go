package broker

import (
	"photofinish/pkg/domain/broker"
	"photofinish/pkg/domain/errors"
)

type ServiceImpl struct {
	repo broker.Repo
}

func NewService(brokerRepo broker.Repo) *ServiceImpl {
	s := new(ServiceImpl)
	s.repo = brokerRepo
	return s
}

func (s *ServiceImpl) CheckExist(outboxId string) bool {
	err := s.repo.CheckExist(outboxId)
	return err != nil
}

func (s *ServiceImpl) FindOutboxList() (*[]broker.Outbox, error) {
	return s.repo.FindNotCompletedOutboxList(10)
}

func (s *ServiceImpl) Delete(outboxId string) error {
	err := s.repo.CheckExist(outboxId)
	if err != nil {
		return errors.ErrNotExists
	}
	err = s.repo.UpdateStatus(outboxId, broker.OutboxDone)
	return err
}
