package tarmak

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {

	var p *schema.Provider
	p = &schema.Provider{
		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:     schema.TypeString,
				Required: true,
			},

			"environment": {
				Type:     schema.TypeString,
				Required: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"tarmak_vault_init_token": resourceVaultInitToken(),
			"tarmak_tunnel":           resourceTarmakTunnel(),
			"tarmak_bastion_instance": resourceTarmakBastionInstance(),
			//"tarmak_vault_kubernetes_cluster": resourceVaultKubernetesCluster(),
		},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

// Config is the configuration structure used to instantiate a
// new Azure management client.
type Config struct {
	ManagementURL string

	Cluster     string
	Environment string
}

func (c *Config) validate() error {
	var err *multierror.Error

	if c.Cluster == "" {
		err = multierror.Append(err, fmt.Errorf("Cluster must be configured for the Tarmak provider"))
	}
	if c.Environment == "" {
		err = multierror.Append(err, fmt.Errorf("Environment must be configured for the Tarmak provider"))
	}

	return err.ErrorOrNil()
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := &Config{
			Cluster:     d.Get("cluster").(string),
			Environment: d.Get("environment").(string),
		}

		if err := config.validate(); err != nil {
			return nil, err
		}

		return config, nil
	}
}

func registerProviderWithSubscription(providerName string, client resources.ProvidersClient) error {
	_, err := client.Register(providerName)
	if err != nil {
		return fmt.Errorf("Cannot register provider %s with Azure Resource Manager: %s.", providerName, err)
	}

	return nil
}

var providerRegistrationOnce sync.Once

// registerAzureResourceProvidersWithSubscription uses the providers client to register
// all Azure resource providers which the Terraform provider may require (regardless of
// whether they are actually used by the configuration or not). It was confirmed by Microsoft
// that this is the approach their own internal tools also take.
func registerAzureResourceProvidersWithSubscription(providerList []resources.Provider, client resources.ProvidersClient) error {
	var err error
	providerRegistrationOnce.Do(func() {
		providers := map[string]struct{}{
			"Microsoft.Compute":           struct{}{},
			"Microsoft.Cache":             struct{}{},
			"Microsoft.ContainerRegistry": struct{}{},
			"Microsoft.ContainerService":  struct{}{},
			"Microsoft.Network":           struct{}{},
			"Microsoft.Cdn":               struct{}{},
			"Microsoft.Storage":           struct{}{},
			"Microsoft.Sql":               struct{}{},
			"Microsoft.Search":            struct{}{},
			"Microsoft.Resources":         struct{}{},
			"Microsoft.ServiceBus":        struct{}{},
			"Microsoft.KeyVault":          struct{}{},
			"Microsoft.EventHub":          struct{}{},
		}

		// filter out any providers already registered
		for _, p := range providerList {
			if _, ok := providers[*p.Namespace]; !ok {
				continue
			}

			if strings.ToLower(*p.RegistrationState) == "registered" {
				log.Printf("[DEBUG] Skipping provider registration for namespace %s\n", *p.Namespace)
				delete(providers, *p.Namespace)
			}
		}

		var wg sync.WaitGroup
		wg.Add(len(providers))
		for providerName := range providers {
			go func(p string) {
				defer wg.Done()
				log.Printf("[DEBUG] Registering provider with namespace %s\n", p)
				if innerErr := registerProviderWithSubscription(p, client); err != nil {
					err = innerErr
				}
			}(providerName)
		}
		wg.Wait()
	})

	return err
}
