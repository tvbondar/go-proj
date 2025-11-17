// HTTP-эндпоинты
// Логи ошибок помогут мониторить API. Encode-ошибка теперь логируется — редкая, но важная.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/tvbondar/go-server/internal/usecases"
)

type HTTPHandler struct {
	usecase *usecases.GetOrderUseCase
}

func NewHTTPHandler(usecase *usecases.GetOrderUseCase) *HTTPHandler {
	return &HTTPHandler{usecase: usecase}
}

func (h *HTTPHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/order/")
	order, err := h.usecase.Execute(r.Context(), id)
	if err != nil {
		log.Printf("HTTP error for order %s: %v", id, err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("Encode error for order %s: %v", id, err)
	}
}
