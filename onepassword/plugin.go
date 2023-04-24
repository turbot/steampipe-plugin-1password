package onepassword

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-1password",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		// DefaultIgnoreConfig: &plugin.IgnoreConfig{
		// 	ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"404"}),
		// },
		TableMap: map[string]*plugin.Table{
			"onepassword_vault":      tableOnepasswordVault(ctx),
			"onepassword_item":       tableOnepasswordItem(ctx),
			"onepassword_item_file":  tableOnepasswordItemFile(ctx),
			"onepassword_item_login": tableOnepasswordItemLogin(ctx),
		},
	}
	return p
}
