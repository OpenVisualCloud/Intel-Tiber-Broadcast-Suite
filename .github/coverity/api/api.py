import requests
import json
from requests.auth import HTTPBasicAuth
import os
import logging

# coverity REST API
# https://documentation.blackduck.com/bundle/coverity-docs/page/cim-api-docs/openapi/cim-openapi.html


def prepare_query() -> dict:
    """
    Prepares a query dictionary for Coverity API based on environment variables.
    This function retrieves the necessary environment variables to construct a query
    dictionary for the Coverity API. If any of the required environment variables are
    not set, it logs an error and returns an empty dictionary.
    Returns:
      dict: A dictionary containing the following keys if all required environment
          variables are set:
          - "base_url": The base URL for the Coverity API.
          - "project_name": The name of the Coverity project.
          - "password": The Coverity API token.
          - "user": The Coverity user.
          - "stream": The Coverity project name (same as "project_name").
          - "columns": A list of columns to be included in the query.
          If any required environment variables are not set, returns an empty dictionary.
    """

    log = logging.getLogger(__name__)

    BASE_URL = os.getenv("COVERITY_BASE_URL", default=None)
    TOKEN = os.getenv("COVERITY_TOKEN", default=None)
    PROJECT_NAME = os.getenv("COVERITY_PROJECT_NAME", default=None)
    USER = os.getenv("COVERITY_USER", default=None)

    if None in [BASE_URL, TOKEN, PROJECT_NAME, USER]:
        log.error("Environment variables not set")
        log.error(
            "Please set the following environment variables: COVERITY_BASE_URL, COVERITY_TOKEN, COVERITY_PROJECT_NAME, COVERITY_USER"
        )
        return {}
    return {
        "base_url": BASE_URL,
        "project_name": PROJECT_NAME,
        "password": TOKEN,
        "user": USER,
        "stream": PROJECT_NAME,
        "columns": [
            "cid",
            "classification",
            "severity",
            "displayType",
            "fileLanguage",
            "displayImpact",
            "displayFile",
            "displayFunction",
            "displayCategory",
        ],
    }

def get_view_id(config: dict) -> int:
    """
    Retrieve the ID of a view based on the provided view name.

    Args:
        config (dict): Configuration dictionary containing 'base_url', 'view_name', 'user', and 'password'.

    Returns:
        int: The view ID if found, otherwise None.
    """
    endpoint = f"{config['base_url']}/api/view/v1/views"

    try:
        response = requests.get(endpoint, auth=HTTPBasicAuth(config["user"], config["password"]))
        response.raise_for_status()

        # Parse the JSON response
        views = response.json()

        # Find the view with the specified name
        for view in views:
            if view["name"] == config["view_name"]:
                return view["id"]

        logging.error(f'View with name "{config["view_name"]}" not found.')
        return None

    except requests.exceptions.HTTPError as http_err:
        logging.error(f'HTTP error occurred: {http_err}')
    except Exception as err:
        logging.error(f'An error occurred: {err}')

def fetch_view_issues(config: dict, view_id: int) -> dict:
    """
    Fetch issues from a specific view based on the provided view ID.

    Args:
        config (dict): Configuration dictionary containing 'base_url', 'user', and 'password'.
        view_id (int): The ID of the view to fetch issues from.

    Returns:
        dict: A dictionary containing the list of issues for the view.
    """
    endpoint = f"{config['base_url']}/api/view/v1/views/{view_id}/issues"
    issues ={}
    try:
        response = requests.get(endpoint, auth=HTTPBasicAuth(config["user"], config["password"]))
        response.raise_for_status()
        issues = response.json()

    except requests.exceptions.HTTPError as http_err:
        logging.error(f'HTTP error occurred: {http_err}')
    except Exception as err:
        logging.error(f'An error occurred: {err}')
    finally:
        return issues






def get_snapshot(config: dict, description: str, version: str) -> int:
    """
    Retrieve a snapshot ID based on the provided description and version.

    Args:
        config (dict): Configuration dictionary containing 'base_url', 'user', 'password', and 'stream'.
        description (str): Description of the snapshot to search for.
        version (str): Version of the snapshot to search for.

    Returns:
        str: The snapshot ID if found, otherwise None.
    """
    search_query_url = f"{config['base_url']}/api/v2/snapshots"
    raw = get_snapshots_list(config)
    ids = list(map(lambda s: s["id"], raw["snapshotsForStream"]))
    snapshot_id = 0
    #TODO: replace with threads
    for id in ids:
        res = requests.get(
            f"{search_query_url}/{id}",
            auth=HTTPBasicAuth(config["user"], config["password"]),
        )
        res.raise_for_status()
        json_res = json.loads(res.text)
        if (
            json_res["description"] == description
            and json_res["sourceVersion"] == version
        ):
            snapshot_id = json_res["snapshotId"]
            break
    return snapshot_id


def get_snapshots_list(config: dict) -> dict:
    """
    Retrieve a list of snapshots for a given stream.

    Args:
        config (dict): Configuration dictionary containing 'base_url', 'user', 'password', and 'stream'.

    Returns:
        dict: A dictionary containing the list of snapshots for the stream.
    """
    search_query_url = f"{config['base_url']}/api/v2/streams/stream/snapshots"
    # Get snapshots in stream
    response = requests.get(
        search_query_url,
        params={"idType": "byName", "name": config["stream"]},
        auth=HTTPBasicAuth(config["user"], config["password"]),
    )
    response.raise_for_status()
    return json.loads(response.text)


def get_snapshot_issues(config: dict) -> dict:
    """
    Retrieve a list of issues for a given snapshot.

    Args:
        config (dict): Configuration dictionary containing 'base_url', 'user', 'password', 'stream', 'project_name', 'columns', and 'snapshot'.
            Optional keys: 'extra_filters'.

    Returns:
        dict: A dictionary containing the list of issues for the snapshot.
    """
    search_query_url = f"{config['base_url']}/api/v2/issues/search"

    search_query_params = {
        "includeColumnLabels": "true",
        "locale": "en_us",
        "offset": 0,
        "queryType": "bySnapshot",
        "rowCount": 10000,
        "sortOrder": "asc",
    }
    filters = [
        {
            "columnKey": "project",
            "matchMode": "oneOrMoreMatch",
            "matchers": [
                {
                    "class": "Project",
                    "name": config["project_name"],
                    "type": "nameMatcher",
                }
            ],
        },
        {
            "columnKey": "streams",
            "matchMode": "oneOrMoreMatch",
            "matchers": [
                {"class": "Stream", "name": config["stream"], "type": "nameMatcher"}
            ],
        },
    ]
    if "extra_filters" in config.keys():
        filters.extend(config["extra_filters"])

    search_query_data = {
        "filters": filters,
        "columns": config["columns"],
        "snapshotScope": {
            "show": {"scope": config["snapshot"], "includeOutdatedSnapshots": "false"}
        },
    }
    # Get issues from snapshot
    response = requests.post(
        search_query_url,
        params=search_query_params,
        json=search_query_data,
        auth=HTTPBasicAuth(config["user"], config["password"]),
    )
    response.raise_for_status()
    return json.loads(response.text)
