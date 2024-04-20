package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	Key     string
	Created time.Time
}

type OTPMap map[string]OTP

func NewOTPMap(ctx context.Context, period time.Duration) OTPMap {
	om := make(OTPMap)
	go om.Retention(ctx, period)
	return om
}

func (om OTPMap) NewOTP() OTP {
	o := OTP{
		Key:     uuid.NewString(),
		Created: time.Now(),
	}
	om[o.Key] = o
	log.Println(om)
	return o
}

func (om OTPMap) VerifyOTP(key string) bool {
	if _, ok := om[key]; !ok {
		return false
	}
	log.Println(key)
	delete(om, key)
	return true
}

func (om OTPMap) Retention(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(400 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			for _, otp := range om {
				if otp.Created.Add(period).Before(time.Now()) {
					delete(om, otp.Key)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
