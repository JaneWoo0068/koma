package provider

import (
	"github.com/metafates/mangal/provider/generic"
	"github.com/metafates/mangal/provider/mangadex"
	"github.com/metafates/mangal/provider/mangapill"
	"github.com/metafates/mangal/provider/weebcentral"
	"github.com/metafates/mangal/source"
)

const CustomProviderExtension = ".lua"

var builtinProviders = []*Provider{
	{
		ID:   mangadex.ID,
		Name: mangadex.Name,
		CreateSource: func() (source.Source, error) {
			return mangadex.New(), nil
		},
	},
	{
		ID:   weebcentral.ID,
		Name: weebcentral.Name,
		CreateSource: func() (source.Source, error) {
			return weebcentral.New(), nil
		},
	},
}

func init() {
	for _, conf := range []*generic.Configuration{
		mangapill.Config,
	} {
		conf := conf
		builtinProviders = append(builtinProviders, &Provider{
			ID:   conf.ID(),
			Name: conf.Name,
			CreateSource: func() (source.Source, error) {
				return generic.New(conf), nil
			},
		})
	}
}
