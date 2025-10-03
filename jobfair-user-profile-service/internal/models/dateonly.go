package models

import (
	"time"
)

type DateOnly time.Time

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	str := string(b)
	// Hapus tanda kutip
	str = str[1 : len(str)-1]
	if str == "" || str == "null" {
		return nil
	}
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

func (d *DateOnly) ToTime() *time.Time {
	if d == nil {
		return nil // biar aman kalau EndDate memang kosong
	}
	t := time.Time(*d)
	return &t
}
