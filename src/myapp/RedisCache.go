package myapp
import (
    "gopkg.in/redis.v3"
    "fmt"
    "github.com/jinzhu/gorm"
    "time"
)
/*var client *redis.Client*/

//RedisCache used to store the client
type RedisCache struct {
    client *redis.Client
}

//Init initiates the connection to the redis server and creates client
func (rc *RedisCache) Init(options *redis.Options) {
    rc.client = redis.NewClient(options)
    err := rc.client.Ping().Err()
    if err != nil {
        panic(err)
    }

}
//InitSubscriptionsCache used to load subscription data
//Uses gorm as ORM to get All the subscription and load
// This is a sample how to load cache from DB
// TODO : Number of rows can be  high so need to do pagination approach
func (rc *RedisCache) InitSubscriptionsCache(db *gorm.DB) {
    var subs []Subscription
    constKey := "Subscription"
    db.Find(&subs)
    length := len(subs)
    for i:=0; i<length; i++ {
        key := fmt.Sprintf("%s-%d",constKey,subs[i].ID);
        err := rc.client.HSet(key, "UserID", string(subs[i].UserID)).Err()
        if err != nil {
            panic(err)
        }
        err = rc.client.HSet(key, "AppIDs", string(subs[i].AppIDs)).Err()
        if err != nil {
            panic(err)
        }
        err = rc.client.HSet(key, "Balance", fmt.Sprintf("%0.2f",subs[i].Balance)).Err()
        if err != nil {
            panic(err)
        }
        err = rc.client.HSet(key, "StartDate", subs[i].StartDate.String()).Err()
        if err != nil {
            panic(err)
        }
    }
}

//AddSubscription to set to cache. Using the HSet to store all fields as in Relational DB
// Application creates Subscription object
//TODO If any additional validations
func (rc *RedisCache) AddSubscription(sub Subscription) {
    constKey := "Subscription-"
    client := rc.client;
    subsID := client.Incr("subscription:id").Val()
    err := client.HMSet(fmt.Sprintf("%s%d",constKey,subsID),"UserID", fmt.Sprintf("%d",sub.UserID),"AppIDs", sub.AppIDs,"Balance", fmt.Sprintf("%0.2f",sub.Balance),"StartDate", sub.StartDate.String()).Err()
    if err != nil {
        panic(err)
    }
}

//SetSubscription to set to cache. Using the HSet to store all fields as in Relational DB
// Application creates Subscription object
//TODO If any additional validations
func (rc *RedisCache) SetSubscription(sub Subscription) {
    constKey := "Subscription"
    key := fmt.Sprintf("%s-%d",constKey,sub.ID);
    err := rc.client.HSet(key, "UserID", string(sub.UserID)).Err()
    if err != nil {
        panic(err)
    }
    err = rc.client.HSet(key, "AppIDs", sub.AppIDs).Err()
    if err != nil {
        panic(err)
    }
    err = rc.client.HSet(key, "Balance", fmt.Sprintf("%0.2f",sub.Balance)).Err()
    if err != nil {
        panic(err)
    }
    err = rc.client.HSet(key, "StartDate", sub.StartDate.String()).Err()
    if err != nil {
        panic(err)
    }
}

//SetApp to set to cache. Using the HSet to store all fields as in Relational DB
// Application creates Subscription object
//TODO If any additional validations
func (rc *RedisCache) SetApp(app App) {
    constKey := "App"
    key := fmt.Sprintf("%s-%d",constKey,app.AppID);
    err := rc.client.HSet(key, "AppName", string(app.AppName)).Err()
    if err != nil {
        panic(err)
    }
    err = rc.client.HSet(key, "AppURL", app.AppURL).Err()
    if err != nil {
        panic(err)
    }
    err = rc.client.HSet(key, "Description", app.Description).Err()
    if err != nil {
        panic(err)
    }
    err = rc.client.HSet(key, "UnitPrice", fmt.Sprintf("%0.6f",app.UnitPrice)).Err()
    if err != nil {
        panic(err)
    }
}

//AddApp to set to cache. Using the HSet to store all fields as in Relational DB
// Application creates Subscription object
//TODO If any additional validations
func (rc *RedisCache) AddApp(app App) {
    constKey := "App-"
    client := rc.client;
    appID := client.Incr("app:id").Val()
    err := client.HMSet(fmt.Sprintf("%s%d",constKey,appID),"AppName",app.AppName,"AppURL", app.AppURL, "Description",app.Description,"UnitPrice", fmt.Sprintf("%0.6f",app.UnitPrice)).Err()
    if err != nil {
        panic(err)
    }
}

//AddPayment to set to cache. Using the HSet to store all fields as in Relational DB
// Application creates Payment object
//TODO If any additional validations
func (rc *RedisCache) AddPayment(pay Payment) {
    constKey := "Payment-"
    client := rc.client;
    payID := client.Incr("payment:id").Val()
    err := client.HMSet(fmt.Sprintf("%s%d",constKey,payID),"SubscriptionID",fmt.Sprintf("%d",pay.SubscriptionID),"DatePaid",pay.DatePaid.String(),"Amount",fmt.Sprintf("%0.6f",pay.Amount)).Err()
    if err != nil {
        panic(err)
    }
}

//SetPayment to set to cache. Using the HSet to store all fields as in Relational DB
// Application creates Payment object
//TODO If any additional validations
func (rc *RedisCache) SetPayment(pay Payment) {
    constKey := "Payment"
    key := fmt.Sprintf("%s-%d",constKey,pay.PaymentID);
    err := rc.client.HSet(key, "SubscriptionID", string(pay.SubscriptionID)).Err()
    if err != nil {
        panic(err)
    }
    err = rc.client.HSet(key, "DatePaid", pay.DatePaid.String()).Err()
    if err != nil {
        panic(err)
    }
    err = rc.client.HSet(key, "Amount", fmt.Sprintf("%0.6f",pay.Amount)).Err()
    if err != nil {
        panic(err)
    }
}


//SetAppUsage to set to cache. Using the HSet to store all fields as in Relational DB
//Assume that AppUsage is  cerated every 30secs or configured time period
//TODO If any additional validations
func (rc *RedisCache) SetAppUsage(appUsage AppUsage) {
    constKey := "AppUsage"
    key := fmt.Sprintf("%s-%d-%d",constKey,appUsage.SubscriptionID,appUsage.AppID);
    value := fmt.Sprintf("%d-%0.6f%s",appUsage.MinsUsed,appUsage.Amount,appUsage.LastUpdated.String())

    err := rc.client.RPush(key, value).Err()
    if err != nil {
        panic(err)
    }

    key = fmt.Sprintf("%s-%d","App",appUsage.AppID)
    unitprice,err := rc.client.HGet(key,"UnitPrice").Float64()

    key = fmt.Sprintf("%s-%d","Subscription",appUsage.SubscriptionID)
    credit,err := rc.client.HGet(key,"Balance").Float64()


    credit =  credit - (unitprice * float64(appUsage.MinsUsed))

    err = rc.client.HSet(key, "Balance", fmt.Sprintf("%0.6f",credit)).Err()
    if err != nil {
        panic(err)
    }

}

//AddAppUsage to set to cache. Using the HSet to store all fields as in Relational DB
//Assume that AppUsage is  cerated every 30secs or configured time period
//TODO If any additional validations
func (rc *RedisCache) AddAppUsage(appUsage AppUsage) {
    constKey := "AppUsage-"
    client := rc.client;
    appUsageID := client.Incr("appusage:id").Val()
    err := client.HMSet(fmt.Sprintf("%s%d",constKey,appUsageID),"SubscriptionID",fmt.Sprintf("%d",appUsage.SubscriptionID),"AppID",fmt.Sprintf("%d",appUsage.AppID),"MinsUsed",fmt.Sprintf("%d",appUsage.MinsUsed)).Err()
    if err != nil {
        panic(err)
    }
    key := fmt.Sprintf("%s-%d","App",appUsage.AppID)
    unitprice,err := client.HGet(key,"UnitPrice").Float64()

    key = fmt.Sprintf("%s-%d","Subscription",appUsage.SubscriptionID)
    credit,err := rc.client.HGet(key,"Balance").Float64()
    credit =  credit - (unitprice * float64(appUsage.MinsUsed))

    err = rc.client.HSet(key, "Balance", fmt.Sprintf("%0.6f",credit)).Err()
    if err != nil {
        panic(err)
    }

}

//SetUser to set to cache. Using the HSet to store all fields as in Relational DB
// Application creates User object on registration
//TODO If any additional validations
func (rc *RedisCache) SetUser(usr User, isNew bool) {
    constKey := "User"
    key := fmt.Sprintf("%s-%d",constKey,usr.UserID);
    err := rc.client.HSet(key, "Name", string(usr.UserID)).Err()
    if err != nil {
        panic(err)
    }
    err = rc.client.HSet(key, "AddressID", string(usr.AddressID)).Err()
    if err != nil {
        panic(err)
    }
    err = rc.client.HSet(key, "Email", usr.Email).Err()
    if err != nil {
        panic(err)
    }
    timeNowStr := time.Now().String()
    if isNew {
        err = rc.client.HSet(key, "Created", timeNowStr).Err()
        if err != nil {
            panic(err)
        }
    }
    err = rc.client.HSet(key, "LastUpdated", timeNowStr).Err()
    if err != nil {
        panic(err)
    }
}


//AddUser to set to cache. Using the HSet to store all fields as in Relational DB
// Application creates User object on registration
//TODO If any additional validations
func (rc *RedisCache) AddUser(usr User) {
    client := rc.client;
    userID := client.Incr("user:id").Val()
    timeNow := time.Now().String()
    err := client.HMSet(fmt.Sprintf("User-%d",userID),"Name",usr.Name,"AddressID",fmt.Sprintf("%d",usr.AddressID), "Email",usr.Email,"Created", timeNow,"Updated",timeNow).Err()
    err = client.Set("user:name:"+usr.Name,userID,0).Err()
    if err != nil {
        panic(err)
    }
}

//AddAddress to set to cache. Using the HSet to store all fields as in Relational DB
// Application creates User object on registration
//TODO If any additional validations
func (rc *RedisCache) AddAddress(addr Address) {
    client := rc.client;
    addrID := client.Incr("addr:id").Val()
    err := client.HMSet(fmt.Sprintf("Address-%d",addrID),"UserID",fmt.Sprintf("%d",addr.UserID),"Line1",addr.Line1,"Line2",addr.Line2, "Country",addr.Country,"City", addr.City,"PostCode",addr.PostCode).Err()
    if err != nil {
        panic(err)
    }
}

//SetAddress to set to cache. Using the HSet to store all fields as in Relational DB
// Application creates User object on registration
//TODO If any additional validations
func (rc *RedisCache) SetAddress(addr Address) {
    constKey := "Address"
    key := fmt.Sprintf("%s-%d",constKey,addr.UserID);
    err := rc.client.HSet(key, "UserID", string(addr.UserID)).Err()
    if err != nil {
        panic(err)
    }
    err = rc.client.HSet(key, "Line1", string(addr.Line1)).Err()
    if err != nil {
        panic(err)
    }
    err = rc.client.HSet(key, "Line2", addr.Line2).Err()
    if err != nil {
        panic(err)
    }
    err = rc.client.HSet(key, "Country", addr.Country).Err()
        if err != nil {
            panic(err)
    }
    err = rc.client.HSet(key, "City", addr.City).Err()
    if err != nil {
        panic(err)
    }
    err = rc.client.HSet(key, "PostCode", addr.PostCode).Err()
    if err != nil {
        panic(err)
    }
}
