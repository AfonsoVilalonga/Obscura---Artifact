import sys

def bps_to_kbps(bps):
    return (bps * 8) / 1000

if len(sys.argv) != 2:
    print("Usage: python3 helper.py <filename>")
    sys.exit(1)

file_path = sys.argv[1]

try:
    with open(file_path, 'r') as f:
        lines = f.readlines()

    if len(lines) != 5:
        print("The file should have exactly 5 lines.")
        sys.exit(1)

    kbps_values = []
    for line in lines:
        line = line.strip()
        try:
            bps = float(line)
            kbps = bps_to_kbps(bps)
            kbps_values.append(kbps)
        except ValueError:
            print(f"Invalid number in line: {line}")
            sys.exit(1)

    average_kbps = sum(kbps_values) / len(kbps_values)

    print("Kbps values for each line:")
    for i, kbps in enumerate(kbps_values, start=1):
        print(f"Line {i}: {kbps:.2f} Kbps")

    print(f"\nAverage Kbps: {average_kbps:.2f} Kbps")

except FileNotFoundError:
    print(f"File '{file_path}' not found.")
    sys.exit(1)
