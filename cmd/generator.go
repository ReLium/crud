package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/ReLium/crud/internal/mongodb"
	"github.com/ReLium/crud/internal/repository"
)

const (
	workerCount      = 10
	requestPerSecond = 300
	count            = 10000
	buffer           = 100
)

type Generator struct {
	repository repository.Repository
	speedLimit <-chan time.Time
	queue      chan repository.Cat
	count      int
}

func Generate() error {
	mongoDBClient, err := mongodb.NewClient(DefaultMongoUrl, DefaultMongoTimeoutMilliseconds)
	if err != nil {
		return err
	}
	repo := repository.NewMongoDBRepo(mongoDBClient)
	generator := &Generator{
		count:      count,
		repository: repo,
		queue:      make(chan repository.Cat, buffer),
		speedLimit: time.Tick(time.Second / requestPerSecond),
	}
	return generator.Process()
}

func (r *Generator) Process() error {
	fmt.Printf("Sending %d сats to target...\n", r.count)
	wg := sync.WaitGroup{}
	for k := 0; k < workerCount; k++ {
		wg.Add(1)
		go func() {
			r.senderWorker()
			wg.Done()
		}()
	}

	for i := 0; i < r.count; i++ {
		r.queue <- repository.Cat{
			Name:       fmt.Sprintf("Cat-%d", i),
			Gender:     "male",
			Color:      "black",
			Vaccinated: true,
		}
	}

	close(r.queue)
	wg.Wait()
	fmt.Println("Done!")
	return nil
}

func (r *Generator) senderWorker() {
	for cat := range r.queue {
		<-r.speedLimit
		r.repository.Insert(&cat)
	}
}
