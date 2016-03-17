package myapp //same package name as source file

import (
    "testing" //import go package for testing related functionality
    "fmt"
    "time"
    "strconv"
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
    "gopkg.in/redis.v3"
)

//Benchmarking  function starts with "Benchmark" and takes a pointer to type testing.B
var (
    db *gorm.DB
	err error
    cache *RedisCache
    options redis.Options
    dbinfo string
    client *redis.Client
)
const (
    dBUSER     = "postgres"
    dBPASSWORD = "postgres"
    dBNAME     = "sample"
    isRedis    = true // This flag needs to be turn off to run Postgresql benchmark
    RedisAddress  = "localhost:6379"
    RedisPass = ""
    RedisDB = 0
)

func setup() {
    if(isRedis) {
        cache = &RedisCache{}
        options = redis.Options{
            Addr:     RedisAddress,
            Password: RedisPass,
            DB:       RedisDB,
        }
        cache.Init(&options)
        client = cache.client
    } else {
        dbinfo = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
            dBUSER, dBPASSWORD, dBNAME)
        db, err = gorm.Open("postgres", dbinfo)
        checkDBErr(err,db)
        //defer db.Close()
        if db.HasTable(&User{}) == false {
            db.AutoMigrate(&User{}, &Address{}, &App{},&AppUsage{},&Payment{},&Subscription{})
        }
    }
}

func teardown() {
    if (isRedis) {
        cache.Close()
    } else {
         db.Close()
    }
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

func testLoadDataFromDB(t *testing.T) {
    //createTestLoadRedisFromDB(db)
}

//BenchmarkSetUser test
func BenchmarkAddUser(b *testing.B) {
    b.StopTimer() //stop the performance timer temporarily while doing initialization
    setup()
    defer teardown()
    b.StartTimer() //restart timer
   for i := 0; i < b.N; i++ {
       addrLine1 := fmt.Sprintf("addr1-%d",i)
       addrLine2 := fmt.Sprintf("addr2-%d",i)
       emailID := fmt.Sprintf("test.%d@g.com",i)
       name := fmt.Sprintf("test_%d",i)
       pCode:= 560000 + i;
       address1 := Address{UserID: i,Line1:addrLine1,Line2:addrLine2,Country:"India",City:"Bangalore",PostCode:strconv.Itoa(pCode)}
       user1 := User{Name:name,AddressID: i,Email:emailID,Created:time.Now(),LastUpdated:time.Now()};
       if(isRedis) {
           cache.AddAddress(address1)
           cache.AddUser(user1)
       } else {
           db.Create(&address1)
           db.Create(&user1)
       }
   }
    b.Log("BenchmarkAddUser passed")
}

func BenchmarkAddSubscriptions(b *testing.B) {
    b.StopTimer() //stop the performance timer temporarily while doing initialization
    setup()
    defer teardown()
    b.StartTimer() //restart timer
    var usrName, email string
    var usrID int
    for i := 0; i < b.N; i++ {
        usrName = fmt.Sprintf("test_%d",i)
        email = fmt.Sprintf("test.%d@g.com",i)
        var usr0 User
        if(isRedis) {
            userID := client.Get("user:username:"+usrName).Val()
            usrID,err = strconv.Atoi(userID)
            subs1 := Subscription{UserID:usrID,AppIDs: strconv.Itoa(i),Balance:0.0}
            cache.AddSubscription(subs1)
        } else {
            db.Where(&User{Name: usrName, Email: email}).First(&usr0)
            subs1 := Subscription{UserID:usr0.UserID,AppIDs: strconv.Itoa(i),Balance:0.0}
            db.Create(&subs1)
        }
    }
    b.Log("BenchmarkAddSubscriptions passed")
}

func BenchmarkAddApps(b *testing.B) {
    b.StopTimer() //stop the performance timer temporarily while doing initialization
    setup()
    defer teardown()
    b.StartTimer() //restart timer

    var appName,appURL, appDesc string
    var appUnitPrice float64
    for  i:=1; i<=b.N;i++  {
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
        app := App{AppName:appName,AppURL:appURL,Description:appDesc,UnitPrice:appUnitPrice}
        if(isRedis) {
           cache.AddApp(app)
        } else {
           db.Create(&app)
        }
    }
    b.Log("BenchmarkAddApps passed")
}
