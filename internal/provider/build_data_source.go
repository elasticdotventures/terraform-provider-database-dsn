// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &buildDataSource{}

func NewBuildDataSource() datasource.DataSource {
	return &buildDataSource{}
}

type buildDataSource struct{}

type buildModel struct {
	Driver   types.String `tfsdk:"driver"`
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
	Host     types.String `tfsdk:"host"`
	Port     types.Int64  `tfsdk:"port"`
	Name     types.String `tfsdk:"name"`
	Params   types.Map    `tfsdk:"params"`
	DSN      types.String `tfsdk:"dsn"`
	ID       types.String `tfsdk:"id"`
}

func (d *buildDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_build"
}

func (d *buildDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Builds a database DSN from component parts.",
		Attributes: map[string]schema.Attribute{
			"driver": schema.StringAttribute{
				MarkdownDescription: "Database driver (e.g., postgres, mysql, sqlserver)",
				Required:            true,
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "Database username",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Database password",
				Optional:            true,
				Sensitive:           true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Database host",
				Required:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Database port",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Required:            true,
			},
			"params": schema.MapAttribute{
				MarkdownDescription: "Additional connection parameters",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"dsn": schema.StringAttribute{
				MarkdownDescription: "The constructed DSN",
				Computed:            true,
				Sensitive:           true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "SHA1 hash of the DSN for resource identification",
				Computed:            true,
			},
		},
	}
}

func (d *buildDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data buildModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	u := &url.URL{
		Scheme: data.Driver.ValueString(),
		Host:   fmt.Sprintf("%s:%d", data.Host.ValueString(), data.Port.ValueInt64()),
	}

	if !data.User.IsNull() && !data.Password.IsNull() {
		u.User = url.UserPassword(data.User.ValueString(), data.Password.ValueString())
	} else if !data.User.IsNull() {
		u.User = url.User(data.User.ValueString())
	}

	u.Path = "/" + data.Name.ValueString()

	if !data.Params.IsNull() && !data.Params.IsUnknown() {
		var params map[string]string
		resp.Diagnostics.Append(data.Params.ElementsAs(ctx, &params, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		q := url.Values{}
		keys := make([]string, 0, len(params))
		for k := range params {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			q.Set(k, params[k])
		}
		u.RawQuery = q.Encode()
	}

	dsn := u.String()
	data.DSN = types.StringValue(dsn)
	
	sum := fmt.Sprintf("%x", sha1.Sum([]byte(dsn)))
	data.ID = types.StringValue(sum)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}