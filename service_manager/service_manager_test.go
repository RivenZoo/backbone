package service_manager

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type counterSvc struct {
	count int
	exit  bool
}

func newCounterSvc() *counterSvc {
	return &counterSvc{}
}

func (svc *counterSvc) Run() error {
	fmt.Println("counterSvc run")
	for i := 0; i < 10000; i++ {
		if !svc.exit {
			time.Sleep(time.Millisecond * 200)
		} else {
			break
		}
		svc.count += 1
		fmt.Println("count: ", svc.count)
	}
	return nil
}

func (svc *counterSvc) Stop() error {
	fmt.Println("counterSvc stop")
	svc.exit = true
	return nil
}

func TestServiceContainer_RunServices(t *testing.T) {
	sc := NewServiceContainer()
	key := "countsvc"
	sc.RegisterCreator(key, NewServiceCreator(newCounterSvc, nil))
	sc.Init()
	go func() {
		svc := sc.GetService(key)
		assert.NotNil(t, svc)
		cntSvc := svc.(*counterSvc)
		for {
			if cntSvc.count < 20 {
				continue
			}
			sc.Close()
			break
		}
		t.Log("svc closed")
	}()
	sc.RunServices()
	t.Log("svc exited")
}
