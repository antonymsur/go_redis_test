package myapp
import "time"
// User represents model User
type User struct {
    UserID    int        `sql:"AUTO_INCREMENT" gorm:"primary_key"`
    Name      string     `sql:"size:255;unique;index"`
    /*Address   Address */   
    AddressID int
    Email     string
    Created   time.Time `sql:"DEFAULT:current_timestamp"`
    LastUpdated time.Time `sql:"DEFAULT:current_timestamp"`
}
