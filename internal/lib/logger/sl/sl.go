package sl

import (
	"log/slog"
)

//Оставляем удобная штучка для вывода подробностей об ошибке

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
