package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const IntentsURL string = "http://localhost:8181/onos/v1/intents"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &intentResource{}
	_ resource.ResourceWithConfigure = &intentResource{}
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &intentResource{}
	_ resource.ResourceWithConfigure = &intentResource{}
)

// intentResource is the resource implementation.
type intentResource struct {
	client *http.Client
}

// NewIntentResource is a helper function to simplify the provider implementation.
func NewIntentResource() resource.Resource {
	return &intentResource{}
}

type intentResourceUnmarshal struct {
	Intents []intentsUnmarshal `tfsdk:"intents"`
}

type intentsUnmarshal struct {
	AppID     string   `json:"appid"`
	ID        string   `json:"id"`
	Key       string   `json:"key"`
	State     string   `json:"state"`
	Type      string   `json:"type"`
	Resources []string `json:"resources"`
}

type intentResourceModel struct {
	Intents []intentsModel `tfsdk:"intents"`
}

type intentsModel struct {
	AppID     types.String `tfsdk:"appid"`
	ID        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
	State     types.String `tfsdk:"state"`
	Type      types.String `tfsdk:"type"`
	Resources types.List   `tfsdk:"resources"`
}

// Metadata returns the resource type name.
func (r *intentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_intent"
}

// Schema defines the schema for the data source.
func (r *intentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"intents": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"appid": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"key": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"resources": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
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

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *intentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

// Read refreshes the Terraform state with the latest data.
func (r *intentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	httpReq, err := http.NewRequest("GET", IntentsURL, nil)
	if err != nil {
		fmt.Println(nil)
	}
	httpReq.SetBasicAuth("onos", "rocks")

	res, err := r.client.Do(httpReq)
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
	intents := intentResourceUnmarshal{}

	err = json.Unmarshal(body, &intents)
	if err != nil {
		fmt.Println(err)
	}

	var state intentResourceModel

	for _, intent := range intents.Intents {
		resourceElements := []attr.Value{types.StringValue(intent.Resources[0]), types.StringValue(intent.Resources[1])}
		intentState := intentsModel{

			AppID:     types.StringValue(intent.AppID),
			ID:        types.StringValue(intent.ID),
			Key:       types.StringValue(intent.Key),
			State:     types.StringValue(intent.State),
			Type:      types.StringValue(intent.Type),
			Resources: types.ListNull(types.StringType),
		}
		intentState.Resources, _ = types.ListValue(types.StringType, resourceElements)
		state.Intents = append(state.Intents, intentState)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *intentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *intentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
