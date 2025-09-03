package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zozole1818/nanny-contract/internal/reports/model"
	"time"
)

type ReportCache struct {
	rdb         *redis.Client
	transformer model.Transformer
}

func NewReportCache(addr string, transformer model.Transformer) ReportCache {
	if addr == "" {
		addr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return ReportCache{
		rdb:         rdb,
		transformer: transformer,
	}
}

func (r ReportCache) GetReport(ctx context.Context, grossIncome float64, sicknessContribution bool) (model.ReportResponse, error) {
	key := constructKey(grossIncome, sicknessContribution)
	result, err := r.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return model.ReportResponse{}, fmt.Errorf("not found")
	}
	resp, err := r.transformer.FromByte(result)
	if err != nil {
		return model.ReportResponse{}, fmt.Errorf("from byte error: %v", err)
	}
	fmt.Println("CACHE: found key!", key)
	return resp, nil
}

func (r ReportCache) AddReport(ctx context.Context, report model.ReportResponse) error {
	key := constructKey(report.GrossIncome, report.SicknessContribution)
	b, err := r.transformer.ToByte(report)
	if err != nil {
		return fmt.Errorf("to byte error: %v", err)
	}
	err = r.rdb.Set(ctx, key, b, 30*time.Second).Err()
	if err != nil {
		return fmt.Errorf("cache error: %v", err)
	}
	fmt.Println("CACHE: added key!", key)
	return nil
}

func constructKey(grossIncome float64, sicknessContribution bool) string {
	return fmt.Sprintf("%f:%t", grossIncome, sicknessContribution)
}
