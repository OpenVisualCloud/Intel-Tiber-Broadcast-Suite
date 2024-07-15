import json
import csv
import sys

def count_vulnerabilities(file_path):
    with open(file_path, 'r') as file:
        data = json.load(file)
    file_vulnerabilities = {}
    file_vulnerabilities['Dockerfile'] = {'CRITICAL': 0, 'HIGH': 0}
    for result in data.get('Results',[]):
        target = result['Target']
        if 'Vulnerabilities' in result:
            if target not in file_vulnerabilities:
                file_vulnerabilities[target] = {'CRITICAL': 0, 'HIGH': 0}
            for vulnerability in result['Vulnerabilities']:
                severity = vulnerability['Severity']
                file_vulnerabilities[target][severity] += 1
    return file_vulnerabilities

def create_csv_summary(file_path, output_file):
    file_vulnerabilities = count_vulnerabilities(file_path)
    with open(output_file, 'w', newline='') as csv_file:
        fieldnames = ['File', 'CRITICAL', 'HIGH', 'TOTAL']
        writer = csv.DictWriter(csv_file, fieldnames=fieldnames)
        writer.writeheader()
        for file, counts in file_vulnerabilities.items():
            counts['TOTAL'] = sum(counts.values())
            row = {'File': file}
            row.update(counts)
            writer.writerow(row)

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python trivy_source_code_summary.py <trivy_report_file> <output_file>")
        sys.exit(1)
    report_file = sys.argv[1]
    output_file = sys.argv[2]
    create_csv_summary(report_file, output_file)
