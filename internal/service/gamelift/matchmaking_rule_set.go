package gamelift

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/gamelift"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceMatchmakingRuleSet() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceMatchmakingRuleSetCreate,
		ReadWithoutTimeout:   resourceMatchmakingRuleSetRead,
		UpdateWithoutTimeout: resourceMatchmakingRuleSetUpdate,
		DeleteWithoutTimeout: resourceMatchmakingRuleSetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"rule_set_body": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 65535),
					validation.StringIsJSON,
				),
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
				DiffSuppressFunc: verify.SuppressEquivalentJSONDiffs,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaTrulyComputed(),
		},
	}
}

func resourceMatchmakingRuleSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).GameLiftConn(ctx)
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(ctx, d.Get("tags").(map[string]interface{})))

	input := gamelift.CreateMatchmakingRuleSetInput{
		Name:        aws.String(d.Get("name").(string)),
		RuleSetBody: aws.String(d.Get("rule_set_body").(string)),
		Tags:        Tags(tags.IgnoreAWS()),
	}
	log.Printf("[INFO] Creating GameLift Matchmaking Rule Set: %s", input)
	out, err := conn.CreateMatchmakingRuleSet(&input)
	if err != nil {
		return diag.Errorf("error creating GameLift Matchmaking Rule Set: %s", err)
	}

	d.SetId(aws.StringValue(out.RuleSet.RuleSetName))

	return resourceMatchmakingRuleSetRead(ctx, d, meta)
}

func resourceMatchmakingRuleSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).GameLiftConn(ctx)
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig

	log.Printf("[INFO] Describing GameLift Matchmaking Rule Set: %s", d.Id())
	out, err := conn.DescribeMatchmakingRuleSets(&gamelift.DescribeMatchmakingRuleSetsInput{
		Names: aws.StringSlice([]string{d.Id()}),
	})
	if err != nil {
		if tfawserr.ErrStatusCodeEquals(err, 400) || tfawserr.ErrCodeEquals(err, gamelift.ErrCodeNotFoundException) {
			log.Printf("[WARN] GameLift Matchmaking Rule Set (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error reading GameLift Matchmaking Rule Set (%s): %s", d.Id(), err)
	}
	ruleSets := out.RuleSets

	if len(ruleSets) < 1 {
		log.Printf("[WARN] GameLift Matchmaking Rule Set (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if len(ruleSets) != 1 {
		return diag.Errorf("expected exactly 1 GameLift Matchmaking Rule Set, found %d under %q",
			len(ruleSets), d.Id())
	}
	ruleSet := ruleSets[0]

	arn := aws.StringValue(ruleSet.RuleSetArn)
	d.Set("arn", arn)
	d.Set("name", ruleSet.RuleSetName)
	d.Set("rule_set_body", ruleSet.RuleSetBody)

	tags, err := listTags(ctx, conn, arn)

	if err != nil {
		return diag.Errorf("error listing tags for GameLift Matchmaking Rule Set (%s): %s", arn, err)
	}

	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return diag.Errorf("error setting tags: %v", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return diag.Errorf("error setting tags_all: %v", err)
	}

	return nil
}

func resourceMatchmakingRuleSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).GameLiftConn(ctx)

	log.Printf("[INFO] Updating GameLift Matchmaking Rule Set: %s", d.Id())

	arn := d.Id()
	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := updateTags(ctx, conn, arn, o, n); err != nil {
			return diag.Errorf("error updating GameLift Matchmaking Rule Set (%s) tags: %s", arn, err)
		}
	}

	return resourceMatchmakingRuleSetRead(ctx, d, meta)
}

func resourceMatchmakingRuleSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).GameLiftConn(ctx)
	log.Printf("[INFO] Deleting GameLift Matchmaking Rule Set: %s", d.Id())
	_, err := conn.DeleteMatchmakingRuleSet(&gamelift.DeleteMatchmakingRuleSetInput{
		Name: aws.String(d.Id()),
	})
	if tfawserr.ErrCodeEquals(err, gamelift.ErrCodeNotFoundException) {
		return nil
	}
	if err != nil {
		return diag.Errorf("error deleting GameLift Matchmaking Rule Set (%s): %s", d.Id(), err)
	}

	return nil
}
