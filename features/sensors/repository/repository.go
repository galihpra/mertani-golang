package repository

import (
	"context"
	"fmt"
	"io"
	"mertani-golang/features/sensors"
	"mertani-golang/utils/cloudinary"

	"gorm.io/gorm"
)

type Sensor struct {
	Id          uint      `gorm:"column:id; primaryKey;"`
	Name        string    `gorm:"column:name; type:varchar(200);"`
	Description string    `gorm:"column:description; type:text;"`
	ImageUrl    string    `gorm:"column:image; type:text;"`
	ImageRaw    io.Reader `gorm:"-"`

	UserId uint `gorm:"column:user_id;"`
	User   User `gorm:"foreignKey:UserId"`

	Details []Detail `gorm:"foreignKey:SensorId"`
}

type User struct {
	Id uint
}

type Detail struct {
	Id   uint   `gorm:"column:id; primaryKey;"`
	Spec string `gorm:"column:spec; type:text"`

	SensorId uint
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
	url, err := repo.cloud.Upload(context.Background(), "sensors", data.ImageRaw)
	if err != nil {
		return err
	}

	var inputDB = new(Sensor)
	inputDB.Name = data.Name
	inputDB.Description = data.Description
	inputDB.UserId = data.UserId
	inputDB.ImageUrl = *url
	inputDB.Details = make([]Detail, len(data.Details))

	if err := repo.db.Create(inputDB).Error; err != nil {
		return err
	}

	for i, detail := range data.Details {
		newDetail := Detail{
			Spec:     detail.Spec,
			SensorId: inputDB.Id,
		}

		inputDB.Details[i] = newDetail
	}
	fmt.Println(inputDB.Details)

	return nil
}
