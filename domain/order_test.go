package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateIfOderIsEmpty(t *testing.T) {
	order := domain.NewOrder()
	err := order.Validate()
	require.Error(t, err)
}

func TestOrderIdIsNotAUuid(t *testing.T) {
	order := domain.NewOrder()
	order.ID = "abc"                                      // Isso deve ser inválido, pois não é um UUID
	order.UserID = "123e4567-e89b-12d3-a456-426614174000" // Exemplo de UUID válido
	err := order.Validate()
	require.Error(t, err) // Espera-se que haja um erro, pois o ID não é um UUID válido
}

func TestorderValidation(t *testing.T) {
	order := domain.NewOrder()

	order.ID = uuid.NewV4().String()                      // Gera um novo UUID para o ID
	order.UserID = "123e4567-e89b-12d3-a456-426614174000" // Exemplo de UUID válido
	order.Status = "pending"
	order.DeliveryAddress = "123 Main St"
	order.PaymentMethod = "Credit Card"
	order.TotalCost = 100.00
	order.CreatedAt = time.Now()
	order.ExpectedDeliveryDate = time.Now().Add(24 * time.Hour)
	order.SellerID = "seller123"
	err := order.Validate()
	require.Nil(t, err) // Espera-se que não haja erro, pois todos os campos são válidos
}
