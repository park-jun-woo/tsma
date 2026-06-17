//ff:func feature=gate type=test
//ff:what cleanupBacking 단위테스트: backing(gen-<item>.go)과 overlay.json 둘 다 지워지는지(존재→삭제), 그리고 둘 다 없을 때도 무음으로 통과(삭제 실패 무시)하는지 직접 호출해 귀속 검증한다.

package tsmagate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanupBacking_RemovesBackingAndOverlay(t *testing.T) {
	root := t.TempDir()
	backingRel := filepath.Join(".tsma", "test", "gen-item.go")
	overlayRel := filepath.Join(".tsma", "test", "overlay.json")
	if err := writeTestFile(root, backingRel, "package pkg\n"); err != nil {
		t.Fatal(err)
	}
	if err := writeTestFile(root, overlayRel, "{}"); err != nil {
		t.Fatal(err)
	}

	cleanupBacking(root, backingRel)

	if _, err := os.Stat(filepath.Join(root, backingRel)); !os.IsNotExist(err) {
		t.Fatalf("backing must be removed, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, overlayRel)); !os.IsNotExist(err) {
		t.Fatalf("overlay.json must be removed, stat err = %v", err)
	}
}

func TestCleanupBacking_MissingScratchIsSilent(t *testing.T) {
	root := t.TempDir()
	// Neither the backing nor the overlay exist: the Remove errors are ignored,
	// so the call is a harmless no-op (scratch, not source of truth).
	cleanupBacking(root, filepath.Join(".tsma", "test", "absent.go"))
}
