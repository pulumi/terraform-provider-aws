package meta

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type dataSourceAwsIPRangesResult struct {
	CreateDate   string
	Prefixes     []dataSourceAwsIPRangesPrefix
	Ipv6Prefixes []dataSourceAwsIPRangesIpv6Prefix `json:"ipv6_prefixes"`
	SyncToken    string
}

type dataSourceAwsIPRangesPrefix struct {
	IpPrefix string `json:"ip_prefix"`
	Region   string
	Service  string
}

type dataSourceAwsIPRangesIpv6Prefix struct {
	Ipv6Prefix string `json:"ipv6_prefix"`
	Region     string
	Service    string
}

func DataSourceIPRanges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIPRangesRead,

		Schema: map[string]*schema.Schema{
			"cidr_blocks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_cidr_blocks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"regions": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"services": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sync_token": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://ip-ranges.amazonaws.com/ip-ranges.json",
			},
		},
	}
}

func dataSourceIPRangesRead(d *schema.ResourceData, meta interface{}) error {

	conn := cleanhttp.DefaultClient()
	url := d.Get("url").(string)

	log.Printf("[DEBUG] Reading IP ranges from %s", url)

	res, err := conn.Get(url)

	if err != nil {
		return fmt.Errorf("Error listing IP ranges from (%s): %w", url, err)
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("Error reading response body from (%s): %w", url, err)
	}

	result := new(dataSourceAwsIPRangesResult)

	if err := json.Unmarshal(data, result); err != nil {
		return fmt.Errorf("Error parsing result from (%s): %w", url, err)
	}

	if err := d.Set("create_date", result.CreateDate); err != nil {
		return fmt.Errorf("Error setting create date: %w", err)
	}

	syncToken, err := strconv.Atoi(result.SyncToken)

	if err != nil {
		return fmt.Errorf("Error while converting sync token: %w", err)
	}

	d.SetId(result.SyncToken)

	if err := d.Set("sync_token", syncToken); err != nil {
		return fmt.Errorf("Error setting sync token: %w", err)
	}

	get := func(key string) *schema.Set {

		set := d.Get(key).(*schema.Set)

		for _, e := range set.List() {

			s := e.(string)

			set.Remove(s)
			set.Add(strings.ToLower(s))

		}

		return set

	}

	var (
		regions        = get("regions")
		services       = get("services")
		noRegionFilter = regions.Len() == 0
		ipPrefixes     []string
		ipv6Prefixes   []string
	)

	matchFilter := func(region, service string) bool {
		matchRegion := noRegionFilter || regions.Contains(strings.ToLower(region))
		matchService := services.Contains(strings.ToLower(service))
		return matchRegion && matchService
	}

	for _, e := range result.Prefixes {
		if matchFilter(e.Region, e.Service) {
			ipPrefixes = append(ipPrefixes, e.IpPrefix)
		}
	}

	for _, e := range result.Ipv6Prefixes {
		if matchFilter(e.Region, e.Service) {
			ipv6Prefixes = append(ipv6Prefixes, e.Ipv6Prefix)
		}
	}

	sort.Strings(ipPrefixes)

	if err := d.Set("cidr_blocks", ipPrefixes); err != nil {
		return fmt.Errorf("Error setting cidr_blocks: %w", err)
	}

	sort.Strings(ipv6Prefixes)

	if err := d.Set("ipv6_cidr_blocks", ipv6Prefixes); err != nil {
		return fmt.Errorf("Error setting ipv6_cidr_blocks: %w", err)
	}

	return nil

}
