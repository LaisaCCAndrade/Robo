package main

//guest

import (
	rabbitmq "Robo/pkg/rabbitMq"
	"Robo/service"
	"Robo/skimas"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Publishs(ch *amqp.Channel, data skimas.Data) error {
	body, err := json.Marshal(data) // tranforma em JSON

	if err != nil {
		return err
	}

	err = ch.Publish(
		"amq.direct",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	Publishs(ch, skimas.Data{
		Domain:       "teste.com",
		Name:         "teste",
		Email:        "emaildeenvio@gmail.com",
		Phone:        "123456789",
		Country:      "BR",
		Organization: "teste",
		CNPJ:         "123456789",
	})

	// Inicializa Rotas
	r := mux.NewRouter()

	corsHandler := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Referrer-Policy", "no-referrer")
			next.ServeHTTP(w, r)
		})
	})

	//abrir rabbit
	ch, err = rabbitmq.OpenChannel()

	if err != nil {
		panic(err)
	}
	defer ch.Close()

	//abrir canal 
	out := make(chan amqp.Delivery)

	//comunicação rabittMq
	go rabbitmq.Consumer(ch, out, "concurrentInformation")

	//envia de 10 em 10
	for i := 1; i <= 10; i++ {
		go func() {
			for {
				msg := <- out
				
				var body skimas.Data
				err := json.Unmarshal(msg.Body, &body) 
				
				if err != nil {
					panic(err)
				}
				
				service.SendEmail(body)
		
				msg.Ack(false)
			}
		}()
	}

	log.Fatal(http.ListenAndServe(":8084", corsHandler(r)))
}
