package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type response struct {
	Status  bool   `json:"status"`
	Code    int    `Json:"code"`
	Message string `Json:"message"`
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	responseReturnApi := response{Status: true, Code: 200, Message: "Serviço de API ligado"}

	json.NewEncoder(w).Encode(responseReturnApi)
}

func getMusica(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode("colecao")

}

//=======================================================

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb+srv://root:BobMarley1981@cluster0.sgfic.mongodb.net/")
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	log.Print("	=> BD Connected")

	// Selecionar a coleção
	collection := client.Database("jPLay").Collection("musicas")

	// Definir o filtro de consulta
	filter := bson.M{"Artista": "Led Zeppelin"} // Por exemplo, buscar por artista "Led Zeppelin"

	// Realizar a consulta
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	// Iterar pelos resultados
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Resultado: %v\n", result)
	}

	// Verificar se houve erro durante a iteração
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	//Router

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getStatus(w, r)
		}
	})

	http.HandleFunc("/song", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			getMusica(w, r)
		}
	})

	log.Print("- Server ON -")
	http.ListenAndServe(":6000", nil)

}
