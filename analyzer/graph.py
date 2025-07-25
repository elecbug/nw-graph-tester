import json
import networkx as nx
import matplotlib.pyplot as plt
import argparse

def parse_propagation_graph(file_path, message_id="1", layout="shell"):
    G = nx.DiGraph()
    with open(file_path, 'r') as f:
        for line in f:
            try:
                node = json.loads(line)
                dst = node["id"]
                for src in node.get("receive_map", {}).get(message_id, []):
                    G.add_edge(src, dst)
            except json.JSONDecodeError:
                continue

    print(f"Total nodes: {G.number_of_nodes()}, edges: {G.number_of_edges()}")

    plt.figure(figsize=(14, 10))
    if layout == "shell":
        pos = nx.shell_layout(G)
    elif layout == "kamada":
        pos = nx.kamada_kawai_layout(G)
    elif layout == "circular":
        pos = nx.circular_layout(G)
    else:
        pos = nx.spring_layout(G, seed=42)

    # draw without labels, smaller nodes and arrows
    nx.draw(G, pos, with_labels=False, node_size=30, node_color='skyblue',
            arrowsize=6, edge_color='gray', width=0.3)

    plt.title(f"Propagation Graph (Message ID {message_id})")
    plt.tight_layout()
    plt.savefig(f"results/propagation_graph_{message_id}_{layout}.png")
    print(f"Saved as results/propagation_graph_{message_id}_{layout}.png")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Parse and visualize propagation graph from JSONL log.")
    parser.add_argument("file", help="Path to the log file (JSONL format)")
    parser.add_argument("--msgid", default="1", help="Message ID to parse (default: 1)")
    parser.add_argument("--layout", default="shell", help="Layout type: shell, kamada, circular, spring")
    args = parser.parse_args()

    parse_propagation_graph(args.file, args.msgid, args.layout)
