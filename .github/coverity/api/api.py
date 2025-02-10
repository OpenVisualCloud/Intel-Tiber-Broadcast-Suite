import requests
import json
from requests.auth import HTTPBasicAuth

# coverity REST API
# https://documentation.blackduck.com/bundle/coverity-docs/page/cim-api-docs/openapi/cim-openapi.html


def get_snapshot(config, description, version):
  search_querry_url = f"{config['base_url']}/api/v2/snapshots"
  raw = get_snapshots_list(config)
  ids = list(map(lambda s: s['id'],raw["snapshotsForStream"]))
  snapshot = None

  for id in ids:
    res = requests.get(f"{search_querry_url}/{id}",
        auth=HTTPBasicAuth(config["user"],config["password"]))
    res.raise_for_status()
    json_res= json.loads(res.text)
    if json_res["description"] == description and json_res["sourceVersion"] == version:
      snapshot = json_res["snapshotId"]
      break
  return snapshot

def get_snapshots_list(config):
  search_querry_url = f"{config['base_url']}/api/v2/streams/stream/snapshots"
    # Get snapshots in stream
  response = requests.get(search_querry_url,
        params={"idType": "byName", "name" : config["stream"]}, 
        auth=HTTPBasicAuth(config["user"],config["password"]))
  response.raise_for_status()
  return  json.loads(response.text)


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
  filters=[
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
