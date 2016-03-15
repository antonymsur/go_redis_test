package myapp
import (
	"time"
	"github.com/jinzhu/gorm"
)
// Subscription represents Model Subscription
type Subscription struct {
	ID uint       //Avoid gorm.Model tto avoid dependency in the Model struct
	UserID int     //TODO make it as associated relation, foreign key with user table
	AppIDs string //Apps []App TODO change to appArray if ORM can handle
	Balance float64 // May be to be renamed to CreditBalance?
	StartDate time.Time  `sql:"DEFAULT:current_timestamp"`
}
//GetSubscriptionByUID find The subscriptin by user id
func (arg *Subscription) GetSubscriptionByUID(uid int, db *gorm.DB) {
	db.Where(Subscription{UserID: uid}).First(arg)
}
