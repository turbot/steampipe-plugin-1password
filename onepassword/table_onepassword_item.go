package onepassword

import (
	"context"

	"github.com/1Password/connect-sdk-go/onepassword"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOnepasswordItem(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "onepassword_item",
		Description: "Retrieve information about your items.",
		List: &plugin.ListConfig{
			ParentHydrate: listVaults,
			Hydrate:       listItems,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "vault_id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getItem,
			KeyColumns: plugin.AllColumns([]string{"id", "vault_id"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the Item.",
			},
			{
				Name:        "vault_id",
				Type:        proto.ColumnType_STRING,
				Description: "The parent vault ID of the Item.",
				Transform:   transform.FromField("Vault.ID"),
			},
			{
				Name:        "favorite",
				Type:        proto.ColumnType_BOOL,
				Description: "Is the item favorite.",
			},
			{
				Name:        "version",
				Type:        proto.ColumnType_INT,
				Description: "The version of the item.",
			},
			{
				Name:        "category",
				Type:        proto.ColumnType_STRING,
				Description: "The category of the item.",
			},

			{
				Name:        "sections",
				Type:        proto.ColumnType_JSON,
				Description: "The category of the item.",
			},
			{
				Name:        "files",
				Type:        proto.ColumnType_JSON,
				Description: "The category of the item.",
			},

			{
				Name:        "last_edited_by",
				Type:        proto.ColumnType_STRING,
				Description: "UUID of the user that last edited the item.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Item created at.",
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Item updated at.",
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

func listItems(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vault := h.Item.(onepassword.Vault)
	vault_id := d.EqualsQuals["vault_id"].GetStringValue()

	// check if the provided vault_id is not matching with the parentHydrate
	if vault_id != "" && vault_id != vault.ID {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_Item.listItems", "connection_error", err)
		return nil, err
	}

	items, err := client.GetItems(vault.ID)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item.listItems", "api_error", err)
		return nil, err
	}

	for _, item := range items {
		d.StreamListItem(ctx, item)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getItem(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetStringValue()
	vault_id := d.EqualsQuals["vault_id"].GetStringValue()

	// Check if id is empty
	if id == "" || vault_id == "" {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item.getItem", "connection_error", err)
		return nil, err
	}

	item, err := client.GetItem(id, vault_id)
	if err != nil {
		plugin.Logger(ctx).Error("onepassword_item.getItem", "api_error", err)
		return nil, err
	}

	return item, nil
}
