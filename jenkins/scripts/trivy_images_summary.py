import csv
import json
import glob
import os
import sys

def parse_trivy_scan_results(file_path):
    with open(file_path, 'r') as file:
        data = json.load(file)
    image_name = data['ArtifactName']
    image_name = os.path.basename(image_name).replace('.tar', '')
    sha256 = data['Metadata']['ImageID'].replace('sha256:', '')
    vulnerabilities = {'CRITICAL': 0, 'HIGH': 0}
    for result in data.get('Results',[]):
        if 'Vulnerabilities' in result:
            for vulnerability in result['Vulnerabilities']:
                severity = vulnerability['Severity']
                vulnerabilities[severity] += 1
    vulnerabilities['TOTAL'] = sum(vulnerabilities.values())
    return image_name, sha256, vulnerabilities

def create_csv_summary(directory_path, output_path):
    headers = ['Image name', 'Critical', 'High', 'Total', 'Status', 'sha256']
    rows = []
    for file_name in glob.glob(os.path.join(directory_path, 'Trivy_*.json')):
        image_name, sha256, vulnerabilities = parse_trivy_scan_results(file_name)
        status = 'PASS' if vulnerabilities['TOTAL'] == 0 else 'FAIL'
        rows.append({
            'Image name': image_name,
            'Critical': vulnerabilities['CRITICAL'],
            'High': vulnerabilities['HIGH'],
            'Total': vulnerabilities['TOTAL'],
            'Status': status,
            'sha256': sha256
        })
    rows.sort(key=lambda x: x['Image name'])

    with open(output_path, 'w', newline='') as csv_file:
        writer = csv.DictWriter(csv_file, fieldnames=headers)
        writer.writeheader()
        writer.writerows(rows)

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python trivy_images_summary.py <directory_path> <output_path>")
        sys.exit(1)
    directory = sys.argv[1]
    output_path = sys.argv[2]
    create_csv_summary(directory, output_path)
