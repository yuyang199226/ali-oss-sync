package main

import (
	"ali-oss-sync/uploader"
	"flag"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"github.com/colinmarc/hdfs"
	"os"
	"strings"
	//    "strconv"
	"path/filepath"
	"time"
)

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func main() {
	var (
		endpoint     = flag.String("endpoint", "", "ali oss endpoint")
		bucket       = flag.String("bucket", "", "ali oss bucket name")
		localprefix  = flag.String("localprefix", "", "local upload prefix dir")
		objectprefix = flag.String("objectprefix", "", "object upload prefix dir")
		remain       = flag.Int("remain", 7, "scan and save days")
		done         = flag.String("done", "__done__", "done file name to indicate to uupload")
	)
	flag.Parse()
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ACCESS_KEY_SECRET")
	fmt.Println("OSS Go SDK Version: ", oss.Version, endpoint)
	client, err := oss.New(*endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		log.Fatal(err)
	}
    hdfsClient, err := hdfs.New("manage13.aibee.cn:8020")
    if err != nil {
		log.Fatal(err)
    }
    hdfsDir, err := hdfsClient.ReadDir("/prod")
    if err != nil {
		log.Fatal(err)
    }
    fmt.Println(hdfsDir)
	bucketObj, err := client.Bucket(*bucket)
	uploader := uploader.Uploader{client, bucketObj, *localprefix, *objectprefix, *remain, *done}

	for i := 0; i <= *remain; i++ {
		now := time.Now()
		d := now.AddDate(0, 0, -i)
		d_str := d.Format("20060102")
		srcDir := filepath.Join(*localprefix, d_str)
		if *done != "" {
			doneFilePath := filepath.Join(srcDir, *done)
			if FileExist(doneFilePath) {
				err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
						return err
					}
					if info.IsDir() {
						fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
						return nil
					} else {
						fmt.Printf("visited file : %q\n", path)
						suffixPath := strings.Replace(path, srcDir, "", 1)
						objectKey := filepath.Join(*objectprefix, d_str, suffixPath)
						fmt.Printf("UpLoad %s -> %s\n", path, objectKey)
						err := uploader.Upload(path, objectKey)
						return err
					}
				})
				if err != nil {
					fmt.Printf("error walking the path %q: %v\n", srcDir, err)
					return
				} else {
					err := os.Remove(doneFilePath)
					if err != nil {
						fmt.Println("del done file failed")
					}
				}
			}
		}

	}
	if uploader.ShouldUpload("/tmp/test.txt", "test.txt") {
		fmt.Println("upload")
		uploader.SimpleUpload("/tmp/test.txt", "test.txt")
	}
}
