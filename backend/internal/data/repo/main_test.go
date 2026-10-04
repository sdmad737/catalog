package repo

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/sdmad737/catalog/backend/internal/core/services/reporting/eventbus"
	"github.com/sdmad737/catalog/backend/internal/data/ent"
	_ "github.com/sdmad737/catalog/backend/pkgs/cgofreesqlite"
	"github.com/sdmad737/catalog/backend/pkgs/faker"
)

var (
	fk   = faker.NewFaker()
	tbus = eventbus.New()

	tClient *ent.Client
	tRepos  *AllRepos
	tUser   UserOut
	tGroup  Group
)

func bootstrap() {
	var (
		err error
		ctx = context.Background()
	)

	tGroup, err = tRepos.Groups.GroupCreate(ctx, "test-group")
	if err != nil {
		log.Fatal(err)
	}

	tUser, err = tRepos.Users.Create(ctx, userFactory())
	if err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	testRoot, err := os.MkdirTemp("", "catalog-repo-*")
	if err != nil {
		log.Fatalf("failed creating test directory: %v", err)
	}

	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	go func() {
		_ = tbus.Run(context.Background())
	}()

	err = client.Schema.Create(context.Background())
	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	tClient = client
	tRepos = New(tClient, tbus, testRoot)

	bootstrap()

	code := m.Run()
	_ = client.Close()
	_ = os.RemoveAll(testRoot)
	os.Exit(code)
}
