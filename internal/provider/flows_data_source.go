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

const HostURL string = "http://localhost:8181/onos/v1"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &flowsDataSource{}
	_ datasource.DataSourceWithConfigure = &flowsDataSource{}
)

// NewFlowsDataSource is a helper function to simplify the provider implementation.
func NewFlowsDataSource() datasource.DataSource {
	return &flowsDataSource{}
}

// flowsDataSource is the data source implementation.
type flowsDataSource struct {
	client *http.Client
}

type flowsDataSourceUnmarshal struct {
	Flows []flowsUnmarshal `tfsdk:"flows"`
}

type flowsUnmarshal struct {
	AppID       string                  `tfsdk:"appid"`
	Bytes       int                     `tfsdk:"bytes"`
	DeviceID    string                  `tfsdk:"deviceid"`
	GroupID     int                     `tfsdk:"groupid"`
	ID          string                  `tfsdk:"id"`
	IsPermanent bool                    `tfsdk:"ispermanent"`
	LastSeen    int                     `tfsdk:"lastseen"`
	Life        int                     `tfsdk:"life"`
	LiveType    string                  `tfsdk:"livetype"`
	Packets     int                     `tfsdk:"packets"`
	Priority    int                     `tfsdk:"priority"`
	State       string                  `tfsdk:"state"`
	TableID     int                     `tfsdk:"tableid"`
	TableName   string                  `tfsdk:"tablename"`
	Timeout     int                     `tfsdk:"timeout"`
	Selector    flowsSelectorUnmarshal  `tfsdk:"selector"`
	Treatment   flowsTreatmentUnmarshal `tfsdk:"treatment"`
}

// flowsSelectorUnmarshal maps flow ingredients data
type flowsSelectorUnmarshal struct {
	Criteria []flowsSelectorCriteriaUnmarshal `tfsdk:"criteria"`
}

type flowsSelectorCriteriaUnmarshal struct {
	EthType string `tfsdk:"ethtype"`
	Mac     string `tfsdk:"mac"`
	Port    int    `tfsdk:"port"`
	Type    string `tfsdk:"type"`
}

type flowsTreatmentUnmarshal struct {
	ClearDeferred bool                                  `tfsdk:"cleardeferred"`
	Deferred      []flowsTreatmentInstructionsUnmarshal `tfsdk:"deferred"` //for deferred instructions
	Instructions  []flowsTreatmentInstructionsUnmarshal `tfsdk:"instructions"`
}

type flowsTreatmentInstructionsUnmarshal struct {
	Port string `tfsdk:"port"`
	Type string `tfsdk:"type"`
}

type flowsDataSourceModel struct {
	Flows []flowsModel `tfsdk:"flows"`
}

type flowsModel struct {
	AppID       types.String        `tfsdk:"appid"`
	Bytes       types.Int64         `tfsdk:"bytes"`
	DeviceID    types.String        `tfsdk:"deviceid"`
	GroupID     types.Int64         `tfsdk:"groupid"`
	ID          types.String        `tfsdk:"id"`
	IsPermanent types.Bool          `tfsdk:"ispermanent"`
	LastSeen    types.Int64         `tfsdk:"lastseen"`
	Life        types.Int64         `tfsdk:"life"`
	LiveType    types.String        `tfsdk:"livetype"`
	Packets     types.Int64         `tfsdk:"packets"`
	Priority    types.Int64         `tfsdk:"priority"`
	State       types.String        `tfsdk:"state"`
	TableID     types.Int64         `tfsdk:"tableid"`
	TableName   types.String        `tfsdk:"tablename"`
	Timeout     types.Int64         `tfsdk:"timeout"`
	Selector    flowsSelectorModel  `tfsdk:"selector"`
	Treatment   flowsTreatmentModel `tfsdk:"treatment"`
}

// flowsSelectorModel maps flow ingredients data
type flowsSelectorModel struct {
	Criteria []flowsSelectorCriteriaModel `tfsdk:"criteria"`
}

type flowsSelectorCriteriaModel struct {
	EthType types.String `tfsdk:"ethtype"`
	Mac     types.String `tfsdk:"mac"`
	Port    types.Int64  `tfsdk:"port"`
	Type    types.String `tfsdk:"type"`
}

type flowsTreatmentModel struct {
	ClearDeferred types.Bool                        `tfsdk:"cleardeferred"`
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
						"appid": schema.StringAttribute{
							Computed: true,
						},
						"bytes": schema.Int64Attribute{
							Computed: true,
						},
						"deviceid": schema.StringAttribute{
							Computed: true,
						},
						"groupid": schema.Int64Attribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"ispermanent": schema.BoolAttribute{
							Computed: true,
						},
						"lastseen": schema.Int64Attribute{
							Computed: true,
						},
						"life": schema.Int64Attribute{
							Computed: true,
						},
						"livetype": schema.StringAttribute{
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
						"tableid": schema.Int64Attribute{
							Computed: true,
						},
						"tablename": schema.StringAttribute{
							Computed: true,
						},
						"timeout": schema.Int64Attribute{
							Computed: true,
						},
						"selector": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"criteria": schema.ListNestedAttribute{
									Computed: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"ethtype": schema.StringAttribute{
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
						"treatment": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"cleardeferred": schema.BoolAttribute{
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
	}
}

func (d *flowsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

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

	httpReq, err := http.NewRequest("GET", HostURL, nil)
	if err != nil {
		fmt.Println(nil)
	}
	httpReq.SetBasicAuth("onos", "rocks")

	res, err := d.client.Do(httpReq)
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
	flows := flowsDataSourceUnmarshal{}

	err = json.Unmarshal(body, &flows)
	if err != nil {
		fmt.Println(err)
	}

	var state flowsDataSourceModel

	for _, flow := range flows.Flows {
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
		state.Flows = append(state.Flows, flowState)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
