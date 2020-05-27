package uploader
import (
    "fmt"
    "os"
    "strconv"
    "time"
    "github.com/aliyun/aliyun-oss-go-sdk/oss"
)
type Uploader struct {
    Client *oss.Client
    Bucket *oss.Bucket
    LocalPrefix string
    ObjectPrefix string
    Remain int
    Done string
}

func handleError(err error) {
    fmt.Println("Error: ", err)
    os.Exit(-1)
}

func (uploader *Uploader) Sync() {
    fmt.Println("sync")


}

func (uploader *Uploader) SimpleUpload(srcFilePath string, dstFilePath string) {
    err := uploader.Bucket.PutObjectFromFile(dstFilePath, srcFilePath)
     if err != nil {
        handleError(err)
     }

}

func (uploader*Uploader) ShouldUpload(srcFilePath string, dstFilePath string) bool {
    isExist, err := uploader.Bucket.IsObjectExist(dstFilePath)
    if err != nil {
        handleError(err)
    }
    if !isExist {
        fmt.Println("object not exist: ", dstFilePath)
        return true
    } else {
        props, err := uploader.Bucket.GetObjectDetailedMeta(dstFilePath)
        if err != nil {
        handleError(err)
        }
        dstSize, err := strconv.ParseInt(props["Content-Length"][0], 10, 64)
        if err != nil {
            handleError(err)
        }
        srcSize := GetLocalFileSize(srcFilePath)
        dstModTimeTime, err := time.Parse(time.RFC1123, props["Last-Modified"][0])
        if err != nil {
            handleError(err)
        }
        dstModTime := dstModTimeTime.Unix() 
        srcModTime := GetLocalFileModTime(srcFilePath)
        if (srcModTime > dstModTime) || (srcModTime == dstModTime && srcSize != dstSize) {
            return true
        } else {
            return false
        }
    }

}
func (uploader *Uploader) Upload (srcFilePath string, dstFilePath string) error {
    err := uploader.Bucket.UploadFile(dstFilePath, srcFilePath, 100*1024, oss.Routines(3), oss.Checkpoint(true, ""))
    return err

}
func GetLocalFileSize(filePath string) int64 {
         f, err := os.Open(filePath)
         defer f.Close()
         if err != nil {
             handleError(err)
         }
         info, err := f.Stat()
         if err != nil {
             handleError(err)
         }
         return info.Size()
}


func GetLocalFileModTime(filePath string) int64 {
         f, err := os.Open(filePath)
         defer f.Close()
         if err != nil {
             handleError(err)
         }
         info, err := f.Stat()
         if err != nil {
             handleError(err)
         }
         return info.ModTime().Unix()
}


