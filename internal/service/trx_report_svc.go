package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bsnack/internal/model"
	"bsnack/internal/repository"
	customerror "bsnack/pkg/custom_error"

	"github.com/redis/go-redis/v9"
)

type reportService struct {
	reportRepo repository.ReportRepository
	redis      *redis.Client
}

type ReportService interface {
	GetTransactionReport(ctx context.Context, start string, end string) (*model.TransactionReportResponse, error)
}

func NewReportService(reportRepo repository.ReportRepository, redis *redis.Client) ReportService {
	return &reportService{
		reportRepo: reportRepo,
		redis:      redis,
	}
}

func (s *reportService) GetTransactionReport(ctx context.Context, start string, end string) (*model.TransactionReportResponse, error) {

	if start == "" || end == "" {
		return nil, customerror.ErrInvalidDateRange
	}

	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil, customerror.ErrInvalidDateRange
	}

	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		return nil, customerror.ErrInvalidDateRange
	}

	if endDate.Before(startDate) {
		return nil, customerror.ErrInvalidDateRange
	}

	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	cacheKey := fmt.Sprintf("report:transactions:%s:%s", start, end)

	if s.redis != nil {
		cached, err := s.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var resp model.TransactionReportResponse
			if err := json.Unmarshal([]byte(cached), &resp); err == nil {
				return &resp, nil
			}
		}
	}

	resp, err := s.reportRepo.GetTransactionReport(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	if s.redis != nil {
		if b, err := json.Marshal(resp); err == nil {
			_ = s.redis.Set(ctx, cacheKey, b, time.Minute).Err()
		}
	}

	return resp, nil
}
