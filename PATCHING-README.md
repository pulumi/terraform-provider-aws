## About

This playbook describes how to upgrade [pulumi/terraform-provider-aws](https://github.com/pulumi/terraform-provider-aws) to a new version of the upstream Terraform provider [hashicorp/terraform-provider-aws](https://github.com/hashicorp/terraform-provider-aws).

## Update playbook

When a new tag is pushed from upstream we need to re-apply our changes on top of that branch.

- The changes are a series of changes to the `internal/service`
- Each commit should be a single logical change we're maintaining
- The second-to-last commit sets up the automation
- The last commit is the application of the automation

### New Tag Procedure

1. Checkout previous patched branch (e.g. `patched-v4.38.0`)
2. Create a new patch branch (e.g. `patched-v4.39.0`):

    ```
    git checkout -b patched-v4.39.0
    ```

3. Rebase the new branch on top of the new tag:

    ```
    git rebase -i v4.39.0
    ```

4. Drop the commit "Apply automated updates", ensure the last commit is "Add automation tooling"
5. Resolve any conflicts with the rebase
6. Restore & run patcher tool:

    ```
    yarn
    yarn tf-patch apply
    ```

7. Review output and address missed replacements

- Remove or update in the `replacements.json` or `prereplacements.json` as needed.

8. Check for additional required replacements:

    ```
    yarn tf-patch check
    ```

9. If new replacements added, then search for all instances of "TODO" in the replacements file and replace with a valid alternative.
10. Re-run patcher tool once complete:

    ```
    yarn tf-patch apply
    ```

11. Stage changes to the `replacements.json` and `prereplacements.json` and amend the last commit:

    ```
    git add *.json
    git commit --amend --no-edit
    ```

12. Stage and commit all other changes:

    ```
    git add .
    git commit -m "Apply automated updates"
    ```

13. Push the completed branch to the Pulumi remote on GitHub (`git push --set-upstream pulumi patched-v...`)

### Adding new changes

We need to add a commit to our stack of patch commits. To do this, we just commit our new manual change, then rebase to sort out the ordering.

If adding a new patch _between_ new tags, just add a new commit at the end of the current `patched-` branch, then sort out the ordering during the next tag update rebase.
