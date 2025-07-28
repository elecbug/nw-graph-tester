import json
import matplotlib.pyplot as plt
import numpy as np
from collections import defaultdict
import argparse
import os

def load_jsonl_data(file_path):
    data = []
    with open(file_path, 'r') as f:
        for line in f:
            line = line.strip()
            if not line:  # Skip empty lines
                continue
            data.append(json.loads(line))
    return data

def analyze_broadcast_metrics(data):
    broadcast_metrics = defaultdict(lambda: {'duplicate_counts': [], 'receiving_rates': []})
    delay_broadcast_metrics = defaultdict(lambda: defaultdict(lambda: {'duplicate_counts': [], 'receiving_rates': []}))
    
    for entry in data:
        broadcast = entry['broadcast']
        delay = entry['delay']
        duplicate_count = entry['duplicate_rate']  # Actually duplicate count, not rate
        receiving_rate = entry['receiving_rate'] * 100  # Convert to percentage
        
        broadcast_metrics[broadcast]['duplicate_counts'].append(duplicate_count)
        broadcast_metrics[broadcast]['receiving_rates'].append(receiving_rate)
        
        delay_broadcast_metrics[delay][broadcast]['duplicate_counts'].append(duplicate_count)
        delay_broadcast_metrics[delay][broadcast]['receiving_rates'].append(receiving_rate)
    
    avg_metrics = {}
    for broadcast, metrics in broadcast_metrics.items():
        avg_metrics[broadcast] = {
            'avg_duplicate_count': np.mean(metrics['duplicate_counts']),
            'avg_receiving_rate': np.mean(metrics['receiving_rates'])
        }
    
    delay_avg_metrics = {}
    for delay, broadcast_data in delay_broadcast_metrics.items():
        delay_avg_metrics[delay] = {}
        for broadcast, metrics in broadcast_data.items():
            delay_avg_metrics[delay][broadcast] = {
                'avg_duplicate_count': np.mean(metrics['duplicate_counts']),
                'avg_receiving_rate': np.mean(metrics['receiving_rates'])
            }
    
    return avg_metrics, delay_avg_metrics

def create_graphs(avg_metrics, delay_avg_metrics, output_dir):
    # Original graphs (overall averages)
    broadcasts = sorted(avg_metrics.keys())
    
    def sort_key(broadcast):
        if broadcast == 'BasicPublish':
            return 1000
        elif broadcast.startswith('WavePublish-'):
            return int(broadcast.split('-')[1])
        return 999
    
    broadcasts = sorted(broadcasts, key=sort_key)
    
    duplicate_counts = [avg_metrics[b]['avg_duplicate_count'] for b in broadcasts]
    receiving_rates = [avg_metrics[b]['avg_receiving_rate'] for b in broadcasts]
    
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(15, 6))
    
    ax1.plot(range(len(broadcasts)), duplicate_counts, marker='o', linewidth=2, markersize=6, color='skyblue')
    ax1.set_xlabel('Broadcast Method')
    ax1.set_ylabel('Average Duplicate Count')
    ax1.set_title('Average Duplicate Count by Broadcast Method (Overall)')
    ax1.set_xticks(range(len(broadcasts)))
    ax1.set_xticklabels(broadcasts, rotation=45, ha='right')
    ax1.grid(axis='y', alpha=0.3)
    
    for i, v in enumerate(duplicate_counts):
        ax1.text(i, v + max(duplicate_counts) * 0.02, f'{v:.1f}', ha='center', va='bottom', fontsize=8)
    
    ax2.plot(range(len(broadcasts)), receiving_rates, marker='s', linewidth=2, markersize=6, color='lightcoral')
    ax2.set_xlabel('Broadcast Method')
    ax2.set_ylabel('Average Receiving Rate (%)')
    ax2.set_title('Average Receiving Rate by Broadcast Method (Overall)')
    ax2.set_xticks(range(len(broadcasts)))
    ax2.set_xticklabels(broadcasts, rotation=45, ha='right')
    ax2.grid(axis='y', alpha=0.3)
    
    # Dynamic y-axis limits for receiving rate
    if receiving_rates:
        min_receiving = min(receiving_rates)
        max_receiving = max(receiving_rates)
        margin = (max_receiving - min_receiving) * 0.1 if max_receiving > min_receiving else 0.1
        ax2.set_ylim(min_receiving - margin, max_receiving + margin)
    
    for i, v in enumerate(receiving_rates):
        margin = (max(receiving_rates) - min(receiving_rates)) * 0.01 if receiving_rates else 0.01
        ax2.text(i, v + margin, f'{v:.2f}%', ha='center', va='bottom', fontsize=8)
    
    plt.tight_layout()
    plt.savefig(f'{output_dir}/broadcast_metrics_comparison.png', 
                dpi=300, bbox_inches='tight')
    plt.show()
    
    # Delay-specific graphs
    create_delay_specific_graphs(delay_avg_metrics, output_dir)

def create_delay_specific_graphs(delay_avg_metrics, output_dir):
    delays = sorted(delay_avg_metrics.keys())
    
    # Get all broadcast methods
    all_broadcasts = set()
    for delay_data in delay_avg_metrics.values():
        all_broadcasts.update(delay_data.keys())
    
    def sort_key(broadcast):
        if broadcast == 'BasicPublish':
            return 1000
        elif broadcast.startswith('WavePublish-'):
            return int(broadcast.split('-')[1])
        return 999
    
    broadcasts = sorted(all_broadcasts, key=sort_key)
    
    # Create color map for different delays
    colors = plt.cm.tab10(np.linspace(0, 1, len(delays)))
    markers = ['o', 's', '^', 'D', 'v', '<', '>', 'p', '*', 'h']
    
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(18, 8))
    
    # Duplicate Count Graph
    for i, delay in enumerate(delays):
        duplicate_counts = []
        for broadcast in broadcasts:
            if broadcast in delay_avg_metrics[delay]:
                duplicate_counts.append(delay_avg_metrics[delay][broadcast]['avg_duplicate_count'])
            else:
                duplicate_counts.append(None)
        
        # Filter out None values for plotting
        x_vals = []
        y_vals = []
        for j, val in enumerate(duplicate_counts):
            if val is not None:
                x_vals.append(j)
                y_vals.append(val)
        
        ax1.plot(x_vals, y_vals, marker=markers[i % len(markers)], 
                linewidth=2, markersize=6, color=colors[i], 
                label=f'Delay {delay}ms', alpha=0.8)
    
    ax1.set_xlabel('Broadcast Method')
    ax1.set_ylabel('Average Duplicate Count')
    ax1.set_title('Average Duplicate Count by Broadcast Method (by Delay)')
    ax1.set_xticks(range(len(broadcasts)))
    ax1.set_xticklabels(broadcasts, rotation=45, ha='right')
    ax1.grid(axis='y', alpha=0.3)
    ax1.legend(bbox_to_anchor=(1.05, 1), loc='upper left')
    
    # Receiving Rate Graph
    for i, delay in enumerate(delays):
        receiving_rates = []
        for broadcast in broadcasts:
            if broadcast in delay_avg_metrics[delay]:
                receiving_rates.append(delay_avg_metrics[delay][broadcast]['avg_receiving_rate'])
            else:
                receiving_rates.append(None)
        
        # Filter out None values for plotting
        x_vals = []
        y_vals = []
        for j, val in enumerate(receiving_rates):
            if val is not None:
                x_vals.append(j)
                y_vals.append(val)
        
        ax2.plot(x_vals, y_vals, marker=markers[i % len(markers)], 
                linewidth=2, markersize=6, color=colors[i], 
                label=f'Delay {delay}ms', alpha=0.8)
    
    ax2.set_xlabel('Broadcast Method')
    ax2.set_ylabel('Average Receiving Rate (%)')
    ax2.set_title('Average Receiving Rate by Broadcast Method (by Delay)')
    ax2.set_xticks(range(len(broadcasts)))
    ax2.set_xticklabels(broadcasts, rotation=45, ha='right')
    ax2.grid(axis='y', alpha=0.3)
    ax2.legend(bbox_to_anchor=(1.05, 1), loc='upper left')
    
    plt.tight_layout()
    plt.savefig(f'{output_dir}/broadcast_metrics_by_delay.png', 
                dpi=300, bbox_inches='tight')
    plt.show()

def create_combined_graph(avg_metrics, delay_avg_metrics, output_dir):
    broadcasts = sorted(avg_metrics.keys())
    
    def sort_key(broadcast):
        if broadcast == 'BasicPublish':
            return 1000
        elif broadcast.startswith('WavePublish-'):
            return int(broadcast.split('-')[1])
        return 999
    
    broadcasts = sorted(broadcasts, key=sort_key)
    
    duplicate_counts = [avg_metrics[b]['avg_duplicate_count'] for b in broadcasts]
    receiving_rates = [avg_metrics[b]['avg_receiving_rate'] for b in broadcasts]
    
    print(f"DEBUG: duplicate_counts = {duplicate_counts}")
    print(f"DEBUG: receiving_rates = {receiving_rates}")
    
    fig, ax1 = plt.subplots(figsize=(12, 8))
    
    color1 = 'tab:blue'
    ax1.set_xlabel('Broadcast Method')
    ax1.set_ylabel('Average Duplicate Count', color=color1)
    line1 = ax1.plot(range(len(broadcasts)), duplicate_counts, 
                     marker='o', linewidth=2, markersize=6, color=color1, label='Duplicate Count')
    ax1.tick_params(axis='y', labelcolor=color1)
    ax1.set_xticks(range(len(broadcasts)))
    ax1.set_xticklabels(broadcasts, rotation=45, ha='right')
    
    ax2 = ax1.twinx()
    color2 = 'tab:red'
    ax2.set_ylabel('Average Receiving Rate (%)', color=color2)
    line2 = ax2.plot(range(len(broadcasts)), receiving_rates, 
                     marker='s', linewidth=2, markersize=6, color=color2, label='Receiving Rate')
    ax2.tick_params(axis='y', labelcolor=color2)
    
    # Dynamic y-axis limits for receiving rate
    min_receiving = min(receiving_rates) if receiving_rates else 99
    max_receiving = max(receiving_rates) if receiving_rates else 100
    margin = (max_receiving - min_receiving) * 0.1 if max_receiving > min_receiving else 0.1
    ax2.set_ylim(min_receiving - margin, max_receiving + margin)

    for i, v in enumerate(duplicate_counts):
        ax1.text(i, v + max(duplicate_counts) * 0.02, f'{v:.1f}', ha='center', va='bottom', fontsize=8)
    
    for i, v in enumerate(receiving_rates):
        ax2.text(i, v + margin * 0.2, f'{v:.2f}%', ha='center', va='bottom', fontsize=8)
    
    plt.title('Broadcast Methods: Duplicate Count vs Receiving Rate')
    ax1.grid(axis='y', alpha=0.3)
    
    # Combine legends from both axes
    lines1 = line1
    lines2 = line2
    labels1 = ['Duplicate Count']
    labels2 = ['Receiving Rate']
    ax1.legend(lines1 + lines2, labels1 + labels2, loc='upper left')
    
    plt.tight_layout()
    plt.savefig(f'{output_dir}/broadcast_metrics_combined.png', 
                dpi=300, bbox_inches='tight')
    plt.show()

def main():
    parser = argparse.ArgumentParser(description="Analyze broadcast metrics and generate graphs.")
    parser.add_argument('--input', type=str, required=True, help="Path to the input JSONL file.")
    parser.add_argument('--output', type=str, required=True, help="Directory to save the output graphs.")
    args = parser.parse_args()

    data_file = args.input
    output_dir = args.output

    # Create output directory if it doesn't exist
    os.makedirs(output_dir, exist_ok=True)

    # Check if input file exists
    if not os.path.exists(data_file):
        print(f"Error: Input file '{data_file}' not found.")
        return

    data = load_jsonl_data(data_file)
    
    if not data:
        print("Error: No data loaded from the input file.")
        return
    
    print(f"Loaded {len(data)} entries from {data_file}")
    
    avg_metrics, delay_avg_metrics = analyze_broadcast_metrics(data)

    if not avg_metrics:
        print("Error: No metrics calculated from the data.")
        return

    print("Broadcast Method Analysis Results:")
    print("=" * 50)
    for broadcast, metrics in sorted(avg_metrics.items()):
        print(f"{broadcast}:")
        print(f"  Average Duplicate Count: {metrics['avg_duplicate_count']:.2f}")
        print(f"  Average Receiving Rate: {metrics['avg_receiving_rate']:.4f}%")
        print()

    print(f"Found {len(delay_avg_metrics)} different delay values: {sorted(delay_avg_metrics.keys())}")
    print()

    create_graphs(avg_metrics, delay_avg_metrics, output_dir)
    create_combined_graph(avg_metrics, delay_avg_metrics, output_dir)

if __name__ == "__main__":
    main()