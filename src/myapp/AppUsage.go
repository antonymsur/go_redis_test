package myapp
import "time"
//AppUsage represents Model AppUsage
type AppUsage struct {
	SubscriptionID int
	AppID int
	MinsUsed int // TODO change to Duration from Time Object
	Amount float64
	LastUpdated time.Time `sql:"DEFAULT:current_timestamp"`
}
