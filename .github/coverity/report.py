import api
import sys
import logging


def main():
    logging.basicConfig(level=logging.INFO)
    log = logging.getLogger(__name__)

    try:
        branch = sys.argv[1]
        commit = sys.argv[2]
    except IndexError:
        print("Usage: python report.py <branch> <commit>")
        sys.exit(1)

    query = api.prepare_query()
    if not query:
        log.error("failed to init query")
        sys.exit(1)

    log.info(f"fetching snapshot for branch {branch} and commit {commit}")
    query["snapshot"] = api.get_snapshot(query, branch, commit)
    if query["snapshot"] == 0:
        log.error(
            "No snapshot found, for the branch:{branch} and commit:{commit}, the analysis might not be done yet"
        )
        sys.exit(1)

    log.info("snapshot found")
    log.info(f"fetching issues for snapshot from {branch}/{commit}")

    issues = api.get_snapshot_issues(query)
    df = api.issues_to_pandas(issues)

    log.info("generating reports")
    df_grpc = api.filter_grpc_issues(df)
    df_launcher = api.filter_launcher_issues(df)
    df_grpc.to_csv("grpc_report.csv", index=False)
    df_launcher.to_csv("launcher_report.csv", index=False)

    log.info("fetching outstanding view issues")
    outstanding_issues = api.fetch_outstanding_view_issues(query)
    df_outstanding = api.issues_to_pandas(outstanding_issues)
    df_outstanding.to_csv("outstanding_issues.csv", index=False)
    log.info("Reports generated")


if __name__ == "__main__":
    main()
