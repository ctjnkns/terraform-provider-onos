package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &flowsDataSource{}
)

// NewFlowsDataSource is a helper function to simplify the provider implementation.
func NewFlowsDataSource() datasource.DataSource {
	return &flowsDataSource{}
}

type onosClient struct {
	HTTPClient *http.Client
	Host       string
	Username   string
	Password   string
}

// flowsDataSource is the data source implementation.
type flowsDataSource struct {
	client *onosClient
}

// flowsModel maps flows schema data.
type flowsDataSourceModel struct {
	Flows []flowsModel `tfsdk:"flows"`
}

type flowsModel struct {
	AppID       types.String        `tfsdk:"appId"`
	Bytes       types.Int64         `tfsdk:"bytes"`
	DeviceID    types.String        `tfsdk:"deviceId"`
	GroupID     types.Int64         `tfsdk:"groupId"`
	ID          types.String        `tfsdk:"id"`
	IsPermanent types.Bool          `tfsdk:"isPermanent"`
	LastSeen    types.Int64         `tfsdk:"lastSeen"`
	Life        types.Int64         `tfsdk:"life"`
	LiveType    types.String        `tfsdk:"liveType"`
	Packets     types.Int64         `tfsdk:"packets"`
	Priority    types.Int64         `tfsdk:"priority"`
	State       types.String        `tfsdk:"state"`
	TableID     types.Int64         `tfsdk:"tableId"`
	TableName   types.String        `tfsdk:"tableName"`
	Timeout     types.Int64         `tfsdk:"timeout"`
	Selector    flowsSelectorModel  `tfsdk:"selector"`
	Treatment   flowsTreatmentModel `tfsdk:"treatment"`
}

// flowsSelectorModel maps flow ingredients data
type flowsSelectorModel struct {
	Criteria []flowsSelectorCriteriaModel `tfsdk:"criteria"`
}

type flowsSelectorCriteriaModel struct {
	EthType types.String `tfsdk:"ethType"`
	Mac     types.String `tfsdk:"mac"`
	Port    types.Int64  `tfsdk:"port"`
	Type    types.String `tfsdk:"type"`
}

type flowsTreatmentModel struct {
	ClearDeferred bool                              `tfsdk:"clearDeferred"`
	Deferred      []flowsTreatmentInstructionsModel `tfsdk:"deferred"` //for deferred instructions
	Instructions  []flowsTreatmentInstructionsModel `tfsdk:"instructions"`
}

type flowsTreatmentInstructionsModel struct {
	Port types.String `tfsdk:"port"`
	Type types.String `tfsdk:"type"`
}

// Metadata returns the data source type name.
func (d *flowsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flows"
}

// Schema defines the schema for the data source.
func (d *flowsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"flows": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"appId": schema.StringAttribute{
							Computed: true,
						},
						"bytes": schema.Int64Attribute{
							Computed: true,
						},
						"deviceId": schema.StringAttribute{
							Computed: true,
						},
						"groupId": schema.Int64Attribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"isPermanent": schema.BoolAttribute{
							Computed: true,
						},
						"lastSeen": schema.Int64Attribute{
							Computed: true,
						},
						"life": schema.Int64Attribute{
							Computed: true,
						},
						"liveType": schema.StringAttribute{
							Computed: true,
						},
						"packets": schema.Int64Attribute{
							Computed: true,
						},
						"priority": schema.Int64Attribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
						"tableId": schema.Int64Attribute{
							Computed: true,
						},
						"tableName": schema.StringAttribute{
							Computed: true,
						},
						"timeout": schema.Int64Attribute{
							Computed: true,
						},
						"selector": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"criteria": schema.ListNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"ethType": schema.StringAttribute{
													Computed: true,
												},
												"mac": schema.StringAttribute{
													Computed: true,
												},
												"port": schema.Int64Attribute{
													Computed: true,
												},
												"type": schema.StringAttribute{
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"treatment": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"clearDeferred": schema.BoolAttribute{
										Computed: true,
									},
									"deferred": schema.ListNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"port": schema.StringAttribute{
													Computed: true,
												},
												"type": schema.StringAttribute{
													Computed: true,
												},
											},
										},
									},
									"instructions": schema.ListNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"port": schema.StringAttribute{
													Computed: true,
												},
												"type": schema.StringAttribute{
													Computed: true,
												},
											},
										},
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

func (d *flowsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*onosClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client

}

// Read refreshes the Terraform state with the latest data.
func (d *flowsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state flowsDataSourceModel

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/onos/v1/flows", d.client.Host), nil)
	if err != nil {
		fmt.Println(nil)
	}
	req.SetBasicAuth(d.client.Username, d.client.Password)

	res, err := d.client.HTTPClient.do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println("status: %d, body: %s", res.StatusCode, body)
	}
	flows := flowsDataSourceModel{}

	err = json.Unmarshal(body, &flows)
	if err != nil {
		fmt.Println(err)
	}

	var state flowsDataSourceModel

	for _, flow := range flows.Flows {
		fmt.Println(flow)
		flowState := flowsModel{
			AppID:       types.StringValue(flow.AppID),
			Bytes:       types.Int64Value(int64(flow.Bytes)),
			DeviceID:    types.StringValue(flow.DeviceID),
			GroupID:     types.Int64Value(int64(flow.GroupID)),
			ID:          types.StringValue(flow.ID),
			IsPermanent: types.BoolValue(flow.IsPermanent),
			LastSeen:    types.Int64Value(int64(flow.LastSeen)),
			Life:        types.Int64Value(int64(flow.Life)),
			LiveType:    types.StringValue(flow.LiveType),
			Packets:     types.Int64Value(int64(flow.Packets)),
			Priority:    types.Int64Value(int64(flow.Priority)),
			State:       types.StringValue(flow.State),
			TableID:     types.Int64Value(int64(flow.TableID)),
			TableName:   types.StringValue(flow.TableName),
			Timeout:     types.Int64Value(int64(flow.Timeout)),
			Selector: flowsSelectorModel{
				Criteria: []flowsSelectorCriteriaModel{},
			},
			Treatment: flowsTreatmentModel{
				ClearDeferred: types.BoolValue(flow.Treatment.ClearDeferred),
				Deferred:      []flowsTreatmentInstructionsModel{},
				Instructions:  []flowsTreatmentInstructionsModel{},
			},
		}

		for _, criteria := range flow.Selector.Criteria {
			flowState.Selector.Criteria = append(flowState.Selector.Criteria, flowsSelectorCriteriaModel{
				EthType: types.StringValue(criteria.EthType),
				Mac:     types.StringValue(criteria.Mac),
				Port:    types.Int64Value(int64(criteria.Port)),
				Type:    types.StringValue(criteria.Type),
			})
		}

		for _, instruction := range flow.Treatment.Instructions {
			flowState.Treatment.Instructions = append(flowState.Treatment.Instructions, flowsTreatmentInstructionsModel{
				Port: types.StringValue(instruction.Port),
				Type: types.StringValue(instruction.Type),
			})
		}
		for _, deferred := range flow.Treatment.Deferred {
			flowState.Treatment.Deferred = append(flowState.Treatment.Deferred, flowsTreatmentInstructionsModel{
				Port: types.StringValue(deferred.Port),
				Type: types.StringValue(deferred.Type),
			})
		}
		//fmt.Println("This state:", state)
		state.Flows = append(state.Flows, flowState)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

/*

	"treatment": schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"clearDeferred": schema.BoolAttribute{
				Computed: true,
			},
			"deferred": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"port": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"instructions": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"port": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	},
*/
