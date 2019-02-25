# Custom Metrics Provided

## Clusters

```
"name": "clusters/memory.limit",
"name": "clusters/memory.request",
"name": "clusters/memory.usage",
"name": "clusters/cpu.limit",
"name": "clusters/cpu.request",
"name": "clusters/cpu.usage_rate",
```

## Namespaces
```
"name": "namespaces/memory.limit",
"name": "namespaces/memory.request",
"name": "namespaces/memory.usage",
"name": "namespaces/cpu.limit",
"name": "namespaces/cpu.request",
"name": "namespaces/cpu.usage_rate",
```

## Nodes
```
"name": "nodes/uptime",

"name": "nodes/memory.request",
"name": "nodes/memory.node_capacity",
"name": "nodes/memory.node_allocatable",
"name": "nodes/memory.rss",
"name": "nodes/memory.working_set",
"name": "nodes/memory.page_faults_rate",
"name": "nodes/memory.major_page_faults",
"name": "nodes/memory.usage",
"name": "nodes/memory.limit",
"name": "nodes/memory.page_faults",
"name": "nodes/memory.major_page_faults_rate",
"name": "nodes/memory.node_reservation",
"name": "nodes/memory.node_utilization",

"name": "nodes/cpu.node_allocatable",
"name": "nodes/cpu.node_capacity",
"name": "nodes/cpu.usage",
"name": "nodes/cpu.request",
"name": "nodes/cpu.node_reservation",
"name": "nodes/cpu.usage_rate",
"name": "nodes/cpu.node_utilization",
"name": "nodes/cpu.limit",

"name": "nodes/network.rx",
"name": "nodes/network.rx_errors_rate",
"name": "nodes/network.rx_errors",
"name": "nodes/network.rx_rate",
"name": "nodes/network.tx",
"name": "nodes/network.tx_errors_rate",
"name": "nodes/network.tx_errors",
"name": "nodes/network.tx_rate",

"name": "nodes/filesystem.usage",
"name": "nodes/filesystem.inodes_free",
"name": "nodes/filesystem.limit",
"name": "nodes/filesystem.available",
"name": "nodes/filesystem.inodes",

"name": "nodes/ephemeral_storage.usage",
"name": "nodes/ephemeral_storage.request",
"name": "nodes/ephemeral_storage.node_utilization",
"name": "nodes/ephemeral_storage.limit",
"name": "nodes/ephemeral_storage.node_capacity",
"name": "nodes/ephemeral_storage.node_reservation",
"name": "nodes/ephemeral_storage.node_allocatable",
```

## Pods
```
"name": "pods/uptime",
"name": "pods/restart_count",

"name": "pods/memory.usage",
"name": "pods/memory.page_faults",
"name": "pods/memory.limit",
"name": "pods/memory.request",
"name": "pods/memory.page_faults_rate",
"name": "pods/memory.major_page_faults_rate",
"name": "pods/memory.working_set",
"name": "pods/memory.major_page_faults",
"name": "pods/memory.rss",

"name": "pods/cpu.usage",
"name": "pods/cpu.request",
"name": "pods/cpu.usage_rate",
"name": "pods/cpu.limit",

"name": "pods/network.rx",
"name": "pods/network.rx_errors_rate",
"name": "pods/network.rx_errors",
"name": "pods/network.rx_rate",
"name": "pods/network.tx",
"name": "pods/network.tx_errors_rate",
"name": "pods/network.tx_errors",
"name": "pods/network.tx_rate",

"name": "pods/filesystem.limit",
"name": "pods/filesystem.usage",
"name": "pods/filesystem.inodes",
"name": "pods/filesystem.available",
"name": "pods/filesystem.inodes_free"

"name": "pods/ephemeral_storage.usage",
"name": "pods/ephemeral_storage.request",
"name": "pods/ephemeral_storage.limit",

```

## Pod Containers
```
"name": "pod_containers/uptime",
"name": "pod_containers/restart_count",

"name": "pod_containers/memory.major_page_faults",
"name": "pod_containers/memory.working_set",
"name": "pod_containers/memory.page_faults_rate",
"name": "pod_containers/memory.usage",
"name": "pod_containers/memory.request",
"name": "pod_containers/memory.page_faults",
"name": "pod_containers/memory.rss",
"name": "pod_containers/memory.limit",
"name": "pod_containers/memory.major_page_faults_rate",

"name": "pod_containers/cpu.limit",
"name": "pod_containers/cpu.usage_rate",
"name": "pod_containers/cpu.usage",
"name": "pod_containers/cpu.request",

"name": "pod_containers/filesystem.inodes_free",
"name": "pod_containers/filesystem.limit",
"name": "pod_containers/filesystem.usage",
"name": "pod_containers/filesystem.available",
"name": "pod_containers/filesystem.inodes",

"name": "pod_containers/ephemeral_storage.request",
"name": "pod_containers/ephemeral_storage.usage",
"name": "pod_containers/ephemeral_storage.limit",
```

## Sys Containers
```
"name": "sys_containers/uptime",

"name": "sys_containers/memory.working_set",
"name": "sys_containers/memory.major_page_faults_rate",
"name": "sys_containers/memory.usage",
"name": "sys_containers/memory.rss",
"name": "sys_containers/memory.major_page_faults",
"name": "sys_containers/memory.page_faults_rate",
"name": "sys_containers/memory.page_faults",

"name": "sys_containers/cpu.usage",
"name": "sys_containers/cpu.usage_rate",
```
