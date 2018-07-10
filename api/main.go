package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/datastore"
	"github.com/go-chi/chi"
	"golang.org/x/oauth2/google"
)

func main() {
	// Use oauth2.NoContext if there isn't a good context to pass in.
	ctx := context.Background()

	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/datastore")
	if err != nil {
		log.Fatal().Err(err).Msg("oops")
	}
	ds, err := datastore.NewClient(ctx, creds.ProjectID)
	// client := oauth2.NewClient(ctx, creds.TokenSource)
	if err != nil {
		log.Fatal().Err(err).Msg("oops")
	}
	// datastoreService, err := datastore.New(client)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("oops")
	// }

	// THIS ERRORS OUT
	// firestoreService, err := firestore.New(client)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("oops")
	// }
	// resp, err := firestoreService.Projects.Databases.Documents.RunQuery(fmt.Sprintf("projects/%v/databases/%v/documents", creds.ProjectID, "(default)"), &firestore.RunQueryRequest{
	// 	StructuredQuery: &firestore.StructuredQuery{
	// 		From: []*firestore.CollectionSelector{&firestore.CollectionSelector{
	// 			CollectionId: "items",
	// 		}},
	// 		Where: &firestore.Filter{
	// 			FieldFilter: &firestore.FieldFilter{
	// 				Field: &firestore.FieldReference{
	// 					FieldPath: "type",
	// 				},
	// 				Op: "EQUAL",
	// 				Value: &firestore.Value{
	// 					StringValue: "foo",
	// 				},
	// 			},
	// 		},
	// 	},
	// }).Do()
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("oops2")
	// }
	// log.Info().Interface("resp", resp).Msg("response")

	query := datastore.NewQuery("items")
	// Filter("type =", "foo")
	// Filter("Priority >=", 4).
	// Order("-Priority")

	it := ds.Run(ctx, query)
	for {
		var task Item
		_, err := it.Next(&task)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal().Err(err).Msg("Error fetching next task")
		}
		fmt.Printf("Task %v, Priority %v\n", task.Name, task.Typ)
	}

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":3000", r)
}

type Item struct {
	Name string
	Typ  string `json:"type" datastore:"type"`
	Bnum int    `datastore:"bnum"`
}
