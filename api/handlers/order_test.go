package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/orders_api/internal/models"
	"github.com/orders_api/internal/repository"
	"github.com/orders_api/internal/service"
	mock_service "github.com/orders_api/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_GetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockServiceOrder(ctrl)
	orderHandler := NewOrderHandler(mockService)

	app := fiber.New()
	app.Get("/orders/:order_uid", orderHandler.GetOrderByUID)

	tests := []struct {
		Name           string
		ID             string
		ExpectedStatus int
		ExpectedBody   string
		MockSetup      func(ms *mock_service.MockServiceOrder)
	}{
		{
			Name:           "Error_wrong_UUID",
			ID:             "wrong_uuid",
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody: `{
		"code": 400,
		"msg":  "Неверный формат uuid"
	}`,
			MockSetup: func(ms *mock_service.MockServiceOrder) {
				ms.EXPECT().GetOrderByUID("wrong_uuid").Return(nil, service.ErrInvalidUUID)
			},
		},
		{
			Name:           "Error_order_not_found",
			ID:             "f47ac10b-58cc-4372-a567-0e02b2c3d400",
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody: `{
					"code": 404,
					"msg":  "заказ не найден"
				}`,
			MockSetup: func(ms *mock_service.MockServiceOrder) {
				ms.EXPECT().GetOrderByUID("f47ac10b-58cc-4372-a567-0e02b2c3d400").Return(nil, repository.ErrOrderNotFoundByUUID)
			},
		},
		{
			Name:           "Success_found_order",
			ID:             "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			ExpectedStatus: http.StatusOK,
			ExpectedBody: `{
			   "order_uid": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			   "track_number": "WBILMTESTTRACK",
			   "entry": "WBIL",
			   "delivery": {
			      "name": "Test Testov",
			      "phone": "+9720000000",
			      "zip": "2639809",
			      "city": "Kiryat Mozkin",
			      "address": "Ploshad Mira 15",
			      "region": "Kraiot",
			      "email": "test@gmail.com"
			   },
			   "payment": {
			      "transaction": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			      "request_id": "",
			      "currency": "USD",
			      "provider": "wbpay",
			      "amount": 1817,
			      "payment_dt": 1637907727,
			      "bank": "alpha",
			      "delivery_cost": 1500,
			      "goods_total": 317,
			      "custom_fee": 0
			   },
			   "items": [
			      {
			         "chrt_id": 9934930,
			         "track_number": "WBILMTESTTRACK",
			         "price": 453,
			         "rid": "ab4219087a764ae0btest",
			         "name": "Mascaras",
			         "sale": 30,
			         "size": "0",
			         "total_price": 317,
			         "nm_id": 2389212,
			         "brand": "Vivienne Sabo",
			         "status": 202
			      }
			   ],
			   "locale": "en",
			   "internal_signature": "",
			   "customer_id": "test",
			   "delivery_service": "meest",
			   "shardkey": "9",
			   "sm_id": 99,
			   "date_created": "2021-11-26T06:22:19Z",
			   "oof_shard": "1"
			}`,

			MockSetup: func(ms *mock_service.MockServiceOrder) {
				// представим время в виде time.Time
				dateCreated, err := time.Parse(time.RFC3339, "2021-11-26T06:22:19Z")
				if err != nil {
					t.Fatal(err)
				}

				// представим uuid в виде строки:
				readyUUID, err := uuid.FromString("f47ac10b-58cc-4372-a567-0e02b2c3d479")
				if err != nil {
					t.Fatal(err)
				}
				ms.EXPECT().GetOrderByUID("f47ac10b-58cc-4372-a567-0e02b2c3d479").
					Return(
						&models.Order{
							OrderUID:    readyUUID,
							TrackNumber: "WBILMTESTTRACK",
							Entry:       "WBIL",
							Delivery: models.Delivery{
								Name:    "Test Testov",
								Phone:   "+9720000000",
								Zip:     "2639809",
								City:    "Kiryat Mozkin",
								Address: "Ploshad Mira 15",
								Region:  "Kraiot",
								Email:   "test@gmail.com",
							},
							Payment: models.Payment{
								Transaction:  readyUUID,
								RequestID:    "",
								Currency:     "USD",
								Provider:     "wbpay",
								Amount:       1817,
								PaymentDt:    1637907727,
								Bank:         "alpha",
								DeliveryCost: 1500,
								GoodsTotal:   317,
								CustomFee:    0,
							},
							Items: []models.Item{
								{
									ChrtID:      9934930,
									TrackNumber: "WBILMTESTTRACK",
									Price:       453,
									Rid:         "ab4219087a764ae0btest",
									Name:        "Mascaras",
									Sale:        30,
									Size:        "0",
									TotalPrice:  317,
									NmID:        2389212,
									Brand:       "Vivienne Sabo",
									Status:      202,
								},
							},
							Locale:            "en",
							InternalSignature: "",
							CustomerID:        "test",
							DeliveryService:   "meest",
							Shardkey:          "9",
							SmID:              99,
							DateCreated:       dateCreated,
							OofShard:          "1",
						}, nil,
					)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/orders/%s", tt.ID), nil)
			req.Header.Set("Content-Type", "application/json")

			tt.MockSetup(mockService)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, tt.ExpectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.JSONEq(t, tt.ExpectedBody, string(body))
		})
	}

}
