package repositories

import (
	"crm-backend/internal/models"

	"gorm.io/gorm"
)

// UserRepository define a interface para operações de usuário no banco de dados
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	EmailExists(email string) (bool, error)
}

// userRepository implementa UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository cria uma nova instância do repositório de usuários
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create cria um novo usuário no banco de dados
func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetByID busca um usuário pelo ID
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail busca um usuário pelo email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update atualiza um usuário existente
func (r *userRepository) Update(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

// Delete remove um usuário do banco de dados (soft delete)
func (r *userRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

// EmailExists verifica se um email já está em uso
func (r *userRepository) EmailExists(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

