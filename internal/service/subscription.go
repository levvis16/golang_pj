package service

import (
	"errors"
	"time"

	"subscription-service/internal/logger"
	"subscription-service/internal/models"
	"subscription-service/internal/repository"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo   *repository.SubscriptionRepository
	logger *logger.Logger
}

func NewSubscriptionService(repo *repository.SubscriptionRepository, logger *logger.Logger) *SubscriptionService {
	return &SubscriptionService{
		repo:   repo,
		logger: logger,
	}
}

func (s *SubscriptionService) Create(req *models.CreateSubscriptionRequest) (*models.Subscription, error) {
	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		s.logger.Error("failed to parse start date", "error", err)
		return nil, errors.New("invalid start date format, expected MM-YYYY")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		s.logger.Error("failed to parse user id", "error", err)
		return nil, errors.New("invalid user id format")
	}

	var endDate *time.Time
	if req.EndDate != "" {
		parsed, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			s.logger.Error("failed to parse end date", "error", err)
			return nil, errors.New("invalid end date format, expected MM-YYYY")
		}
		endDate = &parsed
	}

	now := time.Now()
	subscription := &models.Subscription{
		ID:          uuid.New(),
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(subscription); err != nil {
		s.logger.Error("failed to create subscription", "error", err)
		return nil, err
	}

	s.logger.Info("subscription created", "id", subscription.ID)
	return subscription, nil
}

func (s *SubscriptionService) GetByID(id string) (*models.Subscription, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		s.logger.Error("failed to parse id", "error", err)
		return nil, errors.New("invalid id format")
	}

	subscription, err := s.repo.GetByID(uuidID)
	if err != nil {
		s.logger.Error("failed to get subscription", "error", err)
		return nil, err
	}

	if subscription == nil {
		return nil, errors.New("subscription not found")
	}

	return subscription, nil
}

func (s *SubscriptionService) Update(id string, req *models.UpdateSubscriptionRequest) (*models.Subscription, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		s.logger.Error("failed to parse id", "error", err)
		return nil, errors.New("invalid id format")
	}

	existing, err := s.repo.GetByID(uuidID)
	if err != nil {
		s.logger.Error("failed to get subscription", "error", err)
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("subscription not found")
	}

	if req.ServiceName != "" {
		existing.ServiceName = req.ServiceName
	}
	if req.Price != nil {
		existing.Price = *req.Price
	}
	if req.StartDate != "" {
		startDate, err := time.Parse("01-2006", req.StartDate)
		if err != nil {
			s.logger.Error("failed to parse start date", "error", err)
			return nil, errors.New("invalid start date format")
		}
		existing.StartDate = startDate
	}
	if req.EndDate != nil {
		if *req.EndDate == "" {
			existing.EndDate = nil
		} else {
			endDate, err := time.Parse("01-2006", *req.EndDate)
			if err != nil {
				s.logger.Error("failed to parse end date", "error", err)
				return nil, errors.New("invalid end date format")
			}
			existing.EndDate = &endDate
		}
	}

	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		s.logger.Error("failed to update subscription", "error", err)
		return nil, err
	}

	s.logger.Info("subscription updated", "id", id)
	return existing, nil
}

func (s *SubscriptionService) Delete(id string) error {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		s.logger.Error("failed to parse id", "error", err)
		return errors.New("invalid id format")
	}

	if err := s.repo.Delete(uuidID); err != nil {
		s.logger.Error("failed to delete subscription", "error", err)
		return err
	}

	s.logger.Info("subscription deleted", "id", id)
	return nil
}

func (s *SubscriptionService) List(filter *models.SubscriptionFilter) ([]models.Subscription, error) {
	subscriptions, err := s.repo.List(filter)
	if err != nil {
		s.logger.Error("failed to list subscriptions", "error", err)
		return nil, err
	}

	return subscriptions, nil
}

func (s *SubscriptionService) GetTotalCost(filter *models.SubscriptionFilter) (int, error) {
	total, err := s.repo.GetTotalCost(filter)
	if err != nil {
		s.logger.Error("failed to get total cost", "error", err)
		return 0, err
	}

	return total, nil
}
