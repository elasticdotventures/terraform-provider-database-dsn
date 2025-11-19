// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/xo/dburl"
)

var _ datasource.DataSource = &parseDataSource{}

func NewParseDataSource() datasource.DataSource {
	return &parseDataSource{}
}

type parseDataSource struct{}

type parseModel struct {
	DSN      types.String `tfsdk:"dsn"`
	Driver   types.String `tfsdk:"driver"`
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
	Host     types.String `tfsdk:"host"`
	Port     types.Int64  `tfsdk:"port"`
	Name     types.String `tfsdk:"name"`
	Params   types.Map    `tfsdk:"params"`
	ID       types.String `tfsdk:"id"`
}

func (d *parseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_parse"
}

func (d *parseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Parses a database DSN into component parts.",
		Attributes: map[string]schema.Attribute{
			"dsn": schema.StringAttribute{
				MarkdownDescription: "Database DSN to parse",
				Required:            true,
				Sensitive:           true,
			},
			"driver": schema.StringAttribute{
				MarkdownDescription: "Database driver",
				Computed:            true,
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "Database username",
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Database password",
				Computed:            true,
				Sensitive:           true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Database host",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Database port",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Computed:            true,
			},
			"params": schema.MapAttribute{
				MarkdownDescription: "Additional connection parameters",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "SHA1 hash of the DSN for resource identification",
				Computed:            true,
			},
		},
	}
}

func (d *parseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data parseModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	u, err := dburl.Parse(data.DSN.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"DSN Parse Error",
			fmt.Sprintf("Unable to parse DSN: %s", err),
		)
		return
	}

	data.Driver = types.StringValue(u.Scheme)
	
	if u.User != nil {
		data.User = types.StringValue(u.User.Username())
		if password, ok := u.User.Password(); ok {
			data.Password = types.StringValue(password)
		} else {
			data.Password = types.StringNull()
		}
	} else {
		data.User = types.StringNull()
		data.Password = types.StringNull()
	}

	host := u.Hostname()
	if host == "" {
		data.Host = types.StringNull()
	} else {
		data.Host = types.StringValue(host)
	}

	portStr := u.Port()
	if portStr == "" {
		data.Port = types.Int64Null()
	} else {
		port, err := strconv.ParseInt(portStr, 10, 64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Port Parse Error",
				fmt.Sprintf("Unable to parse port '%s': %s", portStr, err),
			)
			return
		}
		data.Port = types.Int64Value(port)
	}

	path := strings.TrimPrefix(u.Path, "/")
	if path == "" {
		data.Name = types.StringNull()
	} else {
		data.Name = types.StringValue(path)
	}

	if len(u.Query()) > 0 {
		params := make(map[string]string)
		for k, v := range u.Query() {
			if len(v) > 0 {
				params[k] = v[0]
			} else {
				params[k] = ""
			}
		}
		paramsValue, diags := types.MapValueFrom(ctx, types.StringType, params)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Params = paramsValue
	} else {
		data.Params = types.MapNull(types.StringType)
	}

	sum := fmt.Sprintf("%x", sha1.Sum([]byte(data.DSN.ValueString())))
	data.ID = types.StringValue(sum)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}