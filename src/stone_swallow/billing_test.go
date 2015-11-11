package stone_swallow

import (
	"github.com/favclip/testerator"

	"google.golang.org/appengine/datastore"

	"encoding/json"
	"testing"
)

func TestPutBilling(t *testing.T) {
	_, c, err := testerator.SpinUp()
	defer testerator.SpinDown()

	var b Billing
	key := datastore.NewKey(c, "Billing", "Hoge", 0, nil)
	_, err = datastore.Put(c, key, &b)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMultiPutBilling(t *testing.T) {
	_, c, _ := testerator.SpinUp()
	defer testerator.SpinDown()

	data :=
		`[
{
	"accountId": "001298-BBDE11-A32617",
	"lineItemId": "com.google.cloud/services/app-engine/BackendInstances",
	"startTime": "2015-11-07T00:00:00-08:00",
	"endTime": "2015-11-08T00:00:00-08:00",
	"projectNumber": "481141327602",
	"measurements": [
		{
			"measurementId": "com.google.cloud/services/app-engine/BackendInstances",
			"sum": "0",
			"unit": "seconds"
		}
	],
	"cost": {
		"amount": "0",
		"currency": "USD"
	}
}]`

	var billings []BillingJson
	err := json.Unmarshal([]byte(data), &billings)
	if err != nil {
		t.Fatal(err)
	}

	err = putMultiBilling(c, billings)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecodeJson(t *testing.T) {
	data :=
		`[
    {
        "accountId": "001298-BBDE11-A32617",
        "lineItemId": "com.google.cloud/services/app-engine/BackendInstances",
        "startTime": "2015-11-07T00:00:00-08:00",
        "endTime": "2015-11-08T00:00:00-08:00",
        "projectNumber": "481141327602",
        "measurements": [
            {
                "measurementId": "com.google.cloud/services/app-engine/BackendInstances",
                "sum": "0",
                "unit": "seconds"
            }
        ],
        "cost": {
            "amount": "0",
            "currency": "USD"
        }
    }]`

	var billings []BillingJson
	err := json.Unmarshal([]byte(data), &billings)
	if err != nil {
		t.Fatal(err)
	}

	if len(billings) != 1 {
		t.Errorf("billing length != 1. len = %d", len(billings))
	}
	if billings[0].AccountID != "001298-BBDE11-A32617" {
		t.Errorf("unexpected billing[0].AccountID. valu = %s", billings[0].AccountID)
	}
}

func TestQueryBilling(t *testing.T) {
	_, c, _ := testerator.SpinUp()
	defer testerator.SpinDown()

	data :=
		`[
{
"accountId": "001298-BBDE11-A32617",
"lineItemId": "com.google.cloud/services/app-engine/BackendInstances",
"startTime": "2015-11-07T00:00:00-08:00",
"endTime": "2015-11-08T00:00:00-08:00",
"projectNumber": "481141327602",
"measurements": [
	{
		"measurementId": "com.google.cloud/services/app-engine/BackendInstances",
		"sum": "0",
		"unit": "seconds"
	}
],
"cost": {
	"amount": "0",
	"currency": "USD"
}
}]`

	var billings []BillingJson
	err := json.Unmarshal([]byte(data), &billings)
	if err != nil {
		t.Fatal(err)
	}

	err = putMultiBilling(c, billings)
	if err != nil {
		t.Fatal(err)
	}

	res, err := queryBilling(c)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", res)
}
