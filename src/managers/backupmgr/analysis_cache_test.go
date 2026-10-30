package backupmgr

import (
	"context"
	"sync/atomic"
	"testing"
)

func TestSaveAnalyzerCachesUnchangedSave(t *testing.T) {
	path := analysisFixture(t)
	analyzer := NewSaveAnalyzer(1, 2)
	var calls atomic.Int32
	analyzer.analyze = func(ctx context.Context, path string) (SaveAnalysis, error) {
		calls.Add(1)
		return AnalyzeSave(ctx, path)
	}

	first, err := analyzer.Analyze(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	second, err := analyzer.Analyze(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 1 {
		t.Fatalf("expected one deep scan, got %d", calls.Load())
	}
	if first != second {
		t.Fatalf("cached result changed: first=%+v second=%+v", first, second)
	}
}
