import os
import json
import re
from collections import defaultdict
import glob

workspace_dirs = [
    ".",
    "quantasona-demo/backend",
    "sovereign-connect",
    "sovereign-kernel",
    "sovereign-mesh",
    "sovereign-mesh/amln-sen",
    "sovereign-mesh/proto",
    "sovereign-node-go",
    "sovereign-node-go/metaprogramming-go/metaprogramming-go-master"
]

manifest = []
internal_modules = set()

# First pass: map module paths to directories
mod_to_dir = {}
dir_to_mod = {}

for d in workspace_dirs:
    gomod_path = os.path.join(d, "go.mod")
    if os.path.exists(gomod_path):
        with open(gomod_path, "r") as f:
            for line in f:
                if line.startswith("module "):
                    mod_name = line.strip().split(" ")[1]
                    mod_to_dir[mod_name] = d
                    dir_to_mod[d] = mod_name
                    internal_modules.add(mod_name)
                    break

# Extract dependencies and usage
usage_counts = defaultdict(int)
cisl_links = []

for d in workspace_dirs:
    gomod_path = os.path.join(d, "go.mod")
    if not os.path.exists(gomod_path):
        continue
    
    mod_name = dir_to_mod.get(d)
    deps = []
    
    in_require = False
    with open(gomod_path, "r") as f:
        for line in f:
            line = line.strip()
            if line.startswith("require ("):
                in_require = True
                continue
            elif line.startswith("require "):
                parts = line.split(" ")
                if len(parts) >= 3:
                    dep_name = parts[1]
                    dep_ver = parts[2]
                    dep_type = "internal" if dep_name in internal_modules else "external"
                    if dep_type == "internal":
                        usage_counts[dep_name] += 1
                    deps.append({"name": dep_name, "version": dep_ver, "type": dep_type})
                continue
            elif line == ")" and in_require:
                in_require = False
                continue
            
            if in_require and line and not line.startswith("//"):
                parts = line.split(" ")
                if len(parts) >= 2:
                    dep_name = parts[0]
                    dep_ver = parts[1]
                    dep_type = "internal" if dep_name in internal_modules else "external"
                    if dep_type == "internal":
                        usage_counts[dep_name] += 1
                    deps.append({"name": dep_name, "version": dep_ver, "type": dep_type})

    # Find internal package imports to detect CISLs
    # We walk all .go files in this dir
    for root, _, files in os.walk(d):
        # Skip sub-modules
        is_submod = False
        for wd in workspace_dirs:
            if wd != d and wd != "." and root.startswith(wd):
                is_submod = True
                break
        if is_submod: continue

        for file in files:
            if file.endswith(".go"):
                file_path = os.path.join(root, file)
                with open(file_path, "r", encoding="utf-8", errors="ignore") as f:
                    content = f.read()
                    # regex to find imports
                    # single line: import "foo/bar"
                    # multiline: import ( \n "foo/bar" \n )
                    imports = re.findall(r'"(.*?)"', content)
                    for imp in imports:
                        # Check if it matches any internal module
                        for int_mod in internal_modules:
                            if imp == int_mod or imp.startswith(int_mod + "/"):
                                if int_mod != mod_name: # cross-module
                                    cisl_links.append({
                                        "source_module": mod_name,
                                        "target_module": int_mod,
                                        "imported_package": imp,
                                        "file": file_path
                                    })
                                    # Also ensure we count usage at the package level if needed
                                    usage_counts[int_mod] += 1

    manifest.append({
        "module_path": mod_name,
        "directory": d,
        "dependencies": deps,
    })

# Add usage counts and CISLs
out_manifest = {
    "modules": [],
    "critical_inter_service_links": cisl_links
}

for m in manifest:
    m["usage_count"] = usage_counts[m["module_path"]]
    out_manifest["modules"].append(m)

with open("DIAGRAM_MANIFEST.json", "w") as f:
    json.dump(out_manifest, f, indent=4)

print("Created DIAGRAM_MANIFEST.json")
