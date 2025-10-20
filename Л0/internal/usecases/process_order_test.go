// тест  ProcessOrderUseCase (валидный JSON)
package usecases

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockrepos "github.com/tvbondar/go-server/internal/repositories/mocks"
)

func TestProcessOrderUseCase_SaveSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockrepos.NewMockOrderRepository(ctrl)
	mockCache := mockrepos.NewMockOrderRepository(ctrl)

	// Подготовим корректный JSON (минимально)
	raw := []byte(`{
		"order_uid":"o1",
		"track_number":"t1",
		"delivery":{"name":"n","phone":"p"},
		"payment":{"transaction":"tx","currency":"USD","provider":"p","amount":1},
		"items":[{"chrt_id":1,"track_number":"t1","price":1,"name":"n","total_price":1,"nm_id":2,"brand":"b"}],
		"locale":"en",
		"customer_id":"c",
		"delivery_service":"ds",
		"date_created":"2021-11-26T06:22:19Z"
	}`)

	// ожидания: DB.SaveOrder вызовется, затем Cache.SaveOrder вызовется
	mockDB.EXPECT().SaveOrder(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockCache.EXPECT().SaveOrder(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	u := NewProcessOrderUseCase(mockDB, mockCache)
	err := u.Execute(context.Background(), raw)
	require.NoError(t, err)
}
