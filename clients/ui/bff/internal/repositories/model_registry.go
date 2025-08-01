package repositories

import (
	"context"
	"fmt"

	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"

	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type ModelRegistryRepository struct {
}

func NewModelRegistryRepository() *ModelRegistryRepository {
	return &ModelRegistryRepository{}
}

func (m *ModelRegistryRepository) GetAllModelRegistries(sessionCtx context.Context, client k8s.KubernetesClientInterface, namespace string) ([]models.ModelRegistryModel, error) {
	// Default to non-federated mode for backward compatibility
	return m.GetAllModelRegistriesWithMode(sessionCtx, client, namespace, false)
}

// GetAllModelRegistriesWithMode fetches all model registries with support for federated mode
func (m *ModelRegistryRepository) GetAllModelRegistriesWithMode(sessionCtx context.Context, client k8s.KubernetesClientInterface, namespace string, isFederatedMode bool) ([]models.ModelRegistryModel, error) {

	// TODO: In default mode fetch Routes for external access.
	resources, err := client.GetServiceDetails(sessionCtx, namespace)

	if err != nil {
		return nil, fmt.Errorf("error fetching model registries: %w", err)
	}

	var registries = []models.ModelRegistryModel{}
	for _, s := range resources {
		serverAddress := m.ResolveServerAddress(s.ClusterIP, s.HTTPPort, s.IsHTTPS, s.ExternalAddressRest, isFederatedMode)
		registry := models.ModelRegistryModel{
			Name:          s.Name,
			Description:   s.Description,
			DisplayName:   s.DisplayName,
			ServerAddress: serverAddress,
			IsHTTPS:       s.IsHTTPS,
		}
		registries = append(registries, registry)
	}

	return registries, nil
}

func (m *ModelRegistryRepository) GetModelRegistry(sessionCtx context.Context, client k8s.KubernetesClientInterface, namespace string, modelRegistryID string) (models.ModelRegistryModel, error) {
	// Default to non-federated mode for backward compatibility
	return m.GetModelRegistryWithMode(sessionCtx, client, namespace, modelRegistryID, false)
}

// GetModelRegistryWithMode fetches a specific model registry with support for federated mode
func (m *ModelRegistryRepository) GetModelRegistryWithMode(sessionCtx context.Context, client k8s.KubernetesClientInterface, namespace string, modelRegistryID string, isFederatedMode bool) (models.ModelRegistryModel, error) {

	s, err := client.GetServiceDetailsByName(sessionCtx, namespace, modelRegistryID)
	if err != nil {
		return models.ModelRegistryModel{}, fmt.Errorf("error fetching model registry: %w", err)
	}

	modelRegistry := models.ModelRegistryModel{
		Name:          s.Name,
		Description:   s.Description,
		DisplayName:   s.DisplayName,
		ServerAddress: m.ResolveServerAddress(s.ClusterIP, s.HTTPPort, s.IsHTTPS, s.ExternalAddressRest, isFederatedMode),
		IsHTTPS:       s.IsHTTPS,
	}

	return modelRegistry, nil
}

func (m *ModelRegistryRepository) ResolveServerAddress(clusterIP string, httpPort int32, isHTTPS bool, externalAddressRest string, isFederatedMode bool) string {
	// Default behavior - use cluster IP and port
	protocol := "http"
	if isHTTPS {
		protocol = "https"
	}
	// In federated mode, if external address is available, use it
	if isFederatedMode && externalAddressRest != "" {
		// External address is assumed to be HTTPS
		url := fmt.Sprintf("%s://%s/api/model_registry/v1alpha3", protocol, externalAddressRest)
		return url
	}

	url := fmt.Sprintf("%s://%s:%d/api/model_registry/v1alpha3", protocol, clusterIP, httpPort)
	return url
}
