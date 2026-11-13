package reports

import (
	"strings"
	"testing"
)

func TestBarcodeSVGCode128(t *testing.T) {
	svg := BarcodeSVG("R42", 2, 56)
	if !strings.HasPrefix(svg, "<svg") || !strings.HasSuffix(svg, "</svg>") {
		t.Fatalf("expected svg document, got %q", svg[:40])
	}
	if !strings.Contains(svg, `fill="#000000"`) {
		t.Fatal("expected black bars")
	}
	if !strings.Contains(svg, `aria-label="Штрих-код R42"`) {
		t.Fatal("expected accessible label")
	}
	if svg != BarcodeSVG("R42", 2, 56) {
		t.Fatal("barcode rendering must be deterministic")
	}
}

func TestBarcodeSVGFallbacks(t *testing.T) {
	for _, value := range []string{"", "кириллица", "R\u0001"} {
		svg := BarcodeSVG(value, 0, 0)
		if !strings.Contains(svg, "<svg") {
			t.Fatalf("value %q produced no svg", value)
		}
	}
}

func TestTalonVariantValid(t *testing.T) {
	for _, variant := range []string{"base", "vzrm", "tula-prihod", "tula-prodazha", "tula-second", "tula-defect"} {
		if !TalonVariantValid(variant) {
			t.Fatalf("variant %q must be valid", variant)
		}
	}
	if TalonVariantValid("unknown") {
		t.Fatal("unknown variant must be rejected")
	}
}
