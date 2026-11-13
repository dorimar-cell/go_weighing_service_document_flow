package service

import "testing"

func TestSettingsDefaultsKnownScopes(t *testing.T) {
	for _, scope := range SettingsProfileScopes() {
		values, ok := settingsDefaults(scope)
		if !ok {
			t.Fatalf("scope %q is listed but has no defaults", scope)
		}
		if len(values) == 0 {
			t.Fatalf("scope %q has empty defaults", scope)
		}
	}
	if _, ok := settingsDefaults("unknown"); ok {
		t.Fatal("unknown scope must not have defaults")
	}
}

func TestMaskSecretsKeepsPasswordOutOfPortal(t *testing.T) {
	values := map[string]any{"userMail": "a@b.c", "userPassword": "secret"}
	masked := maskSecrets("email", values)
	if masked["userPassword"] != secretMask {
		t.Fatalf("password must be masked, got %v", masked["userPassword"])
	}
	if values["userPassword"] != "secret" {
		t.Fatal("maskSecrets must not mutate the source map")
	}

	restored := restoreSecrets("email", masked, values)
	if restored["userPassword"] != "secret" {
		t.Fatalf("masked password must be restored for save, got %v", restored["userPassword"])
	}

	changed := restoreSecrets("email", map[string]any{"userPassword": "new-pass"}, values)
	if changed["userPassword"] != "new-pass" {
		t.Fatalf("new password must be kept, got %v", changed["userPassword"])
	}
}

func TestValidateProfileEmailRequiresSmtp(t *testing.T) {
	if err := validateProfile("email", map[string]any{"enable": true}); err == nil {
		t.Fatal("enabled e-mail without SMTP host must be rejected")
	}
	if err := validateProfile("email", map[string]any{
		"enable": true, "smtpHost": "smtp.local", "userMail": "a@b.c",
	}); err != nil {
		t.Fatalf("valid e-mail settings rejected: %v", err)
	}
}

func TestValidateProfileRanges(t *testing.T) {
	cases := []map[string]any{
		{"port": float64(70000)},
		{"smtpPort": float64(-1)},
		{"volume": float64(101)},
	}
	for _, values := range cases {
		if err := validateProfile("email", values); err == nil {
			t.Fatalf("values %v must be rejected", values)
		}
	}
}
