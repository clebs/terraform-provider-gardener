package gardener

import (
	"errors"

	gardner_apis "github.com/gardener/gardener/pkg/client/garden/clientset/versioned/typed/garden/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
)

type Bindings struct {
	AwsSecretBinding       string
	GcpSecretBinding       string
	AzureSecretBinding     string
	OpenStackSecretBinding string
	AliCloudSecretBinding  string
}
type Config struct {
	Profile        string
	KubePath       string
	SecretBindings *Bindings
}
type GardenerClient struct {
	NameSpace         string
	DNSBase           string
	GardenerClientSet *gardner_apis.GardenV1beta1Client
	SecretBindings    *Bindings
}

// Client configures and returns a fully initialized GardnerClient
func (c *Config) Client() (interface{}, error) {
	if c.SecretBindings.AwsSecretBinding == "" && c.SecretBindings.GcpSecretBinding == "" && c.SecretBindings.AzureSecretBinding == "" &&
		c.SecretBindings.OpenStackSecretBinding == "" && c.SecretBindings.AliCloudSecretBinding == "" {
		return nil, errors.New("at least one binding needs to be defined")
	}

	config, err := clientcmd.BuildConfigFromFlags("", c.KubePath)
	if err != nil {
		return nil, err
	}
	clientset, err := gardner_apis.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client := &GardenerClient{
		NameSpace:         "garden-" + c.Profile,
		DNSBase:           c.Profile + ".shoot.canary.k8s-hana.ondemand.com",
		GardenerClientSet: clientset,
		SecretBindings:    c.SecretBindings,
	}
	return client, nil
}
