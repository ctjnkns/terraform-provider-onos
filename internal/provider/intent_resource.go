package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ctjnkns/onosclient"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const IntentsURL string = "http://localhost:8181/onos/v1/intents"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &intentResource{}
	_ resource.ResourceWithConfigure   = &intentResource{}
	_ resource.ResourceWithImportState = &intentResource{}
)

// intentResource is the resource implementation.
type intentResource struct {
	client *onosclient.Client
}

// NewIntentResource is a helper function to simplify the provider implementation.
func NewIntentResource() resource.Resource {
	return &intentResource{}
}

type intentResourceModel struct {
	Intent      intentModel  `tfsdk:"intent"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

type intentModel struct {
	ID       types.String `tfsdk:"id"`
	AppID    types.String `tfsdk:"appid"`
	Key      types.String `tfsdk:"key"`
	Type     types.String `tfsdk:"type"`
	Priority types.Int64  `tfsdk:"priority"`
	One      types.String `tfsdk:"one"`
	Two      types.String `tfsdk:"two"`
	//ID          string             `tfsdk:"id"`
	//State       string             `tfsdk:"state"`
	//Resources []string `tfsdk:"resources"`
	//Selector    selectorModel      `tfsdk:"selector"`
	//Treatment   treatmentModel     `tfsdk:"treatment"`
	//Constraints []constraintsModel `tfsdk:"constraints"`

}

/*
type selectorModel struct {
	Criteria []criteriaModel `tfsdk:"criteria"`
}

type criteriaModel struct {
	EthType string `tfsdk:"ethtype"`
	Mac     string `tfsdk:"mac"`
	Port    int    `tfsdk:"port"`
	Type    string `tfsdk:"type"`
}

type treatmentModel struct {
	ClearDeferred bool                `tfsdk:"cleardeferred"`
	Deferred      []instructionsModel `tfsdk:"deferred"` //for deferred instructions
	Instructions  []instructionsModel `tfsdk:"instructions"`
}

type instructionsModel struct {
	Port string `tfsdk:"port"`
	Type string `tfsdk:"type"`
}

type constraintsModel struct {
	Inclusive bool     `tfsdk:"inclusive"`
	Types     []string `tfsdk:"types"`
	Type      string   `tfsdk:"type"`
}
*/

// Metadata returns the resource type name.
func (r *intentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_intent"
}

// Schema defines the schema for the data source.
func (r *intentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"intent": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"appid": schema.StringAttribute{
						Required: true,
					},
					"key": schema.StringAttribute{
						Required: true,
					},
					"type": schema.StringAttribute{
						Required: true,
					}, /*
						"resources": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
						},*/
					"priority": schema.Int64Attribute{
						Required: true,
					},
					"one": schema.StringAttribute{
						Required: true,
					},
					"two": schema.StringAttribute{
						Required: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *intentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *intentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan intentResourceModel
	//	Intents []intentModel `tfsdk:"intents"`
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	intent := onosclient.Intent{
		AppID:    string(plan.Intent.AppID.ValueString()),
		Key:      string(plan.Intent.Key.ValueString()),
		Type:     string(plan.Intent.Type.ValueString()),
		Priority: int(plan.Intent.Priority.ValueInt64()),
		One:      string(plan.Intent.One.ValueString()),
		Two:      string(plan.Intent.Two.ValueString()),
	}

	// Create new intent
	intent, err := r.client.CreateIntent(intent)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating intent",
			"Could not create intent, unexpected error: "+err.Error(),
		)
		return
	}
	// Map response body to schema and populate Computed attribute values
	plan.Intent = intentModel{
		ID:       types.StringValue(intent.ID),
		AppID:    types.StringValue(intent.AppID),
		Key:      types.StringValue(intent.Key),
		Type:     types.StringValue(intent.Type),
		Priority: types.Int64Value(int64(intent.Priority)),
		One:      types.StringValue(intent.One),
		Two:      types.StringValue(intent.Two),
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *intentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Get current state
	var state intentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed intent value from Onos
	intent := onosclient.Intent{
		AppID: string(state.Intent.AppID.ValueString()),
		Key:   string(state.Intent.Key.ValueString()),
	}
	intent, err := r.client.GetIntent(intent)
	if err != nil {
		//set state to empty since the intent wasn't found.
		state.Intent = intentModel{}
	} else {

		// Overwrite intent with refreshed state
		state.Intent = intentModel{
			ID:       types.StringValue(intent.ID),
			AppID:    types.StringValue(intent.AppID),
			Key:      types.StringValue(intent.Key),
			Type:     types.StringValue(intent.Type),
			Priority: types.Int64Value(int64(intent.Priority)),
			One:      types.StringValue(intent.One),
			Two:      types.StringValue(intent.Two),
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *intentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan intentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	intent := onosclient.Intent{
		AppID:    string(plan.Intent.AppID.ValueString()),
		Key:      string(plan.Intent.Key.ValueString()),
		Type:     string(plan.Intent.Type.ValueString()),
		Priority: int(plan.Intent.Priority.ValueInt64()),
		One:      string(plan.Intent.One.ValueString()),
		Two:      string(plan.Intent.Two.ValueString()),
	}

	// Update existing intent
	intent, err := r.client.UpdateIntent(intent)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Onos Intent",
			"Could not update intent, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from Getintent as Updateintent items are not
	// populated.
	// NA

	// Update resource state with updated items and timestamp
	plan.Intent = intentModel{
		ID:       types.StringValue(intent.ID),
		AppID:    types.StringValue(intent.AppID),
		Key:      types.StringValue(intent.Key),
		Type:     types.StringValue(intent.Type),
		Priority: types.Int64Value(int64(intent.Priority)),
		One:      types.StringValue(intent.One),
		Two:      types.StringValue(intent.Two),
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *intentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state intentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//convert to model
	intent := onosclient.Intent{
		AppID: string(state.Intent.AppID.ValueString()),
		Key:   string(state.Intent.Key.ValueString()),
	}

	// Delete existing intent
	err := r.client.DeleteIntent(intent)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting onos intent",
			"Could not delete intent, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *intentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: AppID,Key. Got: %q", req.ID),
		)
		return
	}

	// Create the intent model with the passed in import string data
	var plan intentResourceModel
	plan.Intent = intentModel{
		AppID: types.StringValue(idParts[0]),
		Key:   types.StringValue(idParts[1]),
	}
	// Set the resp
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
