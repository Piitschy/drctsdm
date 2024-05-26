package directus_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Piitschy/drctsdm/internal/directus"
	"github.com/testcontainers/testcontainers-go"
)

func NewDirectusContainerWithCollection(t *testing.T, version string, collection *directus.Collection) (context.Context, testcontainers.Container, *directus.Directus) {
	ctx, container, d := NewDirectusContainer(t, version)

	err := d.Login(adminEmail, adminPassword)
	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("Failed to login: %s", err)
	}

	err = d.CreateCollection(collection)
	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("Failed to create collection: %s", err)
	}

	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("Failed to create collection: %s", err)
	}
	return ctx, container, d
}

func LoadTestCollection(t *testing.T, name string) *directus.Collection {
	f, err := os.Open(filepath.Join("..", "..", "test", "testdata", name))
	if err != nil {
		t.Fatalf("Failed to open file: %s", err)
	}
	defer f.Close()
	bytes, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("Failed to read file: %s", err)
	}

	c, err := directus.UnmarshalCollection(bytes)
	if err != nil {
		t.Fatalf("Failed to unmarshal collection: %s", err)
	}
	return c
}

func TestCreateCollection(t *testing.T) {
	ctx, container, _ := NewDirectusContainerWithCollection(t, "latest", LoadTestCollection(t, "article.json"))
	defer container.Terminate(ctx)
}

func TestGetCollection(t *testing.T) {
	ctx, container, d := NewDirectusContainerWithCollection(t, "latest", LoadTestCollection(t, "article.json"))
	defer container.Terminate(ctx)

	err := d.Login(adminEmail, adminPassword)
	if err != nil {
		t.Fatalf("Failed to login: %s", err)
	}

	// fmt.Println(d.Url)
	// time.Sleep(5 * time.Minute)
	c, err := d.GetCollection("article")
	if err != nil {
		t.Fatalf("Failed to get collection: %s", err)
	}

	if c.Collection != "article" {
		t.Fatalf("Collection name should be 'article'")
	}
}