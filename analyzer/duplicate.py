import json
import matplotlib.pyplot as plt
import numpy as np
from collections import defaultdict
import argparse

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
    broadcast_metrics = defaultdict(lambda: {'duplicate_rates': [], 'receiving_rates': []})
    
    for entry in data:
        broadcast = entry['broadcast']
        duplicate_rate = entry['duplicate_rate']
        receiving_rate = entry['receiving_rate']
        
        broadcast_metrics[broadcast]['duplicate_rates'].append(duplicate_rate)
        broadcast_metrics[broadcast]['receiving_rates'].append(receiving_rate)
    
    avg_metrics = {}
    for broadcast, metrics in broadcast_metrics.items():
        avg_metrics[broadcast] = {
            'avg_duplicate_rate': np.mean(metrics['duplicate_rates']),
            'avg_receiving_rate': np.mean(metrics['receiving_rates'])
        }
    
    return avg_metrics

def create_graphs(avg_metrics, output_dir):
    broadcasts = sorted(avg_metrics.keys())
    
    def sort_key(broadcast):
        if broadcast == 'BasicPublish':
            return 0
        elif broadcast.startswith('WavePublish-'):
            return int(broadcast.split('-')[1])
        return 999
    
    broadcasts = sorted(broadcasts, key=sort_key)
    
    duplicate_rates = [avg_metrics[b]['avg_duplicate_rate'] for b in broadcasts]
    receiving_rates = [avg_metrics[b]['avg_receiving_rate'] for b in broadcasts]
    
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(15, 6))
    
    ax1.bar(range(len(broadcasts)), duplicate_rates, color='skyblue', alpha=0.7)
    ax1.set_xlabel('Broadcast Method')
    ax1.set_ylabel('Average Duplicate Rate (%)')
    ax1.set_title('Average Duplicate Rate by Broadcast Method')
    ax1.set_xticks(range(len(broadcasts)))
    ax1.set_xticklabels(broadcasts, rotation=45, ha='right')
    ax1.grid(axis='y', alpha=0.3)
    
    for i, v in enumerate(duplicate_rates):
        ax1.text(i, v + 0.1, f'{v:.2f}', ha='center', va='bottom', fontsize=8)
    
    ax2.bar(range(len(broadcasts)), receiving_rates, color='lightcoral', alpha=0.7)
    ax2.set_xlabel('Broadcast Method')
    ax2.set_ylabel('Average Receiving Rate (%)')
    ax2.set_title('Average Receiving Rate by Broadcast Method')
    ax2.set_xticks(range(len(broadcasts)))
    ax2.set_xticklabels(broadcasts, rotation=45, ha='right')
    ax2.grid(axis='y', alpha=0.3)
    ax2.set_ylim(0.99, 1.0005)  # Set y-axis limits for better visibility
    
    for i, v in enumerate(receiving_rates):
        ax2.text(i, v + 0.0001, f'{v:.4f}', ha='center', va='bottom', fontsize=8)
    
    plt.tight_layout()
    plt.savefig(f'{output_dir}/broadcast_metrics_comparison.png', 
                dpi=300, bbox_inches='tight')
    plt.show()

def create_combined_graph(avg_metrics, output_dir):
    broadcasts = sorted(avg_metrics.keys())
    
    def sort_key(broadcast):
        if broadcast == 'BasicPublish':
            return 0
        elif broadcast.startswith('WavePublish-'):
            return int(broadcast.split('-')[1])
        return 999
    
    broadcasts = sorted(broadcasts, key=sort_key)
    
    duplicate_rates = [avg_metrics[b]['avg_duplicate_rate'] for b in broadcasts]
    receiving_rates = [avg_metrics[b]['avg_receiving_rate'] for b in broadcasts]
    
    fig, ax1 = plt.subplots(figsize=(12, 8))
    
    color1 = 'tab:blue'
    ax1.set_xlabel('Broadcast Method')
    ax1.set_ylabel('Average Duplicate Rate (%)', color=color1)
    bars1 = ax1.bar([x - 0.2 for x in range(len(broadcasts))], duplicate_rates, 
                    width=0.4, color=color1, alpha=0.7, label='Duplicate Rate')
    ax1.tick_params(axis='y', labelcolor=color1)
    ax1.set_xticks(range(len(broadcasts)))
    ax1.set_xticklabels(broadcasts, rotation=45, ha='right')
    
    ax2 = ax1.twinx()
    color2 = 'tab:red'
    ax2.set_ylabel('Average Receiving Rate (%)', color=color2)
    bars2 = ax2.bar([x + 0.2 for x in range(len(broadcasts))], receiving_rates, 
                    width=0.4, color=color2, alpha=0.7, label='Receiving Rate')
    ax2.tick_params(axis='y', labelcolor=color2)
    ax2.set_ylim(0.99, 1.0005)

    for i, v in enumerate(duplicate_rates):
        ax1.text(i - 0.2, v + 0.1, f'{v:.2f}', ha='center', va='bottom', fontsize=8)
    
    for i, v in enumerate(receiving_rates):
        ax2.text(i + 0.2, v + 0.0001, f'{v:.4f}', ha='center', va='bottom', fontsize=8)
    
    plt.title('Broadcast Methods: Duplicate Rate vs Receiving Rate')
    plt.grid(axis='y', alpha=0.3)
    
    lines1, labels1 = ax1.get_legend_handles_labels()
    lines2, labels2 = ax2.get_legend_handles_labels()
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

    data = load_jsonl_data(data_file)
    avg_metrics = analyze_broadcast_metrics(data)

    print("Broadcast Method Analysis Results:")
    print("=" * 50)
    for broadcast, metrics in sorted(avg_metrics.items()):
        print(f"{broadcast}:")
        print(f"  Average Duplicate Rate: {metrics['avg_duplicate_rate']:.4f}%")
        print(f"  Average Receiving Rate: {metrics['avg_receiving_rate']:.6f}")
        print()

    create_graphs(avg_metrics, output_dir)
    create_combined_graph(avg_metrics, output_dir)

if __name__ == "__main__":
    main()