package onepassword

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOnepasswordVault(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "onepassword_vault",
		Description: "ACLs for the service version.",
		List: &plugin.ListConfig{
			Hydrate: listVaults,
		},
		// Get: &plugin.GetConfig{
		// 	Hydrate:    getVault,
		// 	KeyColumns: plugin.SingleColumn("name"),
		// },
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of this vault.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of this vault.",
			},
			{
				Name:        "attr_version",
				Type:        proto.ColumnType_INT,
				Description: "The vault version.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Timestamp of when the vault was created.",
			},
			{
				Name:        "content_version",
				Type:        proto.ColumnType_INT,
				Description: "The version of the vault contents.",
				Transform:   transform.FromField("ContentVersoin"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description of this vault.",
			},
			{
				Name:        "items",
				Type:        proto.ColumnType_INT,
				Description: "Number of active items in the vault.",
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of this vault.",
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Timestamp of when the vault was updated.",
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		},
	}
}

func listVaults(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_vault.listVaults", "connection_error", err)
		return nil, err
	}

	vaults, err := client.GetVaults()
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_vault.listVaults", "query_error", err)
		return nil, err
	}

	for _, vault := range vaults {
		d.StreamListItem(ctx, vault)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
