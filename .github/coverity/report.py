import api
import os
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

    querry = api.prepare_querry()
    if not querry:
        log.error("failed to init querry")
        sys.exit(1)

    log.info(f"fetching snapshot for branch {branch} and commit {commit}")
    querry["snapshot"] = api.get_snapshot(querry, branch, commit)
    if querry["snapshot"] == 0:
        log.error("No snapshot found, for the branch:{branch} and commit:{commit}")
        sys.exit(1)

    log.info("spapshot found")
    log.info(f"fetching issues for snapshot from {branch}/{commit}")

    issues = api.get_snapshot_issues(querry)
    df = api.issues_to_pandas(issues)

    log.info("generating reports")
    df_grpc = api.filter_grpc_issues(df)
    df_launcher = api.filter_launcher_issues(df)
    df_grpc.to_csv("grpc_report.csv", index=False)
    df_launcher.to_csv("launcher_report.csv", index=False)
    log.info("Reports generated")


if __name__ == "__main__":
    main()
