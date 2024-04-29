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

func TestUserRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	user := &domain.User{}
	user, err := domain.NewUser(user)
	user.ID = uuid.NewV4().String()
	user.Name = "Ivan Amado Cardoso"
	user.CreatedAt = time.Now()

	repoUser := repositories.UserRepositoryDb{Db: db}
	repoUser.Insert(user)

	userInserted, err := repoUser.Find(user.ID)
	require.NotEmpty(t, userInserted.ID)
	require.Nil(t, err)
	require.Equal(t, userInserted.ID, user.ID)
	require.Equal(t, userInserted.Name, user.Name)

}

func TestUserRepositoryDbUpdate(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	user := &domain.User{}
	user, err := domain.NewUser(user)
	user.ID = uuid.NewV4().String()
	user.Name = "Ivan Amado"
	user.CreatedAt = time.Now()

	repoUser := repositories.UserRepositoryDb{Db: db}
	repoUser.Update(user)

	UserInserted, err := repoUser.Find(user.ID)
	require.NotEmpty(t, UserInserted.ID)
	require.Nil(t, err)
	require.Equal(t, UserInserted.Name, user.Name)
}
