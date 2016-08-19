//This is a little script to convert all tiff images in a directory to jpgs and
//upload them to Google Cloud Storage.
package main

import (
	"context"
	"flag"
	"image/jpeg"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"golang.org/x/image/tiff"
)

const tiffExt = ".tif"
const jpgExt = ".jpg"
const defaultBucket = "areed"

func main() {
	dirname := flag.String("directory", ".", "path to directory with the tiff source files")
	bucketName := flag.String("bucket", defaultBucket, "your Cloud Storage bucket where the images will be uploaded")
	bucketDir := flag.String("bucketdir", "", "directory in your bucket where theimages will be uploaded")
	flag.Parse()
	dirInfo, err := os.Stat(*dirname)
	if err != nil {
		log.Fatal(err)
	}
	if !dirInfo.IsDir() {
		log.Fatalf("%s is not a directory", *dirname)
	}
	dir, err := os.Open(*dirname)
	if err != nil {
		log.Fatal(err)
	}
	//prepare the google storage client
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	bucket := client.Bucket(*bucketName)
	contents, err := dir.Readdirnames(0)
	if err != nil {
		log.Fatal(err)
	}
	for _, n := range contents {
		if filepath.Ext(n) != tiffExt {
			continue
		}
		convert(bucket, path.Join(*bucketDir, strings.Replace(n, tiffExt, jpgExt, 1)), filepath.Join(*dirname, n))
	}
}

func convert(bucket *storage.BucketHandle, dst, src string) {
	println(dst)
	f, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, err := tiff.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	w := bucket.Object(dst).NewWriter(context.Background())
	defer w.Close()
	w.ContentType = "image/jpeg"
	w.ACL = []storage.ACLRule{{storage.AllUsers, storage.RoleReader}}
	err = jpeg.Encode(w, img, &jpeg.Options{Quality: 60})
	if err != nil {
		log.Fatal(err)
	}
}
