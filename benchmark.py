import threading
import requests
import random
import string
import time
import queue

# Config
URLS = [
    "http://127.0.0.1:1024/Scores/",
    "http://127.0.0.1:1025/Scores/"
]
NUM_THREADS = 128
DURATION = 10  # seconds

def generate_random_key(length=8):
    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))

def worker(stop_event, request_count, error_count):
    session = requests.Session()
    while not stop_event.is_set():
        url = random.choice(URLS) + generate_random_key()
        try:
            resp = session.get(url)
            # We don't care about the result, just load
            request_count.put(1)
        except Exception as e:
            error_count.put(1)

def main():
    stop_event = threading.Event()
    request_count = queue.Queue()
    error_count = queue.Queue()
    threads = []

    print(f"Starting QPS test with {NUM_THREADS} threads for {DURATION} seconds...")
    start_time = time.time()

    for _ in range(NUM_THREADS):
        t = threading.Thread(target=worker, args=(stop_event, request_count, error_count))
        threads.append(t)
        t.start()

    try:
        time.sleep(DURATION)
    except KeyboardInterrupt:
        print("\nStopping...")
    
    stop_event.set()

    for t in threads:
        t.join()

    duration = time.time() - start_time
    total_requests = request_count.qsize()
    total_errors = error_count.qsize()
    qps = total_requests / duration

    print("\n--- Test Results ---")
    print(f"Total Requests: {total_requests}")
    print(f"Total Errors:   {total_errors}")
    print(f"Duration:       {duration:.2f} seconds")
    print(f"QPS:            {qps:.2f}")

if __name__ == "__main__":
    main()
