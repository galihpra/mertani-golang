package repository

import (
	"context"
	"errors"
	"io"
	"mertani-golang/features/sensors"
	"mertani-golang/utils/cloudinary"
	"strings"

	"gorm.io/gorm"
)

type Sensor struct {
	Id          uint      `gorm:"column:id; primaryKey;"`
	Name        string    `gorm:"column:name; type:varchar(200);unique;"`
	Description string    `gorm:"column:description; type:text;"`
	ImageUrl    string    `gorm:"column:image; type:text;"`
	ImageRaw    io.Reader `gorm:"-"`

	UserId uint `gorm:"column:user_id;"`
	User   User `gorm:"foreignKey:UserId"`
}

type User struct {
	Id uint
}

type sensorRepository struct {
	db    *gorm.DB
	cloud cloudinary.Cloud
}

func NewSensorRepository(db *gorm.DB, cloud cloudinary.Cloud) sensors.Repository {
	return &sensorRepository{
		db:    db,
		cloud: cloud,
	}
}

func (repo *sensorRepository) Create(data sensors.Sensor) error {
	var inputDB = new(Sensor)
	if data.ImageRaw != nil {
		url, err := repo.cloud.Upload(context.Background(), "sensors", data.ImageRaw)
		if err != nil {
			return err
		}
		if url != nil {
			inputDB.ImageUrl = *url
		}
	}

	inputDB.Name = data.Name
	inputDB.Description = data.Description
	inputDB.UserId = data.UserId

	if err := repo.db.Create(inputDB).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			return errors.New("duplicate: sensor name already exist")
		}

		return err
	}

	return nil
}
