//go:build ignores

package cognitoidp_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	tfopensearch "github.com/hashicorp/terraform-provider-aws/internal/service/opensearch"
)

const (
	openSearchDomainMaxLen = 28

	openSearchDomainPrefix    = "tf-acc-"
	openSearchDomainPrefixLen = len(openSearchDomainPrefix)

	openSearchDomainRemainderLen = openSearchDomainMaxLen - openSearchDomainPrefixLen
)

func randomOpenSearchDomainName() string {
	return fmt.Sprintf(openSearchDomainPrefix+"%s", sdkacctest.RandString(openSearchDomainRemainderLen))
}

func TestAccCognitoIDPManagedUserPoolClient_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var client cognitoidentityprovider.UserPoolClientType
	rName := randomOpenSearchDomainName()
	resourceName := "aws_cognito_managed_user_pool_client.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckIdentityProvider(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentityprovider.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserPoolClientDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedUserPoolClientConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(fmt.Sprintf(`^AmazonOpenSearchService-%s`, rName))),
					resource.TestCheckResourceAttr(resourceName, "access_token_validity", "0"),
					resource.TestCheckResourceAttr(resourceName, "allowed_oauth_flows.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "allowed_oauth_flows.0", "code"),
					resource.TestCheckResourceAttr(resourceName, "allowed_oauth_flows_user_pool_client", "true"),
					resource.TestCheckResourceAttr(resourceName, "allowed_oauth_scopes.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "allowed_oauth_scopes.0", "email"),
					resource.TestCheckResourceAttr(resourceName, "allowed_oauth_scopes.1", "openid"),
					resource.TestCheckResourceAttr(resourceName, "allowed_oauth_scopes.2", "phone"),
					resource.TestCheckResourceAttr(resourceName, "allowed_oauth_scopes.3", "profile"),
					resource.TestCheckResourceAttr(resourceName, "analytics_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "auth_session_validity", "3"),
					resource.TestCheckResourceAttr(resourceName, "callback_urls.#", "1"),
					resource.TestMatchResourceAttr(resourceName, "callback_urls.0", regexp.MustCompile(fmt.Sprintf(`https://search-%s-\w+.%s.es.amazonaws.com/_dashboards/app/home`, rName, acctest.Region()))),
					resource.TestMatchResourceAttr(resourceName, "client_secret", regexp.MustCompile(`\w+`)),
					resource.TestCheckResourceAttr(resourceName, "default_redirect_uri", ""),
					resource.TestCheckResourceAttr(resourceName, "enable_propagate_additional_user_context_data", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_token_revocation", "true"),
					resource.TestCheckNoResourceAttr(resourceName, "explicit_auth_flows"),
					resource.TestCheckResourceAttr(resourceName, "id_token_validity", "0"),
					resource.TestCheckResourceAttr(resourceName, "logout_urls.#", "1"),
					resource.TestMatchResourceAttr(resourceName, "logout_urls.0", regexp.MustCompile(fmt.Sprintf(`https://search-%s-\w+.%s.es.amazonaws.com/_dashboards/app/home`, rName, acctest.Region()))),
					resource.TestCheckResourceAttr(resourceName, "prevent_user_existence_errors", ""),
					resource.TestCheckNoResourceAttr(resourceName, "read_attributes"),
					resource.TestCheckResourceAttr(resourceName, "refresh_token_validity", "30"),
					resource.TestCheckResourceAttr(resourceName, "supported_identity_providers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "supported_identity_providers.0", "COGNITO"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.#", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "user_pool_id", "aws_cognito_user_pool.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "write_attributes.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
		},
	})
}

func TestAccCognitoIDPManagedUserPoolClient_namePattern(t *testing.T) {
	ctx := acctest.Context(t)
	var client cognitoidentityprovider.UserPoolClientType
	rName := randomOpenSearchDomainName()
	resourceName := "aws_cognito_managed_user_pool_client.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckIdentityProvider(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentityprovider.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserPoolClientDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedUserPoolClientConfig_namePattern(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(fmt.Sprintf(`^AmazonOpenSearchService-%s`, rName))),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_pattern",
				},
			},
		},
	})
}

func TestAccCognitoIDPManagedUserPoolClient_enableRevocation(t *testing.T) {
	ctx := acctest.Context(t)
	var client cognitoidentityprovider.UserPoolClientType
	rName := randomOpenSearchDomainName()

	resourceName := "aws_cognito_managed_user_pool_client.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckIdentityProvider(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentityprovider.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserPoolClientDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedUserPoolClientConfig_revocation(rName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "enable_token_revocation", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
			{
				Config: testAccManagedUserPoolClientConfig_revocation(rName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "enable_token_revocation", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
		},
	})
}

func TestAccCognitoIDPManagedUserPoolClient_accessTokenValidity(t *testing.T) {
	ctx := acctest.Context(t)
	var client cognitoidentityprovider.UserPoolClientType
	rName := randomOpenSearchDomainName()
	resourceName := "aws_cognito_managed_user_pool_client.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckIdentityProvider(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentityprovider.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserPoolClientDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedUserPoolClientConfig_accessTokenValidity(rName, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "access_token_validity", "5"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
			{
				Config: testAccManagedUserPoolClientConfig_accessTokenValidity(rName, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "access_token_validity", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
		},
	})
}

func TestAccCognitoIDPManagedUserPoolClient_idTokenValidity(t *testing.T) {
	ctx := acctest.Context(t)
	var client cognitoidentityprovider.UserPoolClientType
	rName := randomOpenSearchDomainName()
	resourceName := "aws_cognito_managed_user_pool_client.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckIdentityProvider(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentityprovider.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserPoolClientDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedUserPoolClientConfig_idTokenValidity(rName, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "id_token_validity", "5"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
			{
				Config: testAccManagedUserPoolClientConfig_idTokenValidity(rName, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "id_token_validity", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
		},
	})
}

func TestAccCognitoIDPManagedUserPoolClient_refreshTokenValidity(t *testing.T) {
	ctx := acctest.Context(t)
	var client cognitoidentityprovider.UserPoolClientType
	rName := randomOpenSearchDomainName()
	resourceName := "aws_cognito_managed_user_pool_client.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckIdentityProvider(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentityprovider.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserPoolClientDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedUserPoolClientConfig_refreshTokenValidity(rName, 60),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "refresh_token_validity", "60"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
			{
				Config: testAccManagedUserPoolClientConfig_refreshTokenValidity(rName, 120),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "refresh_token_validity", "120"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
		},
	})
}

func TestAccCognitoIDPManagedUserPoolClient_tokenValidityUnits(t *testing.T) {
	ctx := acctest.Context(t)
	var client cognitoidentityprovider.UserPoolClientType
	rName := randomOpenSearchDomainName()
	resourceName := "aws_cognito_managed_user_pool_client.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckIdentityProvider(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentityprovider.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserPoolClientDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedUserPoolClientConfig_tokenValidityUnits(rName, "days"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.access_token", "days"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.id_token", "days"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.refresh_token", "days"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
			{
				Config: testAccManagedUserPoolClientConfig_tokenValidityUnits(rName, "hours"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.access_token", "hours"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.id_token", "hours"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.refresh_token", "hours"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
		},
	})
}

func TestAccCognitoIDPManagedUserPoolClient_tokenValidityUnits_AccessToken(t *testing.T) {
	ctx := acctest.Context(t)
	var client cognitoidentityprovider.UserPoolClientType
	rName := randomOpenSearchDomainName()
	resourceName := "aws_cognito_managed_user_pool_client.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckIdentityProvider(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentityprovider.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserPoolClientDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedUserPoolClientConfig_tokenValidityUnits_Unit(rName, "access_token", "days"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.access_token", "days"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.id_token", "hours"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.refresh_token", "days"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
			{
				Config: testAccManagedUserPoolClientConfig_tokenValidityUnits(rName, "hours"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.access_token", "hours"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.id_token", "hours"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.refresh_token", "hours"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
		},
	})
}

func TestAccCognitoIDPManagedUserPoolClient_tokenValidityUnitsWTokenValidity(t *testing.T) {
	ctx := acctest.Context(t)
	var client cognitoidentityprovider.UserPoolClientType
	rName := randomOpenSearchDomainName()
	resourceName := "aws_cognito_managed_user_pool_client.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckIdentityProvider(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentityprovider.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserPoolClientDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedUserPoolClientConfig_tokenValidityUnitsTokenValidity(rName, "days"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.access_token", "days"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.id_token", "days"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.refresh_token", "days"),
					resource.TestCheckResourceAttr(resourceName, "id_token_validity", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
			{
				Config: testAccManagedUserPoolClientConfig_tokenValidityUnitsTokenValidity(rName, "hours"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.access_token", "hours"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.id_token", "hours"),
					resource.TestCheckResourceAttr(resourceName, "token_validity_units.0.refresh_token", "hours"),
					resource.TestCheckResourceAttr(resourceName, "id_token_validity", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
		},
	})
}

func TestAccCognitoIDPManagedUserPoolClient_allFields(t *testing.T) {
	ctx := acctest.Context(t)
	var client cognitoidentityprovider.UserPoolClientType
	rName := randomOpenSearchDomainName()
	resourceName := "aws_cognito_managed_user_pool_client.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckIdentityProvider(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentityprovider.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserPoolClientDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedUserPoolClientConfig_allFields(rName, 300),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(fmt.Sprintf(`^AmazonOpenSearchService-%s`, rName))),
					resource.TestCheckResourceAttr(resourceName, "explicit_auth_flows.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "explicit_auth_flows.*", "CUSTOM_AUTH_FLOW_ONLY"),
					resource.TestCheckTypeSetElemAttr(resourceName, "explicit_auth_flows.*", "USER_PASSWORD_AUTH"),
					resource.TestCheckTypeSetElemAttr(resourceName, "explicit_auth_flows.*", "ADMIN_NO_SRP_AUTH"),
					resource.TestCheckResourceAttr(resourceName, "read_attributes.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "read_attributes.*", "email"),
					resource.TestCheckResourceAttr(resourceName, "write_attributes.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "write_attributes.*", "email"),
					resource.TestCheckResourceAttr(resourceName, "refresh_token_validity", "300"),
					resource.TestCheckResourceAttr(resourceName, "allowed_oauth_flows.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_oauth_flows.*", "code"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_oauth_flows.*", "implicit"),
					resource.TestCheckResourceAttr(resourceName, "allowed_oauth_flows_user_pool_client", "true"),
					resource.TestCheckResourceAttr(resourceName, "allowed_oauth_scopes.#", "5"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_oauth_scopes.*", "openid"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_oauth_scopes.*", "email"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_oauth_scopes.*", "phone"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_oauth_scopes.*", "aws.cognito.signin.user.admin"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_oauth_scopes.*", "profile"),
					resource.TestCheckResourceAttr(resourceName, "callback_urls.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "callback_urls.*", "https://www.example.com/callback"),
					resource.TestCheckTypeSetElemAttr(resourceName, "callback_urls.*", "https://www.example.com/redirect"),
					resource.TestCheckResourceAttr(resourceName, "default_redirect_uri", "https://www.example.com/redirect"),
					resource.TestCheckResourceAttr(resourceName, "logout_urls.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "logout_urls.*", "https://www.example.com/login"),
					resource.TestCheckResourceAttr(resourceName, "prevent_user_existence_errors", "LEGACY"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
		},
	})
}

func TestAccCognitoIDPManagedUserPoolClient_allFieldsUpdatingOneField(t *testing.T) {
	ctx := acctest.Context(t)
	var client cognitoidentityprovider.UserPoolClientType
	rName := randomOpenSearchDomainName()
	resourceName := "aws_cognito_managed_user_pool_client.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckIdentityProvider(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentityprovider.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserPoolClientDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedUserPoolClientConfig_allFields(rName, 300),
			},
			{
				Config: testAccManagedUserPoolClientConfig_allFields(rName, 299),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(fmt.Sprintf(`^AmazonOpenSearchService-%s`, rName))),
					resource.TestCheckResourceAttr(resourceName, "explicit_auth_flows.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "explicit_auth_flows.*", "CUSTOM_AUTH_FLOW_ONLY"),
					resource.TestCheckTypeSetElemAttr(resourceName, "explicit_auth_flows.*", "USER_PASSWORD_AUTH"),
					resource.TestCheckTypeSetElemAttr(resourceName, "explicit_auth_flows.*", "ADMIN_NO_SRP_AUTH"),
					resource.TestCheckResourceAttr(resourceName, "read_attributes.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "read_attributes.*", "email"),
					resource.TestCheckResourceAttr(resourceName, "write_attributes.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "write_attributes.*", "email"),
					resource.TestCheckResourceAttr(resourceName, "refresh_token_validity", "299"),
					resource.TestCheckResourceAttr(resourceName, "allowed_oauth_flows.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_oauth_flows.*", "code"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_oauth_flows.*", "implicit"),
					resource.TestCheckResourceAttr(resourceName, "allowed_oauth_flows_user_pool_client", "true"),
					resource.TestCheckResourceAttr(resourceName, "allowed_oauth_scopes.#", "5"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_oauth_scopes.*", "openid"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_oauth_scopes.*", "email"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_oauth_scopes.*", "phone"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_oauth_scopes.*", "aws.cognito.signin.user.admin"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_oauth_scopes.*", "profile"),
					resource.TestCheckResourceAttr(resourceName, "callback_urls.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "callback_urls.*", "https://www.example.com/callback"),
					resource.TestCheckTypeSetElemAttr(resourceName, "callback_urls.*", "https://www.example.com/redirect"),
					resource.TestCheckResourceAttr(resourceName, "default_redirect_uri", "https://www.example.com/redirect"),
					resource.TestCheckResourceAttr(resourceName, "logout_urls.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "logout_urls.*", "https://www.example.com/login"),
					resource.TestCheckResourceAttr(resourceName, "prevent_user_existence_errors", "LEGACY"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
		},
	})
}

func TestAccCognitoIDPManagedUserPoolClient_analyticsApplicationID(t *testing.T) {
	ctx := acctest.Context(t)
	var client cognitoidentityprovider.UserPoolClientType
	rName := randomOpenSearchDomainName()
	resourceName := "aws_cognito_managed_user_pool_client.test"
	pinpointResourceName := "aws_pinpoint_app.analytics"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			testAccPreCheckIdentityProvider(ctx, t)
			testAccPreCheckPinpointApp(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentityprovider.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserPoolClientDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedUserPoolClientConfig_analyticsApplicationID(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "analytics_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "analytics_configuration.0.application_id", pinpointResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "analytics_configuration.0.external_id", rName),
					resource.TestCheckResourceAttrPair(resourceName, "analytics_configuration.0.role_arn", "aws_iam_role.analytics", "arn"),
					resource.TestCheckResourceAttr(resourceName, "analytics_configuration.0.user_data_shared", "false"),
					resource.TestCheckNoResourceAttr(resourceName, "analytics_configuration.0.application_arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
			{
				Config: testAccManagedUserPoolClientConfig_analyticsShareData(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "analytics_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "analytics_configuration.0.application_id", pinpointResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "analytics_configuration.0.external_id", rName),
					resource.TestCheckResourceAttr(resourceName, "analytics_configuration.0.user_data_shared", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
			{
				Config: testAccManagedUserPoolClientConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "analytics_configuration.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
		},
	})
}

func TestAccCognitoIDPManagedUserPoolClient_analyticsWithARN(t *testing.T) {
	t.Skip("this test hangs on deletion")

	ctx := acctest.Context(t)
	var client cognitoidentityprovider.UserPoolClientType
	rName := randomOpenSearchDomainName()
	resourceName := "aws_cognito_managed_user_pool_client.test"
	pinpointResourceName := "aws_pinpoint_app.analytics"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			testAccPreCheckIdentityProvider(ctx, t)
			testAccPreCheckPinpointApp(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentityprovider.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserPoolClientDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedUserPoolClientConfig_analyticsARN(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "analytics_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "analytics_configuration.0.application_arn", pinpointResourceName, "arn"),
					acctest.CheckResourceAttrGlobalARN(resourceName, "analytics_configuration.0.role_arn", "iam", "role/aws-service-role/cognito-idp.amazonaws.com/AWSServiceRoleForAmazonCognitoIdp"),
					resource.TestCheckResourceAttr(resourceName, "analytics_configuration.0.user_data_shared", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
			{
				Config: testAccManagedUserPoolClientConfig_analyticsARNShareData(rName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "analytics_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "analytics_configuration.0.application_arn", pinpointResourceName, "arn"),
					acctest.CheckResourceAttrGlobalARN(resourceName, "analytics_configuration.0.role_arn", "iam", "role/aws-service-role/cognito-idp.amazonaws.com/AWSServiceRoleForAmazonCognitoIdp"),
					resource.TestCheckResourceAttr(resourceName, "analytics_configuration.0.user_data_shared", "true"),
					resource.TestCheckNoResourceAttr(resourceName, "analytics_configuration.0.application_id"),
					resource.TestCheckNoResourceAttr(resourceName, "analytics_configuration.0.external_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
		},
	})
}

func TestAccCognitoIDPManagedUserPoolClient_authSessionValidity(t *testing.T) {
	ctx := acctest.Context(t)
	var client cognitoidentityprovider.UserPoolClientType
	rName := randomOpenSearchDomainName()
	resourceName := "aws_cognito_managed_user_pool_client.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckIdentityProvider(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentityprovider.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserPoolClientDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedUserPoolClientConfig_authSessionValidity(rName, 15),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "auth_session_validity", "15"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
			{
				Config: testAccManagedUserPoolClientConfig_authSessionValidity(rName, 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					resource.TestCheckResourceAttr(resourceName, "auth_session_validity", "10"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccUserPoolClientImportStateIDFunc(ctx, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
		},
	})
}

func TestAccCognitoIDPManagedUserPoolClient_Disappears_OpenSearchDomain(t *testing.T) {
	ctx := acctest.Context(t)
	var client cognitoidentityprovider.UserPoolClientType
	rName := randomOpenSearchDomainName()
	resourceName := "aws_cognito_managed_user_pool_client.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckIdentityProvider(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cognitoidentityprovider.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserPoolClientDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedUserPoolClientConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckUserPoolClientExists(ctx, resourceName, &client),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfopensearch.ResourceDomain(), "aws_opensearch_domain.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccManagedUserPoolClientBaseConfig(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_cognito_user_pool" "test" {
  name = %[1]q
}

resource "aws_cognito_user_pool_domain" "test" {
  domain       = %[1]q
  user_pool_id = aws_cognito_user_pool.test.id
}

resource "aws_cognito_identity_pool" "test" {
  identity_pool_name               = %[1]q
  allow_unauthenticated_identities = false

  lifecycle {
    ignore_changes = [cognito_identity_providers]
  }
}

resource "aws_opensearch_domain" "test" {
  domain_name = %[1]q

  engine_version = "OpenSearch_1.1"

  cognito_options {
    enabled          = true
    user_pool_id     = aws_cognito_user_pool.test.id
    identity_pool_id = aws_cognito_identity_pool.test.id
    role_arn         = aws_iam_role.test.arn
  }

  ebs_options {
    ebs_enabled = true
    volume_size = 10
  }

  depends_on = [
    aws_cognito_user_pool_domain.test,
    aws_iam_role_policy_attachment.test,
  ]

  timeouts {
    delete = "20m"
  }
}

resource "aws_iam_role" "test" {
  name               = %[1]q
  path               = "/service-role/"
  assume_role_policy = data.aws_iam_policy_document.test.json
}

data "aws_iam_policy_document" "test" {
  statement {
    sid     = ""
    actions = ["sts:AssumeRole"]
    effect  = "Allow"

    principals {
      type = "Service"
      identifiers = [
        "es.${data.aws_partition.current.dns_suffix}",
      ]
    }
  }
}

resource "aws_iam_role_policy_attachment" "test" {
  role       = aws_iam_role.test.name
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonESCognitoAccess"
}
`, rName)
}

func testAccManagedUserPoolClientConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_cognito_managed_user_pool_client" "test" {
  name_prefix  = "AmazonOpenSearchService-%[1]s"
  user_pool_id = aws_cognito_user_pool.test.id

  depends_on = [
    aws_opensearch_domain.test,
  ]
}
`, rName))
}

func testAccManagedUserPoolClientConfig_namePattern(rName string) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_cognito_managed_user_pool_client" "test" {
  name_pattern = %[1]q
  user_pool_id = aws_cognito_user_pool.test.id

  depends_on = [
    aws_opensearch_domain.test,
  ]
}
`, rName))
}

func testAccManagedUserPoolClientConfig_revocation(rName string, revoke bool) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_cognito_managed_user_pool_client" "test" {
  name_prefix  = "AmazonOpenSearchService-%[1]s"
  user_pool_id = aws_cognito_user_pool.test.id

  depends_on = [
    aws_opensearch_domain.test,
  ]

  enable_token_revocation = %[2]t
}
`, rName, revoke))
}

func testAccManagedUserPoolClientConfig_accessTokenValidity(rName string, validity int) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_cognito_managed_user_pool_client" "test" {
  name_prefix  = "AmazonOpenSearchService-%[1]s"
  user_pool_id = aws_cognito_user_pool.test.id

  depends_on = [
    aws_opensearch_domain.test,
  ]

  access_token_validity = %[2]d
}
`, rName, validity))
}

func testAccManagedUserPoolClientConfig_idTokenValidity(rName string, validity int) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_cognito_managed_user_pool_client" "test" {
  name_prefix  = "AmazonOpenSearchService-%[1]s"
  user_pool_id = aws_cognito_user_pool.test.id

  depends_on = [
    aws_opensearch_domain.test,
  ]

  id_token_validity = %[2]d
}
`, rName, validity))
}

func testAccManagedUserPoolClientConfig_refreshTokenValidity(rName string, refreshTokenValidity int) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_cognito_managed_user_pool_client" "test" {
  name_prefix  = "AmazonOpenSearchService-%[1]s"
  user_pool_id = aws_cognito_user_pool.test.id

  depends_on = [
    aws_opensearch_domain.test,
  ]

  refresh_token_validity = %[2]d
}
`, rName, refreshTokenValidity))
}

func testAccManagedUserPoolClientConfig_tokenValidityUnits(rName, units string) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_cognito_managed_user_pool_client" "test" {
  name_prefix  = "AmazonOpenSearchService-%[1]s"
  user_pool_id = aws_cognito_user_pool.test.id

  depends_on = [
    aws_opensearch_domain.test,
  ]

  token_validity_units {
    access_token  = %[2]q
    id_token      = %[2]q
    refresh_token = %[2]q
  }
}
`, rName, units))
}

func testAccManagedUserPoolClientConfig_tokenValidityUnits_Unit(rName, unit, value string) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_cognito_managed_user_pool_client" "test" {
  name_prefix  = "AmazonOpenSearchService-%[1]s"
  user_pool_id = aws_cognito_user_pool.test.id

  depends_on = [
    aws_opensearch_domain.test,
  ]

  token_validity_units {
    %[2]s = %[3]q
  }
}
`, rName, unit, value))
}

func testAccManagedUserPoolClientConfig_tokenValidityUnitsTokenValidity(rName, units string) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_cognito_managed_user_pool_client" "test" {
  name_prefix  = "AmazonOpenSearchService-%[1]s"
  user_pool_id = aws_cognito_user_pool.test.id

  depends_on = [
    aws_opensearch_domain.test,
  ]

  id_token_validity = 1

  token_validity_units {
    access_token  = %[2]q
    id_token      = %[2]q
    refresh_token = %[2]q
  }
}
`, rName, units))
}

func testAccManagedUserPoolClientConfig_allFields(rName string, refreshTokenValidity int) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_cognito_managed_user_pool_client" "test" {
  name_prefix  = "AmazonOpenSearchService-%[1]s"
  user_pool_id = aws_cognito_user_pool.test.id

  depends_on = [
    aws_opensearch_domain.test,
  ]

  explicit_auth_flows = ["ADMIN_NO_SRP_AUTH", "CUSTOM_AUTH_FLOW_ONLY", "USER_PASSWORD_AUTH"]

  read_attributes  = ["email"]
  write_attributes = ["email"]

  refresh_token_validity        = %[2]d
  prevent_user_existence_errors = "LEGACY"

  allowed_oauth_flows                  = ["code", "implicit"]
  allowed_oauth_flows_user_pool_client = "true"
  allowed_oauth_scopes                 = ["phone", "email", "openid", "profile", "aws.cognito.signin.user.admin"]

  callback_urls        = ["https://www.example.com/redirect", "https://www.example.com/callback"]
  default_redirect_uri = "https://www.example.com/redirect"
  logout_urls          = ["https://www.example.com/login"]
}
`, rName, refreshTokenValidity))
}

func testAccManagedUserPoolClientAnalyticsBaseConfig(rName string) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientBaseConfig(rName),
		fmt.Sprintf(`
data "aws_caller_identity" "current" {}

resource "aws_pinpoint_app" "analytics" {
  name = %[1]q
}

resource "aws_iam_role" "analytics" {
  name = "%[1]s-analytics"

  depends_on = [

  ]

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "cognito-idp.${data.aws_partition.current.dns_suffix}"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "analytics" {
  name = %[1]q
  role = aws_iam_role.analytics.id

  policy = <<-EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "mobiletargeting:UpdateEndpoint",
        "mobiletargeting:PutEvents"
      ],
      "Effect": "Allow",
      "Resource": "arn:${data.aws_partition.current.partition}:mobiletargeting:*:${data.aws_caller_identity.current.account_id}:apps/${aws_pinpoint_app.analytics.application_id}*"
    }
  ]
}
EOF
}
`, rName))
}

func testAccManagedUserPoolClientConfig_analyticsApplicationID(rName string) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientAnalyticsBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_cognito_managed_user_pool_client" "test" {
  name_prefix  = "AmazonOpenSearchService-%[1]s"
  user_pool_id = aws_cognito_user_pool.test.id

  depends_on = [
    aws_opensearch_domain.test,
  ]

  analytics_configuration {
    application_id = aws_pinpoint_app.analytics.application_id
    external_id    = %[1]q
    role_arn       = aws_iam_role.analytics.arn
  }
}
`, rName))
}

func testAccManagedUserPoolClientConfig_analyticsShareData(rName string) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientAnalyticsBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_cognito_managed_user_pool_client" "test" {
  name_prefix  = "AmazonOpenSearchService-%[1]s"
  user_pool_id = aws_cognito_user_pool.test.id

  depends_on = [
    aws_opensearch_domain.test,
  ]

  analytics_configuration {
    application_id   = aws_pinpoint_app.analytics.application_id
    external_id      = %[1]q
    role_arn         = aws_iam_role.analytics.arn
    user_data_shared = true
  }
}
`, rName))
}

func testAccManagedUserPoolClientConfig_analyticsARN(rName string) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientAnalyticsBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_cognito_managed_user_pool_client" "test" {
  name_prefix  = "AmazonOpenSearchService-%[1]s"
  user_pool_id = aws_cognito_user_pool.test.id

  depends_on = [
    aws_opensearch_domain.test,
  ]

  analytics_configuration {
    application_arn = aws_pinpoint_app.analytics.arn
  }
}
`, rName))
}

func testAccManagedUserPoolClientConfig_analyticsARNShareData(rName string, share bool) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientAnalyticsBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_cognito_user_pool_client" "test" {
  name_prefix  = "AmazonOpenSearchService-%[1]s"
  user_pool_id = aws_cognito_user_pool.test.id

  depends_on = [
    aws_opensearch_domain.test,
  ]

  analytics_configuration {
    application_arn  = aws_pinpoint_app.test.arn
    user_data_shared = %[2]t
  }
}
`, rName, share))
}

func testAccManagedUserPoolClientConfig_authSessionValidity(rName string, validity int) string {
	return acctest.ConfigCompose(
		testAccManagedUserPoolClientBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_cognito_managed_user_pool_client" "test" {
  name_prefix  = "AmazonOpenSearchService-%[1]s"
  user_pool_id = aws_cognito_user_pool.test.id

  depends_on = [
    aws_opensearch_domain.test,
  ]

  auth_session_validity = %[2]d
}
`, rName, validity))
}
