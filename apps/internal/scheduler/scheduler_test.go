package scheduler

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"squzy/apps/internal/config"
	"squzy/apps/internal/job"
	"testing"
	"time"
)

type storageMock struct {

}

func (s storageMock) Write(id string, log job.CheckError) error {
	return nil
}

type jb struct {
	count int
}

func (j *jb) Do() job.CheckError {
	j.count += 1
	return nil
}

func TestNew(t *testing.T) {
	var cfg config.Config
	j := &jb{count: 0}
	t.Run("Tests: Scheduler.New()", func(t *testing.T) {
		t.Run("Should: create new app without error", func(t *testing.T) {
			_, err := New(cfg, time.Second, j, &storageMock{})
			assert.Equal(t, nil, err)
		})
		t.Run("Should: create new app with 'intervalLessHalfSecondError' error", func(t *testing.T) {
			_, err := New(cfg, time.Millisecond, j, &storageMock{})
			assert.Equal(t, intervalLessHalfSecondError, err)
		})
	})
}

func TestSchl_Run(t *testing.T) {
	t.Run("Tests: Scheduler.Run()", func(t *testing.T) {
		t.Run("Should: run without error ", func(t *testing.T) {
			var cfg config.Config
			j := &jb{count: 0}
			i, _ := New(cfg, time.Second, j, &storageMock{})
			err := i.Run()
			assert.Equal(t, nil, err)
			i.Stop()
		})
		t.Run("Should: run with 'alreadyRunError' error", func(t *testing.T) {
			var cfg config.Config
			j := &jb{count: 0}
			i, _ := New(cfg, time.Second, j, &storageMock{})
			i.Run()
			err := i.Run()
			assert.Equal(t, alreadyRunError, err)
			i.Stop()
		})
		t.Run("Should: run job every second ", func(t *testing.T) {
			var cfg config.Config
			j := &jb{count: 0}
			i, _ := New(cfg, time.Second, j, &storageMock{})
			i.Run()
			time.AfterFunc(time.Second, func() {
				assert.Equal(t, 1, j.count)
			})
			time.AfterFunc(time.Second, func() {
				assert.Equal(t, 2, j.count)
			})
			i.Stop()
		})
	})
}

func TestSchl_Stop(t *testing.T) {
	t.Run("Tests: Scheduler.Stop()", func(t *testing.T) {
		t.Run("Should: stop without error ", func(t *testing.T) {
			var cfg config.Config
			j := &jb{count: 0}
			i, _ := New(cfg, time.Second, j, &storageMock{})
			_ = i.Run()
			err := i.Stop()
			assert.Equal(t, nil, err)
		})
		t.Run("Should: stop with 'alreadyStopError' error", func(t *testing.T) {
			var cfg config.Config
			j := &jb{count: 0}
			i, _ := New(cfg, time.Second, j, &storageMock{})
			i.Run()
			i.Stop()
			err := i.Stop()
			assert.Equal(t, alreadyStopError, err)
		})
	})
}

func TestSchl_IsRun(t *testing.T) {
	t.Run("Tests: Scheduler.IsRun()", func(t *testing.T) {
		t.Run("Should: return true ", func(t *testing.T) {
			var cfg config.Config
			j := &jb{count: 0}
			i, _ := New(cfg, time.Second, j, &storageMock{})
			_ = i.Run()
			assert.Equal(t, true, i.IsRun())
			_ = i.Stop()

		})
		t.Run("Should: return false", func(t *testing.T) {
			t.Run("Suite: after creation", func(t *testing.T) {
				var cfg config.Config
				j := &jb{count: 0}
				i, _ := New(cfg, time.Second, j, &storageMock{})
				assert.Equal(t, false, i.IsRun())
			})
			t.Run("Suite: after stop", func(t *testing.T) {
				var cfg config.Config
				j := &jb{count: 0}
				i, _ := New(cfg, time.Second, j, &storageMock{})
				i.Run()
				i.Stop()
				assert.Equal(t, false, i.IsRun())
			})
		})
	})
}

func TestSchl_GetId(t *testing.T) {
	t.Run("Should: return string(uuid) id", func(t *testing.T) {
		var cfg config.Config
		j := &jb{count: 0}
		i, _ := New(cfg, time.Second, j, &storageMock{})
		id:= i.GetId()
		_, err := uuid.Parse(id)
		assert.IsType(t, "", id)
		assert.Equal(t, nil, err)
	})
}
