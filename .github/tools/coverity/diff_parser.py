import re


def read_diff_file(file_path):
    with open(file_path, 'r') as file:
        return file.read()

def parse_git_diff(diff_text):
    # Regular expression to match the diff header lines
    file_regex = re.compile(r'^diff --git a/(.+) b/(.+)$')
    hunk_regex = re.compile(r'^@@ -\d+,\d+ \+(\d+),(\d+) @@')

    changes = []
    current_file = None

    for line in diff_text.splitlines():
        file_match = file_regex.match(line)
        if file_match:
            current_file = file_match.group(2)
            changes.append({"file_changed": current_file, "lines_changed": []})
            continue

        hunk_match = hunk_regex.match(line)
        if hunk_match and current_file:
            start_line = int(hunk_match.group(1))
            line_count = int(hunk_match.group(2))
            if line_count == 0:
                continue
            end_line = start_line + line_count - 1
            changes[-1]["lines_changed"].append((start_line, end_line))

    return changes

def main():
    diff_file_path = '0001-hwupload_async.diff'  # Replace with the actual path to your diff file
    diff_text = read_diff_file(diff_file_path)
    changes = parse_git_diff(diff_text)
    
    for change in changes:
        print(f"File: {change['file_changed']}")
        print(f"Changed lines: {change['lines_changed']}")

if __name__ == "__main__":
    main()