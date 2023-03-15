/**
 * @file mongo_test.go
 * @author dezhenzhao
 * @brief
 * @version 0.1
 * @date 2023-03-14
 * @copyright Copyright (c) 2021 The powermock Authors. All rights reserved.
**/

package mongo

import (
	"context"
	"encoding/base64"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
	"time"
)

func TestPlugin_GetAnnouncement(t *testing.T) {
	// 创建一条mock 数据，发送load 请求
	m := Plugin{}
	m.GetAnnouncement() <- struct{}{}
}

var MongoDBClient *mongo.Database

func TestPlugin_New(t *testing.T) {
	user := "root"
	password := "223238"
	host := "localhost"
	port := "27017"
	dbName := "runner_go"
	timeOut := 2
	maxNum := 50

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?w=majority", user, password, host, port)
	// 设置连接超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut))
	defer cancel()
	// 通过传进来的uri连接相关的配置
	o := options.Client().ApplyURI(uri)
	// 设置最大连接数 - 默认是100 ，不设置就是最大 max 64
	o.SetMaxPoolSize(uint64(maxNum))
	// 发起链接
	client, err := mongo.Connect(ctx, o)
	if err != nil {
		fmt.Println("ConnectToDB", err)
		return
	}
	// 判断服务是不是可用
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		fmt.Println("ConnectToDB", err)
		return
	}
	// 返回 client
	dbName = "runner_go"

	MongoDBClient = client.Database(dbName)
	find()
}

func insert() {
	ash := MockData{
		Path:      "/v1/hello",
		Method:    []string{"post", "get"},
		UniqueKey: "uuid",
		Mock:      "eyJ1bmlxdWVLZXkiOiJhZHZhbmNlZF9leGFtcGxlIiwicGF0aCI6Ii9hZHYvSGVsbG8iLCJtZXRob2QiOiJHRVQiLCJjYXNlcyI6W3siY29uZGl0aW9uIjp7InNpbXBsZSI6eyJpdGVtcyI6W3sib3BlcmFuZFgiOiIkcmVxdWVzdC5ib2R5LnVpZCIsIm9wZXJhdG9yIjoiPD0iLCJvcGVyYW5kWSI6IjEwMDAifV19fSwicmVzcG9uc2UiOnsic2ltcGxlIjp7ImJvZHkiOnsieC11bml0LWlkIjoiMyIsIngtdW5pdC1yZWdpb24iOiJzaCJ9LCJ0cmFpbGVyIjp7IngtYXBpLXZlcnNpb24iOiIxLjMuMiJ9LCJib2R5Ijoie1widGltZXN0YW1wXCI6IFwiMTExMVwiLCBcIm1lc3NhZ2VcIjogXCJUaGlzIG1lc3NhZ2Ugd2lsbCBvbmx5IGJlIHJldHVybmVkIHdoZW4gdWlkIDw9IDEwMDBcIiwgXCJhbW91bnRcIjogXCJ7eyAkbW9jay5wcmljZSB9fVwifVxuIn19fSx7ImNvbmRpdGlvbiI6eyJzaW1wbGUiOnsiaXRlbXMiOlt7Im9wZXJhbmRYIjoiJHJlcXVlc3QuYm9keS51aWQiLCJvcGVyYXRvciI6Ij4iLCJvcGVyYW5kWSI6IjEwMDAifV19fSwicmVzcG9uc2UiOnsic2NyaXB0Ijp7ImxhbmciOiJqYXZhc2NyaXB0IiwiY29udGVudCI6IihmdW5jdGlvbigpe1xuICAgIGZ1bmN0aW9uIHJhbmRvbShtaW4sIG1heCl7XG4gICAgICAgIHJldHVybiBwYXJzZUludChNYXRoLnJhbmRvbSgpKihtYXgtbWluKzEpK21pbiwxMCk7XG4gICAgfVxuICAgIHJldHVybiB7XG4gICAgICAgIGNvZGU6IDAsXG4gICAgICAgIGJvZHk6IHtcbiAgICAgICAgICAgIFwieC11bml0LWlkXCI6IChyZXF1ZXN0LmJvZHlbXCJ1aWRcIl0gJSA1KS50b1N0cmluZygpLFxuICAgICAgICAgICAgXCJ4LXVuaXQtcmVnaW9uXCI6IFwiYmpcIixcbiAgICAgICAgfSxcbiAgICAgICAgdHJhaWxlcjoge1xuICAgICAgICAgICAgXCJ4LWFwaS12ZXJzaW9uXCI6IFwiMS4zLjJcIixcbiAgICAgICAgfSxcbiAgICAgICAgYm9keToge1xuICAgICAgICAgICAgdGltZXN0YW1wOiBNYXRoLmNlaWwobmV3IERhdGUoKS5nZXRUaW1lKCkgLyAxMDAwKSxcbiAgICAgICAgICAgIG1lc3NhZ2U6IFwidGhpcyBtZXNzYWdlIGlzIGdlbmVyYXRlZCBieSBqYXZhc2NyaXB0LCB5b3VyIHVpZCBpczogXCIgKyByZXF1ZXN0LmJvZHlbXCJ1aWRcIl0sXG4gICAgICAgICAgICBhbW91bnQ6IHJhbmRvbSgwLCA1MDAwKSxcbiAgICAgICAgfSxcbiAgICB9XG59KSgpIn19fV19",
	}
	insertResult, err := MongoDBClient.Collection("mock").InsertOne(context.TODO(), ash)
	if err != nil {
		fmt.Println(err)
	}
	println("Inserted a single document: ", insertResult.InsertedID)
}

func find() {
	findOptions := options.Find()
	findOptions.SetLimit(10)
	cur, err := MongoDBClient.Collection("mock").Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		fmt.Println(err)
	}
	var results []*MockData
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem MockData
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		fmt.Println(err)
	}
	for _, v := range results {
		fmt.Println(v.UniqueKey)
	}
}

func TestDecode(t *testing.T) {
	b := "eyJ1bmlxdWVLZXkiOiJhZHZhbmNlZF9leGFtcGxlIiwicGF0aCI6Ii9hZHYvSGVsbG8iLCJtZXRob2QiOiJHRVQiLCJjYXNlcyI6W3siY29uZGl0aW9uIjp7InNpbXBsZSI6eyJpdGVtcyI6W3sib3BlcmFuZFgiOiIkcmVxdWVzdC5ib2R5LnVpZCIsIm9wZXJhdG9yIjoiPD0iLCJvcGVyYW5kWSI6IjEwMDAifV19fSwicmVzcG9uc2UiOnsic2ltcGxlIjp7ImJvZHkiOnsieC11bml0LWlkIjoiMyIsIngtdW5pdC1yZWdpb24iOiJzaCJ9LCJ0cmFpbGVyIjp7IngtYXBpLXZlcnNpb24iOiIxLjMuMiJ9LCJib2R5Ijoie1widGltZXN0YW1wXCI6IFwiMTExMVwiLCBcIm1lc3NhZ2VcIjogXCJUaGlzIG1lc3NhZ2Ugd2lsbCBvbmx5IGJlIHJldHVybmVkIHdoZW4gdWlkIDw9IDEwMDBcIiwgXCJhbW91bnRcIjogXCJ7eyAkbW9jay5wcmljZSB9fVwifVxuIn19fSx7ImNvbmRpdGlvbiI6eyJzaW1wbGUiOnsiaXRlbXMiOlt7Im9wZXJhbmRYIjoiJHJlcXVlc3QuYm9keS51aWQiLCJvcGVyYXRvciI6Ij4iLCJvcGVyYW5kWSI6IjEwMDAifV19fSwicmVzcG9uc2UiOnsic2NyaXB0Ijp7ImxhbmciOiJqYXZhc2NyaXB0IiwiY29udGVudCI6IihmdW5jdGlvbigpe1xuICAgIGZ1bmN0aW9uIHJhbmRvbShtaW4sIG1heCl7XG4gICAgICAgIHJldHVybiBwYXJzZUludChNYXRoLnJhbmRvbSgpKihtYXgtbWluKzEpK21pbiwxMCk7XG4gICAgfVxuICAgIHJldHVybiB7XG4gICAgICAgIGNvZGU6IDAsXG4gICAgICAgIGJvZHk6IHtcbiAgICAgICAgICAgIFwieC11bml0LWlkXCI6IChyZXF1ZXN0LmJvZHlbXCJ1aWRcIl0gJSA1KS50b1N0cmluZygpLFxuICAgICAgICAgICAgXCJ4LXVuaXQtcmVnaW9uXCI6IFwiYmpcIixcbiAgICAgICAgfSxcbiAgICAgICAgdHJhaWxlcjoge1xuICAgICAgICAgICAgXCJ4LWFwaS12ZXJzaW9uXCI6IFwiMS4zLjJcIixcbiAgICAgICAgfSxcbiAgICAgICAgYm9keToge1xuICAgICAgICAgICAgdGltZXN0YW1wOiBNYXRoLmNlaWwobmV3IERhdGUoKS5nZXRUaW1lKCkgLyAxMDAwKSxcbiAgICAgICAgICAgIG1lc3NhZ2U6IFwidGhpcyBtZXNzYWdlIGlzIGdlbmVyYXRlZCBieSBqYXZhc2NyaXB0LCB5b3VyIHVpZCBpczogXCIgKyByZXF1ZXN0LmJvZHlbXCJ1aWRcIl0sXG4gICAgICAgICAgICBhbW91bnQ6IHJhbmRvbSgwLCA1MDAwKSxcbiAgICAgICAgfSxcbiAgICB9XG59KSgpIn19fV19="
	str, _ := base64.StdEncoding.DecodeString(b)
	fmt.Println(string(str))
}
