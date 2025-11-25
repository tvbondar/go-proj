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

	// Подготовим корректный JSON (полный)
	raw := []byte(`{
		"order_uid":"o1",
		"track_number":"t1",
		"entry":"WBIL",
		"delivery":{
			"name":"n",
			"phone":"p",
			"zip":"00000",
			"city":"c",
			"address":"addr",
			"region":"r",
			"email":"test@example.com"
		},
		"payment":{
			"transaction":"tx",
			"request_id":"",
			"currency":"USD",
			"provider":"p",
			"amount":1,
			"payment_dt":0,
			"bank":"b",
			"delivery_cost":0,
			"goods_total":1,
			"custom_fee":0
		},
		"items":[
			{
				"chrt_id":1,
				"track_number":"t1",
				"price":1,
				"rid":"r1",
				"name":"n",
				"sale":0,
				"size":"S",
				"total_price":1,
				"nm_id":2,
				"brand":"b",
				"status":1
			}
		],
		"locale":"en",
		"internal_signature":"",
		"customer_id":"c",
		"delivery_service":"ds",
		"shardkey":"1",
		"sm_id":1,
		"date_created":"2021-11-26T06:22:19Z",
		"oof_shard":"1"
	}`)

	// ожидания: DB.SaveOrder вызовется, затем Cache.SaveOrder вызовется
	mockDB.EXPECT().SaveOrder(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockCache.EXPECT().SaveOrder(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	u := NewProcessOrderUseCase(mockDB, mockCache)
	err := u.Execute(context.Background(), raw)
	require.NoError(t, err)
}
