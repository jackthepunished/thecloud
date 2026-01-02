package services

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/poyraz/cloud/internal/core/domain"
	"github.com/poyraz/cloud/internal/core/ports"
)

type LBWorker struct {
	lbRepo       ports.LBRepository
	proxyAdapter ports.LBProxyAdapter
}

func NewLBWorker(lbRepo ports.LBRepository, proxyAdapter ports.LBProxyAdapter) *LBWorker {
	return &LBWorker{
		lbRepo:       lbRepo,
		proxyAdapter: proxyAdapter,
	}
}

func (w *LBWorker) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("Load Balancer Worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Load Balancer Worker stopping")
			return
		case <-ticker.C:
			w.processCreatingLBs(ctx)
			w.processDeletingLBs(ctx)
		}
	}
}

func (w *LBWorker) processCreatingLBs(ctx context.Context) {
	lbs, err := w.lbRepo.List(ctx)
	if err != nil {
		log.Printf("Worker: failed to list LBs: %v", err)
		return
	}

	for _, lb := range lbs {
		if lb.Status == domain.LBStatusCreating {
			w.deployLB(ctx, lb)
		}
	}
}

func (w *LBWorker) processDeletingLBs(ctx context.Context) {
	lbs, err := w.lbRepo.List(ctx)
	if err != nil {
		return
	}

	for _, lb := range lbs {
		if lb.Status == domain.LBStatusDeleted {
			w.cleanupLB(ctx, lb)
		}
	}
}

func (w *LBWorker) deployLB(ctx context.Context, lb *domain.LoadBalancer) {
	log.Printf("Worker: deploying LB %s", lb.ID)

	targets, err := w.lbRepo.ListTargets(ctx, lb.ID)
	if err != nil {
		log.Printf("Worker: failed to list targets for LB %s: %v", lb.ID, err)
		return
	}

	_, err = w.proxyAdapter.DeployProxy(ctx, lb, targets)
	if err != nil {
		log.Printf("Worker: failed to deploy proxy for LB %s: %v", lb.ID, err)
		return
	}

	lb.Status = domain.LBStatusActive
	if err := w.lbRepo.Update(ctx, lb); err != nil {
		log.Printf("Worker: failed to update status for LB %s: %v", lb.ID, err)
	} else {
		log.Printf("Worker: LB %s is now ACTIVE", lb.ID)
	}
}

func (w *LBWorker) cleanupLB(ctx context.Context, lb *domain.LoadBalancer) {
	log.Printf("Worker: cleaning up LB %s", lb.ID)

	err := w.proxyAdapter.RemoveProxy(ctx, lb.ID)
	if err != nil {
		log.Printf("Worker: failed to remove proxy for LB %s: %v", lb.ID, err)
		// We might still want to delete from DB if proxy is gone or error is "not found"
	}

	if err := w.lbRepo.Delete(ctx, lb.ID); err != nil {
		log.Printf("Worker: failed to delete LB %s from DB: %v", lb.ID, err)
	} else {
		log.Printf("Worker: LB %s fully removed", lb.ID)
	}
}
