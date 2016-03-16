package myapp
// App reperesents the model APP
type App struct {
	AppID   int     `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	AppName string  `sql:"size:255;index"`
	AppURL	string 	`sql:"size:255;index"`
	Description string
	UnitPrice  float64    //Per Minute Usage
}
