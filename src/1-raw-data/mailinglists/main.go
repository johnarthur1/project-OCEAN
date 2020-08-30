// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
This package is for loading different mailing list data types into Cloud Storage.
*/

package main

import (
	"1-raw-data/gcs"
	"1-raw-data/mailinglists/mailman"
	"1-raw-data/mailinglists/pipermail"
	"context"
	"flag"
	"log"
)

var (
	projectID      = flag.String("project-id", "", "project id")
	bucketName     = flag.String("bucket-name", "", "bucket name to store files")
	mailingList    = flag.String("mailinglist", "piper", "Choose which mailing list to process either piper (default), mailman")
	mailingListURL = flag.String("mailinglist-url", "", "mailing list url to pull files from")
	startDate      = flag.String("start-date", "", "Start date in format of year-month-date and 4dig-2dig-2dig")
	endDate        = flag.String("end-date", "", "End date in format of year-month-date and 4dig-2dig-2dig")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gcs := gcs.StorageConnection{
		BucketName: *bucketName,
		ProjectID:  *projectID,
	}
	gcs.Ctx = ctx

	if err := gcs.ConnectGCSClient(); err != nil {
		log.Fatalf("Connect GCS failes: %v", err)
	}

	if err := gcs.CreateGCSBucket(); err != nil {
		log.Fatalf("Create GCS Bucket failed: %v", err)
	}

	switch *mailingList {
	case "piper":
		if err := pipermail.GetMailingListData(gcs, *mailingListURL); err != nil {
			log.Fatalf("Mailman load failed: %v", err)
		}
	case "mailman":
		if err := mailman.GetMailmanData(gcs, *mailingListURL, *startDate, *endDate); err != nil {
			log.Fatalf("Mailman load failed: %v", err)
		}
	default:
		log.Fatalf("Mailing list %v is not an option. Change the option submitted.: ", mailingList)
	}
}
