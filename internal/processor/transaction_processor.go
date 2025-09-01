package processor

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"insider-go-backend/internal/services"
)

// TxOp: işlem tipi
type TxOp string

const (
	OpCredit   TxOp = "credit"
	OpDebit    TxOp = "debit"
	OpTransfer TxOp = "transfer"
)

// TxJob: iş olarak kuyruğa konulacak işlem
type TxJob struct {
	Op       TxOp
	UserID   int
	ToUserID int     // transfer için hedef kullanıcı
	Amount   float64 // miktar (>0)
}

// TxStats: atomik sayaçlar
type TxStats struct {
	enqueued  int64
	processed int64
	succeeded int64
	failed    int64
}

func (s *TxStats) Snapshot() (enq, proc, ok, fail int64) {
	return atomic.LoadInt64(&s.enqueued), atomic.LoadInt64(&s.processed), atomic.LoadInt64(&s.succeeded), atomic.LoadInt64(&s.failed)
}

// TransactionProcessor: worker pool + channel kuyruğu
type TransactionProcessor struct {
	jobs    chan TxJob
	workers int

	wg   sync.WaitGroup
	quit chan struct{}

	stats TxStats
}

// default processor (opsiyonel global kullanım için)
var defaultProc *TransactionProcessor

// StartDefault: varsayılan işlemciyi başlat
func StartDefault(workers, queueCapacity int) *TransactionProcessor {
	if defaultProc != nil {
		return defaultProc
	}
	p := NewTransactionProcessor(workers, queueCapacity)
	p.Start()
	defaultProc = p
	return defaultProc
}

// GetDefault: varsayılan işlemciyi döner (nil olabilir)
func GetDefault() *TransactionProcessor { return defaultProc }

// StopDefault: varsayılan işlemciyi durdurur
func StopDefault() {
	if defaultProc != nil {
		defaultProc.Stop()
		defaultProc = nil
	}
}

// NewTransactionProcessor: workers ve kuyruk kapasitesi ile oluşturur
func NewTransactionProcessor(workers, queueCapacity int) *TransactionProcessor {
	if workers <= 0 {
		workers = 1
	}
	if queueCapacity <= 0 {
		queueCapacity = 64
	}
	return &TransactionProcessor{
		jobs:    make(chan TxJob, queueCapacity),
		workers: workers,
		quit:    make(chan struct{}),
	}
}

// Start: worker'ları başlatır
func (p *TransactionProcessor) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func(id int) {
			defer p.wg.Done()
			slog.Info("txproc.worker.start", "id", id)
			for {
				select {
				case <-p.quit:
					slog.Info("txproc.worker.stop", "id", id)
					return
				case job, ok := <-p.jobs:
					if !ok {
						slog.Info("txproc.worker.queue_closed", "id", id)
						return
					}
					p.handle(job)
				}
			}
		}(i + 1)
	}
}

// Stop: worker'lara dur işareti gönderir ve bekler (kuyruğu kapatmaz)
func (p *TransactionProcessor) Stop() {
	close(p.quit)
	p.wg.Wait()
}

// CloseQueue: kuyruğu kapatır (Start edilmişse worker'lar kapanır)
func (p *TransactionProcessor) CloseQueue() { close(p.jobs) }

// Enqueue: işi kuyruğa ekler (bloklayıcı)
func (p *TransactionProcessor) Enqueue(job TxJob) error {
	if job.Amount <= 0 {
		return errors.New("amount must be > 0")
	}
	atomic.AddInt64(&p.stats.enqueued, 1)
	p.jobs <- job
	return nil
}

// TryEnqueue: kuyruğa iş eklemeyi non-blocking dener
func (p *TransactionProcessor) TryEnqueue(job TxJob) bool {
	if job.Amount <= 0 {
		return false
	}
	select {
	case p.jobs <- job:
		atomic.AddInt64(&p.stats.enqueued, 1)
		return true
	default:
		return false
	}
}

// Stats: atomik sayaçların anlık değerleri
func (p *TransactionProcessor) Stats() (enq, proc, ok, fail int64) { return p.stats.Snapshot() }

// handle: tek bir işi işler ve sayaçları günceller
func (p *TransactionProcessor) handle(job TxJob) {
	atomic.AddInt64(&p.stats.processed, 1)
	var err error
	start := time.Now()
	switch job.Op {
	case OpCredit:
		_, err = services.Credit(job.UserID, job.Amount)
	case OpDebit:
		_, err = services.Debit(job.UserID, job.Amount)
	case OpTransfer:
		_, _, err = services.Transfer(job.UserID, job.ToUserID, job.Amount)
	default:
		err = errors.New("unknown op")
	}
	if err != nil {
		slog.Error("txproc.job.failed", "op", string(job.Op), "user", job.UserID, "to", job.ToUserID, "amount", job.Amount, "err", err, "took", time.Since(start))
		atomic.AddInt64(&p.stats.failed, 1)
		return
	}
	slog.Info("txproc.job.ok", "op", string(job.Op), "user", job.UserID, "to", job.ToUserID, "amount", job.Amount, "took", time.Since(start))
	atomic.AddInt64(&p.stats.succeeded, 1)
}

// ProcessBatchConcurrently: geçici bir worker pool ile verilen işleri eşzamanlı işler ve tamamlanınca döner
func ProcessBatchConcurrently(ctx context.Context, jobs []TxJob, concurrency int) (okCount, failCount int64) {
	if concurrency <= 0 {
		concurrency = 4
	}
	jobCh := make(chan TxJob)
	var wg sync.WaitGroup
	var ok, fail int64

	// workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case j, okc := <-jobCh:
					if !okc {
						return
					}
					var err error
					switch j.Op {
					case OpCredit:
						_, err = services.Credit(j.UserID, j.Amount)
					case OpDebit:
						_, err = services.Debit(j.UserID, j.Amount)
					case OpTransfer:
						_, _, err = services.Transfer(j.UserID, j.ToUserID, j.Amount)
					default:
						err = errors.New("unknown op")
					}
					if err != nil {
						atomic.AddInt64(&fail, 1)
					} else {
						atomic.AddInt64(&ok, 1)
					}
				}
			}
		}(i + 1)
	}

	// feed jobs
	go func() {
		defer close(jobCh)
		for _, j := range jobs {
			select {
			case <-ctx.Done():
				return
			case jobCh <- j:
			}
		}
	}()

	wg.Wait()
	return atomic.LoadInt64(&ok), atomic.LoadInt64(&fail)
}
