/* Copyright 2025. McKinsey & Company */

package apiserver

import (
	"context"
	"errors"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"

	"mckinsey.com/ark/internal/storage"
)

type fakeBackend struct {
	pingErr error
}

func (f *fakeBackend) Create(context.Context, string, string, string, runtime.Object) error {
	return nil
}

func (f *fakeBackend) Get(context.Context, string, string, string) (runtime.Object, error) {
	return nil, nil
}

func (f *fakeBackend) List(context.Context, string, string, storage.ListOptions) ([]runtime.Object, string, error) {
	return nil, "", nil
}

func (f *fakeBackend) Update(context.Context, string, string, string, runtime.Object) error {
	return nil
}

func (f *fakeBackend) UpdateStatus(context.Context, string, string, string, runtime.Object) error {
	return nil
}
func (f *fakeBackend) Delete(context.Context, string, string, string) error { return nil }
func (f *fakeBackend) Watch(context.Context, string, string, storage.WatchOptions) (watch.Interface, error) {
	return nil, nil
}

func (f *fakeBackend) GetResourceVersion(context.Context, string, string, string) (int64, error) {
	return 0, nil
}
func (f *fakeBackend) Ping(context.Context) error { return f.pingErr }
func (f *fakeBackend) Close() error               { return nil }

func TestServerReadyz(t *testing.T) {
	t.Run("backend not initialized", func(t *testing.T) {
		s := New(Config{})
		if err := s.Readyz(nil); err == nil {
			t.Fatal("expected error when backend is nil")
		}
	})

	t.Run("backend healthy", func(t *testing.T) {
		s := New(Config{})
		s.backend = &fakeBackend{}
		if err := s.Readyz(nil); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("backend unhealthy", func(t *testing.T) {
		s := New(Config{})
		s.backend = &fakeBackend{pingErr: errors.New("db down")}
		if err := s.Readyz(nil); err == nil {
			t.Fatal("expected error when backend ping fails")
		}
	})
}
