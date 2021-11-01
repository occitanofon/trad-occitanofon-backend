package activity

import (
	"testing"
	"time"
)

func TestActiveTranslators(t *testing.T) {
	activeTranslators := NewService()

	translatorIDs := []string{
		"614d86064145d514a16cf2cf",
		"615eab14659a9c407bf70c39",
		"615154c6c23db9a67b545b32",
	}

	for _, translatorID := range translatorIDs {
		activeTranslators.AddOrKeepActive(translatorID)
	}

	if activeTranslators.Total() > 1 {
		t.Logf("%d active translators\n", activeTranslators.Total())
	} else {
		t.Logf("%d active translator\n", activeTranslators.Total())
	}

	t.Log("Waiting ...")
	time.Sleep(11 * time.Minute)

	activeTranslators.AddOrKeepActive(translatorIDs[1])

	if activeTranslators.Total() > 1 {
		t.Logf("%d active translators\n", activeTranslators.Total())
	} else {
		t.Logf("%d active translator\n", activeTranslators.Total())
	}
}
