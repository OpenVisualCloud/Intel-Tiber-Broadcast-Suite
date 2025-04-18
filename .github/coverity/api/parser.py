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
    if raw_response_data["rows"]:
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
    else:
        df = pd.DataFrame(columns=raw_response_data["columns"])
        df.loc[0] = [" "] * len(raw_response_data["columns"])
        df.iloc[0,0] = "no issues found"
    return df

def is_report_clean(df):
    return df.iloc[0,0] == "no issues found"
def is_grpc_record(x):
    return "gRPC" in x or "grpc" in x
def is_launcher_record(x):
    return "Go" in x
def filter_grpc_issues(df): 
    if is_report_clean(df):
        return df
    else:
        return df[df["displayFile"].apply(is_grpc_record)]
    
def filter_launcher_issues(df):
    if is_report_clean(df):
        return df
    else:
        return df[df["fileLanguage"].apply(is_launcher_record)]       
