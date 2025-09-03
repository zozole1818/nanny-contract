package reports

import (
	"fmt"
	"github.com/signintech/gopdf"
	"math"
	"strings"
	"testing"
)

func Test_calculateContribution(t *testing.T) {
	type args struct {
		base       float64
		percentage float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "jeden",
			args: args{base: 2333, percentage: 1.675},
			want: 1.68,
		},
		{
			name: "jeden",
			args: args{base: 2333, percentage: 1.674},
			want: 1.67,
		},
		{
			name: "jeden",
			args: args{base: 2333, percentage: 1.676},
			want: 1.68,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := math.Round(tt.args.percentage*100) / 100; got != tt.want {
				t.Errorf("calculateContribution() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Map[T any, U any](a []T, fn func(T) U) []U {
	r := make([]U, len(a))
	for i, aa := range a {
		r[i] = fn(aa)
	}
	return r
}

func Test_test(t *testing.T) {
	nums := []int{1, 2, 3, 4}
	squared := Map(nums, func(x int) int {
		return x * x
	})
	fmt.Println("squared:", squared) // Output: [1 4 9 16]

	// Example 2: convert int to string
	asStrings := Map(nums, func(x int) string {
		return fmt.Sprintf("Number: %d", x)
	})
	fmt.Println("asStrings:", asStrings)

	// Example 3: convert string to uppercase
	words := []string{"go", "is", "fun"}
	upper := Map(words, func(s string) string {
		return strings.ToUpper(s)
	})
	fmt.Println("upper:", upper) // Output: [GO IS FUN]
}

func TestName(t *testing.T) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	// Draw a simple grid every 50 units
	pdf.SetLineWidth(0.2)
	for x := 0; x < 595; x += 50 { // A4 width in points ~595
		pdf.Line(float64(x), 0, float64(x), 842) // height in points
	}
	for y := 0; y < 842; y += 50 {
		pdf.Line(0, float64(y), 595, float64(y))
	}

	pdf.WritePdf("C:\\cygwin64\\home\\zuzanna.tomaszewska\\return-to-work\\nanny-contract\\backend\\cmd\\grid.pdf")
}
