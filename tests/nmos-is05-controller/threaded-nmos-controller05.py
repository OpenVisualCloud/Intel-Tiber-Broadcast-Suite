#SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
#
#SPDX-License-Identifier: BSD-3-Clause

import argparse
import subprocess
import json
import requests
import threading
import time

# This class is used to store IS-05 Connection API details for the sender and receiver
class ConnectionDetails:
    def __init__(self, sender_destination_ip, sender_destination_port, sender_source_ip, sender_source_port, receiver_interface_ip):
        self.sender_destination_ip = sender_destination_ip
        self.sender_destination_port = sender_destination_port
        self.sender_source_ip = sender_source_ip
        self.sender_source_port = sender_source_port
        self.receiver_interface_ip = receiver_interface_ip
    

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
    
# index represents the index of the resource receiver or sender (stream) in the list
def get_id(receiver_or_sender, index=0):
    if receiver_or_sender:
        resource = parse_json(receiver_or_sender)
        if resource:
            id = resource[index].strip('/')
            return id
        else:
            print("Failed to parse data")
            return None
    else:
        print("Failed to fetch data")
        return None

def update_json_file_receiver(file_path, sender_id,sdp_url,connection_details):
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
    data['transport_params'][0]['interface_ip'] = connection_details.receiver_interface_ip
   
    # Save the updated JSON data back to the file (PATCH receiver staged endpoint)
    with open(file_path, 'w') as file:
        json.dump(data, file, indent=4)

def update_json_file_sender(file_path, receiver_id, connection_details):
    # Load the existing JSON data (PATCH sender staged endpoint)
    with open(file_path, 'r') as file:
        data = json.load(file)
    data['receiver_id'] = receiver_id
    data['transport_params'][0]['source_ip'] = connection_details.sender_source_ip
    data['transport_params'][0]['source_port'] = connection_details.sender_source_port
    data['transport_params'][0]['destination_ip'] = connection_details.sender_destination_ip
    data['transport_params'][0]['destination_port'] = connection_details.sender_destination_port
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

def process_sender(patched_file,sender_patch_url,receiver_id,connection_details):
    update_json_file_sender(patched_file, receiver_id, connection_details)
    send_patch_request(patched_file, sender_patch_url)

def process_receiver(file_path,staged_receiver_url, sender_id,sdp_url,connection_details):
    update_json_file_receiver(file_path, sender_id, sdp_url, connection_details)
    send_patch_request(file_path, staged_receiver_url)

def main():
    parser = argparse.ArgumentParser(description="IS-04 and IS-05: Provide ip and port of sender and receiver NMOS nodes to connect.")
    parser.add_argument("--receiver_ip", required=True, help="Receiver IP address.")
    parser.add_argument("--receiver_port", type=int, required=True, help="Receiver port.")
    parser.add_argument("--sender_ip", required=True, help="Sender IP address.")
    parser.add_argument("--sender_port", type=int, required=True, help="Sender port.")
    parser.add_argument("--receiver_index", type=int, default=0, help="Index of the receiver in the list.")
    parser.add_argument("--sender_index", type=int, default=0, help="Index of the sender in the list.")
    parser.add_argument("--receiver_interface_ip", help="Receiver interface IP address.")
    parser.add_argument("--sender_destination_ip", help="Sender destination IP address.")
    parser.add_argument("--sender_destination_port", type=int, help="Sender destination port.")
    parser.add_argument("--sender_source_ip", help="Sender source IP address.")
    parser.add_argument("--sender_source_port", type=int, help="Sender source port.")
    input_connection_data = parser.parse_args()
    receiver_ip = input_connection_data.receiver_ip
    receiver_port = input_connection_data.receiver_port
    sender_ip = input_connection_data.sender_ip
    sender_port = input_connection_data.sender_port
    receiver_index = input_connection_data.receiver_index
    sender_index = input_connection_data.sender_index
    
    connection_details = ConnectionDetails(
        sender_destination_ip=input_connection_data.sender_destination_ip,
        sender_destination_port=input_connection_data.sender_destination_port,
        sender_source_ip=input_connection_data.sender_source_ip,
        sender_source_port=input_connection_data.sender_source_port,
        receiver_interface_ip=input_connection_data.receiver_interface_ip
    )
    
    if not receiver_ip or not receiver_port or not sender_ip or not sender_port:
        print("Error: Receiver IP, Receiver Port, Sender IP, or Sender Port is not provided.")
        return

    if receiver_index is None or sender_index is None:
        print("Error: Receiver index or sender index is not provided.")
        return

    receiver_url = f"http://{receiver_ip}:{receiver_port}/x-nmos/connection/v1.1/single/receivers/"
    sender_url = f"http://{sender_ip}:{sender_port}/x-nmos/connection/v1.1/single/senders/"
    receiver_data = fetch_data(receiver_url)
    sender_data = fetch_data(sender_url)
    receiver_id = get_id(receiver_data, receiver_index)
    sender_id = get_id(sender_data, sender_index)

    # logic for updating sender staged endpoint and activate
    patched_file="sender.json"
    sender_patch_url = f"http://{sender_ip}:{sender_port}/x-nmos/connection/v1.1/single/senders/{sender_id}/staged"
    sender_thread = threading.Thread(target=process_sender, args=(patched_file,sender_patch_url,receiver_id,connection_details,))

    # logic for updating receiver staged including transportfile and activate
    file_path = "receiver.json"
    sdp_url = f"http://{sender_ip}:{sender_port}/x-nmos/connection/v1.1/single/senders/{sender_id}/transportfile"
    staged_receiver_url = f"http://{receiver_ip}:{receiver_port}/x-nmos/connection/v1.1/single/receivers/{receiver_id}/staged"
    receiver_thread = threading.Thread(target=process_receiver, args=(file_path, staged_receiver_url, sender_id,sdp_url,connection_details,))

    sender_thread.start()
    time.sleep(2)
    receiver_thread.start()
    sender_thread.join()
    receiver_thread.join()

if __name__ == "__main__":
    main()
