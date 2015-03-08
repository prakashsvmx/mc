/*
 * Mini Object Storage, (C) 2014,2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/codegangsta/cli"
	"github.com/minio-io/mc/pkg/s3"
)

const (
	recvFormat  = "2006-01-02T15:04:05.000Z"
	printFormat = "2006-01-02 15:04:05"
)

const (
	delimiter = '/'
)

func parseTime(t string) string {
	ti, _ := time.Parse(recvFormat, t)
	return ti.Format(printFormat)
}

func parseLastModified(t string) string {
	ti, _ := time.Parse(time.RFC1123, t)
	return ti.Format(printFormat)
}

func printBuckets(v []*s3.Bucket) {
	for _, b := range v {
		fmt.Printf("%s   %s\n", parseTime(b.CreationDate), b.Name)
	}
}

func printObjects(v []*s3.Item) {
	if len(v) > 0 {
		sort.Sort(s3.BySize(v))
		for _, b := range v {
			fmt.Printf("%s   %d %s\n", parseTime(b.LastModified), b.Size, b.Key)
		}
	}
}

func printPrefixes(v []*s3.Prefix) {
	if len(v) > 0 {
		for _, b := range v {
			fmt.Printf("                      DIR %s\n", b.Prefix)
		}
	}
}

func printObject(v int64, date, key string) {
	fmt.Printf("%s   %d %s\n", parseLastModified(date), v, key)
}

func doFsList(c *cli.Context) {
	var items []*s3.Item
	var prefixes []*s3.Prefix

	config, err := getMcConfig()
	if err != nil {
		log.Fatal(err)
	}

	s3c, err := getNewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	fsoptions, err := parseOptions(c)
	if err != nil {
		log.Fatal(err)
	}
	switch true {
	case fsoptions.bucket == "":
		buckets, err := s3c.Buckets()
		if err != nil {
			log.Fatal(err)
		}
		printBuckets(buckets)
	case fsoptions.key == "":
		items, prefixes, err = s3c.GetBucket(fsoptions.bucket, "", "", "", s3.MaxKeys)
		if err != nil {
			log.Fatal(err)
		}
		printPrefixes(prefixes)
		printObjects(items)
	case fsoptions.key != "":
		var date string
		var size int64
		size, date, err = s3c.Stat(fsoptions.key, fsoptions.bucket)
		if err != nil {
			items, prefixes, err = s3c.GetBucket(fsoptions.bucket, "", fsoptions.key, string(delimiter), s3.MaxKeys)
			if err != nil {
				log.Fatal(err)
			}
			printPrefixes(prefixes)
			printObjects(items)
		} else {
			printObject(size, date, fsoptions.key)
		}
	}
}
