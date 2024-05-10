package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
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

var (
	ErrNotFound = errors.New("not found")
)

type redisHandler struct{
	client *redis.Client
}

func NewRedisHandler() *redisHandler{
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB:0,
	})
	rh := redisHandler{client}
	ping,err :=	rh.client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	fmt.Println(ping)
	return &rh
}

func (rh redisHandler) Add(key string, car Car) error{
	c,err:= json.Marshal(car)
	if err != nil {
		return err
	}
	err = rh.client.Set(context.Background(), key,c,0).Err()
	if err != nil {
		fmt.Printf("Fail to Set Pair Value Key:%v, Value:%v\n",key,c)
		return err
	}
	fmt.Printf("Key: %v\n",key)
	return nil
}
func (rh redisHandler) Get(key string) (Car, error){
	c:=Car{}
	val, err := rh.client.Get(context.Background(),key).Result()
	if err != nil {
		fmt.Println("Failed to GET", err)
		return c,ErrNotFound
	}
	
	err = json.Unmarshal([]byte(val), &c)
	if err != nil {
		fmt.Println("Failed to Unmarshal", err)
		return c,err
	}
	return c,nil
}
func (rh redisHandler) Update(key string, car Car) error{
	return rh.Add(key, car)
}
func (rh redisHandler) List() (map[string]Car, error){
	m := make(map[string]Car)
	k:=rh.client.Keys(context.Background(),"*")
	for _, value:=range k.Val(){
		fmt.Println(value)
		val,err:= rh.Get(value)
		if err != nil {
			return nil, err
		}
		m[value]=val
	}
	return m, nil
}
func (rh redisHandler) Remove(key string) error{
	resp:= rh.client.Del(context.Background(),key)
	if resp.Val() == 0 {
		return ErrNotFound
	}
	return nil
}