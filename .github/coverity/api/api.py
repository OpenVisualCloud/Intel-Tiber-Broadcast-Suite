import requests
import json
from requests.auth import HTTPBasicAuth

# coverity REST API
# https://documentation.blackduck.com/bundle/coverity-docs/page/cim-api-docs/openapi/cim-openapi.html

def get_snapshot_issues(config):
  search_querry_url = f"{config['base_url']}/api/v2/issues/search"

  search_querry_params = {
      "includeColumnLabels": "true",
      "locale": "en_us",
      "offset": 0,
      "queryType": "bySnapshot",
      "rowCount": 10000,
      "sortOrder": "asc"
  }
  filters= filters=[
      {
        "columnKey": "project",
        "matchMode": "oneOrMoreMatch",
        "matchers": [
          {
            "class": "Project",
              "name": config["project_name"],
              "type": "nameMatcher"
          }
        ]
      },
      {
        "columnKey": "streams",
        "matchMode": "oneOrMoreMatch",
        "matchers": [
          {
            "class": "Stream",
            "name": config["stream"],
            "type": "nameMatcher"
          }
        ]
      }]
  if "extra_filters" in config.keys():
    filters.extend(config["extra_filters"])

  search_querry_data =  {
    "filters": filters,
    "columns": config["columns"],
    "snapshotScope": {
    "show": {
      "scope": config["snapshot"],
      "includeOutdatedSnapshots": "false"
      }
    }
  }
  # Get issues from snapshot
  response = requests.post(search_querry_url,
        params=search_querry_params, 
        json=search_querry_data, 
        auth=HTTPBasicAuth(config["user"],config["password"]))
  response.raise_for_status()
  return  json.loads(response.text)


def update_issues_triage(config):
  put_triage_attribute_querry_url = f"{config['base_url']}/issues/triage"
  put_triage_attribute_querry_params = {
    "triageStoreName": "Default Triage Store", # buil-in store
    "locale": "en_us",
  }
  put_triage_attribute_querry_data = {
    "cids": config["cids"],
    "attributeValuesList": config["attributeValuesList"]
  }

  response = requests.put(put_triage_attribute_querry_url, 
    params=put_triage_attribute_querry_params,
    json=put_triage_attribute_querry_data, 
    auth=HTTPBasicAuth(config["user"],config["password"]))
  response.raise_for_status()