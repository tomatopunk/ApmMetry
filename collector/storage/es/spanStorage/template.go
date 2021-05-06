package spanStorage

import (
	"context"
	"strings"
)

func (w esSpanWriter) CreateTemplate(ctx context.Context, indexName string, template string) error {
	err := w.client.PutTemplate(ctx, indexName, strings.NewReader(template))
	return err
}
