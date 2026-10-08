package modrinth

import (
	modrinthApi "codeberg.org/theepicblock/go-modrinth/modrinth"
	"github.com/mannie-exe/packwiz-tx/core"
	"strings"
	"testing"
)

func TestVersionEnvironments(t *testing.T) {
	cases := map[string]string{
		"client_only": core.ClientSide, "server_only": core.UniversalSide,
		"dedicated_server_only": core.ServerSide, "client_and_server": core.UniversalSide,
		"server_only_client_optional":   core.UniversalSide,
		"client_only_server_optional":   core.UniversalSide,
		"client_or_server_prefers_both": core.UniversalSide,
		"client_or_server":              "either", "singleplayer_only": core.ClientSide,
		"future_value": "",
	}
	for env, want := range cases {
		t.Run(env, func(t *testing.T) {
			if got := getSide(&modrinthApi.Version{Environment: &env}); got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
		})
	}
	if got := getSide(&modrinthApi.Version{}); got != "" {
		t.Fatalf("missing environment = %q", got)
	}
}
func TestMissingEnvironmentDoesNotPanicBeforeFileValidation(t *testing.T) {
	project, version := "Project1", "Version1"
	err := createFileMeta(&modrinthApi.Project{ID: &project, Slug: &project}, &modrinthApi.Version{ID: &version}, &modrinthApi.File{}, core.Pack{}, nil)
	if err == nil || !strings.Contains(err.Error(), "hash") {
		t.Fatalf("expected missing hash, got %v", err)
	}
}
