package provider

import (
	"context"
	"fmt"

	"github.com/ctjnkns/onosclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &hostsDataSource{}
	_ datasource.DataSourceWithConfigure = &hostsDataSource{}
)

// NewHostsDataSource is a helper function to simplify the provider implementation.
func NewHostsDataSource() datasource.DataSource {
	return &hostsDataSource{}
}

// hostsDataSource is the data source implementation.
type hostsDataSource struct {
	client *onosclient.Client
}

type hostsDataSourceModel struct {
	Hosts []hostsModel `tfsdk:"hosts"`
}

type hostsModel struct {
	ID          types.String          `tfsdk:"id"`
	Mac         types.String          `tfsdk:"mac"`
	Vlan        types.String          `tfsdk:"vlan"`
	InnerVlan   types.String          `tfsdk:"innervlan"`
	OuterTpid   types.String          `tfsdk:"outertpid"`
	Configured  types.Bool            `tfsdk:"configured"`
	Suspended   types.Bool            `tfsdk:"suspended"`
	IPAddresses types.List            `tfsdk:"ipaddresses"`
	Locations   []hostsLocationsModel `tfsdk:"locations"`
}

type hostsLocationsModel struct {
	ElementID types.String `tfsdk:"elementid"`
	Port      types.String `tfsdk:"port"`
}

// Metadata returns the data source type name.
func (d *hostsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hosts"
}

// Schema defines the schema for the data source.
func (d *hostsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"hosts": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"mac": schema.StringAttribute{
							Computed: true,
						},
						"vlan": schema.StringAttribute{
							Computed: true,
						},
						"innervlan": schema.StringAttribute{
							Computed: true,
						},
						"outertpid": schema.StringAttribute{
							Computed: true,
						},
						"configured": schema.BoolAttribute{
							Computed: true,
						},
						"suspended": schema.BoolAttribute{
							Computed: true,
						},
						"ipaddresses": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"locations": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"elementid": schema.StringAttribute{
										Computed: true,
									},
									"port": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *hostsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*onosclient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *onos.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client

}

// Read refreshes the Terraform state with the latest data.
func (d *hostsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state hostsDataSourceModel
	hosts, err := d.client.GetHosts()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Onos Hosts",
			err.Error(),
		)
		return
	}

	for _, host := range hosts.Hosts {

		//https://github.com/buildkite/terraform-provider-buildkite/blob/2b6b85c94482a0b5ff44beb4b24296016ac31659/buildkite/data_source_meta.go#L57
		ips, diag := types.ListValueFrom(ctx, types.StringType, host.IPAddresses)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}

		hostState := hostsModel{
			ID:          types.StringValue(host.ID),
			Mac:         types.StringValue(host.Mac),
			Vlan:        types.StringValue(host.Vlan),
			InnerVlan:   types.StringValue(host.InnerVlan),
			OuterTpid:   types.StringValue(host.OuterTpid),
			Configured:  types.BoolValue(host.Configured),
			Suspended:   types.BoolValue(host.Suspended),
			IPAddresses: ips,
			Locations:   []hostsLocationsModel{},
		}

		for _, location := range host.Locations {
			hostState.Locations = append(hostState.Locations, hostsLocationsModel{
				ElementID: types.StringValue(location.ElementID),
				Port:      types.StringValue(location.Port),
			})
		}
		state.Hosts = append(state.Hosts, hostState)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
