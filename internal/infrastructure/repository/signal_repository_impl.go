// internal/infrastructure/repository/signal_repository_impl.go

package repository

import (
	"context"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

type signalRepositoryImpl struct {
	db *gorm.DB
}

func NewSignalRepository(db *gorm.DB) repository.SignalRepository {
	return &signalRepositoryImpl{db: db}
}

func (r *signalRepositoryImpl) Save(ctx context.Context, signal *model.Signal) error {
	result := r.db.WithContext(ctx).Create(signal)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to save signal")
	}
	return nil
}

func (r *signalRepositoryImpl) FindByID(ctx context.Context, id uint) (*model.Signal, error) {
	var signal model.Signal
	result := r.db.WithContext(ctx).First(&signal, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(result.Error, "failed to find signal by id")
	}
	return &signal, nil
}

func (r *signalRepositoryImpl) FindBySymbol(ctx context.Context, symbol string) ([]*model.Signal, error) {
	var signals []*model.Signal
	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).Find(&signals)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to find signals by symbol")
	}
	return signals, nil
}
