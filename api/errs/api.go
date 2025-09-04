package errs

// ErrorResponse модель возвращаемой ошибки
// @Description Модель описывает возвращаемую ошибку: код и краткое сообщение
type ErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

const (
	BadRequestCode          = 400
	NotFoundCode            = 404
	InternalServerErrorCode = 500
)

var (
	ErrInternalServer = ErrorResponse{
		Code: InternalServerErrorCode,
		Msg:  "внутренняя ошибка сервера при исполнении запроса",
	}

	ErrInvalidJSON = ErrorResponse{
		Code: BadRequestCode,
		Msg:  "Неверный формат данных",
	}

	ErrValidateJSON = ErrorResponse{
		Code: BadRequestCode,
		Msg:  "Неверно указаны данные",
	}

	ErrOrderNotFound = ErrorResponse{
		Code: NotFoundCode,
		Msg:  "заказ не найден",
	}

	ErrInvalidUUID = ErrorResponse{
		Code: BadRequestCode,
		Msg:  "Неверный формат uuid",
	}

	ErrOrderExistsUUID = ErrorResponse{
		Code: BadRequestCode,
		Msg:  "заказ с таким uuid уже существует",
	}

	ErrOrderExistsTrack = ErrorResponse{
		Code: BadRequestCode,
		Msg:  "заказ с таким трек номером уже существует",
	}

	ErrPaymentExists = ErrorResponse{
		Code: BadRequestCode,
		Msg:  "платеж с такой транзакцией уже существует",
	}
)
