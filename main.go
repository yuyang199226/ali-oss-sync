package main
import (
    "fmt"
    "github.com/aliyun/aliyun-oss-go-sdk/oss"
    "ali-oss-sync/uploader"
    "flag"
    "os"
    "log"
)
func main() {
    var (
        endpoint          = flag.String("endpoint", "", "ali oss endpoint")
        bucket            = flag.String("bucket", "", "ali oss bucket name")
    )
    flag.Parse()
    accessKeyId := os.Getenv("ACCESS_KEY_ID")
    accessKeySecret := os.Getenv("ACCESS_KEY_SECRET")
    fmt.Println("OSS Go SDK Version: ", oss.Version, endpoint)
    client , err := oss.New(*endpoint, accessKeyId, accessKeySecret)
    if err != nil {
        log.Fatal(err)
    }
    bucketObj, err := client.Bucket(*bucket)
    uploader := uploader.Uploader{client, bucketObj, "/tmp"}
    if uploader.ShouldUpload("/tmp/test.txt", "test.txt") {
        fmt.Println("upload")
        uploader.SimpleUpload("/tmp/test.txt", "test.txt")
    }
}
