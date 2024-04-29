package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/queue"
	"encoding/json"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

type OrderService struct {
	Db             *gorm.DB
	Domain         domain.Order
	MessageChannel chan amqp.Delivery
	RabbitMQ       *queue.RabbitMQ
}

type oderNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewOrderService(db *gorm.DB, rabbitMQ *queue.RabbitMQ, messageChannel chan amqp.Delivery) *OrderService {
	return &OrderService{
		Db:             db,
		Domain:         domain.Order{},
		MessageChannel: messageChannel,
		RabbitMQ:       rabbitMQ,
	}
}

func (s *OrderService) ProcessOrders(ch *amqp.Channel) {
	msgs, err := s.RabbitMQ.Channel.Consume(
		"order_queue", // nome da fila
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var order domain.Order
			err := json.Unmarshal(d.Body, &order)
			if err != nil {
				log.Printf("Error parsing order: %v", err)
				continue
			}

			switch order.Status {
			case "new":
				_, err = s.InsertAndNotify(&order)
			case "canceled":
				err = s.CancelAndNotify(&order)
			case "finished":
				err = s.FinalizeAndNotify(&order)
			default:
				log.Printf("Unknown order status: %s", order.Status)
				continue
			}

			if err != nil {
				log.Printf("Error processing order: %v", err)
			} else {
				log.Printf("Order processed successfully: %s", order.ID)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (s *OrderService) InsertAndNotify(order *domain.Order) (domain.Order, error) {
	OrderRepository := repositories.OrderRepositoryDb{Db: s.Db}
	insertedOrder, err := OrderRepository.Insert(order)
	if err != nil {
		log.Printf("Erro ao inserir o pedido: %v", err)
		return domain.Order{}, err
	}
	orderJson, err := json.Marshal(insertedOrder)
	if err != nil {
		log.Printf("Erro ao serializar o pedido para JSON: %v", err)
		return domain.Order{}, err // Retorna uma Order vazia junto com o erro
	}
	err = s.notify(orderJson, insertedOrder.ID)
	if err != nil {
		log.Printf("Erro ao enviar a notificação para o RabbitMQ: %v", err)
		return domain.Order{}, err
	}
	return *insertedOrder, nil
}

func (s *OrderService) FinalizeAndNotify(order *domain.Order) error {
	OrderRepository := repositories.OrderRepositoryDb{Db: s.Db}
	updatedOrder, err := OrderRepository.Update(order)
	if err != nil {
		log.Printf("Erro ao finalizar o pedido: %v", err)
		return err
	}
	orderJson, err := json.Marshal(order)
	if err != nil {
		log.Printf("Erro ao serializar o pedido para JSON: %v", err)
		return err
	}
	err = s.notify(orderJson, order.ID)
	if err != nil {
		log.Printf("Erro ao enviar a notificação para o RabbitMQ: %v", err)
		return err
	}
	log.Printf("Pedido finalizado e notificação enviada com sucesso: %v", updatedOrder.ID)

	// atualizar o estoque e sincronizar com o PDV
	return nil
}

func (s *OrderService) CancelAndNotify(order *domain.Order) error {

	OrderRepository := repositories.OrderRepositoryDb{Db: s.Db}
	insertedOrder, err := OrderRepository.Update(order)
	if err != nil {
		log.Printf("Error canceling order: %v", err)
		return err
	}

	orderJson, err := json.Marshal(order)
	if err != nil {
		log.Printf("Erro ao serializar o pedido para JSON: %v", err)
		return err
	}

	err = s.notify(orderJson, order.ID)
	if err != nil {
		log.Printf("Error sending the notification to RabbitMQ: %v", err)
		return err
	}

	log.Printf("Order canceled and notification sent successfully: %v", insertedOrder.ID)
	return nil
}

func (j *OrderService) notify(orderJson []byte, id string) error {
	err := j.RabbitMQ.Notify(
		string(orderJson),
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"),
	)
	if err != nil {
		log.Printf("Error sending notification to RabbitMQ: %v", err)
		return err
	}
	log.Printf("Order inserted and notification sent successfully: %v", id)
	return nil
}
