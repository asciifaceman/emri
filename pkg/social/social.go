package social

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/asciifaceman/emri/pkg/dal/models"
	"go.uber.org/zap"
)

var (
	prefix       = "https://%s"
	aboutURI     = "/api/v2/instance"
	moderatedURI = "/api/v1/instance/domain_blocks"
	peeredURI    = "/api/v1/instance/peers"
	activityURI  = "/api/v1/instance/activity"
	validDomain  = regexp.MustCompile(`^(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)?$`)
)

func buildURL(domain string, uri string) string {
	base := fmt.Sprintf(prefix, domain)
	return fmt.Sprintf("%s%s", base, uri)
}

/*
Validates the given domain is, in fact, a valid domain and exists/responds
*/
func Validate(domain string) error {
	if !validDomain.Match([]byte(domain)) {
		return fmt.Errorf("domain is not valid")
	}

	resp, err := http.Get(buildURL(domain, activityURI))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	activity := []*models.SingleActivityResponse{}
	err = json.NewDecoder(resp.Body).Decode(&activity)

	if err != nil {
		return err
	}

	clear(activity)

	return nil
}

/*
Reads /api/v2/instance and returns a *models.AboutResponse or error
*/
func About(domain string) (*models.AboutResponse, error) {
	l := zap.S().Named("mas.about")
	l.Debugw("querying domain", "domain", domain, "endpoint", aboutURI)

	resp, err := http.Get(buildURL(domain, aboutURI))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	instance := &models.AboutResponse{}
	err = json.NewDecoder(resp.Body).Decode(instance)

	return instance, err
}

/*
Reads /api/v1/instance/domain_blocks and returns a models.DomainBlockResponse or error
*/
func Moderated(domain string) (*models.DomainBlockResponse, error) {
	l := zap.S().Named("mas.moderated")
	l.Debugw("querying domain", "domain", domain, "endpoint", moderatedURI)
	return nil, nil
}

/*
Reads /api/v1/instance/peers and returns a models.PeeredResponse or error
*/
func Peered(domain string) (*models.PeeredResponse, error) {
	l := zap.S().Named("mas.peered")
	l.Debugw("querying domain", "domain", domain, "endpoint", peeredURI)
	return nil, nil
}
