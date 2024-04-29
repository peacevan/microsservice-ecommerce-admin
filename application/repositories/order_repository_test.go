package repositories_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestOrderRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	order := domain.NewOrder()
	order.ID = uuid.NewV4().String()
	order.FilePath = "path"
	order.CreatedAt = time.Now()

	repo := repositories.OrderRepositoryDb{Db: db}
	repo.Insert(order)

	v, err := repo.Find(order.ID)

	require.NotEmpty(t, v.ID)
	require.Nil(t, err)
	require.Equal(t, v.ID, order.ID)
}
