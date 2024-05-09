package models

import "errors"

var (
	NotFoundErr = errors.New("not found")
)

type Car struct {
	Brand string    `json:"brand"`
	Model string    `json:"model"`
	Year  int16     `json:"year"`
	Motor MotorSpec `json:"motor"`
}

type MotorSpec struct {
	Size       float32 `json:"size"`
	Horsepower float32 `json:"horsepower"`
	Torque     float32 `json:"torque"`
	Max_rpm    int     `json:"max_rpm"`
}