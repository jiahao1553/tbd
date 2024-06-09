package internal

import (
	"errors"
	"sync"

	"github.com/gwenwindflower/tbd/shared"
)

func WriteFiles(ts shared.SourceTables, bd string, prefix string, database *string, schema *string) error {
	if len(ts.SourceTables) == 0 {
		return errors.New("no tables to write")
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		WriteYAML(ts, bd, database, schema)
	}()
	go func() {
		defer wg.Done()
		WriteStagingModels(ts, bd, prefix)
	}()
	wg.Wait()
	return nil
}
