package context

import (
	"reflect"
	"testing"
)

func Test_getHashableTurboEnvVarsFromOs(t *testing.T) {
	env := []string{
		"SOME_ENV_VAR=excluded",
		"SOME_OTHER_ENV_VAR=excluded",
		"FIRST_THASH_ENV_VAR=first",
		"TURBO_TOKEN=never",
		"SOME_OTHER_THASH_ENV_VAR=second",
		"TURBO_TEAM=never",
	}
	gotNames, gotPairs := getHashableTurboEnvVarsFromOs(env)
	wantNames := []string{"FIRST_THASH_ENV_VAR", "SOME_OTHER_THASH_ENV_VAR"}
	wantPairs := []string{"FIRST_THASH_ENV_VAR=first", "SOME_OTHER_THASH_ENV_VAR=second"}
	if !reflect.DeepEqual(wantNames, gotNames) {
		t.Errorf("getHashableTurboEnvVarsFromOs() env names got = %v, want %v", gotNames, wantNames)
	}
	if !reflect.DeepEqual(wantPairs, gotPairs) {
		t.Errorf("getHashableTurboEnvVarsFromOs() env pairs got = %v, want %v", gotPairs, wantPairs)
	}
}

func Test_isWorkspaceReference(t *testing.T) {
	tests := []struct {
		name              string
		packageVersion    string
		dependencyVersion string
		want              bool
	}{
		{
			name:              "handles exact match",
			packageVersion:    "1.2.3",
			dependencyVersion: "1.2.3",
			want:              true,
		},
		{
			name:              "handles semver range satisfied",
			packageVersion:    "1.2.3",
			dependencyVersion: "^1.0.0",
			want:              true,
		},
		{
			name:              "handles semver range not-satisfied",
			packageVersion:    "2.3.4",
			dependencyVersion: "^1.0.0",
			want:              false,
		},
		{
			name:              "handles workspace protocol with version",
			packageVersion:    "1.2.3",
			dependencyVersion: "workspace:1.2.3",
			want:              true,
		},
		{
			name:              "handles workspace protocol with relative path",
			packageVersion:    "1.2.3",
			dependencyVersion: "workspace:../other-package/",
			want:              true,
		},
		{
			name:              "handles npm protocol with satisfied semver range",
			packageVersion:    "1.2.3",
			dependencyVersion: "npm:^1.2.3",
			want:              true, // default in yarn is to use the workspace version unless `enableTransparentWorkspaces: true`. This isn't currently being checked.
		},
		{
			name:              "handles npm protocol with non-satisfied semver range",
			packageVersion:    "2.3.4",
			dependencyVersion: "npm:^1.2.3",
			want:              false,
		},
		{
			name:              "handles pre-release versions",
			packageVersion:    "1.2.3",
			dependencyVersion: "1.2.2-alpha-1234abcd.0",
			want:              false,
		},
		{
			name:              "handles non-semver package version",
			packageVersion:    "sometag",
			dependencyVersion: "1.2.3",
			want:              true, // for backwards compatability with the code before versions were verified
		},
		{
			name:              "handles non-semver package version",
			packageVersion:    "1.2.3",
			dependencyVersion: "sometag",
			want:              true, // for backwards compatability with the code before versions were verified
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isWorkspaceReference(tt.packageVersion, tt.dependencyVersion)
			if got != tt.want {
				t.Errorf("isWorkspaceReference() got = %v, want %v", got, tt.want)
			}
		})
	}
}
