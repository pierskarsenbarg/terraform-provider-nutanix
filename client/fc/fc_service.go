package fc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

// Operations ...
type Operations struct {
	client *client.Client
}

// Service ...
type Service interface {
	GetImagedNode(uuid string) (*ImagedNodeDetails, error)
	ListImagedNodes(req *ImagedNodesListInput) (*ImagedNodesListResponse, error)
	GetImagedCluster(uuid string) (*ImagedClusterDetails, error)
	ListImagedClusters(input *ImagedClustersListInput) (*ImagedClustersListResponse, error)
	CreateCluster(input *CreateClusterInput) (*CreateClusterResponse, error)
	UpdateCluster(clusterUUID string, updateData *UpdateClusterData) error
	DeleteCluster(clusterUUID string) error
	CreateAPIKey(input *CreateAPIKeysInput) (*CreateAPIKeysResponse, error)
	GetAPIKey(uuid string) (*CreateAPIKeysResponse, error)
	ListAPIKeys(body *ListMetadataInput) (*ListAPIKeysResponse, error)
}

func (op Operations) GetImagedNode(nodeUUID string) (*ImagedNodeDetails, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/imaged_nodes/%s", nodeUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	imagedNodeDetails := new(ImagedNodeDetails)

	if err != nil {
		return nil, err
	}

	return imagedNodeDetails, op.client.Do(ctx, req, imagedNodeDetails)
}

func (op Operations) ListImagedNodes(input *ImagedNodesListInput) (*ImagedNodesListResponse, error) {
	ctx := context.TODO()
	path := "/imaged_nodes/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	imagedNodesListResponse := new(ImagedNodesListResponse)

	if err != nil {
		return nil, err
	}

	return imagedNodesListResponse, op.client.Do(ctx, req, imagedNodesListResponse)
}

func (op Operations) GetImagedCluster(clusterUUID string) (*ImagedClusterDetails, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/imaged_clusters/%s", clusterUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	imagedClusterDetails := new(ImagedClusterDetails)

	if err != nil {
		return nil, err
	}

	return imagedClusterDetails, op.client.Do(ctx, req, imagedClusterDetails)
}

func (op Operations) ListImagedClusters(input *ImagedClustersListInput) (*ImagedClustersListResponse, error) {
	ctx := context.TODO()
	path := "/imaged_clusters/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	imagedClustersListResponse := new(ImagedClustersListResponse)

	if err != nil {
		return nil, err
	}

	return imagedClustersListResponse, op.client.Do(ctx, req, imagedClustersListResponse)
}

func (op Operations) CreateCluster(input *CreateClusterInput) (*CreateClusterResponse, error) {
	ctx := context.TODO()
	path := "/imaged_clusters"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	createClusterResponse := new(CreateClusterResponse)

	if err != nil {
		return nil, err
	}

	return createClusterResponse, op.client.Do(ctx, req, createClusterResponse)
}

func (op Operations) UpdateCluster(clusterUUID string, updateData *UpdateClusterData) error {
	ctx := context.TODO()
	path := fmt.Sprintf("/imaged_clusters/%s", clusterUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPut, path, updateData)

	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

func (op Operations) DeleteCluster(clusterUUID string) error {
	ctx := context.TODO()
	path := fmt.Sprintf("/imaged_clusters/%s", clusterUUID)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)

	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

//Create a new api key which will be used by remote nodes to authenticate with Foundation Central
func (op Operations) CreateAPIKey(input *CreateAPIKeysInput) (*CreateAPIKeysResponse, error) {
	ctx := context.TODO()
	path := "/api_keys"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)
	if err != nil {
		return nil, err
	}

	createApiResponse := new(CreateAPIKeysResponse)
	return createApiResponse, op.client.Do(ctx, req, createApiResponse)
}

//Get an api key given its UUID.
func (op Operations) GetAPIKey(uuid string) (*CreateAPIKeysResponse, error) {
	ctx := context.TODO()
	path := fmt.Sprintf("/api_keys/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, uuid)
	if err != nil {
		return nil, err
	}

	createApiResponse := new(CreateAPIKeysResponse)
	return createApiResponse, op.client.Do(ctx, req, createApiResponse)
}

//List all the api keys.
func (op Operations) ListAPIKeys(body *ListMetadataInput) (*ListAPIKeysResponse, error) {
	ctx := context.TODO()
	path := "/api_keys/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}

	listApiKeysResponse := new(ListAPIKeysResponse)
	return listApiKeysResponse, op.client.Do(ctx, req, listApiKeysResponse)
}
