package service

import "mertani-golang/features/sensors"

type sensorService struct {
	repo sensors.Repository
}

func NewSensorService(repo sensors.Repository) sensors.Service {
	return &sensorService{
		repo: repo,
	}
}

func (srv *sensorService) Create(data sensors.Sensor) error {
	if err := srv.repo.Create(data); err != nil {
		return err
	}

	return nil
}

func (srv *sensorService) GetAll() ([]sensors.Sensor, error) {
	result, err := srv.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (srv *sensorService) Delete(SensorId, UserId uint) error {
	err := srv.repo.Delete(SensorId, UserId)
	if err != nil {
		return err
	}

	return nil
}

func (srv *sensorService) Update(UserId uint, updateSensor sensors.Sensor) error {
	err := srv.repo.Update(UserId, updateSensor)
	if err != nil {
		return err
	}

	return nil
}
