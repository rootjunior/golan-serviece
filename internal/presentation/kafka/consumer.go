package kafka

import (
	"context"
	"errors"
	"go-service/internal/config"
	"go-service/internal/presentation/kafka/v1"
	"io"
	"log"
	"sync"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type Topic struct {
	Config  kafkago.ReaderConfig
	Handler v1.Handler
}

type Server struct {
	cfg config.Config
	v1  *v1.Controller

	topics []*kafkago.Reader
	mu     sync.Mutex
	once   sync.Once
}

func NewServer(cfg *config.Config, v1 *v1.Controller) *Server {
	return &Server{
		cfg: *cfg,
		v1:  v1,
	}
}

func (s *Server) getTopics() []Topic {
	return []Topic{
		{
			Config: kafkago.ReaderConfig{
				Brokers: []string{"localhost:9092"},
				Topic:   "create-user",
				GroupID: "go-service-create-user",

				StartOffset: kafkago.FirstOffset,

				MinBytes: 1,
				MaxBytes: 10e6,
			},
			Handler: s.v1.CreateUser,
		},
	}
}

func (s *Server) addTopic(topic *kafkago.Reader) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.topics = append(s.topics, topic)
}

func (s *Server) run(ctx context.Context) error {
	readers := s.getTopics()
	errCh := make(chan error, len(readers))
	for _, topicCfg := range readers {
		topicCfg := topicCfg
		topic := kafkago.NewReader(topicCfg.Config)
		s.addTopic(topic)
		go func() {
			errCh <- s.consume(ctx, topic, topicCfg)
		}()
	}

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *Server) consume(
	ctx context.Context,
	topic *kafkago.Reader,
	readerCfg Topic,
) error {
	for {
		msg, err := topic.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, io.EOF) {
				return nil
			}
			log.Printf(
				"failed to fetch kafka message, topic=%s, error=%v\n",
				readerCfg.Config.Topic,
				err,
			)
			time.Sleep(time.Second)
			continue
		}

		if err := readerCfg.Handler(ctx, msg); err != nil {
			log.Printf(
				"failed to handle kafka message, topic=%s, partition=%d, offset=%d, error=%v\n",
				msg.Topic,
				msg.Partition,
				msg.Offset,
				err,
			)
			continue
		}

		if err := topic.CommitMessages(ctx, msg); err != nil {
			log.Printf(
				"failed to commit kafka message, topic=%s, partition=%d, offset=%d, error=%v\n",
				msg.Topic,
				msg.Partition,
				msg.Offset,
				err,
			)
			continue
		}
	}
}

func (s *Server) Start(ctx context.Context) {
	go func() {
		go func() {
			<-ctx.Done()
			log.Printf("Done kafka context")

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := s.Stop(shutdownCtx); err != nil {
				log.Printf("kafka shutdown error: %v\n", err)
			}
		}()

		if err := s.run(ctx); err != nil {
			log.Printf("kafka server error: %v\n", err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	var resultErr error

	done := make(chan struct{})

	go func() {
		s.once.Do(func() {
			s.mu.Lock()
			topics := make([]*kafkago.Reader, 0, len(s.topics))
			topics = append(topics, s.topics...)
			s.mu.Unlock()

			for _, topic := range topics {
				if err := topic.Close(); err != nil {
					resultErr = errors.Join(resultErr, err)
				}
			}
		})

		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return resultErr
	}
}
