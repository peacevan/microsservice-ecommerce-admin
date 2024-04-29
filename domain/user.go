package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type User struct {
	ID        string    `json:"user_id" valid:"uuid" gorm:"type:uuid;primary_key"`
	Name      string    `json:"name" valid:"notnull"`
	Email     string    `json:"email" valid:"notnull"` // Corrigido de "status" para "email"
	Error     string    `valid:"-"`
	CreatedAt time.Time `json:"created_at" valid:"-"`
	UpdatedAt time.Time `json:"updated_at" valid:"-"`
}

func NewUser(user *User) (*User, error) {
	// Chama o método prepare no objeto user passado como argumento
	user.prepare()

	// Valida o objeto user
	err := user.Validate()

	// Se houver um erro de validação, retorna nil e o erro
	if err != nil {
		return nil, err
	}

	// Se tudo estiver bem, retorna o objeto user
	return user, nil
}

func (user *User) prepare() {
	user.ID = uuid.NewV4().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

func (user *User) Validate() error {
	_, err := govalidator.ValidateStruct(user)

	if err != nil {
		return err
	}

	return nil
}
