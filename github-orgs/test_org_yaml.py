"""Validate kubeflow/org.yaml file.

The tests in this file validates:
  - org.yaml is a valid YAML file
  - All team maintainers and members appear under org admins or members
"""

import os.path as osp

import yaml

ORG_YAML = "kubeflow/org.yaml"

def test_team_member_is_in_org():
  path = osp.join(osp.dirname(__file__), ORG_YAML)
  with open(path) as stream:
    org_data = yaml.safe_load(stream)
    
    for org_name, org in org_data["orgs"].items():
      org_members = org["members"]
      org_admins = org["admins"]

      for team_name, team in org["teams"].items():
        # Verify both all groups under team
        for group_name in ["maintainers", "members"]:
          if group_name in team:
            for team_member in team[group_name]:
              assert team_member in org_members or team_member in org_admins, \
                "{} (team {}) not an admin or member of org {}".format(
                  team_member, team_name, org_name)

test_team_member_is_in_org()
