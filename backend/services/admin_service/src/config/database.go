package config

import (
    "context"
    "log"
    "time"
    "os"
    
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// ConnectDatabase kết nối đến MongoDB và trả về database
func ConnectDatabase() error {

    // Lấy giá trị từ các biến môi trường
    mongoURI := os.Getenv("MONGO_URI")
    dbName := os.Getenv("MONGO_DB_NAME")
    //os.Getenv("")

    if mongoURI == "" || dbName == "" {
        log.Fatal("Required environment variables are missing!")
    }

    // Tạo context với timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Sử dụng Connect để kết nối đến MongoDB
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        return err
    }

    // Kiểm tra kết nối
    err = client.Ping(ctx, nil)
    if err != nil {
        return err
    }

    DB = client.Database(dbName)

    log.Println("Connected to MongoDB!")
    return nil
}
