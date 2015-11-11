package stone_swallow

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
)

type demo struct {
	c   context.Context
	w   http.ResponseWriter
	ctx context.Context
}

type Billing struct {
	AccountID     string    `json:"accountId"`
	LineItemID    string    `json:"lineItemId"`
	ProjectNumber string    `json:"projectNumber"`
	Source        string    `json:"-" datastore:",unindexd"`
	Cost          float64   `json:"cost" datastore:",unindexd"`
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
}

type BillingJson struct {
	AccountID     string        `json:"accountId"`
	Cost          BillingCost   `json:"cost"`
	EndTime       time.Time     `json:"endTime"`
	LineItemID    string        `json:"lineItemId"`
	Measurements  []Measurement `json:"measurements"`
	ProjectNumber string        `json:"projectNumber"`
	StartTime     time.Time     `json:"startTime"`
}

type BillingCost struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type Measurement struct {
	MeasurementID string `json:"measurementId"`
	Sum           string `json:"sum"`
	Unit          string `json:"unit"`
}

type BillingSum struct {
	StartTime time.Time          `json:"startTime"`
	Cost      map[string]float64 `json:"cost"`
}

func listBilling4chart(w http.ResponseWriter, r *http.Request, c context.Context) {
	bMap, err := queryBilling(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(bMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json;utf-8")
	w.WriteHeader(http.StatusOK)
}

func queryBilling(c context.Context) (map[string]BillingSum, error) {
	d := time.Now().Add(-90 * 24 * time.Hour)

	bMap := make(map[string]BillingSum)

	q := datastore.NewQuery("Billing").Filter("StartTime >= ", d).Order("-StartTime")
	t := q.Run(c)
	for {
		var bill Billing
		_, err := t.Next(&bill)
		if err == datastore.Done {
			break
		}
		if err != nil {
			log.Errorf(c, "getching nexe Billing %v", err)
			return nil, err
		}

		idArray := strings.Split(bill.LineItemID, "/")
		serviceID := idArray[2]
		if val, ok := bMap[bill.StartTime.String()]; ok {
			if costVal, ok := val.Cost[serviceID]; ok {
				val.Cost[serviceID] = costVal + bill.Cost
			} else {
				val.Cost[serviceID] = bill.Cost
			}
		} else {
			serviceMap := make(map[string]float64)
			serviceMap[serviceID] = bill.Cost
			bMap[bill.StartTime.String()] = BillingSum{
				StartTime: bill.StartTime,
				Cost:      serviceMap,
			}
		}
	}
	return bMap, nil
}

func importBilling(w http.ResponseWriter, r *http.Request, c context.Context) {
	bucket := r.FormValue("bucket")
	fileName := r.FormValue("fileName")
	if len(fileName) < 1 {
		http.Error(w, "required fileName", http.StatusBadRequest)
		return
	}

	c, cancel := context.WithDeadline(c, time.Now().Add(30*time.Second))
	defer cancel()

	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(c, storage.ScopeFullControl),
			Base:   &urlfetch.Transport{Context: c},
		},
	}
	ctx := cloud.NewContext(appengine.AppID(c), hc)

	d := &demo{
		c:   c,
		w:   w,
		ctx: ctx,
	}
	billings, err := d.readFile(bucket, fileName)
	if err != nil {
		log.Errorf(ctx, "read file error = %s", err.Error())
		return
	}

	err = putMultiBilling(c, billings)
	if err != nil {
		log.Errorf(ctx, "putMultiBilling Error. err = %s", err.Error())
		http.Error(w, "putMultiBilling Error", http.StatusInternalServerError)
		return
	}
}

func putMultiBilling(ctx context.Context, billingJsons []BillingJson) error {
	keys := make([]*datastore.Key, 0)
	billings := make([]Billing, 0)
	for _, bj := range billingJsons {
		keyName := fmt.Sprintf("%s-_-%s-_-%s-_-%s", bj.ProjectNumber, bj.LineItemID, bj.StartTime, bj.AccountID)
		log.Infof(ctx, "Key Name = %s", keyName)

		keys = append(keys, datastore.NewKey(ctx, "Billing", keyName, 0, nil))
		b, err := json.Marshal(bj)
		if err != nil {
			log.Errorf(ctx, "billing json marshal error")
			return err
		}
		cost, err := strconv.ParseFloat(bj.Cost.Amount, 64)
		if err != nil {
			log.Errorf(ctx, "cost parse float error. const = %s", bj.Cost.Amount)
			return err
		}
		billings = append(billings, Billing{
			AccountID:     bj.AccountID,
			LineItemID:    bj.LineItemID,
			ProjectNumber: bj.ProjectNumber,
			Source:        string(b),
			Cost:          cost,
			StartTime:     bj.StartTime,
			EndTime:       bj.EndTime,
		})
	}

	_, err := datastore.PutMulti(ctx, keys, billings)
	return err
}

func (d *demo) errorf(format string, args ...interface{}) {
	log.Errorf(d.c, format, args...)
}

// readFile reads the named file in Google Cloud Storage.
func (d *demo) readFile(bucket string, fileName string) ([]BillingJson, error) {
	io.WriteString(d.w, "\nAbbreviated file content (first line and last 1K):\n")

	rc, err := storage.NewReader(d.ctx, bucket, fileName)
	if err != nil {
		d.errorf("readFile: unable to open file from bucket %q, file %q: %v", bucket, fileName, err)
		return nil, err
	}
	defer rc.Close()
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		d.errorf("readFile: unable to read byte from bucket %q, file %q: %v", bucket, fileName, err)
		return nil, err
	}
	log.Infof(d.c, "body = %s", b)

	var billings []BillingJson
	err = json.Unmarshal(b, &billings)
	if err != nil {
		d.errorf("readFile: unable to decode json from bucket %q, file %q: %v", bucket, fileName, err)
		return nil, err
	}

	log.Infof(d.c, "%v", billings)

	return billings, nil
}
