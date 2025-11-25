// тест GetOrderUseCase (cache hit)
package usecases

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/tvbondar/go-server/internal/entities"
	mockrepos "github.com/tvbondar/go-server/internal/repositories/mocks"
)

func TestGetOrderUseCase_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mockrepos.NewMockOrderRepository(ctrl)
	order := entities.Order{OrderUID: "id1"}

	// cacheRepo.GetOrderByID returns order => usecase should return it and not call DB
	mockCache.EXPECT().GetOrderByID(gomock.Any(), "id1").Return(order, nil).Times(1)

	u := NewGetOrderUseCase(mockCache, nil) // cacheRepo, dbRepo
	got, err := u.Execute(context.Background(), "id1")
	require.NoError(t, err)
	require.Equal(t, "id1", got.OrderUID)
}
