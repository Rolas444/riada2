package repository

import (
	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
	"gorm.io/gorm"
)

type gormPersonRepository struct {
	db *gorm.DB
}

func NewGormPersonRepository(db *gorm.DB) ports.PersonRepository {
	return &gormPersonRepository{db: db}
}

func (r *gormPersonRepository) Save(person *domain.Person) error {
	// Save actualiza el registro si tiene una clave primaria, o crea uno nuevo si no la tiene.
	// println("Saving person:", person.ID, person.Name, person.MiddleName)
	return r.db.Save(person).Error
}

func (r *gormPersonRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Person{}, id).Error
}

func (r *gormPersonRepository) FindByID(id uint) (*domain.Person, error) {
	var person domain.Person
	if err := r.db.Preload("Addresses").Preload("Phones").First(&person, id).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *gormPersonRepository) FindByUserID(userID uint) (*domain.Person, error) {
	var person domain.Person
	if err := r.db.Preload("Addresses").Preload("Phones").Where("user_id = ?", userID).First(&person).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *gormPersonRepository) Search(searchTerm string) ([]domain.Person, error) {
	var persons []domain.Person
	query := r.db.Preload("Addresses").Preload("Phones")

	if searchTerm != "" {
		likeTerm := "%" + searchTerm + "%"
		// Busca en la concatenación de nombre, apellido paterno y materno, O en el número de documento.
		// Esta sintaxis es para PostgreSQL.
		query = query.Where("name || ' ' || middle_name || ' ' || last_name ILIKE ? OR doc_number = ?", likeTerm, searchTerm)
	}

	if err := query.Order("id DESC").Limit(300).Find(&persons).Error; err != nil {
		return nil, err
	}
	return persons, nil
}

func (r *gormPersonRepository) FindByDocument(docType domain.DocType, docNumber string) (*domain.Person, error) {
	var person domain.Person
	if err := r.db.Where("type_doc = ? AND doc_number = ?", docType, docNumber).First(&person).Error; err != nil {
		return nil, err
	}
	return &person, nil
}
