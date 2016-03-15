    package main
import (
    "fmt"
    "myapp"
    "time"
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
    "math/rand"
    "strconv"
)
const (
    dBUSER     = "postgres"
    dBPASSWORD = "postgres"
    dBNAME     = "sample"
)
func main() {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
        dBUSER, dBPASSWORD, dBNAME)
    db, err := gorm.Open("postgres", dbinfo)
    checkDBErr(err,db)
    defer db.Close()
    if db.HasTable(&myapp.User{}) == false {
        db.AutoMigrate(&myapp.User{}, &myapp.Address{}, &myapp.App{},&myapp.AppUsage{},&myapp.Payment{},&myapp.Subscription{})
    }
    cache := &myapp.RedisCache{}
    cache.Init()

    //createTestLoadRedisFromDB(db)
    var inRedis bool
    inRedis = false
    start := time.Now()
    createTestData(db,cache,inRedis) //
    elapsed := time.Since(start)
    fmt.Printf("createTestData in DB took %s\n", elapsed)
    inRedis = true
    start = time.Now()
    createTestData(db,cache,inRedis) //
    elapsed = time.Since(start)
    fmt.Printf("createTestData in Redis took %s\n", elapsed)

}
//checkDBErr
func checkDBErr(err error,db *gorm.DB) {
    if err != nil {
        panic(err)
    }
    err = db.DB().Ping()
    if err != nil {
        panic(err)
    }
}
//checkErr
func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
//createTestData
func createTestData(db *gorm.DB,cache *myapp.RedisCache,inRedis bool) {
    createUsers(db,cache,1000,inRedis)
    createApps(db,cache,50,inRedis)
    createSubscriptions(db,cache,1000,inRedis)
    createPayments(db,cache,1000,inRedis)
    //useApps(db)
}
//createTestLoadRedisFromDB
func createTestLoadRedisFromDB(db *gorm.DB,cache *myapp.RedisCache) {
    cache.InitSubscriptionsCache(db)
}

//createUsers
func createUsers (db *gorm.DB,cache *myapp.RedisCache,nUsers int,inRedis bool){
    for i := 1; i <= nUsers; i++ {
        addrLine1 := fmt.Sprintf("addr1-%d",i)
        addrLine2 := fmt.Sprintf("addr2-%d",i)
        emailID := fmt.Sprintf("test.%d@g.com",i)
        name := fmt.Sprintf("test_%d",i)
        pCode:= 560000 + i;
        address1 := myapp.Address{UserID: i,Line1:addrLine1,Line2:addrLine2,Country:"India",City:"Bangalore",PostCode:strconv.Itoa(pCode)}
        user1 := myapp.User{UserID: i,Name:name,AddressID: i,Email:emailID,Created:time.Now(),LastUpdated:time.Now()};
        if(inRedis) {
            cache.SetUser(user1,true)
        } else {
            db.Create(&address1)
            db.Create(&user1)
        }
    }
}
//Create Apps
func createApps(db *gorm.DB,cache *myapp.RedisCache,nApps int,inRedis bool) {
    var appName,appURL, appDesc string
    var appUnitPrice float64
    for  i:=1; i<=nApps;i++  {
        if i%2 == 0 {
            appName = fmt.Sprintf("hello%d",i)
            appDesc = fmt.Sprintf("This is a sample Hello%d App",i)
            appUnitPrice =  0.50

        } else {
            appName = fmt.Sprintf("help%d",i)
            appDesc = fmt.Sprintf("This is a sample Help%d App",i)
            appUnitPrice =  0.55
        }
        appURL = fmt.Sprintf("https://app%d.example.com",i)
        app := myapp.App{AppName:appName,AppURL:appURL,Description:appDesc,UnitPrice:appUnitPrice}
        if(inRedis) {
            cache.SetApp(app)
        } else {
            db.Create(&app)
        }
    }

}
//Create Subscriptions
func createSubscriptions(db *gorm.DB, cache *myapp.RedisCache,nSubs int,inRedis bool) {
    var usrName, email string
    var usrAppID int
    for i:=1;i<=nSubs;i++ {
        usrName = fmt.Sprintf("test_%d",i)
        email = fmt.Sprintf("test.%d@g.com",i)
        var usr0 myapp.User
        var apps0 myapp.App
        db.Where(&myapp.User{Name: usrName, Email: email}).First(&usr0)
        db.Where(&myapp.App{AppID: usrAppID}).First(&apps0)
        subs1 := myapp.Subscription{UserID:usr0.UserID,AppIDs: strconv.Itoa(i),Balance:0.0}
        if(inRedis) {
            cache.SetSubscription(subs1)
        } else {
            db.Create(&subs1)
        }
    }
}
func createPayments(db *gorm.DB,cache *myapp.RedisCache,nSubs int,inRedis bool) {
    var usr0 myapp.User
    var paymentAmt float64
    paymentAmt = 1000.0;
    var usrName, email string
    for i:=1;i<=nSubs;i++ {
        usrName = fmt.Sprintf("test_%d",i)
        email = fmt.Sprintf("test.%d@g.com",i)
        db.Where(&myapp.User{Name: usrName, Email: email}).First(&usr0)
        subs0 := &myapp.Subscription{}
        subs0.GetSubscriptionByUID(usr0.UserID,db)
        payment1 := myapp.Payment{SubscriptionID:subs0.ID,Amount:paymentAmt}
        subs0.Balance = subs0.Balance + paymentAmt
        if(inRedis) {
            // Make it in single transaction
            cache.SetPayment(payment1)
            cache.SetSubscription(*subs0)
        } else {
            // Make it in single transaction
            db.Create(&payment1)
            db.Save(&subs0)
        }
    }
}

func useApps(db *gorm.DB) {
    jobs := make(chan int, 100)
	for w := 1; w <= 10; w++ {
		go useApp(w, jobs,db)
	}
    for j := 1; j <= 1000; j++ {
		jobs <- j
	}
	close(jobs)
}
func useApp(id int, jobs <-chan int, db *gorm.DB) {
    appID := id%50 // Test purpose get appId with in the test data
    var app0 myapp.App
    var subs0 myapp.Subscription
    var amount float64
    db.Where(&myapp.App{AppID: appID}).First(&app0)
    minsUsed := rand.Intn(3)
    if minsUsed == 0 {
        minsUsed = 1
    }
    amount = app0.UnitPrice * float64(minsUsed)
    for j := range jobs {
        fmt.Printf("JobId %d,minsUsed = %d, appUnitPrice = %f , appID = %d\n",j,minsUsed,app0.UnitPrice,appID)
        appUsage := myapp.AppUsage{SubscriptionID:id,AppID:appID,MinsUsed:minsUsed,Amount:amount}
        db.Create(&appUsage)
        db.Where("id = ?",id).First(&subs0)
        subs0.Balance = subs0.Balance - amount
        db.Save(&subs0)
		//time.Sleep(time.Second)
	}
}
