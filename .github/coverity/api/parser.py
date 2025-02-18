import pandas as pd


def issues_to_pandas(raw_response_data: dict) -> pd.DataFrame:
    """
    Convert raw response data from Coverity API to a pandas DataFrame.

    Args:
      raw_response_data (dict): The raw response data from the Coverity API,
                    expected to have a "rows" key containing a list of dictionaries.

    Returns:
      pd.DataFrame: A pandas DataFrame containing the converted response data,
              sorted by the "displayImpact" column with a categorical order of ["High", "Medium", "Low"].
    """
    converted_response = []
    for row in raw_response_data["rows"]:
        converted_dict = {item["key"]: item["value"] for item in row}
        converted_response.append(converted_dict)
    df = pd.DataFrame(converted_response)
    if "displayImpact" in df.columns:
        impact_order = ["High", "Medium", "Low"]
        df["displayImpact"] = pd.Categorical(
            df["displayImpact"], categories=impact_order
        )
        df = df.sort_values(by="displayImpact")
    # remove abs path from displayFile
    df["displayFile"] = df["displayFile"].str.replace(r".*/_work/", "", regex=True)
    return df


is_grpc_record = lambda x: "gRPC" in x or "grpc" in x
is_launcher_record = lambda x: "Go" in x
filter_grpc_issues = lambda df: df[df["displayFile"].apply(is_grpc_record)]
filter_launcher_issues = lambda df: df[df["fileLanguage"].apply(is_launcher_record)]
