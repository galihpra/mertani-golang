package service

import (
	"mertani-golang/features/devices"
)

type deviceService struct {
	repo devices.Repository
}

func NewDeviceService(repo devices.Repository) devices.Service {
	return &deviceService{
		repo: repo,
	}
}

func (srv *deviceService) Create(data devices.Device) error {
	if err := srv.repo.Create(data); err != nil {
		return err
	}

	return nil
}

func (srv *deviceService) GetAll() ([]devices.Device, error) {
	result, err := srv.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (srv *deviceService) Delete(DeviceId, UserId uint) error {
	err := srv.repo.Delete(DeviceId, UserId)
	if err != nil {
		return err
	}

	return nil
}

func (srv *deviceService) Update(UserId uint, updateDevice devices.Device) error {
	err := srv.repo.Update(UserId, updateDevice)
	if err != nil {
		return err
	}

	return nil
}
