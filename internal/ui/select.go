package ui

import (
	"charm.land/huh/v2"
)

func Select(title string, options []string) (string, error) {
	opts := make([]huh.Option[string], len(options))
	for i, o := range options {
		opts[i] = huh.NewOption(o, o)
	}

	var selected string

	sel := huh.NewSelect[string]().
		Title("Pick a task").
		Options(opts...).
		Height(6).
		Filtering(true).
		Value(&selected)

	km := huh.NewDefaultKeyMap()
	km.Select.Filter.SetEnabled(false)                   // start in filter mode
	km.Select.SetFilter.SetEnabled(true)                 // so esc is live from the start
	km.Select.SetFilter.SetHelp("esc", "stop filtering") // relabel (keeps huh's state logic)
	km.Select.ClearFilter.SetHelp("esc", "clear filter") // relabel; stays disabled until a filter is set
	km.Select.Prev.Unbind()                              // single-component form: no shift+tab
	km.Select.Next.Unbind()

	form := huh.NewForm(huh.NewGroup(sel)).WithKeyMap(km)

	if err := form.Run(); err != nil {
		return "", err
	}
	return selected, nil
}
