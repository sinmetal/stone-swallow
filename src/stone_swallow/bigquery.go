package stone_swallow

import (
	"appengine"
	"code.google.com/p/goauth2/appengine/serviceaccount"
	"code.google.com/p/google-api-go-client/bigquery/v2"
	"fmt"
	"net/http"
)

const (
	PROJECT_ID string = "sinpkmnms"
	DATASET_ID string = "pokemonms"
	TABLE_ID   string = "Pokemon"
)

func listBigQuery(w http.ResponseWriter, r *http.Request, c appengine.Context) {
	client, err := serviceaccount.NewClient(c, "https://www.googleapis.com/auth/bigquery")
	if err != nil {
		c.Errorf("failed to create service account client: %#v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	service, err := bigquery.New(client)
	if err != nil {
		c.Errorf("failed to create service account: %#v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := service.Jobs.Query(PROJECT_ID, &bigquery.QueryRequest{
		Kind:  "bigquery#queryRequest",
		Query: "SELECT pokemonName FROM [726962906418:" + DATASET_ID + "." + TABLE_ID + "] LIMIT 10",
	}).Do()
	if err != nil {
		c.Errorf("failed to query: %#v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, row := range response.Rows {
		for _, cell := range row.F {
			fmt.Fprintf(w, "%v\n", cell.V)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
