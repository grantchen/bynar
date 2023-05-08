package scope

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"

const (
	AccountIDClaimKey      = "custom:account_id"
	OrganizationIDClaimKey = "custom:organization_id"
)

type RequestScope struct {
	AccountID      int
	OrganizationID int
}

func ResolveFromToken(token string) (RequestScope, error) {
	claims, claimsErr := utils.ExtractClaimsFromJWT(token)
	if claimsErr != nil {
		return RequestScope{}, nil
	}

	accountID, accountIDErr := utils.ResolveIntClaim(AccountIDClaimKey, claims)
	if accountIDErr != nil {
		return RequestScope{}, accountIDErr
	}

	organizationID, organizationIDErr := utils.ResolveIntClaim(OrganizationIDClaimKey, claims)
	if organizationIDErr != nil {
		return RequestScope{}, organizationIDErr
	}

	return RequestScope{AccountID: accountID, OrganizationID: organizationID}, nil
}
