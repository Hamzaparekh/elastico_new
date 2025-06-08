import pandas as pd
import matplotlib.pyplot as plt

# Load the Elastico benchmark results
csv_path = "/mnt/data/benchmark.csv"
df = pd.read_csv(csv_path)

# Calculate throughput as NumTransactions / Latency(s)
df["Throughput (tx/s)"] = df["NumTransactions"] / df["Latency(s)"]

# Aggregate average metrics by PayloadSize
agg_metrics = df.groupby("PayloadSize").agg(
    Avg_Latency=("Latency(s)", "mean"),
    Avg_CPU=("CPU(%)", "mean"),
    Avg_Memory=("Memory(MB)", "mean"),
    Avg_Throughput=("Throughput (tx/s)", "mean")
).reset_index()

# Display the aggregated metrics DataFrame
import ace_tools as tools
tools.display_dataframe_to_user(
    name="Aggregated Elastico Metrics (Multiple Transactions)",
    dataframe=agg_metrics
)

# Plotting configuration for readability
plt.rcParams.update({'figure.figsize': (8, 5)})

# Plot 1: Average Latency vs. Payload Size
plt.figure()
plt.plot(agg_metrics["PayloadSize"], agg_metrics["Avg_Latency"], marker='o')
plt.title("Elastico: Average Latency vs Payload Size\n(Multiple Transactions per Iteration)")
plt.xlabel("Payload Size (bytes)")
plt.ylabel("Latency (s)")
plt.grid(True)
plt.show()

# Plot 2: Average CPU Usage vs. Payload Size
plt.figure()
plt.plot(agg_metrics["PayloadSize"], agg_metrics["Avg_CPU"], marker='o')
plt.title("Elastico: Average CPU Usage vs Payload Size\n(Multiple Transactions per Iteration)")
plt.xlabel("Payload Size (bytes)")
plt.ylabel("CPU Usage (%)")
plt.grid(True)
plt.show()

# Plot 3: Average Memory (RSS) vs. Payload Size
plt.figure()
plt.plot(agg_metrics["PayloadSize"], agg_metrics["Avg_Memory"], marker='o')
plt.title("Elastico: Average Memory (RSS) vs Payload Size\n(Multiple Transactions per Iteration)")
plt.xlabel("Payload Size (bytes)")
plt.ylabel("Memory (MB)")
plt.grid(True)
plt.show()

# Plot 4: Average Throughput vs. Payload Size
plt.figure()
plt.plot(agg_metrics["PayloadSize"], agg_metrics["Avg_Throughput"], marker='o')
plt.title("Elastico: Average Throughput vs Payload Size\n(Multiple Transactions per Iteration)")
plt.xlabel("Payload Size (bytes)")
plt.ylabel("Throughput (tx/s)")
plt.grid(True)
plt.show()
