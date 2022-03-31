package cmd

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ReLium/crud/internal/mongodb"
	"github.com/ReLium/crud/internal/repository"
	"github.com/cheggaaa/pb/v3"
)

const (
	workerCount      = 10
	requestPerSecond = 100
	count            = 1000
	buffer           = 10
)

type Generator struct {
	repository  repository.Repository
	speedLimit  <-chan time.Time
	queue       chan repository.Cat
	progressBar *pb.ProgressBar
	count       int
}

func Generate() error {
	host, timeout := getMongoDBSettings()
	mongoDBClient, err := mongodb.NewClient(host, timeout)
	if err != nil {
		return err
	}
	repo := repository.NewMongoDBRepo(mongoDBClient)
	generator := &Generator{
		count:       count,
		repository:  repo,
		queue:       make(chan repository.Cat, buffer),
		speedLimit:  time.Tick(time.Second / requestPerSecond),
		progressBar: pb.New(count),
	}
	return generator.Process()
}

func (r *Generator) Process() error {
	fmt.Printf("Sending %d сats to db...\n", r.count)
	wg := sync.WaitGroup{}
	for k := 0; k < workerCount; k++ {
		wg.Add(1)
		go func() {
			r.senderWorker()
			wg.Done()
		}()
	}

	r.progressBar.Start()
	for i := 0; i < r.count; i++ {
		r.queue <- *r.buildCat(i)
	}

	close(r.queue)
	wg.Wait()
	r.progressBar.Finish()
	return nil
}

func (r *Generator) senderWorker() {
	for cat := range r.queue {
		<-r.speedLimit
		r.repository.Insert(&cat)
		r.progressBar.Increment()
	}
}

func (r *Generator) buildCat(idx int) *repository.Cat {
	colors := []string{"white", "black", "silver", "gray", "ginger"}
	genders := []string{"male", "female"}
	names := []string{"Bublina", "Leontynka", "Merlin", "Micka", "Casidy"}

	return &repository.Cat{
		Name:       fmt.Sprintf("%s-%d", names[rand.Intn(len(names))], idx),
		Color:      colors[rand.Intn(len(colors))],
		Gender:     genders[rand.Intn(len(genders))],
		Vaccinated: rand.Float64() > 0.5,
	}
}
