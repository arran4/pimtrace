package plotoutput

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"pimtrace"
)

func BarPlot(d pimtrace.Data, fn string) error {
	p := plot.New()
	if d.Len() == 0 {
		return nil
	}
	first, ok := d.Entry(0).(pimtrace.HasStringArray)
	if !ok {
		return fmt.Errorf("plotter: invalid type")
	}
	headers := first.HeadersStringArray()
	hvm := map[string]plotter.Values{}
	var xName []string
	for i := 0; i < d.Len(); i++ {
		r := d.Entry(i)
		for hi, header := range headers {
			v, err := r.Get(header)
			if err != nil {
				return fmt.Errorf("fetching column %s from row %d: %w", header, i, err)
			}
			if hi == 0 {
				xName = append(xName, v.String())
				continue
			}
			vi := v.Float64()
			if vi != nil {
				vii := *vi
				hvm[header] = append(hvm[header], vii)
			}
		}
	}
	w := vg.Points(20)
	for i, header := range headers {
		if i == 0 {
			continue
		}
		bar, err := plotter.NewBarChart(hvm[header], w)
		if err != nil {
			return fmt.Errorf("plot error: %w", err)
		}
		bar.LineStyle.Width = vg.Length(0)
		bar.Color = plotutil.Color(i)
		bar.Offset = w * vg.Points(float64(i)-2)
		p.Add(bar)
		p.Legend.Add(header, bar)
	}
	p.NominalX(xName...)
	p.Legend.Top = true
	if err := p.Save(16*vg.Inch, 3*vg.Inch, fn); err != nil {
		return fmt.Errorf("bar plot: %w", err)
	}
	return nil
}
