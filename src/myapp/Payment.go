package myapp
import (
  "time"
)
//Payment represents Model Payment
type Payment struct {
	PaymentID  int   `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	SubscriptionID   uint
	Amount  float64
    DatePaid time.Time     `sql:"DEFAULT:current_timestamp"`
}
