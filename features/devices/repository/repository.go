package repository

import (
	"context"
	"errors"
	"io"
	"mertani-golang/features/devices"
	"mertani-golang/utils/cloudinary"
	"strings"

	"gorm.io/gorm"
)

type Device struct {
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

type deviceRepository struct {
	db    *gorm.DB
	cloud cloudinary.Cloud
}

func NewDeviceRepository(db *gorm.DB, cloud cloudinary.Cloud) devices.Repository {
	return &deviceRepository{
		db:    db,
		cloud: cloud,
	}
}

func (repo *deviceRepository) Create(data devices.Device) error {
	var inputDB = new(Device)
	if data.ImageRaw != nil {
		url, err := repo.cloud.Upload(context.Background(), "devices", data.ImageRaw)
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
			return errors.New("duplicate: device name already exist")
		}

		return err
	}

	return nil
}

func (repo *deviceRepository) GetAll() ([]devices.Device, error) {
	var dataDevice []Device

	if err := repo.db.Find(&dataDevice).Error; err != nil {
		return nil, err
	}

	var result []devices.Device
	for _, device := range dataDevice {
		result = append(result, devices.Device{
			Id:          device.Id,
			Name:        device.Name,
			ImageUrl:    device.ImageUrl,
			Description: device.Description,
			UserId:      device.UserId,
		})
	}

	return result, nil
}

func (repo *deviceRepository) Delete(DeviceId, UserId uint) error {
	var data Device

	if err := repo.db.First(&data, DeviceId).Error; err != nil {
		return errors.New("device not found")
	}

	if data.UserId != UserId {
		return errors.New("not authorized: you are not authorized to delete this device")
	}

	if err := repo.db.Delete(&data).Error; err != nil {
		return err
	}

	return nil
}

func (repo *deviceRepository) Update(UserId uint, updateDevice devices.Device) error {
	var data Device

	if err := repo.db.First(&data, updateDevice.Id).Error; err != nil {
		return errors.New("device not found")
	}

	if data.UserId != UserId {
		return errors.New("not authorized: you are not authorized to update this device")
	}

	if updateDevice.ImageRaw != nil {
		url, err := repo.cloud.Upload(context.Background(), "devices", updateDevice.ImageRaw)
		if err != nil {
			return err
		}
		if url != nil {
			updateDevice.ImageUrl = *url
		}
	}

	if err := repo.db.Model(&data).Updates(Device{
		Name:        updateDevice.Name,
		Description: updateDevice.Description,
		ImageUrl:    updateDevice.ImageUrl,
	}).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			return errors.New("duplicate: device name already exist")
		}

		return err
	}

	return nil
}
