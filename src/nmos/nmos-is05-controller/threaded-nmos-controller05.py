#SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
#
#SPDX-License-Identifier: BSD-3-Clause

import argparse
import subprocess
import json
import requests
import threading
import time

def fetch_data(url):
    try:
        response = requests.get(url)
        if response.status_code == 200:
            print("Successful request")
            print(response.text)
            return response.text
        else:
            print("Request Failed")
            print(response.status_code)
            print(response.text)
            return ""
    except Exception as e:
        print(f"Error fetching SDP data from {url}: {e}")
        return ""

def fetch_and_save_sdp(sdp_url, filename):
    retry = 10
    for attempt in range(retry):
        try:
            response = requests.get(sdp_url, timeout=5)
            if response.status_code == 200:
                response.raise_for_status()
                content = response.text
                with open(filename, 'w', encoding='utf-8') as file:
                    file.write(content)
                break  # Exit the loop if the request is successful
            else:
                print(f"Attempt {attempt + 1}/{retry} to fetch SDP transportfile when it is ready")
        except requests.exceptions.RequestException as e:
            print(f"Attempt {attempt + 1}/{retry}: Error occurred - {e}")
        time.sleep(0.5)
    else:
        raise Exception(f"Failed to fetch SDP data from {sdp_url} after {retry} attempts")

def parse_json(data):
    try:
        return json.loads(data)
    except json.JSONDecodeError as e:
        print(f"Error parsing JSON response: {e}")
        return None
def get_id(receiver_or_sender):
    if receiver_or_sender:
        resource = parse_json(receiver_or_sender)
        if resource:
            id = resource[0].strip('/')
            return id
        else:
            print("Failed to parse data")
            return None
    else:
        print("Failed to fetch data")
        return None

def update_json_file_receiver(file_path, sender_id,sdp_url):
    # Load the existing JSON data (PATCH receiver staged endpoint)
    with open(file_path, 'r') as file:
        data = json.load(file)
    # Update the JSON data (PATCH receiver staged endpoint)
    filename="fetched_sender.sdp"
    fetch_and_save_sdp(sdp_url,filename)
    with open(filename, 'r') as spd_f:
        spd_f = spd_f.read()
    data["transport_file"]["data"] = spd_f
    data["sender_id"] = sender_id
    # Save the updated JSON data back to the file (PATCH receiver staged endpoint)
    with open(file_path, 'w') as file:
        json.dump(data, file, indent=4)

def update_json_file_sender(file_path, receiver_id):
    # Load the existing JSON data (PATCH sender staged endpoint)
    with open(file_path, 'r') as file:
        data = json.load(file)
    data['receiver_id'] = receiver_id
    # Save the updated JSON data back to the file (PATCH sender staged endpoint)
    with open(file_path, 'w') as file:
        json.dump(data, file, indent=4)

def send_patch_request(patched_file, url):
    curl_command = [
        "curl", "-X", "PATCH",
        "-H", "Content-Type: application/json",
        url,
        "--data", f"@{patched_file}"
    ]
    try:
        result = subprocess.run(curl_command, capture_output=True, text=True, check=True)
        print(result.stdout)
    except subprocess.CalledProcessError as e:
        print(f"Error executing curl command: {e}")
        print(e.output)

def process_sender(patched_file,sender_patch_url,receiver_id ):
    update_json_file_sender(patched_file, receiver_id)
    send_patch_request(patched_file, sender_patch_url)

def process_receiver(file_path,staged_receiver_url, sender_id,sdp_url):
    update_json_file_receiver(file_path, sender_id, sdp_url)
    send_patch_request(file_path, staged_receiver_url)

def main():
    parser = argparse.ArgumentParser(description="IS-04 and IS-05: Provide ip and port of sender and receiver NMOS nodes to connect.")
    parser.add_argument("--receiver_ip", required=True, help="Receiver IP address.")
    parser.add_argument("--receiver_port", type=int, required=True, help="Receiver port.")
    parser.add_argument("--sender_ip", required=True, help="Sender IP address.")
    parser.add_argument("--sender_port", type=int, required=True, help="Sender port.")
    input_connection_data = parser.parse_args()
    receiver_ip = input_connection_data.receiver_ip
    receiver_port = input_connection_data.receiver_port
    sender_ip = input_connection_data.sender_ip
    sender_port = input_connection_data.sender_port

    receiver_url = f"http://{receiver_ip}:{receiver_port}/x-nmos/connection/v1.1/single/receivers/"
    sender_url = f"http://{sender_ip}:{sender_port}/x-nmos/connection/v1.1/single/senders/"
    receiver_data = fetch_data(receiver_url)
    sender_data = fetch_data(sender_url)
    receiver_id = get_id(receiver_data)
    sender_id = get_id(sender_data)

    # logic for updating sender staged endpoint and activate
    patched_file="sender.json"
    sender_patch_url = f"http://{sender_ip}:{sender_port}/x-nmos/connection/v1.1/single/senders/{sender_id}/staged"
    sender_thread = threading.Thread(target=process_sender, args=(patched_file,sender_patch_url,receiver_id,))

    # logic for updating receiver staged including transportfile and activate
    file_path = "receiver.json"
    sdp_url = f"http://{sender_ip}:{sender_port}/x-nmos/connection/v1.1/single/senders/{sender_id}/transportfile"
    staged_receiver_url = f"http://{receiver_ip}:{receiver_port}/x-nmos/connection/v1.1/single/receivers/{receiver_id}/staged"
    receiver_thread = threading.Thread(target=process_receiver, args=(file_path, staged_receiver_url, sender_id,sdp_url,))

    sender_thread.start()
    time.sleep(2)
    receiver_thread.start()
    sender_thread.join()
    receiver_thread.join()

if __name__ == "__main__":
    main()
